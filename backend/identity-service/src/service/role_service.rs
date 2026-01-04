use crate::db::{AuthorizationRepositoryTrait, RoleRepositoryTrait};
use crate::proto::authorization::role_service_server::RoleService as RoleServiceTrait;
use crate::proto::authorization::{
    AssignPermissionToRoleRequest, AssignPermissionToRoleResponse, CreateRoleRequest,
    CreateRoleResponse, DeleteRoleRequest, DeleteRoleResponse, GetRolePermissionsRequest,
    GetRolePermissionsResponse, GetRoleRequest, GetRoleResponse, ListRolesRequest,
    ListRolesResponse, RemovePermissionFromRoleRequest, RemovePermissionFromRoleResponse,
    UpdateRoleRequest, UpdateRoleResponse,
};
use std::sync::Arc;
use tonic::{Request, Response, Status};
use uuid::Uuid;

pub struct RoleService {
    authorization_repository: Arc<dyn AuthorizationRepositoryTrait>,
    role_repository: Arc<dyn RoleRepositoryTrait>,
}

impl RoleService {
    pub fn new(
        authorization_repository: Arc<dyn AuthorizationRepositoryTrait>,
        role_repository: Arc<dyn RoleRepositoryTrait>,
    ) -> Self {
        Self {
            authorization_repository,
            role_repository,
        }
    }
}

#[tonic::async_trait]
impl RoleServiceTrait for RoleService {
    async fn create_role(
        &self,
        request: Request<CreateRoleRequest>,
    ) -> Result<Response<CreateRoleResponse>, Status> {
        let role_req = request.into_inner();

        let role = self
            .role_repository
            .create_role(&role_req.name, &role_req.description)
            .await
            .map_err(|e| Status::internal(format!("Failed to create role: {}", e)))?;

        Ok(Response::new(CreateRoleResponse {
            success: true,
            role: Some(role.into()),
        }))
    }

    async fn get_role(
        &self,
        request: Request<GetRoleRequest>,
    ) -> Result<Response<GetRoleResponse>, Status> {
        todo!()
    }

    async fn list_roles(
        &self,
        request: Request<ListRolesRequest>,
    ) -> Result<Response<ListRolesResponse>, Status> {
        let list_req = request.into_inner();

        let roles = self
            .role_repository
            .list_roles()
            .await
            .map_err(|e| Status::internal(format!("Failed to list roles: {}", e)))?;

        Ok(Response::new(ListRolesResponse {
            count: roles.len() as i32,
            roles: roles.into_iter().map(|role| role.into()).collect(),
            success: true,
        }))
    }

    async fn update_role(
        &self,
        request: Request<UpdateRoleRequest>,
    ) -> Result<Response<UpdateRoleResponse>, Status> {
        todo!()
    }

    async fn delete_role(
        &self,
        request: Request<DeleteRoleRequest>,
    ) -> Result<Response<DeleteRoleResponse>, Status> {
        let role_id = Uuid::parse_str(&request.into_inner().role_id).map_err(|_| {
            Status::invalid_argument("Failed to parse role ID: must be a valid UUID.")
        })?;

        self.role_repository
            .delete_role(role_id)
            .await
            .map_err(|e| Status::internal(format!("Failed to delete role: {}", e)))?;

        Ok(Response::new(DeleteRoleResponse {
            success: true,
            message: "Role deleted successfully".to_string(),
        }))
    }

    async fn assign_permission_to_role(
        &self,
        request: Request<AssignPermissionToRoleRequest>,
    ) -> Result<Response<AssignPermissionToRoleResponse>, Status> {
        let permission_req = request.into_inner();

        let role_id = Uuid::parse_str(&permission_req.role_id).map_err(|_| {
            Status::invalid_argument("Failed to parse role ID: must be a valid UUID.")
        })?;

        let permission_id = Uuid::parse_str(&permission_req.permission_id).map_err(|_| {
            Status::invalid_argument("Failed to parse permission ID: must be a valid UUID.")
        })?;

        self.authorization_repository
            .assign_permission_to_role(role_id, permission_id)
            .await
            .map_err(|e| Status::internal(format!("Failed to assign permission to role: {}", e)))?;

        Ok(Response::new(AssignPermissionToRoleResponse {
            success: true,
            message: "Permission assigned successfully".to_string(),
        }))
    }

    async fn remove_permission_from_role(
        &self,
        request: Request<RemovePermissionFromRoleRequest>,
    ) -> Result<Response<RemovePermissionFromRoleResponse>, Status> {
        let permission_req = request.into_inner();

        let role_id = Uuid::parse_str(&permission_req.role_id).map_err(|_| {
            Status::invalid_argument("Failed to parse role ID: must be a valid UUID.")
        })?;

        let permission_id = Uuid::parse_str(&permission_req.permission_id).map_err(|_| {
            Status::invalid_argument("Failed to parse permission ID: must be a valid UUID.")
        })?;

        self.role_repository
            .remove_permission_from_role(role_id, permission_id)
            .await
            .map_err(|e| {
                Status::internal(format!("Failed to remove permission from role: {}", e))
            })?;

        Ok(Response::new(RemovePermissionFromRoleResponse {
            success: true,
        }))
    }

    async fn get_role_permissions(
        &self,
        request: Request<GetRolePermissionsRequest>,
    ) -> Result<Response<GetRolePermissionsResponse>, Status> {
        todo!()
    }
}
