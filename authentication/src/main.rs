use crate::proto::authentication::authentication_service_server::AuthenticationServiceServer;
use crate::service::authentication::AuthenticationService;
use jsonwebtoken::{encode, EncodingKey, Header};
use serde::Serialize;
use tonic::{transport::Server, Status};

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
    let addr = "[::1]:50051".parse()?;
    let auth_service = AuthenticationService::default();

    println!("AuthService Server listening on {}", addr);

    Server::builder()
        .add_service(AuthenticationServiceServer::new(auth_service))
        .serve(addr)
        .await?;
    Ok(())
}
