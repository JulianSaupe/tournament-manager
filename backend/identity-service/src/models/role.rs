use chrono::{DateTime, Utc};
use sqlx::FromRow;
use uuid::Uuid;

#[derive(FromRow)]
pub struct Role {
    pub id: Uuid,
    pub name: String,
    pub description: String,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

impl From<Role> for crate::proto::authorization::Role {
    fn from(role: Role) -> crate::proto::authorization::Role {
        crate::proto::authorization::Role {
            id: role.id.to_string(),
            name: role.name,
            description: role.description,
            created_at: Some(prost_types::Timestamp {
                seconds: role.created_at.timestamp(),
                nanos: role.created_at.timestamp_subsec_nanos() as i32,
            }),
            updated_at: Some(prost_types::Timestamp {
                seconds: role.updated_at.timestamp(),
                nanos: role.updated_at.timestamp_subsec_nanos() as i32,
            }),
        }
    }
}
