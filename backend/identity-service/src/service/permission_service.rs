use crate::db::AuthorizationRepositoryTrait;
use crate::proto::authorization::permission_service_server::PermissionService as PermissionServiceTrait;
use crate::proto::authorization::{
    CreatePermissionRequest, CreatePermissionResponse, DeletePermissionRequest,
    DeletePermissionResponse, GetPermissionRequest, GetPermissionResponse, ListPermissionsRequest,
    ListPermissionsResponse, UpdatePermissionRequest, UpdatePermissionResponse,
};
use std::sync::Arc;
use tonic::{Request, Response, Status};

pub struct PermissionService {
    authorization_repository: Arc<dyn AuthorizationRepositoryTrait>,
}

impl PermissionService {
    pub fn new(authorization_repository: Arc<dyn AuthorizationRepositoryTrait>) -> Self {
        Self {
            authorization_repository,
        }
    }
}

#[tonic::async_trait]
impl PermissionServiceTrait for PermissionService {
    async fn create_permission(
        &self,
        request: Request<CreatePermissionRequest>,
    ) -> Result<Response<CreatePermissionResponse>, Status> {
        let permission_req = request.into_inner();

        let permission_id = self
            .authorization_repository
            .create_permission(&permission_req.name)
            .await
            .map_err(|e| Status::internal("Failed to create permission:"))?;

        Ok(Response::new(CreatePermissionResponse {
            permission_id: permission_id.to_string(),
            success: true,
        }))
    }

    async fn get_permission(
        &self,
        request: Request<GetPermissionRequest>,
    ) -> Result<Response<GetPermissionResponse>, Status> {
        todo!()
    }

    async fn list_permissions(
        &self,
        request: Request<ListPermissionsRequest>,
    ) -> Result<Response<ListPermissionsResponse>, Status> {
        todo!()
    }

    async fn update_permission(
        &self,
        request: Request<UpdatePermissionRequest>,
    ) -> Result<Response<UpdatePermissionResponse>, Status> {
        todo!()
    }

    async fn delete_permission(
        &self,
        request: Request<DeletePermissionRequest>,
    ) -> Result<Response<DeletePermissionResponse>, Status> {
        todo!()
    }
}
