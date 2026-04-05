use crate::db::repository_error::RepositoryError;
use crate::db::{AuthorizationRepositoryTrait, RoleRepositoryTrait, UserRepositoryTrait};
use crate::proto::user::user_service_server::UserService as UserServiceTrait;
use crate::proto::user::{
    CreateRequest, CreateResponse, DeleteRequest, DeleteResponse, ResetPasswordRequest,
    ResetPasswordResponse,
};
use std::sync::Arc;
use tonic::{Request, Response, Status};
use uuid::Uuid;

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

        Ok(Response::new(CreateResponse {
            success: true,
            user_id: user_id.to_string(),
        }))
    }

    async fn delete(
        &self,
        request: Request<DeleteRequest>,
    ) -> Result<Response<DeleteResponse>, Status> {
        let delete_req = request.into_inner();

        let user_id = Uuid::parse_str(&delete_req.user_id).map_err(|_| {
            Status::invalid_argument("Failed to parse user ID: must be a valid UUID.")
        })?;

        self.user_repository
            .delete(user_id)
            .await
            .map_err(|e| match e {
                RepositoryError::NotFound => Status::not_found("User not found"),
                _ => Status::internal(format!("Failed to delete user: {}", e)),
            })?;

        Ok(Response::new(DeleteResponse { success: true }))
    }

    async fn reset_password(
        &self,
        request: Request<ResetPasswordRequest>,
    ) -> Result<Response<ResetPasswordResponse>, Status> {
        todo!()
    }
}
