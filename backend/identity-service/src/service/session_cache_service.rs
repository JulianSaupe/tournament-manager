use crate::models::session::Session;
use chrono::Utc;
use std::collections::HashMap;
use std::sync::Arc;
use std::time::Duration;
use tokio::select;
use tokio::sync::RwLock;
use tokio_util::sync::CancellationToken;

pub struct SessionCacheService {
    sessions: Arc<RwLock<HashMap<String, Session>>>,
    max_capacity: usize,
}

#[tonic::async_trait]
pub trait SessionCacheServiceTrait: Send + Sync {
    async fn get(&self, id: String) -> Option<Session>;
    async fn set(&self, id: String, session: Session);
    async fn remove(&self, id: String);
}

impl SessionCacheService {
    pub fn new(
        cleanup_interval: Duration,
        shutdown_token: CancellationToken,
        max_capacity: usize,
    ) -> SessionCacheService {
        let sessions = Arc::new(RwLock::new(HashMap::<String, Session>::new()));

        let cleanup_sessions = sessions.clone();

        tokio::spawn(async move {
            let token = shutdown_token.clone();
            let mut interval = tokio::time::interval(cleanup_interval);
            loop {
                select! {
                    _ = token.cancelled() => break,
                    _ = interval.tick() => {
                        let mut sessions = cleanup_sessions.write().await;
                        let now = Utc::now();
                        let expired_count = sessions.len();
                        sessions.retain(|_, session| session.expires_at > now);
                        let expired_count = expired_count - sessions.len();

                        if expired_count > 0 {
                            tracing::debug!("Cleaned up {} expired sessions", expired_count);
                        }
                    }
                }
            }
        });

        SessionCacheService {
            sessions,
            max_capacity,
        }
    }
}

#[tonic::async_trait]
impl SessionCacheServiceTrait for SessionCacheService {
    async fn get(&self, id: String) -> Option<Session> {
        let sessions = self.sessions.read().await;
        let session = sessions.get(&id).cloned();

        if let Some(ref s) = session {
            if s.expires_at <= Utc::now() {
                drop(sessions);
                self.remove(id).await;
                return None;
            }
        }

        session
    }

    async fn set(&self, id: String, session: Session) {
        let mut sessions = self.sessions.write().await;

        if sessions.len() >= self.max_capacity && !sessions.contains_key(&id) {
            if let Some(oldest_key) = sessions
                .iter()
                .min_by_key(|(_, s)| s.last_accessed_at)
                .map(|(k, _)| k.clone())
            {
                sessions.remove(&oldest_key);
                tracing::debug!("Evicted oldest session due to capacity limit");
            }
        }

        sessions.insert(id, session);
    }

    async fn remove(&self, id: String) {
        let mut sessions = self.sessions.write().await;
        sessions.remove(&id);
    }
}
