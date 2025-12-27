use crate::proto::account::account_service_server::AccountService as AccountServiceTrait;
use crate::proto::account::{
    CreateRequest, CreateResponse, DeleteRequest, DeleteResponse, ResetPasswordRequest,
    ResetPasswordResponse,
};
use tonic::{Request, Response, Status};

#[derive(Debug, Default)]
pub struct AccountService {}

#[tonic::async_trait]
impl AccountServiceTrait for AccountService {
    async fn create(
        &self,
        request: Request<CreateRequest>,
    ) -> Result<Response<CreateResponse>, Status> {
        let create_req = request.into_inner();

        let response = CreateResponse {
            success: true,
            user_id: "".to_string(),
        };

        Ok(Response::new(response))
    }

    async fn delete(
        &self,
        request: Request<DeleteRequest>,
    ) -> Result<Response<DeleteResponse>, Status> {
        todo!()
    }

    async fn reset_password(
        &self,
        request: Request<ResetPasswordRequest>,
    ) -> Result<Response<ResetPasswordResponse>, Status> {
        todo!()
    }
}
