use crate::adapter::driven::database::{
    AuthorizationRepository, AuthorizationRepositoryTrait, CachedSessionRepository, Database,
    PermissionRepository, PermissionRepositoryTrait, RoleRepository, RoleRepositoryTrait,
    SessionRepository, SessionRepositoryTrait, UserRepository, UserRepositoryTrait,
};
use crate::adapter::driving::grpc::authentication_service::GrpcAuthenticationService;
use crate::adapter::driving::grpc::authorization_service::AuthorizationService;
use crate::adapter::driving::grpc::permission_service::PermissionService;
use crate::adapter::driving::grpc::role_service::RoleService;
use crate::adapter::driving::grpc::user_service::UserService;
use crate::adapter::driving::http::authentication_handler::AuthenticationHandler;
use crate::config::Config;
use crate::interceptor::auth_interceptor::AuthInterceptor;
use crate::proto::authentication::authentication_service_server::AuthenticationServiceServer;
use crate::proto::authorization::authorization_service_server::AuthorizationServiceServer;
use crate::proto::authorization::permission_service_server::PermissionServiceServer;
use crate::proto::authorization::role_service_server::RoleServiceServer;
use crate::proto::user::user_service_server::UserServiceServer;
use crate::service::session_cache_service::{SessionCacheService, SessionCacheServiceTrait};
use crate::service::{AuthenticationService, AuthenticationServiceTrait};
use std::sync::Arc;
use std::time::Duration;
use tokio_util::sync::CancellationToken;
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
        let shutdown_token = CancellationToken::new();
        let database = Database::new(&self.config.db.url).await?;

        // Cache Services
        let session_cache_service: Arc<dyn SessionCacheServiceTrait> = Arc::new(
            SessionCacheService::new(Duration::from_mins(1), shutdown_token.clone(), 1000),
        );

        // Repositories
        let user_repository: Arc<dyn UserRepositoryTrait> =
            Arc::new(UserRepository::new(database.clone()));

        let session_repository: Arc<dyn SessionRepositoryTrait> =
            Arc::new(CachedSessionRepository::new(
                Arc::new(SessionRepository::new(database.clone())),
                session_cache_service,
            ));

        let authorization_repository: Arc<dyn AuthorizationRepositoryTrait> =
            Arc::new(AuthorizationRepository::new(database.clone()));

        let permission_repository: Arc<dyn PermissionRepositoryTrait> =
            Arc::new(PermissionRepository::new(database.clone()));

        let role_repository: Arc<dyn RoleRepositoryTrait> = Arc::new(RoleRepository::new(database));

        // Services
        let authentication_service: Arc<dyn AuthenticationServiceTrait> = Arc::new(
            AuthenticationService::new(user_repository.clone(), session_repository.clone()),
        );

        // Grpc Services
        let grpc_authentication_service =
            GrpcAuthenticationService::new(authentication_service.clone());

        let user_service = UserService::new(
            user_repository.clone(),
            authorization_repository.clone(),
            role_repository.clone(),
        );

        let authorization_service = AuthorizationService::new(authorization_repository.clone());
        let permission_service = PermissionService::new(permission_repository.clone());
        let role_service =
            RoleService::new(authorization_repository.clone(), role_repository.clone());

        let interceptor = AuthInterceptor::new(self.config.auth.token.clone());

        let addr = format!("[::1]:{}", self.config.server_port).parse()?;
        println!("GRPC server listening on {}", addr);

        let grpc_server = Server::builder()
            .layer(async_interceptor(interceptor))
            .add_service(AuthenticationServiceServer::new(
                grpc_authentication_service,
            ))
            .add_service(UserServiceServer::new(user_service))
            .add_service(AuthorizationServiceServer::new(authorization_service))
            .add_service(PermissionServiceServer::new(permission_service))
            .add_service(RoleServiceServer::new(role_service))
            .serve(addr);

        let http_handler = AuthenticationHandler::new(authentication_service.clone());
        let http_addr_str = format!("[::1]:{}", self.config.server_port + 1);
        let http_addr: std::net::SocketAddr = http_addr_str.parse()?;
        println!("HTTP Server listening on {}", http_addr);

        let listener = tokio::net::TcpListener::bind(http_addr).await?;
        let http_server = axum::serve(listener, http_handler.router());

        tokio::select! {
            res = grpc_server => res?,
            res = http_server => res?,
            _ = shutdown_token.cancelled() => {},
        }

        shutdown_token.cancel();

        Ok(())
    }
}
