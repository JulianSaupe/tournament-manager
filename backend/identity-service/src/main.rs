use crate::db::{AccountRepository, Database};
use crate::proto::account::account_service_server::AccountServiceServer;
use crate::proto::authentication::authentication_service_server::AuthenticationServiceServer;
use crate::service::account_service::AccountService;
use crate::service::authentication_service::AuthenticationService;
use jsonwebtoken::{encode, EncodingKey, Header};
use serde::Serialize;
use tonic::{transport::Server, Status};

mod db;
mod models;
mod proto;
mod service;
mod utils;

const JWT_SECRET: &[u8] = b"secret";
const TOKEN_EXPIRATION: i64 = 10_000_000_000;

#[derive(Debug, Serialize)]
struct Claims {
    sub: String,
    exp: i64,
}

fn generate_token(username: &str) -> Result<String, Status> {
    let claims = Claims {
        sub: username.to_owned(),
        exp: TOKEN_EXPIRATION,
    };
    encode(
        &Header::default(),
        &claims,
        &EncodingKey::from_secret(JWT_SECRET),
    )
    .map_err(|_| Status::internal("Failed to create token"))
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let database_url = db::get_database_url().map_err(|e| {
        sqlx::Error::Configuration(Box::new(std::io::Error::new(
            std::io::ErrorKind::InvalidInput,
            e,
        )))
    })?;

    let database = Database::new(&database_url).await?;
    let account_repository = AccountRepository::new(database).await;

    let addr = "[::1]:5000".parse()?;
    let authentication_service = AuthenticationService::default();
    let account_service = AccountService::new(account_repository);

    println!("Server listening on {}", addr);

    Server::builder()
        .add_service(AuthenticationServiceServer::new(authentication_service))
        .add_service(AccountServiceServer::new(account_service))
        .serve(addr)
        .await?;
    Ok(())
}
