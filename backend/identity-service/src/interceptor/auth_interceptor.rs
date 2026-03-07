use std::future::Future;
use std::pin::Pin;
use tonic::{Request, Status};
use tonic_async_interceptor::AsyncInterceptor;

#[derive(Clone)]
pub struct AuthInterceptor {
    token: String,
}

impl AuthInterceptor {
    pub fn new(token: String) -> Self {
        Self { token }
    }

    fn validate_token(token: &str, expected_token: &str) -> bool {
        if !token.starts_with("Bearer ") {
            return false;
        }

        let token = token.trim_start_matches("Bearer ");

        token == expected_token
    }
}

impl AsyncInterceptor for AuthInterceptor {
    type Future = Pin<Box<dyn Future<Output = Result<Request<()>, Status>> + Send>>;

    fn call(&mut self, request: Request<()>) -> Self::Future {
        let expected_token = self.token.clone();
        Box::pin(async move {
            let token = request
                .metadata()
                .get("authorization")
                .and_then(|v| v.to_str().ok());

            match token {
                Some(t) if Self::validate_token(t, &expected_token) => Ok(request),
                _ => Err(Status::unauthenticated("Invalid or missing token")),
            }
        })
    }
}
