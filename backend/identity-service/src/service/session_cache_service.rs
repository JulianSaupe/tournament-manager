use crate::models::session::Session;
use chrono::Utc;
use sqlx::__rt::sleep;
use std::collections::HashMap;
use std::sync::Arc;
use std::time::Duration;
use tokio::select;
use tokio::sync::RwLock;
use tokio_util::sync::CancellationToken;

pub struct SessionCacheService {
    sessions: Arc<RwLock<HashMap<String, Session>>>,
}

#[tonic::async_trait]
pub trait SessionCacheServiceTrait: Send + Sync {
    async fn get(&self, id: String) -> Option<Session>;
    async fn set(&self, id: String, session: Session);
    async fn remove(&self, id: String);
}

impl SessionCacheService {
    pub fn new(
        cleanup_duration: Duration,
        shutdown_token: CancellationToken,
    ) -> SessionCacheService {
        let sessions = Arc::new(RwLock::new(HashMap::<String, Session>::new()));

        let cleanup_sessions = sessions.clone();

        tokio::spawn(async move {
            let token = shutdown_token.clone();
            loop {
                select! {
                    _ = token.cancelled() => break,
                    _ = sleep(cleanup_duration) => {
                        let mut sessions = cleanup_sessions.write().await;
                        sessions.retain(|_, session| session.expires_at > Utc::now());
                    }
                }
            }
        });

        SessionCacheService { sessions }
    }
}

#[tonic::async_trait]
impl SessionCacheServiceTrait for SessionCacheService {
    async fn get(&self, id: String) -> Option<Session> {
        let sessions = self.sessions.read().await;
        sessions.get(&id).cloned()
    }

    async fn set(&self, id: String, session: Session) {
        let mut sessions = self.sessions.write().await;
        sessions.insert(id, session);
    }

    async fn remove(&self, id: String) {
        let mut sessions = self.sessions.write().await;
        sessions.remove(&id);
    }
}
