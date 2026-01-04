use crate::db::PermissionRepositoryTrait;
use crate::proto::authorization::permission_service_server::PermissionService as PermissionServiceTrait;
use crate::proto::authorization::{
    CreatePermissionRequest, CreatePermissionResponse, DeletePermissionRequest,
    DeletePermissionResponse, GetPermissionByNameRequest, GetPermissionResponse,
    ListPermissionsRequest, ListPermissionsResponse, UpdatePermissionRequest,
    UpdatePermissionResponse,
};
use std::sync::Arc;
use tonic::{Request, Response, Status};

pub struct PermissionService {
    permission_repository: Arc<dyn PermissionRepositoryTrait>,
}

impl PermissionService {
    pub fn new(permission_repository: Arc<dyn PermissionRepositoryTrait>) -> Self {
        Self {
            permission_repository,
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
            .permission_repository
            .create_permission(&permission_req.name)
            .await
            .map_err(|e| Status::internal("Failed to create permission:"))?;

        Ok(Response::new(CreatePermissionResponse {
            permission_id: permission_id.to_string(),
            success: true,
        }))
    }

    async fn get_permission_by_name(
        &self,
        request: Request<GetPermissionByNameRequest>,
    ) -> Result<Response<GetPermissionResponse>, Status> {
        let permission_req = request.into_inner();

        let permission = self
            .permission_repository
            .get_permission_by_name(&permission_req.permission_name)
            .await
            .map_err(|e| Status::internal("Failed to get permission by name:"))?;

        Ok(Response::new(GetPermissionResponse {
            permission: Some(permission.into()),
            success: true,
        }))
    }

    async fn list_permissions(
        &self,
        request: Request<ListPermissionsRequest>,
    ) -> Result<Response<ListPermissionsResponse>, Status> {
        let permissions = self
            .permission_repository
            .list_permissions()
            .await
            .map_err(|e| Status::internal(format!("Failed to list permissions: {}", e)))?;

        Ok(Response::new(ListPermissionsResponse {
            count: permissions.len() as i32,
            permissions: permissions.into_iter().map(Into::into).collect(),
            success: true,
        }))
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
