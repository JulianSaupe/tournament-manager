use crate::domain::errors::service_error::ServiceError;
use crate::proto::authentication::authentication_service_server::AuthenticationService as GrpcAuthenticationServiceTrait;
use crate::proto::authentication::{
    LoginRequest, LoginResponse, LogoutRequest, LogoutResponse, ValidateSessionRequest,
    ValidateSessionResponse,
};
use crate::service::AuthenticationServiceTrait;
use prost_types::Timestamp;
use std::sync::Arc;
use tonic::{Request, Response, Status};
use uuid::Uuid;

pub struct GrpcAuthenticationService {
    authentication_service: Arc<dyn AuthenticationServiceTrait>,
}

impl GrpcAuthenticationService {
    pub fn new(authentication_service: Arc<dyn AuthenticationServiceTrait>) -> Self {
        GrpcAuthenticationService {
            authentication_service,
        }
    }
}

#[tonic::async_trait]
impl GrpcAuthenticationServiceTrait for GrpcAuthenticationService {
    async fn login(
        &self,
        request: Request<LoginRequest>,
    ) -> Result<Response<LoginResponse>, Status> {
        let login_req = request.into_inner();

        let ip_address = if login_req.ip_address.is_empty() {
            None
        } else {
            Some(login_req.ip_address)
        };

        let user_agent = if login_req.user_agent.is_empty() {
            None
        } else {
            Some(login_req.user_agent)
        };

        let session = self
            .authentication_service
            .login(
                &login_req.email,
                &login_req.password,
                ip_address,
                user_agent,
            )
            .await
            .map_err(|e| match e {
                ServiceError::Unauthorized(_) => Status::unauthenticated("Invalid credentials"),
                _ => Status::internal(format!("Failed to login: {}", e)),
            })?;

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
            .authentication_service
            .validate_session(session_id)
            .await
            .map_err(|e| match e {
                ServiceError::NotFound(_) => Status::unauthenticated("Session not found"),
                _ => Status::internal(format!("Failed to validate session: {}", e)),
            })?;

        let expires_at = Timestamp {
            seconds: session.expires_at.timestamp(),
            nanos: session.expires_at.timestamp_subsec_nanos() as i32,
        };

        Ok(Response::new(ValidateSessionResponse {
            valid: true,
            user_id: session.user_id.to_string(),
            expires_at: Some(expires_at),
            message: "Session is valid".to_string(),
        }))
    }

    async fn logout(
        &self,
        request: Request<LogoutRequest>,
    ) -> Result<Response<LogoutResponse>, Status> {
        let logout_req = request.into_inner();

        let session_id = Uuid::parse_str(&logout_req.session_id)
            .map_err(|_| Status::invalid_argument("Invalid session ID format"))?;

        self.authentication_service
            .logout(session_id)
            .await
            .map_err(|e| match e {
                ServiceError::NotFound(_) => Status::not_found("Session not found"),
                _ => Status::internal(format!("Failed to logout: {}", e)),
            })?;

        Ok(Response::new(LogoutResponse {
            success: true,
            message: "Logout successful".to_string(),
        }))
    }
}
