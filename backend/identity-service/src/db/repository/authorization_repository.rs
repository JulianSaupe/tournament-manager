use crate::db::Database;
use crate::models::role::Role;
use uuid::Uuid;

pub struct AuthorizationRepository {
    database: Database,
}

impl AuthorizationRepository {
    pub fn new(database: Database) -> Self {
        Self { database }
    }
}

#[tonic::async_trait]
pub trait AuthorizationRepositoryTrait: Send + Sync {
    async fn assign_role(&self, user_id: Uuid, role_id: Uuid) -> Result<(), String>;
    async fn revoke_role(&self, user_id: Uuid, role_id: Uuid) -> Result<(), String>;
    async fn get_roles_for_user(&self, user_id: Uuid) -> Result<Vec<Role>, String>;
    async fn get_user_permissions(&self, user_id: Uuid) -> Result<Vec<String>, String>;
    async fn assign_permission_to_role(
        &self,
        role_id: Uuid,
        permission_id: Uuid,
    ) -> Result<(), String>;
}

#[tonic::async_trait]
impl AuthorizationRepositoryTrait for AuthorizationRepository {
    async fn assign_role(&self, user_id: Uuid, role_id: Uuid) -> Result<(), String> {
        sqlx::query(r#"INSERT INTO users_roles (user_id, role_id) VALUES ($1, $2)"#)
            .bind(user_id)
            .bind(role_id)
            .execute(self.database.pool())
            .await
            .map_err(|e| format!("Failed to assign role to user: {}", e))?;

        Ok(())
    }

    async fn revoke_role(&self, user_id: Uuid, role_id: Uuid) -> Result<(), String> {
        sqlx::query(r#"DELETE FROM users_roles WHERE user_id = $1 AND role_id = $2"#)
            .bind(user_id)
            .bind(role_id)
            .execute(self.database.pool())
            .await
            .map_err(|e| format!("Failed to revoke role from user: {}", e))?;

        Ok(())
    }

    async fn get_roles_for_user(&self, user_id: Uuid) -> Result<Vec<Role>, String> {
        let roles: Vec<Role> = sqlx::query_as(
            r#"SELECT r.id, r.name, r.description, r.created_at, r.updated_at
                   FROM users_roles ur INNER JOIN roles r ON ur.role_id = r.id
                   WHERE ur.user_id = $1"#,
        )
        .bind(user_id)
        .fetch_all(self.database.pool())
        .await
        .map_err(|e| format!("Failed to get roles for user: {}", e))?;

        Ok(roles)
    }

    async fn get_user_permissions(&self, user_id: Uuid) -> Result<Vec<String>, String> {
        let permissions: Vec<String> = sqlx::query_scalar(
            r#"SELECT p.name FROM permissions p
                    INNER JOIN role_permissions rp ON p.id = rp.permission_id
                    INNER JOIN users_roles ur ON rp.role_id = ur.role_id
                    WHERE ur.user_id = $1"#,
        )
        .bind(user_id)
        .fetch_all(self.database.pool())
        .await
        .map_err(|e| format!("Failed to get user permissions: {}", e))?;

        Ok(permissions)
    }

    async fn assign_permission_to_role(
        &self,
        role_id: Uuid,
        permission_id: Uuid,
    ) -> Result<(), String> {
        sqlx::query(r#"INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)"#)
            .bind(role_id)
            .bind(permission_id)
            .execute(self.database.pool())
            .await
            .map_err(|e| format!("Failed to assign permission to role: {}", e))?;

        Ok(())
    }
}
