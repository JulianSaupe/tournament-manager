use crate::db::Database;
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
pub trait UserRepositoryTrait: Send + Sync {
    async fn create_user(
        &self,
        username: String,
        email: String,
        password: String,
    ) -> Result<Uuid, String>;

    async fn delete(&self, id: &String) -> Result<(), String>;

    async fn find_by_email_and_password(&self, email: &str) -> Option<String>;
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

    async fn delete(&self, id: &String) -> Result<(), String> {
        sqlx::query(r#"DELETE FROM users WHERE id = $1"#)
            .bind(uuid::Uuid::parse_str(id).unwrap())
            .execute(self.database.pool())
            .await
            .map_err(|_| "Failed to delete user")?;

        Ok(())
    }

    async fn find_by_email_and_password(&self, email: &str) -> Option<String> {
        let password_hash: String =
            sqlx::query_scalar(r#"SELECT password FROM users WHERE email = $1"#)
                .bind(email)
                .fetch_one(self.database.pool())
                .await
                .ok()?;

        Some(password_hash)
    }
}
