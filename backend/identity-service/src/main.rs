use crate::db::{AccountRepository, AccountRepositoryTrait, Database};
use crate::proto::account::account_service_server::AccountServiceServer;
use crate::proto::authentication::authentication_service_server::AuthenticationServiceServer;
use crate::service::account_service::AccountService;
use crate::service::authentication_service::AuthenticationService;
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
    let account_repository: Arc<dyn AccountRepositoryTrait> =
        Arc::new(AccountRepository::new(database));

    let addr = "[::1]:5000".parse()?;
    let authentication_service = AuthenticationService::new(account_repository.clone());
    let account_service = AccountService::new(account_repository.clone());

    println!("Server listening on {}", addr);

    Server::builder()
        .add_service(AuthenticationServiceServer::new(authentication_service))
        .add_service(AccountServiceServer::new(account_service))
        .serve(addr)
        .await?;
    Ok(())
}
