use crate::domain::errors::service_error::ServiceError;
use crate::service::AuthenticationServiceTrait;
use axum::{
    Json, Router,
    extract::{Path, State},
    http::StatusCode,
    response::{IntoResponse, Response},
    routing::{delete, get, post},
};
use serde::Deserialize;
use std::sync::Arc;
use uuid::Uuid;

pub struct AuthenticationHandler {
    authentication_service: Arc<dyn AuthenticationServiceTrait>,
}

impl AuthenticationHandler {
    pub fn new(authentication_service: Arc<dyn AuthenticationServiceTrait>) -> Self {
        Self {
            authentication_service,
        }
    }

    pub fn router(self) -> Router {
        Router::new()
            .route("/login", post(login))
            .route("/logout/:session_id", delete(logout))
            .route("/validate/:session_id", get(validate_session))
            .with_state(Arc::new(self))
    }
}

#[derive(Deserialize)]
pub struct LoginRequest {
    pub email: String,
    pub password: String,
    pub ip_address: Option<String>,
    pub user_agent: Option<String>,
}

async fn login(
    State(handler): State<Arc<AuthenticationHandler>>,
    Json(payload): Json<LoginRequest>,
) -> impl IntoResponse {
    match handler
        .authentication_service
        .login(
            &payload.email,
            &payload.password,
            payload.ip_address,
            payload.user_agent,
        )
        .await
    {
        Ok(session) => (StatusCode::OK, Json(session)).into_response(),
        Err(e) => map_error(e),
    }
}

async fn logout(
    State(handler): State<Arc<AuthenticationHandler>>,
    Path(session_id): Path<Uuid>,
) -> impl IntoResponse {
    match handler.authentication_service.logout(session_id).await {
        Ok(_) => StatusCode::NO_CONTENT.into_response(),
        Err(e) => map_error(e),
    }
}

async fn validate_session(
    State(handler): State<Arc<AuthenticationHandler>>,
    Path(session_id): Path<Uuid>,
) -> impl IntoResponse {
    match handler
        .authentication_service
        .validate_session(session_id)
        .await
    {
        Ok(session) => (StatusCode::OK, Json(session)).into_response(),
        Err(e) => map_error(e),
    }
}

fn map_error(error: ServiceError) -> Response {
    match error {
        ServiceError::NotFound(msg) => (StatusCode::NOT_FOUND, Json(msg)).into_response(),
        ServiceError::Unauthorized(msg) => (StatusCode::UNAUTHORIZED, Json(msg)).into_response(),
        ServiceError::Internal(msg) => {
            (StatusCode::INTERNAL_SERVER_ERROR, Json(msg)).into_response()
        }
    }
}
