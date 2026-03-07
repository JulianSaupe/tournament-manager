use std::env;

pub struct Config {
    pub db: DatabaseConfig,
    pub auth: AuthConfig,
    pub server_port: u16,
}

pub struct DatabaseConfig {
    pub url: String,
}

pub struct AuthConfig {
    pub token: String,
}

impl Config {
    pub fn from_env() -> Result<Self, Box<dyn std::error::Error>> {
        if let Err(e) = dotenvy::from_path("../../docker/.identity-service.env") {
            eprintln!("Warning: Could not load .env file: {}", e);
        }

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

        let token = env::var("AUTH_TOKEN").unwrap_or_else(|_| "".to_string());

        println!("token: {}", token);
        assert!(!token.is_empty(), "AUTH_TOKEN must be set");

        Ok(Config {
            db: DatabaseConfig { url: db_url },
            auth: AuthConfig { token },
            server_port,
        })
    }
}
