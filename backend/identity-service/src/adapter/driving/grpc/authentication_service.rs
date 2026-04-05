use crate::adapter::driven::database::repository_error::RepositoryError;
use crate::adapter::driven::database::{SessionRepositoryTrait, UserRepositoryTrait};
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
            .map_err(|e| match e {
                RepositoryError::NotFound => Status::unauthenticated("Session not found"),
                _ => Status::internal(format!("Failed to validate session: {}", e)),
            })?;

        let _ = self
            .session_repository
            .update_last_accessed(session_id)
            .await;

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

#[cfg(test)]
mod tests {
    use super::*;
    use crate::adapter::driven::repository::session_repository::MockSessionRepositoryTrait;
    use crate::adapter::driven::repository::user_repository::MockUserRepositoryTrait;
    use crate::models::session::Session;
    use chrono::Utc;
    use mockall::predicate::*;

    #[tokio::test]
    async fn test_login_success() {
        let mut mock_user_repo = MockUserRepositoryTrait::new();
        let mut mock_session_repo = MockSessionRepositoryTrait::new();

        let user_id = Uuid::now_v7();
        let email = "test@example.com".to_string();
        let password = "password123".to_string();
        let password_hash = crate::utils::hash_string(&password).unwrap();

        let email_clone = email.clone();
        mock_user_repo
            .expect_find_by_email()
            .with(eq(email_clone))
            .times(1)
            .returning(move |_| {
                let password_hash = password_hash.clone();
                Box::pin(async move { Some((user_id, password_hash)) })
            });

        let session_id = Uuid::now_v7();
        let now = Utc::now();
        mock_session_repo
            .expect_create_session()
            .with(eq(user_id), eq(None), eq(None), eq(SESSION_DURATION_HOURS))
            .times(1)
            .returning(move |u_id, _, _, _| {
                Box::pin(async move {
                    Ok(Session {
                        session_id,
                        user_id: u_id,
                        ip_address: None,
                        user_agent: None,
                        created_at: now,
                        expires_at: now + chrono::Duration::hours(SESSION_DURATION_HOURS),
                        last_accessed_at: now,
                    })
                })
            });

        let service =
            AuthenticationService::new(Arc::new(mock_user_repo), Arc::new(mock_session_repo));

        let request = Request::new(LoginRequest {
            email,
            password,
            ip_address: "".to_string(),
            user_agent: "".to_string(),
        });

        let response = service.login(request).await.unwrap();
        let inner = response.into_inner();

        assert!(inner.success);
        assert_eq!(inner.session_id, session_id.to_string());
        assert_eq!(inner.message, "Login successful");
    }

    #[tokio::test]
    async fn test_validate_session_success() {
        let mock_user_repo = MockUserRepositoryTrait::new();
        let mut mock_session_repo = MockSessionRepositoryTrait::new();

        let session_id = Uuid::now_v7();
        let user_id = Uuid::now_v7();
        let now = Utc::now();

        mock_session_repo
            .expect_validate_session()
            .with(eq(session_id))
            .times(1)
            .returning(move |s_id| {
                Box::pin(async move {
                    Ok(Session {
                        session_id: s_id,
                        user_id,
                        ip_address: None,
                        user_agent: None,
                        created_at: now,
                        expires_at: now + chrono::Duration::hours(1),
                        last_accessed_at: now,
                    })
                })
            });

        mock_session_repo
            .expect_update_last_accessed()
            .with(eq(session_id))
            .times(1)
            .returning(|_| Box::pin(async move { Ok(()) }));

        let service =
            AuthenticationService::new(Arc::new(mock_user_repo), Arc::new(mock_session_repo));

        let request = Request::new(ValidateSessionRequest {
            session_id: session_id.to_string(),
        });

        let response = service.validate_session(request).await.unwrap();
        let inner = response.into_inner();

        assert!(inner.valid);
        assert_eq!(inner.user_id, user_id.to_string());
        assert_eq!(inner.message, "Session is valid");
    }

    #[tokio::test]
    async fn test_logout_success() {
        let mock_user_repo = MockUserRepositoryTrait::new();
        let mut mock_session_repo = MockSessionRepositoryTrait::new();

        let session_id = Uuid::now_v7();

        mock_session_repo
            .expect_delete_session()
            .with(eq(session_id))
            .times(1)
            .returning(|_| Box::pin(async move { Ok(()) }));

        let service =
            AuthenticationService::new(Arc::new(mock_user_repo), Arc::new(mock_session_repo));

        let request = Request::new(LogoutRequest {
            session_id: session_id.to_string(),
        });

        let response = service.logout(request).await.unwrap();
        let inner = response.into_inner();

        assert!(inner.success);
        assert_eq!(inner.message, "Logout successful");
    }
}
