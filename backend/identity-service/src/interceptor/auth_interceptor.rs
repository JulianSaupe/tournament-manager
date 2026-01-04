use std::future::Future;
use std::pin::Pin;
use tonic::{Request, Status};
use tonic_async_interceptor::AsyncInterceptor;

#[derive(Clone)]
pub struct AuthInterceptor;

impl AsyncInterceptor for AuthInterceptor {
    type Future = Pin<Box<dyn Future<Output = Result<Request<()>, Status>> + Send>>;

    fn call(&mut self, request: Request<()>) -> Self::Future {
        Box::pin(async move {
            let token = request
                .metadata()
                .get("authorization")
                .and_then(|v| v.to_str().ok());

            match token {
                Some(t) if validate_token(t) => Ok(request),
                _ => Err(Status::unauthenticated("Invalid or missing token")),
            }
        })
    }
}

fn validate_token(token: &str) -> bool {
    // TODO: Implement actual token validation
    token.starts_with("Bearer ")
}
