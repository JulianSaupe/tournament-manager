use crate::generate_token;
use crate::proto::authentication::authentication_service_server::AuthenticationService as AuthenticationServiceTrait;
use crate::proto::authentication::{
    ExtendLifetimeRequest, ExtendLifetimeResponse, LoginRequest, LoginResponse,
};
use tonic::{Request, Response, Status};

const ADMIN_USER: &str = "admin";
const ADMIN_PASS: &str = "password";

#[derive(Debug, Default)]
pub struct AuthenticationService {}

#[tonic::async_trait]
impl AuthenticationServiceTrait for AuthenticationService {
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

        let response = LoginResponse {
            success,
            token,
            message,
        };

        Ok(Response::new(response))
    }

    async fn extend_lifetime(
        &self,
        request: Request<ExtendLifetimeRequest>,
    ) -> Result<Response<ExtendLifetimeResponse>, Status> {
        todo!()
    }
}
