use crate::utils::hash_string;
use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Account {
    id: String,
    username: String,
    email: String,
    password_hash: String,
    created_at: i64,
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
