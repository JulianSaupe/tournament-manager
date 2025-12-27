use sqlx::postgres::{PgPool, PgPoolOptions};
use std::env;

pub type DbPool = PgPool;

fn get_database_url() -> Result<String, String> {
    if let Ok(url) = env::var("DATABASE_URL") {
        return Ok(url);
    }

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

pub async fn init_pool() -> Result<DbPool, sqlx::Error> {
    let database_url = get_database_url().map_err(|e| {
        sqlx::Error::Configuration(Box::new(std::io::Error::new(
            std::io::ErrorKind::InvalidInput,
            e,
        )))
    })?;

    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&database_url)
        .await?;

    sqlx::migrate!().run(&pool).await?;

    println!("Database connection pool established");
    Ok(pool)
}
