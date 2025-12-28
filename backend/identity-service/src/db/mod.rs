mod database;
mod repository;

pub use database::Database;
pub use repository::*;

use sqlx::postgres::{PgPool, PgPoolOptions};
use std::env;

pub type DbPool = PgPool;

pub fn get_database_url() -> Result<String, String> {
    let host = env::var("DB_HOST").unwrap_or_else(|_| "localhost".to_string());
    let port = env::var("DB_PORT").unwrap_or_else(|_| "5433".to_string());
    let user = env::var("DB_USER").unwrap_or_else(|_| "postgres".to_string());
    let password = env::var("DB_PASSWORD").unwrap_or_else(|_| "postgres".to_string());
    let name = env::var("DB_NAME").unwrap_or_else(|_| "identity-service".to_string());

    Ok(format!(
        "postgres://{}:{}@{}:{}/{}",
        user, password, host, port, name
    ))
}
