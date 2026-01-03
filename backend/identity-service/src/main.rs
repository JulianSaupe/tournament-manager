use crate::db::{Database, UserRepository, UserRepositoryTrait};
use crate::proto::authentication::authentication_service_server::AuthenticationServiceServer;
use crate::proto::user::user_service_server::UserServiceServer;
use crate::service::authentication_service::AuthenticationService;
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
    let user_repository: Arc<dyn UserRepositoryTrait> = Arc::new(UserRepository::new(database));

    let addr = "[::1]:5000".parse()?;
    let authentication_service = AuthenticationService::new(user_repository.clone());
    let user_service = UserService::new(user_repository.clone());

    println!("Server listening on {}", addr);

    Server::builder()
        .add_service(AuthenticationServiceServer::new(authentication_service))
        .add_service(UserServiceServer::new(user_service))
        .serve(addr)
        .await?;
    Ok(())
}
