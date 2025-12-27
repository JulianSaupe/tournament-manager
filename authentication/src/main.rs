use auth::auth_service_server::{AuthService, AuthServiceServer};
use auth::{LoginRequest, LoginResponse};
use jsonwebtoken::{encode, EncodingKey, Header};
use serde::Serialize;
use tonic::{transport::Server, Request, Response, Status};

pub mod auth {
    tonic::include_proto!("authentication");
}

#[derive(Debug, Default)]
pub struct MyAuthService {}

#[derive(Debug, Serialize)]
struct Claims {
    sub: String,
    exp: i64,
}

#[tonic::async_trait]
impl AuthService for MyAuthService {
    async fn login(
        &self,
        request: Request<LoginRequest>,
    ) -> Result<Response<LoginResponse>, Status> {
        println!("Got a login request: {:?}", request);

        let req = request.into_inner();

        let success = req.username == "admin" && req.password == "password";

        let message = if success {
            "Login successful".to_string()
        } else {
            "Invalid credentials".to_string()
        };

        let claims = Claims {
            sub: "admin".to_owned(),
            exp: 10000000000,
        };

        let token = encode(
            &Header::default(),
            &claims,
            &EncodingKey::from_secret("secret".as_ref()),
        )
        .map_err(|_| Status::internal("Failed to create token"))?;

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
    let auth_service = MyAuthService::default();

    println!("AuthService Server listening on {}", addr);

    Server::builder()
        .add_service(AuthServiceServer::new(auth_service))
        .serve(addr)
        .await?;

    Ok(())
}
