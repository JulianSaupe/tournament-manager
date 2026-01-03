use crate::db::Database;
use crate::models::session::Session;
use chrono::{Duration, Utc};
use uuid::Uuid;

pub struct SessionRepository {
    database: Database,
}

impl SessionRepository {
    pub fn new(database: Database) -> Self {
        Self { database }
    }
}

#[tonic::async_trait]
pub trait SessionRepositoryTrait: Send + Sync {
    async fn create_session(
        &self,
        user_id: Uuid,
        ip_address: Option<String>,
        user_agent: Option<String>,
        duration_hours: i64,
    ) -> Result<Session, String>;

    async fn validate_session(&self, session_id: Uuid) -> Result<Option<Session>, String>;

    async fn delete_session(&self, session_id: Uuid) -> Result<(), String>;

    async fn delete_user_sessions(&self, user_id: Uuid) -> Result<i64, String>;

    async fn cleanup_expired_sessions(&self) -> Result<i64, String>;

    async fn update_last_accessed(&self, session_id: Uuid) -> Result<(), String>;
}

#[tonic::async_trait]
impl SessionRepositoryTrait for SessionRepository {
    async fn create_session(
        &self,
        user_id: Uuid,
        ip_address: Option<String>,
        user_agent: Option<String>,
        duration_hours: i64,
    ) -> Result<Session, String> {
        let expires_at = Utc::now() + Duration::hours(duration_hours);

        let session: Session = sqlx::query_as(
            r#"
                INSERT INTO sessions (user_id, ip_address, user_agent, expires_at)
                VALUES ($1, $2, $3, $4)
                RETURNING session_id, user_id, ip_address, user_agent, created_at, expires_at, last_accessed_at
            "#
        )
            .bind(user_id)
            .bind(ip_address)
            .bind(user_agent)
            .bind(expires_at)
            .fetch_one(self.database.pool())
            .await
            .map_err(|e| format!("Failed to create session: {}", e))?;

        Ok(session)
    }

    async fn validate_session(&self, session_id: Uuid) -> Result<Option<Session>, String> {
        let session = sqlx::query_as(
            r#"
                SELECT session_id, user_id, ip_address, user_agent, created_at, expires_at, last_accessed_at
                FROM sessions
                WHERE session_id = $1 AND expires_at > NOW()
            "#
        )
            .bind(session_id)
            .fetch_optional(self.database.pool())
            .await
            .map_err(|e| format!("Failed to validate session: {}", e))?;

        Ok(session)
    }

    async fn delete_session(&self, session_id: Uuid) -> Result<(), String> {
        sqlx::query(r#"DELETE FROM sessions WHERE session_id = $1"#)
            .bind(session_id)
            .execute(self.database.pool())
            .await
            .map_err(|e| format!("Failed to delete session: {}", e))?;

        Ok(())
    }

    async fn delete_user_sessions(&self, user_id: Uuid) -> Result<i64, String> {
        let result = sqlx::query(r#"DELETE FROM sessions WHERE user_id = $1"#)
            .bind(user_id)
            .execute(self.database.pool())
            .await
            .map_err(|e| format!("Failed to delete user sessions: {}", e))?;

        Ok(result.rows_affected() as i64)
    }

    async fn cleanup_expired_sessions(&self) -> Result<i64, String> {
        let result = sqlx::query(r#"DELETE FROM sessions WHERE expires_at < NOW()"#)
            .execute(self.database.pool())
            .await
            .map_err(|e| format!("Failed to cleanup expired sessions: {}", e))?;

        Ok(result.rows_affected() as i64)
    }

    async fn update_last_accessed(&self, session_id: Uuid) -> Result<(), String> {
        sqlx::query(r#"UPDATE sessions SET last_accessed_at = NOW() WHERE session_id = $1"#)
            .bind(session_id)
            .execute(self.database.pool())
            .await
            .map_err(|e| format!("Failed to update last accessed time: {}", e))?;

        Ok(())
    }
}
