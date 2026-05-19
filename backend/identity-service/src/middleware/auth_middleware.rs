use axum::http::HeaderName;
use axum::{
    extract::Request,
    http::{HeaderMap, StatusCode},
    middleware::Next,
    response::Response,
};

pub(crate) async fn authenticate(
    headers: HeaderMap,
    request: Request,
    next: Next,
) -> Result<Response, StatusCode> {
    match get_token(&headers) {
        Some(token) if token_is_valid(token) => {
            let response = next.run(request).await;
            Ok(response)
        }
        _ => Err(StatusCode::UNAUTHORIZED),
    }
}

fn get_token(headers: &HeaderMap) -> Option<&str> {
    let value = headers.get(HeaderName::from_static("authorization"));
    value.and_then(|v| v.to_str().ok())
}

fn token_is_valid(token: &str) -> bool {
    if !token.starts_with("Bearer") {
        return false;
    }

    token == "Bearer"
}
