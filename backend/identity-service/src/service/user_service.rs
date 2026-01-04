use crate::db::{AuthorizationRepositoryTrait, RoleRepositoryTrait, UserRepositoryTrait};
use crate::proto::user::user_service_server::UserService as UserServiceTrait;
use crate::proto::user::{
    CreateRequest, CreateResponse, DeleteRequest, DeleteResponse, ResetPasswordRequest,
    ResetPasswordResponse,
};
use std::sync::Arc;
use tonic::{Request, Response, Status};

pub struct UserService {
    user_repository: Arc<dyn UserRepositoryTrait>,
    authorization_repository: Arc<dyn AuthorizationRepositoryTrait>,
    role_repository: Arc<dyn RoleRepositoryTrait>,
}

impl UserService {
    pub fn new(
        user_repository: Arc<dyn UserRepositoryTrait>,
        authorization_repository: Arc<dyn AuthorizationRepositoryTrait>,
        role_repository: Arc<dyn RoleRepositoryTrait>,
    ) -> Self {
        Self {
            user_repository,
            authorization_repository,
            role_repository,
        }
    }
}

#[tonic::async_trait]
impl UserServiceTrait for UserService {
    async fn create(
        &self,
        request: Request<CreateRequest>,
    ) -> Result<Response<CreateResponse>, Status> {
        let create_req = request.into_inner();

        let user_id = self
            .user_repository
            .create_user(create_req.username, create_req.email, create_req.password)
            .await
            .map_err(|_| Status::internal("Failed to create user"))?;

        let role_id = self
            .role_repository
            .get_role_by_name("user")
            .await
            .unwrap()
            .id;

        self.authorization_repository
            .assign_role(user_id, role_id)
            .await
            .unwrap();

        let response = CreateResponse {
            success: true,
            user_id: user_id.to_string(),
        };

        Ok(Response::new(response))
    }

    async fn delete(
        &self,
        request: Request<DeleteRequest>,
    ) -> Result<Response<DeleteResponse>, Status> {
        let delete_req = request.into_inner();

        self.user_repository
            .delete(&delete_req.user_id)
            .await
            .map_err(|_| Status::internal("Failed to delete user"))?;

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
