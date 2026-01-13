use crate::config::Config;
use crate::db::{
    AuthorizationRepository, AuthorizationRepositoryTrait, Database, PermissionRepository,
    PermissionRepositoryTrait, RoleRepository, RoleRepositoryTrait, SessionRepository,
    SessionRepositoryTrait, UserRepository, UserRepositoryTrait,
};
use crate::interceptor::auth_interceptor::AuthInterceptor;
use crate::proto::authentication::authentication_service_server::AuthenticationServiceServer;
use crate::proto::authorization::authorization_service_server::AuthorizationServiceServer;
use crate::proto::authorization::permission_service_server::PermissionServiceServer;
use crate::proto::authorization::role_service_server::RoleServiceServer;
use crate::proto::user::user_service_server::UserServiceServer;
use crate::service::authentication_service::AuthenticationService;
use crate::service::authorization_service::AuthorizationService;
use crate::service::permission_service::PermissionService;
use crate::service::role_service::RoleService;
use crate::service::user_service::UserService;
use std::sync::Arc;
use tonic::transport::Server;
use tonic_async_interceptor::async_interceptor;

pub struct App {
    config: Config,
}

impl App {
    pub fn new(config: Config) -> Self {
        Self { config }
    }

    pub async fn run(&self) -> Result<(), Box<dyn std::error::Error>> {
        let database = Database::new(&self.config.db.url).await?;

        let user_repository: Arc<dyn UserRepositoryTrait> =
            Arc::new(UserRepository::new(database.clone()));

        let session_repository: Arc<dyn SessionRepositoryTrait> =
            Arc::new(SessionRepository::new(database.clone()));

        let authorization_repository: Arc<dyn AuthorizationRepositoryTrait> =
            Arc::new(AuthorizationRepository::new(database.clone()));

        let permission_repository: Arc<dyn PermissionRepositoryTrait> =
            Arc::new(PermissionRepository::new(database.clone()));

        let role_repository: Arc<dyn RoleRepositoryTrait> = Arc::new(RoleRepository::new(database));

        let authentication_service =
            AuthenticationService::new(user_repository.clone(), session_repository.clone());

        let user_service = UserService::new(
            user_repository.clone(),
            authorization_repository.clone(),
            role_repository.clone(),
        );

        let authorization_service = AuthorizationService::new(authorization_repository.clone());
        let permission_service = PermissionService::new(permission_repository.clone());
        let role_service =
            RoleService::new(authorization_repository.clone(), role_repository.clone());

        let addr = format!("[::1]:{}", self.config.server_port).parse()?;
        println!("Server listening on {}", addr);

        Server::builder()
            .layer(async_interceptor(AuthInterceptor))
            .add_service(AuthenticationServiceServer::new(authentication_service))
            .add_service(UserServiceServer::new(user_service))
            .add_service(AuthorizationServiceServer::new(authorization_service))
            .add_service(PermissionServiceServer::new(permission_service))
            .add_service(RoleServiceServer::new(role_service))
            .serve(addr)
            .await?;

        Ok(())
    }
}
