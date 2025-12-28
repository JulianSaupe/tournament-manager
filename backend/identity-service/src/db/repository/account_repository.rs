use crate::db::Database;
use crate::utils::hash_string;
use tonic::Status;
use uuid::Uuid;

pub struct AccountRepository {
    database: Database,
}

impl AccountRepository {
    pub async fn new(database: Database) -> Self {
        AccountRepository { database }
    }

    pub async fn create_account(
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

    pub async fn delete(&self, id: &String) -> Result<(), String> {
        sqlx::query(r#"DELETE FROM accounts WHERE id = $1"#)
            .bind(uuid::Uuid::parse_str(id).unwrap())
            .execute(self.database.pool())
            .await
            .map_err(|_| "Failed to delete account")?;

        Ok(())
    }
}
