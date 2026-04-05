use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct User {
    id: Uuid,
    username: String,
    email: String,
    password_hash: String,
    created_at: i64,
}
