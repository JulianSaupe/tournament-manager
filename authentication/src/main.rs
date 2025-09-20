use tonic::{transport::Server, Request, Response, Status};

pub mod auth {
    tonic::include_proto!("authentication");
}

use auth::auth_service_server::{AuthService, AuthServiceServer};
use auth::{LoginRequest, LoginResponse};

#[derive(Debug, Default)]
pub struct MyAuthService {}

#[tonic::async_trait]
impl AuthService for MyAuthService {
    async fn login(
        &self,
        request: Request<LoginRequest>,
    ) -> Result<Response<LoginResponse>, Status> {
        println!("Got a login request: {:?}", request);

        let req = request.into_inner();

        let success = req.username == "admin" && req.password == "password";
        let token = if success {
            "sample_jwt_token".to_string()
        } else {
            String::new()
        };
        let message = if success {
            "Login successful".to_string()
        } else {
            "Invalid credentials".to_string()
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
    let auth_service = MyAuthService::default();

    println!("AuthService Server listening on {}", addr);

    Server::builder()
        .add_service(AuthServiceServer::new(auth_service))
        .serve(addr)
        .await?;

    Ok(())
}
