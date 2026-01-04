use crate::db::Database;
use crate::models::role::Role;
use uuid::Uuid;

pub struct RoleRepository {
    database: Database,
}

impl RoleRepository {
    pub fn new(database: Database) -> Self {
        Self { database }
    }
}

#[tonic::async_trait]
pub trait RoleRepositoryTrait: Send + Sync {
    async fn list_roles(&self) -> Result<Vec<Role>, String>;
    async fn create_role(&self, name: &str, description: &str) -> Result<Role, String>;
    async fn update_role(&self, id: Uuid, new_name: &str) -> Result<(), String>;
    async fn delete_role(&self, id: Uuid) -> Result<(), String>;
    async fn get_role_by_name(&self, name: &str) -> Result<Role, String>;
    async fn remove_permission_from_role(
        &self,
        role_id: Uuid,
        permission_id: Uuid,
    ) -> Result<(), String>;
}

#[tonic::async_trait]
impl RoleRepositoryTrait for RoleRepository {
    async fn list_roles(&self) -> Result<Vec<Role>, String> {
        let roles: Vec<Role> =
            sqlx::query_as(r#"SELECT id, name, description, created_at, updated_at FROM roles"#)
                .fetch_all(self.database.pool())
                .await
                .map_err(|e| format!("Failed to list roles: {}", e))?;

        Ok(roles)
    }

    async fn create_role(&self, name: &str, description: &str) -> Result<Role, String> {
        let role: Role = sqlx::query_as(
            r#"INSERT INTO roles (name, description) VALUES ($1, $2)
                    RETURNING id, name, description, created_at, updated_at"#,
        )
        .bind(name)
        .bind(description)
        .fetch_one(self.database.pool())
        .await
        .map_err(|e| format!("Failed to create role: {}", e))?;

        Ok(role)
    }

    async fn update_role(&self, id: Uuid, new_name: &str) -> Result<(), String> {
        sqlx::query(r#"UPDATE roles SET name = $1 WHERE id = $2"#)
            .bind(new_name)
            .bind(id)
            .execute(self.database.pool())
            .await
            .map_err(|e| format!("Failed to update role: {}", e))?;

        Ok(())
    }

    async fn delete_role(&self, id: Uuid) -> Result<(), String> {
        sqlx::query(r#"DELETE FROM roles WHERE id = $1"#)
            .bind(id)
            .execute(self.database.pool())
            .await
            .map_err(|e| format!("Failed to delete role: {}", e))?;

        Ok(())
    }

    async fn get_role_by_name(&self, name: &str) -> Result<Role, String> {
        let role: Role = sqlx::query_as(
            r#"SELECT id, name, description, created_at, updated_at FROM roles WHERE name = $1"#,
        )
        .bind(name)
        .fetch_one(self.database.pool())
        .await
        .map_err(|e| format!("Failed to get role by name: {}", e))?;

        Ok(role)
    }

    async fn remove_permission_from_role(
        &self,
        role_id: Uuid,
        permission_id: Uuid,
    ) -> Result<(), String> {
        sqlx::query(r#"DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2"#)
            .bind(role_id)
            .bind(permission_id)
            .execute(self.database.pool())
            .await
            .map_err(|e| format!("Failed to remove permission from role: {}", e))?;

        Ok(())
    }
}
