use crate::db::{SessionRepositoryTrait, UserRepositoryTrait};
use crate::proto::authentication::authentication_service_server::AuthenticationService as AuthenticationServiceTrait;
use crate::proto::authentication::{
    LoginRequest, LoginResponse, LogoutRequest, LogoutResponse, ValidateSessionRequest,
    ValidateSessionResponse,
};
use crate::utils::verify_hash;
use prost_types::Timestamp;
use std::sync::Arc;
use tonic::{Request, Response, Status};
use uuid::Uuid;

const SESSION_DURATION_HOURS: i64 = 24;

pub struct AuthenticationService {
    user_repository: Arc<dyn UserRepositoryTrait>,
    session_repository: Arc<dyn SessionRepositoryTrait>,
}

impl AuthenticationService {
    pub fn new(
        user_repository: Arc<dyn UserRepositoryTrait>,
        session_repository: Arc<dyn SessionRepositoryTrait>,
    ) -> Self {
        AuthenticationService {
            user_repository,
            session_repository,
        }
    }
}

#[tonic::async_trait]
impl AuthenticationServiceTrait for AuthenticationService {
    async fn login(
        &self,
        request: Request<LoginRequest>,
    ) -> Result<Response<LoginResponse>, Status> {
        let login_req = request.into_inner();

        let user_data = self.user_repository.find_by_email(&login_req.email).await;

        let (user_id, password_hash) = match user_data {
            Some((id, hash)) => (id, hash),
            None => {
                return Ok(Response::new(LoginResponse {
                    success: false,
                    session_id: String::new(),
                    expires_at: None,
                    message: "Invalid credentials".to_string(),
                }));
            }
        };

        let success = verify_hash(&login_req.password, &password_hash);

        if !success {
            return Ok(Response::new(LoginResponse {
                success: false,
                session_id: String::new(),
                expires_at: None,
                message: "Invalid credentials".to_string(),
            }));
        }

        let session = self
            .session_repository
            .create_session(
                user_id,
                if login_req.ip_address.is_empty() {
                    None
                } else {
                    Some(login_req.ip_address)
                },
                if login_req.user_agent.is_empty() {
                    None
                } else {
                    Some(login_req.user_agent)
                },
                SESSION_DURATION_HOURS,
            )
            .await
            .map_err(|e| Status::internal(format!("Failed to create session: {}", e)))?;

        let expires_at = Timestamp {
            seconds: session.expires_at.timestamp(),
            nanos: session.expires_at.timestamp_subsec_nanos() as i32,
        };

        let response = LoginResponse {
            success: true,
            session_id: session.session_id.to_string(),
            expires_at: Some(expires_at),
            message: "Login successful".to_string(),
        };

        Ok(Response::new(response))
    }

    async fn validate_session(
        &self,
        request: Request<ValidateSessionRequest>,
    ) -> Result<Response<ValidateSessionResponse>, Status> {
        let validate_req = request.into_inner();

        let session_id = Uuid::parse_str(&validate_req.session_id)
            .map_err(|_| Status::invalid_argument("Invalid session ID format"))?;

        let session = self
            .session_repository
            .validate_session(session_id)
            .await
            .map_err(|e| Status::internal(format!("Failed to validate session: {}", e)))?;

        match session {
            Some(sess) => {
                let _ = self
                    .session_repository
                    .update_last_accessed(session_id)
                    .await;

                let expires_at = Timestamp {
                    seconds: sess.expires_at.timestamp(),
                    nanos: sess.expires_at.timestamp_subsec_nanos() as i32,
                };

                Ok(Response::new(ValidateSessionResponse {
                    valid: true,
                    user_id: sess.user_id.to_string(),
                    expires_at: Some(expires_at),
                    message: "Session is valid".to_string(),
                }))
            }
            None => Ok(Response::new(ValidateSessionResponse {
                valid: false,
                user_id: String::new(),
                expires_at: None,
                message: "Session is invalid or expired".to_string(),
            })),
        }
    }

    async fn logout(
        &self,
        request: Request<LogoutRequest>,
    ) -> Result<Response<LogoutResponse>, Status> {
        let logout_req = request.into_inner();

        let session_id = Uuid::parse_str(&logout_req.session_id)
            .map_err(|_| Status::invalid_argument("Invalid session ID format"))?;

        self.session_repository
            .delete_session(session_id)
            .await
            .map_err(|e| Status::internal(format!("Failed to logout: {}", e)))?;

        Ok(Response::new(LogoutResponse {
            success: true,
            message: "Logout successful".to_string(),
        }))
    }
}
