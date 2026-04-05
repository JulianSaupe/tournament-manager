use crate::adapter::driven::database::Database;
use crate::adapter::driven::database::repository_error::RepositoryError;
use crate::domain::models::session::Session;
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
#[mockall::automock]
pub trait SessionRepositoryTrait: Send + Sync {
    async fn create_session(
        &self,
        user_id: Uuid,
        ip_address: Option<String>,
        user_agent: Option<String>,
        duration_hours: i64,
    ) -> Result<Session, RepositoryError>;

    async fn validate_session(&self, session_id: Uuid) -> Result<Session, RepositoryError>;

    async fn delete_session(&self, session_id: Uuid) -> Result<(), RepositoryError>;

    async fn delete_user_sessions(&self, user_id: Uuid) -> Result<(), RepositoryError>;

    async fn cleanup_expired_sessions(&self) -> Result<i64, RepositoryError>;

    async fn update_last_accessed(&self, session_id: Uuid) -> Result<(), RepositoryError>;
}

#[tonic::async_trait]
impl SessionRepositoryTrait for SessionRepository {
    async fn create_session(
        &self,
        user_id: Uuid,
        ip_address: Option<String>,
        user_agent: Option<String>,
        duration_hours: i64,
    ) -> Result<Session, RepositoryError> {
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
            .map_err(RepositoryError::from)?;

        Ok(session)
    }

    async fn validate_session(&self, session_id: Uuid) -> Result<Session, RepositoryError> {
        let session = sqlx::query_as(
            r#"
                SELECT session_id, user_id, ip_address, user_agent, created_at, expires_at, last_accessed_at
                FROM sessions
                WHERE session_id = $1 AND expires_at > NOW()
            "#
        )
            .bind(session_id)
            .fetch_one(self.database.pool())
            .await
            .map_err(RepositoryError::from)?;

        Ok(session)
    }

    async fn delete_session(&self, session_id: Uuid) -> Result<(), RepositoryError> {
        let result = sqlx::query(r#"DELETE FROM sessions WHERE session_id = $1"#)
            .bind(session_id)
            .execute(self.database.pool())
            .await
            .map_err(RepositoryError::from)?;

        if result.rows_affected() == 0 {
            return Err(RepositoryError::NotFound);
        }

        Ok(())
    }

    async fn delete_user_sessions(&self, user_id: Uuid) -> Result<(), RepositoryError> {
        let result = sqlx::query(r#"DELETE FROM sessions WHERE user_id = $1"#)
            .bind(user_id)
            .execute(self.database.pool())
            .await
            .map_err(RepositoryError::from)?;

        if result.rows_affected() == 0 {
            return Err(RepositoryError::NotFound);
        }

        Ok(())
    }

    async fn cleanup_expired_sessions(&self) -> Result<i64, RepositoryError> {
        let result = sqlx::query(r#"DELETE FROM sessions WHERE expires_at < NOW()"#)
            .execute(self.database.pool())
            .await
            .map_err(RepositoryError::from)?;

        Ok(result.rows_affected() as i64)
    }

    async fn update_last_accessed(&self, session_id: Uuid) -> Result<(), RepositoryError> {
        sqlx::query(r#"UPDATE sessions SET last_accessed_at = NOW() WHERE session_id = $1"#)
            .bind(session_id)
            .execute(self.database.pool())
            .await
            .map_err(RepositoryError::from)?;

        Ok(())
    }
}
