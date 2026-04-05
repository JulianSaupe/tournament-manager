use crate::adapter::driven::database::SessionRepositoryTrait;
use crate::adapter::driven::database::repository_error::RepositoryError;
use crate::domain::models::session::Session;
use crate::service::session_cache_service::SessionCacheServiceTrait;
use std::sync::Arc;
use uuid::Uuid;

pub struct CachedSessionRepository {
    repository: Arc<dyn SessionRepositoryTrait>,
    cache_service: Arc<dyn SessionCacheServiceTrait>,
}

impl CachedSessionRepository {
    pub fn new(
        repository: Arc<dyn SessionRepositoryTrait>,
        cache_service: Arc<dyn SessionCacheServiceTrait>,
    ) -> Self {
        Self {
            repository,
            cache_service,
        }
    }
}

#[tonic::async_trait]
impl SessionRepositoryTrait for CachedSessionRepository {
    async fn create_session(
        &self,
        user_id: Uuid,
        ip_address: Option<String>,
        user_agent: Option<String>,
        duration_hours: i64,
    ) -> Result<Session, RepositoryError> {
        let session = self
            .repository
            .create_session(user_id, ip_address, user_agent, duration_hours)
            .await?;
        self.cache_service
            .set(session.session_id.to_string(), session.clone())
            .await;
        Ok(session)
    }

    async fn validate_session(&self, session_id: Uuid) -> Result<Session, RepositoryError> {
        if let Some(session) = self.cache_service.get(session_id.to_string()).await {
            return Ok(session);
        }

        let session = self.repository.validate_session(session_id).await?;

        self.cache_service
            .set(session_id.to_string(), session.clone())
            .await;

        Ok(session)
    }

    async fn delete_session(&self, session_id: Uuid) -> Result<(), RepositoryError> {
        self.cache_service.remove(session_id.to_string()).await;
        self.repository.delete_session(session_id).await
    }

    async fn delete_user_sessions(&self, user_id: Uuid) -> Result<(), RepositoryError> {
        self.cache_service.remove_user_sessions(user_id).await;
        self.repository.delete_user_sessions(user_id).await
    }

    async fn cleanup_expired_sessions(&self) -> Result<i64, RepositoryError> {
        self.repository.cleanup_expired_sessions().await
    }

    async fn update_last_accessed(&self, session_id: Uuid) -> Result<(), RepositoryError> {
        self.cache_service
            .update_last_accessed(session_id.to_string())
            .await;
        self.repository.update_last_accessed(session_id).await
    }
}
