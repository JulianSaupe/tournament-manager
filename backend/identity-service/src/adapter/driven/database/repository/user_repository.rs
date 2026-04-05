use crate::adapter::driven::database::Database;
use crate::adapter::driven::database::repository_error::RepositoryError;
use crate::utils::hash_string;
use uuid::Uuid;

pub struct UserRepository {
    database: Database,
}

impl UserRepository {
    pub fn new(database: Database) -> Self {
        Self { database }
    }
}

#[tonic::async_trait]
#[mockall::automock]
pub trait UserRepositoryTrait: Send + Sync {
    async fn create_user(
        &self,
        username: String,
        email: String,
        password: String,
    ) -> Result<Uuid, String>;

    async fn delete(&self, id: Uuid) -> Result<(), RepositoryError>;

    async fn find_by_email(&self, email: &str) -> Option<(Uuid, String)>;
}

#[tonic::async_trait]
impl UserRepositoryTrait for UserRepository {
    async fn create_user(
        &self,
        username: String,
        email: String,
        password: String,
    ) -> Result<Uuid, String> {
        let id: Uuid = sqlx::query_scalar(
            r#"
                INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id
            "#,
        )
        .bind(username)
        .bind(email)
        .bind(hash_string(&password).unwrap())
        .fetch_one(self.database.pool())
        .await
        .map_err(|_| "Failed to create user")?;

        Ok(id)
    }

    async fn delete(&self, id: Uuid) -> Result<(), RepositoryError> {
        let result = sqlx::query(r#"DELETE FROM users WHERE id = $1"#)
            .bind(id)
            .execute(self.database.pool())
            .await
            .map_err(RepositoryError::from)?;

        if result.rows_affected() == 0 {
            return Err(RepositoryError::NotFound);
        }

        Ok(())
    }

    async fn find_by_email(&self, email: &str) -> Option<(Uuid, String)> {
        let result: Option<(Uuid, String)> =
            sqlx::query_as(r#"SELECT id, password FROM users WHERE email = $1"#)
                .bind(email)
                .fetch_optional(self.database.pool())
                .await
                .ok()?;

        result
    }
}
