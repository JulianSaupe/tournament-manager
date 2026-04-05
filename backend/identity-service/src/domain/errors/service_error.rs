use crate::adapter::driven::repository_error::RepositoryError;

#[derive(Debug)]
pub enum ServiceError {
    NotFound(String),
    Unauthorized(String),
    Internal(String),
}

impl std::fmt::Display for ServiceError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            ServiceError::NotFound(msg) => write!(f, "Not found: {}", msg),
            ServiceError::Unauthorized(msg) => write!(f, "Unauthorized: {}", msg),
            ServiceError::Internal(msg) => write!(f, "Internal error: {}", msg),
        }
    }
}

impl From<RepositoryError> for ServiceError {
    fn from(e: RepositoryError) -> Self {
        match e {
            RepositoryError::NotFound => ServiceError::NotFound("Record not found".to_string()),
            RepositoryError::DatabaseError(e) => ServiceError::Internal(e.to_string()),
        }
    }
}
