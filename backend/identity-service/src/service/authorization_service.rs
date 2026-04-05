use crate::db::AuthorizationRepositoryTrait;
use crate::proto::authorization::authorization_service_server::AuthorizationService as AuthorizationServiceTrait;
use crate::proto::authorization::{AssignRoleToUserRequest, AssignRoleToUserResponse, CheckPermissionRequest, CheckPermissionResponse, GetRoleUsersRequest, GetRoleUsersResponse, GetUserPermissionsRequest, GetUserPermissionsResponse, GetUserRolesRequest, GetUserRolesResponse, RemoveRoleFromUserRequest, RemoveRoleFromUserResponse};
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

    async fn assign_role_to_user(&self, request: Request<AssignRoleToUserRequest>) -> Result<Response<AssignRoleToUserResponse>, Status> {
        todo!()
    }

    async fn remove_role_from_user(&self, request: Request<RemoveRoleFromUserRequest>) -> Result<Response<RemoveRoleFromUserResponse>, Status> {
        todo!()
    }

    async fn get_user_roles(&self, request: Request<GetUserRolesRequest>) -> Result<Response<GetUserRolesResponse>, Status> {
        todo!()
    }

    async fn get_role_users(&self, request: Request<GetRoleUsersRequest>) -> Result<Response<GetRoleUsersResponse>, Status> {
        todo!()
    }
}
