use crate::db::PermissionRepositoryTrait;
use crate::db::repository_error::RepositoryError;
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
            .map_err(|_| Status::internal("Failed to create permission:"))?;

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
            .map_err(|e| match e {
                RepositoryError::NotFound => Status::not_found("Permission not found"),
                _ => Status::internal("Failed to get permission by name:"),
            })?;

        Ok(Response::new(GetPermissionResponse {
            permission: Some(permission.into()),
            success: true,
        }))
    }

    async fn list_permissions(
        &self,
        request: Request<ListPermissionsRequest>,
    ) -> Result<Response<ListPermissionsResponse>, Status> {
        let request_params = request.into_inner();
        let page = request_params.page;
        let page_size = request_params.page_size;

        if !page.is_some() && page_size.is_some() || page.is_some() && !page_size.is_some() {
            return Err(Status::invalid_argument(
                "Both page and page_size must be provided or neither.".to_string(),
            ));
        }

        if page.is_some() && page.unwrap() < 1 {
            return Err(Status::invalid_argument(
                "Page must be greater than 0.".to_string(),
            ));
        }

        if page_size.is_some() && page_size.unwrap() < 1 {
            return Err(Status::invalid_argument(
                "Page size must be greater than 0.".to_string(),
            ));
        }

        let permissions = self
            .permission_repository
            .list_permissions(page, page_size)
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
        let permission_req = request.into_inner();

        let permission_id = uuid::Uuid::parse_str(&permission_req.permission_id)
            .map_err(|e| Status::invalid_argument(format!("Invalid permission id: {}", e)))?;

        self.permission_repository
            .update_permission(permission_id, &permission_req.name)
            .await
            .map_err(|e| match e {
                RepositoryError::NotFound => Status::not_found("Permission not found"),
                _ => Status::internal("Failed to update permission:"),
            })?;

        Ok(Response::new(UpdatePermissionResponse {
            success: true,
            message: "Permission updated successfully.".to_string(),
        }))
    }

    async fn delete_permission(
        &self,
        request: Request<DeletePermissionRequest>,
    ) -> Result<Response<DeletePermissionResponse>, Status> {
        let permission_req = request.into_inner();

        let permission_id = uuid::Uuid::parse_str(&permission_req.permission_id)
            .map_err(|e| Status::invalid_argument(format!("Invalid permission id: {}", e)))?;

        self.permission_repository
            .delete_permission(permission_id)
            .await
            .map_err(|e| match e {
                RepositoryError::NotFound => Status::not_found("Permission not found"),
                _ => Status::internal("Failed to delete permission:"),
            })?;

        Ok(Response::new(DeletePermissionResponse {
            success: true,
            message: "Permission deleted successfully.".to_string(),
        }))
    }
}
