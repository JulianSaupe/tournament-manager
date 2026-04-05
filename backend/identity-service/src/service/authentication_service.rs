use crate::adapter::driven::repository_error::RepositoryError;
use crate::adapter::driven::{SessionRepositoryTrait, UserRepositoryTrait};
use crate::domain::errors::service_error::ServiceError;
use crate::domain::models::session::Session;
use crate::proto::authentication::LoginResponse;
use crate::utils::verify_hash;
use prost_types::Timestamp;
use std::sync::Arc;
use tonic::{Response, Status};
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
        Self {
            user_repository,
            session_repository,
        }
    }
}

#[tonic::async_trait]
pub trait AuthenticationServiceTrait: Send + Sync {
    async fn login(
        &self,
        email: &str,
        password: &str,
        ip_address: Option<String>,
        user_agent: Option<String>,
    ) -> Result<Session, ServiceError>;
    async fn validate_session(&self, session_id: Uuid) -> Result<Session, ServiceError>;
    async fn logout(&self, session_id: Uuid) -> Result<(), ServiceError>;
}

#[tonic::async_trait]
impl AuthenticationServiceTrait for AuthenticationService {
    async fn login(
        &self,
        email: &str,
        password: &str,
        ip_address: Option<String>,
        user_agent: Option<String>,
    ) -> Result<Session, ServiceError> {
        let user_data = self.user_repository.find_by_email(email).await;

        let (user_id, password_hash) = match user_data {
            Some((id, hash)) => (id, hash),
            None => {
                return Err(ServiceError::Unauthorized(
                    "Invalid credentials".to_string(),
                ));
            }
        };

        let success = verify_hash(password, &password_hash);

        if !success {
            return Err(ServiceError::Unauthorized(
                "Invalid credentials".to_string(),
            ));
        }

        let session = self
            .session_repository
            .create_session(user_id, ip_address, user_agent, SESSION_DURATION_HOURS)
            .await
            .map_err(ServiceError::from)?;

        Ok(session)
    }

    async fn validate_session(&self, session_id: Uuid) -> Result<Session, ServiceError> {
        let session = self.session_repository.validate_session(session_id).await?;

        self.session_repository
            .update_last_accessed(session_id)
            .await
            .map_err(ServiceError::from)?;

        Ok(session)
    }

    async fn logout(&self, session_id: Uuid) -> Result<(), ServiceError> {
        self.session_repository
            .delete_session(session_id)
            .await
            .map_err(ServiceError::from)?;

        Ok(())
    }
}
