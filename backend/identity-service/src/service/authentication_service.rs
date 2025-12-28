use crate::db::{AccountRepository, AccountRepositoryTrait};
use crate::proto::authentication::authentication_service_server::AuthenticationService as AuthenticationServiceTrait;
use crate::proto::authentication::{
    ExtendLifetimeRequest, ExtendLifetimeResponse, LoginRequest, LoginResponse,
};
use crate::utils::{generate_token, hash_string, verify_hash};
use std::sync::Arc;
use tonic::{Request, Response, Status};

pub struct AuthenticationService {
    account_repository: Arc<dyn AccountRepositoryTrait>,
}

impl AuthenticationService {
    pub fn new(account_repository: Arc<dyn AccountRepositoryTrait>) -> Self {
        AuthenticationService { account_repository }
    }
}

#[tonic::async_trait]
impl AuthenticationServiceTrait for AuthenticationService {
    async fn login(
        &self,
        request: Request<LoginRequest>,
    ) -> Result<Response<LoginResponse>, Status> {
        let login_req = request.into_inner();

        let password_hash = self
            .account_repository
            .find_by_email_and_password(&login_req.email)
            .await;

        let success =
            password_hash.is_some() && verify_hash(&login_req.password, &password_hash.unwrap());

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
