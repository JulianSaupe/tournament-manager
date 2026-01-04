use crate::db::AuthorizationRepositoryTrait;
use crate::proto::authorization::authorization_service_server::AuthorizationService as AuthorizationServiceTrait;
use crate::proto::authorization::{
    CheckPermissionRequest, CheckPermissionResponse, GetUserPermissionsRequest,
    GetUserPermissionsResponse, ValidateAccessRequest, ValidateAccessResponse,
};
use std::sync::Arc;
use tonic::{Request, Response, Status};
use uuid::Uuid;

pub struct AuthorizationService {
    authorization_service: Arc<dyn AuthorizationRepositoryTrait>,
}

impl AuthorizationService {
    pub fn new(authorization_service: Arc<dyn AuthorizationRepositoryTrait>) -> Self {
        Self {
            authorization_service,
        }
    }
}

#[tonic::async_trait]
impl AuthorizationServiceTrait for AuthorizationService {
    async fn check_permission(
        &self,
        request: Request<CheckPermissionRequest>,
    ) -> Result<Response<CheckPermissionResponse>, Status> {
        todo!()
    }

    async fn get_user_permissions(
        &self,
        request: Request<GetUserPermissionsRequest>,
    ) -> Result<Response<GetUserPermissionsResponse>, Status> {
        let permission_req = request.into_inner();

        let user_id = Uuid::parse_str(&permission_req.user_id)
            .map_err(|_| Status::invalid_argument("Invalid user ID format"))?;

        let permissions = self
            .authorization_service
            .get_user_permissions(user_id)
            .await
            .map_err(|_| Status::internal("Failed to get user permissions"))?;

        Ok(Response::new(GetUserPermissionsResponse {
            success: true,
            permission_names: permissions,
        }))
    }

    async fn validate_access(
        &self,
        request: Request<ValidateAccessRequest>,
    ) -> Result<Response<ValidateAccessResponse>, Status> {
        todo!()
    }
}
