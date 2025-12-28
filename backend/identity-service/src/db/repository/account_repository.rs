use crate::db::Database;
use crate::utils::hash_string;
use tonic::async_trait;
use uuid::Uuid;

pub struct AccountRepository {
    database: Database,
}

impl AccountRepository {
    pub fn new(database: Database) -> Self {
        Self { database }
    }
}

#[async_trait]
pub trait AccountRepositoryTrait: Send + Sync {
    async fn create_account(
        &self,
        username: String,
        email: String,
        password: String,
    ) -> Result<Uuid, String>;

    async fn delete(&self, id: &String) -> Result<(), String>;

    async fn find_by_email_and_password(&self, email: &str) -> Option<String>;
}

#[async_trait]
impl AccountRepositoryTrait for AccountRepository {
    async fn create_account(
        &self,
        username: String,
        email: String,
        password: String,
    ) -> Result<Uuid, String> {
        let id: Uuid = sqlx::query_scalar(
            r#"
                INSERT INTO accounts (username, email, password) VALUES ($1, $2, $3) RETURNING id
            "#,
        )
        .bind(username)
        .bind(email)
        .bind(hash_string(&password).unwrap())
        .fetch_one(self.database.pool())
        .await
        .map_err(|_| "Failed to create account")?;

        Ok(id)
    }

    async fn delete(&self, id: &String) -> Result<(), String> {
        sqlx::query(r#"DELETE FROM accounts WHERE id = $1"#)
            .bind(uuid::Uuid::parse_str(id).unwrap())
            .execute(self.database.pool())
            .await
            .map_err(|_| "Failed to delete account")?;

        Ok(())
    }

    async fn find_by_email_and_password(&self, email: &str) -> Option<String> {
        let password_hash: String =
            sqlx::query_scalar(r#"SELECT password FROM accounts WHERE email = $1"#)
                .bind(email)
                .fetch_one(self.database.pool())
                .await
                .ok()?;

        Some(password_hash)
    }
}
