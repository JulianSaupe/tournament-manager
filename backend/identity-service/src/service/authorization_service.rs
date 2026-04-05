use crate::db::AuthorizationRepositoryTrait;
use crate::proto::authorization::authorization_service_server::AuthorizationService as AuthorizationServiceTrait;
use crate::proto::authorization::{
    CheckPermissionRequest, CheckPermissionResponse, GetUserPermissionsRequest,
    GetUserPermissionsResponse,
};
use std::sync::Arc;
use tonic::{Request, Response, Status};
use uuid::Uuid;

pub struct AuthorizationService {
    authorization_repository: Arc<dyn AuthorizationRepositoryTrait>,
}

impl AuthorizationService {
    pub fn new(authorization_repository: Arc<dyn AuthorizationRepositoryTrait>) -> Self {
        Self {
            authorization_repository,
        }
    }
}

#[tonic::async_trait]
impl AuthorizationServiceTrait for AuthorizationService {
    async fn check_permission(
        &self,
        request: Request<CheckPermissionRequest>,
    ) -> Result<Response<CheckPermissionResponse>, Status> {
        let permission_req = request.into_inner();

        let user_id = Uuid::parse_str(&permission_req.user_id)
            .map_err(|_| Status::invalid_argument("Invalid user ID format"))?;

        let user_permissions = self
            .authorization_repository
            .get_user_permissions(user_id)
            .await
            .map_err(|_| Status::internal("Failed to get user permissions"))?;

        let has_permission = user_permissions.contains(&permission_req.permission_name);

        let message = if has_permission {
            "Permission granted".to_string()
        } else {
            "Permission denied".to_string()
        };

        Ok(Response::new(CheckPermissionResponse {
            allowed: has_permission,
            message,
        }))
    }

    async fn get_user_permissions(
        &self,
        request: Request<GetUserPermissionsRequest>,
    ) -> Result<Response<GetUserPermissionsResponse>, Status> {
        let permission_req = request.into_inner();

        let user_id = Uuid::parse_str(&permission_req.user_id)
            .map_err(|_| Status::invalid_argument("Invalid user ID format"))?;

        let permissions = self
            .authorization_repository
            .get_user_permissions(user_id)
            .await
            .map_err(|_| Status::internal("Failed to get user permissions"))?;

        Ok(Response::new(GetUserPermissionsResponse {
            success: true,
            permission_names: permissions,
        }))
    }
}
