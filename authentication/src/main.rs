use auth::authentication_service_server::{AuthenticationService, AuthenticationServiceServer};
use auth::{LoginRequest, LoginResponse};
use jsonwebtoken::{encode, EncodingKey, Header};
use serde::Serialize;
use tonic::{transport::Server, Request, Response, Status};

pub mod models;
mod utils;

use models::Account;

pub mod auth {
    tonic::include_proto!("authentication");
}

const JWT_SECRET: &[u8] = b"secret";
const TOKEN_EXPIRATION: i64 = 10_000_000_000;
const ADMIN_USER: &str = "admin";
const ADMIN_PASS: &str = "password";

#[derive(Debug, Default)]
pub struct AuthService {}

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

#[tonic::async_trait]
impl AuthenticationService for AuthService {
    async fn login(
        &self,
        request: Request<LoginRequest>,
    ) -> Result<Response<LoginResponse>, Status> {
        let login_req = request.into_inner();

        let success = login_req.email == ADMIN_USER && login_req.password == ADMIN_PASS;

        let (message, token) = if success {
            (
                "Login successful".to_string(),
                generate_token(&login_req.email)?,
            )
        } else {
            ("Invalid credentials".to_string(), String::new())
        };

        let reply = LoginResponse {
            success,
            token,
            message,
        };

        Ok(Response::new(reply))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;
    let auth_service = AuthService::default();

    println!("AuthService Server listening on {}", addr);

    Server::builder()
        .add_service(AuthenticationServiceServer::new(auth_service))
        .serve(addr)
        .await?;
    Ok(())
}
