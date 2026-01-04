use sqlx::FromRow;
use uuid::Uuid;

#[derive(FromRow)]
pub struct Permission {
    pub id: Uuid,
    pub name: String,
    pub description: String,
    pub created_at: chrono::DateTime<chrono::Utc>,
    pub updated_at: chrono::DateTime<chrono::Utc>,
}

impl From<Permission> for crate::proto::authorization::Permission {
    fn from(permission: Permission) -> Self {
        Self {
            id: permission.id.to_string(),
            name: permission.name,
            description: permission.description,
            created_at: Some(prost_types::Timestamp {
                seconds: permission.created_at.timestamp(),
                nanos: permission.created_at.timestamp_subsec_nanos() as i32,
            }),
            updated_at: Some(prost_types::Timestamp {
                seconds: permission.updated_at.timestamp(),
                nanos: permission.updated_at.timestamp_subsec_nanos() as i32,
            }),
        }
    }
}
