use crate::db::{AccountRepository, AccountRepositoryTrait};
use crate::proto::account::account_service_server::AccountService as AccountServiceTrait;
use crate::proto::account::{
    CreateRequest, CreateResponse, DeleteRequest, DeleteResponse, ResetPasswordRequest,
    ResetPasswordResponse,
};
use crate::utils::hash_string;
use std::sync::Arc;
use tonic::{Request, Response, Status};

pub struct AccountService {
    account_repository: Arc<dyn AccountRepositoryTrait>,
}

impl AccountService {
    pub fn new(account_repository: Arc<dyn AccountRepositoryTrait>) -> Self {
        Self { account_repository }
    }
}

#[tonic::async_trait]
impl AccountServiceTrait for AccountService {
    async fn create(
        &self,
        request: Request<CreateRequest>,
    ) -> Result<Response<CreateResponse>, Status> {
        let create_req = request.into_inner();

        let id = self
            .account_repository
            .create_account(create_req.username, create_req.email, create_req.password)
            .await
            .map_err(|_| Status::internal("Failed to create account"))?;

        let response = CreateResponse {
            success: true,
            user_id: id.to_string(),
        };

        Ok(Response::new(response))
    }

    async fn delete(
        &self,
        request: Request<DeleteRequest>,
    ) -> Result<Response<DeleteResponse>, Status> {
        let delete_req = request.into_inner();

        self.account_repository
            .delete(&delete_req.user_id)
            .await
            .map_err(|_| Status::internal("Failed to delete account"))?;

        let response = DeleteResponse { success: true };

        Ok(Response::new(response))
    }

    async fn reset_password(
        &self,
        request: Request<ResetPasswordRequest>,
    ) -> Result<Response<ResetPasswordResponse>, Status> {
        todo!()
    }
}
