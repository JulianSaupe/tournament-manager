use std::env;

pub struct Config {
    pub db: DatabaseConfig,
    pub server_port: u16,
}

pub struct DatabaseConfig {
    pub url: String,
}

impl Config {
    pub fn from_env() -> Result<Self, Box<dyn std::error::Error>> {
        dotenvy::dotenv().ok();

        let host = env::var("DB_HOST").unwrap_or_else(|_| "localhost".to_string());
        let port = env::var("DB_PORT").unwrap_or_else(|_| "5433".to_string());
        let user = env::var("DB_USER").unwrap_or_else(|_| "postgres".to_string());
        let password = env::var("DB_PASSWORD").unwrap_or_else(|_| "postgres".to_string());
        let name = env::var("DB_NAME").unwrap_or_else(|_| "identity-service".to_string());

        let db_url = format!(
            "postgres://{}:{}@{}:{}/{}",
            user, password, host, port, name
        );

        let server_port = env::var("SERVER_PORT")
            .unwrap_or_else(|_| "5000".to_string())
            .parse()?;

        Ok(Config {
            db: DatabaseConfig { url: db_url },
            server_port,
        })
    }
}
