use crate::db::{
    AuthorizationRepositoryTrait, Database, PermissionRepositoryTrait, RoleRepositoryTrait,
    SessionRepository, SessionRepositoryTrait, UserRepository, UserRepositoryTrait,
};
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

mod db;
mod models;
mod proto;
mod service;
mod utils;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let database_url = db::get_database_url().map_err(|e| {
        sqlx::Error::Configuration(Box::new(std::io::Error::new(
            std::io::ErrorKind::InvalidInput,
            e,
        )))
    })?;

    let database = Database::new(&database_url).await?;

    let user_repository: Arc<dyn UserRepositoryTrait> =
        Arc::new(UserRepository::new(database.clone()));

    let session_repository: Arc<dyn SessionRepositoryTrait> =
        Arc::new(SessionRepository::new(database.clone()));

    let authorization_repository: Arc<dyn AuthorizationRepositoryTrait> =
        Arc::new(db::AuthorizationRepository::new(database.clone()));

    let permission_repository: Arc<dyn PermissionRepositoryTrait> =
        Arc::new(db::PermissionRepository::new(database.clone()));

    let role_repository: Arc<dyn RoleRepositoryTrait> = Arc::new(db::RoleRepository::new(database));

    let authentication_service =
        AuthenticationService::new(user_repository.clone(), session_repository.clone());

    let user_service = UserService::new(
        user_repository.clone(),
        authorization_repository.clone(),
        role_repository.clone(),
    );

    let authorization_service = AuthorizationService::new(authorization_repository.clone());
    let permission_service = PermissionService::new(permission_repository.clone());
    let role_service = RoleService::new(authorization_repository.clone(), role_repository.clone());

    let addr = "[::1]:5000".parse()?;
    println!("Server listening on {}", addr);

    Server::builder()
        .add_service(AuthenticationServiceServer::new(authentication_service))
        .add_service(UserServiceServer::new(user_service))
        .add_service(AuthorizationServiceServer::new(authorization_service))
        .add_service(PermissionServiceServer::new(permission_service))
        .add_service(RoleServiceServer::new(role_service))
        .serve(addr)
        .await?;
    Ok(())
}
