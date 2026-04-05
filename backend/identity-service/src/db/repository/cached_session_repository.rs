use uuid::Uuid;
use crate::db::{SessionRepository, SessionRepositoryTrait};
use crate::models::session::Session;
use crate::service::session_cache_service::{SessionCacheService, SessionCacheServiceTrait};

pub struct CachedSessionRepository {
    repository: SessionRepository,
    cache_service: SessionCacheService,
}

impl CachedSessionRepository {
    pub fn new(repository: SessionRepository, cache_service: SessionCacheService) -> Self {
        Self { repository, cache_service }
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
    ) -> Result<Session, String> {
        let session = self.repository.create_session(user_id, ip_address, user_agent, duration_hours).await?;
        self.cache_service.set(session.session_id.to_string(), session.clone()).await;
        Ok(session)
    }

    async fn validate_session(&self, session_id: Uuid) -> Result<Option<Session>, String> {
        if let Some(session) = self.cache_service.get(session_id.to_string()).await {
            return Ok(Some(session));
        }

        self.repository.validate_session(session_id).await
    }

    async fn delete_session(&self, session_id: Uuid) -> Result<(), String> {
        self.cache_service.remove(session_id.to_string()).await;
        self.repository.delete_session(session_id).await
    }

    async fn delete_user_sessions(&self, user_id: Uuid) -> Result<i64, String> {
        self.repository.delete_user_sessions(user_id).await
    }

    async fn cleanup_expired_sessions(&self) -> Result<i64, String> {
        self.repository.cleanup_expired_sessions().await
    }

    async fn update_last_accessed(&self, session_id: Uuid) -> Result<(), String> {
        self.repository.update_last_accessed(session_id).await
    }
}