use crate::db::Database;
use crate::models::permission::Permission;
use uuid::Uuid;

pub struct PermissionRepository {
    database: Database,
}

impl PermissionRepository {
    pub fn new(database: Database) -> Self {
        Self { database }
    }
}

#[tonic::async_trait]
pub trait PermissionRepositoryTrait: Send + Sync {
    async fn create_permission(&self, name: &str) -> Result<Uuid, String>;
    async fn list_permissions(
        &self,
        page: Option<i32>,
        page_size: Option<i32>,
    ) -> Result<Vec<Permission>, String>;
    async fn update_permission(&self, id: Uuid, new_name: &str) -> Result<(), String>;
    async fn delete_permission(&self, id: Uuid) -> Result<(), String>;
    async fn get_permission_by_name(&self, name: &str) -> Result<Permission, String>;
}

#[tonic::async_trait]
impl PermissionRepositoryTrait for PermissionRepository {
    async fn create_permission(&self, name: &str) -> Result<Uuid, String> {
        let permission_id =
            sqlx::query_scalar(r#"INSERT INTO permissions (name) VALUES ($1) RETURNING id"#)
                .bind(name)
                .fetch_one(self.database.pool())
                .await
                .map_err(|e| format!("Failed to create permission: {}", e))?;

        Ok(permission_id)
    }

    async fn list_permissions(
        &self,
        page: Option<i32>,
        page_size: Option<i32>,
    ) -> Result<Vec<Permission>, String> {
        if let (Some(page), Some(page_size)) = (page, page_size) {
            let offset = (page - 1) * page_size;

            let permissions: Vec<Permission> = sqlx::query_as(
                r#"SELECT id, name, description, created_at, updated_at FROM permissions LIMIT $1 OFFSET $2"#,
            )
                .bind(page_size)
                .bind(offset)
                .fetch_all(self.database.pool())
                .await
                .map_err(|e| format!("Failed to list permissions: {}", e))?;

            return Ok(permissions);
        }

        let permissions: Vec<Permission> = sqlx::query_as(
            r#"SELECT id, name, description, created_at, updated_at FROM permissions"#,
        )
        .fetch_all(self.database.pool())
        .await
        .map_err(|e| format!("Failed to list permissions: {}", e))?;

        Ok(permissions)
    }

    async fn update_permission(&self, id: Uuid, new_name: &str) -> Result<(), String> {
        sqlx::query(
            r#"UPDATE permissions
                    SET name = $2, updated_at = now()
                    WHERE id = $1"#,
        )
        .bind(id)
        .bind(new_name)
        .execute(self.database.pool())
        .await
        .map_err(|e| format!("Failed to update permission: {}", e))?;

        Ok(())
    }

    async fn delete_permission(&self, id: Uuid) -> Result<(), String> {
        sqlx::query(r#"DELETE FROM permissions WHERE id = $1"#)
            .bind(id)
            .execute(self.database.pool())
            .await
            .map_err(|e| format!("Failed to delete permission: {}", e))?;

        Ok(())
    }

    async fn get_permission_by_name(&self, name: &str) -> Result<Permission, String> {
        let permission = sqlx::query_as(
            r#"SELECT id, name, description, created_at, updated_at FROM permissions WHERE name = $1"#,
        )
        .bind(name)
        .fetch_one(self.database.pool())
        .await
        .map_err(|e| format!("Failed to get permission by name: {}", e))?;

        Ok(permission)
    }
}
