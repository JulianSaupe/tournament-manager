use serde::{Deserialize, Serialize};
use uuid::Uuid;
use crate::utils::hash_string;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Account {
    pub id: String,
    pub username: String,
    pub email: String,
    pub password_hash: String,
    pub created_at: i64,
}

impl Account {
    pub fn new(username: String, email: String, password: String) -> Self {
        Self {
            id: Uuid::new_v4().to_string(),
            username,
            email,
            password_hash: hash_string(&password).unwrap(),
            created_at: chrono::Utc::now().timestamp(),
        }
    }
}