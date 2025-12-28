use crate::db::AccountRepository;
use crate::proto::account::account_service_server::AccountService as AccountServiceTrait;
use crate::proto::account::{
    CreateRequest, CreateResponse, DeleteRequest, DeleteResponse, ResetPasswordRequest,
    ResetPasswordResponse,
};
use crate::utils::hash_string;
use tonic::{Request, Response, Status};

pub struct AccountService {
    repository: AccountRepository,
}

impl AccountService {
    pub fn new(repository: AccountRepository) -> Self {
        Self { repository }
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
            .repository
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

        self.repository
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
