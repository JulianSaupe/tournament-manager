#[derive(Debug)]
pub enum RepositoryError {
    NotFound,
    DatabaseError(sqlx::Error),
}

impl std::fmt::Display for RepositoryError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            RepositoryError::NotFound => write!(f, "Record not found"),
            RepositoryError::DatabaseError(e) => write!(f, "Database error: {}", e),
        }
    }
}

impl From<sqlx::Error> for RepositoryError {
    fn from(e: sqlx::Error) -> Self {
        match e {
            sqlx::Error::RowNotFound => RepositoryError::NotFound,
            _ => RepositoryError::DatabaseError(e),
        }
    }
}
