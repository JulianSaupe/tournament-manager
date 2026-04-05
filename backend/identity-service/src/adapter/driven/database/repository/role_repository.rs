use crate::adapter::driven::database::Database;
use crate::adapter::driven::database::repository_error::RepositoryError;
use crate::domain::models::permission::Permission;
use crate::domain::models::role::Role;
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
    async fn list_roles(
        &self,
        page: Option<i32>,
        page_size: Option<i32>,
    ) -> Result<Vec<Role>, RepositoryError>;
    async fn create_role(&self, name: &str, description: &str) -> Result<Role, RepositoryError>;
    async fn update_role(
        &self,
        id: Uuid,
        new_name: &str,
        new_description: &str,
    ) -> Result<(), RepositoryError>;
    async fn delete_role(&self, id: Uuid) -> Result<(), RepositoryError>;
    async fn get_role_by_name(&self, name: &str) -> Result<Role, RepositoryError>;
    async fn remove_permission_from_role(
        &self,
        role_id: Uuid,
        permission_id: Uuid,
    ) -> Result<(), RepositoryError>;

    async fn get_role_permissions(&self, role_id: Uuid)
    -> Result<Vec<Permission>, RepositoryError>;
    async fn get_role_by_id(&self, id: Uuid) -> Result<Role, RepositoryError>;
}

#[tonic::async_trait]
impl RoleRepositoryTrait for RoleRepository {
    async fn list_roles(
        &self,
        page: Option<i32>,
        page_size: Option<i32>,
    ) -> Result<Vec<Role>, RepositoryError> {
        let roles = match (page, page_size) {
            (Some(page), Some(page_size)) => {
                let offset = (page - 1) * page_size;
                sqlx::query_as(
                    r#"SELECT id, name, description, created_at, updated_at FROM roles LIMIT $1 OFFSET $2"#,
                )
                .bind(page_size)
                .bind(offset)
                .fetch_all(self.database.pool())
                .await
            }
            _ => {
                sqlx::query_as(
                    r#"SELECT id, name, description, created_at, updated_at FROM roles"#,
                )
                .fetch_all(self.database.pool())
                .await
            }
        }
        .map_err(RepositoryError::from)?;

        Ok(roles)
    }

    async fn create_role(&self, name: &str, description: &str) -> Result<Role, RepositoryError> {
        let role: Role = sqlx::query_as(
            r#"INSERT INTO roles (name, description) VALUES ($1, $2)
                    RETURNING id, name, description, created_at, updated_at"#,
        )
        .bind(name)
        .bind(description)
        .fetch_one(self.database.pool())
        .await
        .map_err(RepositoryError::from)?;

        Ok(role)
    }

    async fn update_role(
        &self,
        id: Uuid,
        new_name: &str,
        new_description: &str,
    ) -> Result<(), RepositoryError> {
        sqlx::query(r#"UPDATE roles SET name = $1, description = $2 WHERE id = $3"#)
            .bind(new_name)
            .bind(new_description)
            .bind(id)
            .execute(self.database.pool())
            .await
            .map_err(RepositoryError::from)?;

        Ok(())
    }

    async fn delete_role(&self, id: Uuid) -> Result<(), RepositoryError> {
        let result = sqlx::query(r#"DELETE FROM roles WHERE id = $1"#)
            .bind(id)
            .execute(self.database.pool())
            .await
            .map_err(RepositoryError::from)?;

        if result.rows_affected() == 0 {
            return Err(RepositoryError::NotFound);
        }

        Ok(())
    }

    async fn get_role_by_name(&self, name: &str) -> Result<Role, RepositoryError> {
        let role: Role = sqlx::query_as(
            r#"SELECT id, name, description, created_at, updated_at FROM roles WHERE name = $1"#,
        )
        .bind(name)
        .fetch_one(self.database.pool())
        .await
        .map_err(RepositoryError::from)?;

        Ok(role)
    }

    async fn remove_permission_from_role(
        &self,
        role_id: Uuid,
        permission_id: Uuid,
    ) -> Result<(), RepositoryError> {
        let result = sqlx::query(
            r#"DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2"#,
        )
        .bind(role_id)
        .bind(permission_id)
        .execute(self.database.pool())
        .await
        .map_err(RepositoryError::from)?;

        if result.rows_affected() == 0 {
            return Err(RepositoryError::NotFound);
        }

        Ok(())
    }

    async fn get_role_permissions(
        &self,
        role_id: Uuid,
    ) -> Result<Vec<Permission>, RepositoryError> {
        let permissions: Vec<Permission> = sqlx::query_as(
            r#"SELECT permissions.id, permissions.name, permissions.description, permissions.created_at, permissions.updated_at
                    FROM role_permissions
                    INNER JOIN permissions ON role_permissions.permission_id = permissions.id
                    WHERE role_permissions.role_id = $1"#,
        )
        .bind(role_id)
        .fetch_all(self.database.pool())
        .await
        .map_err(RepositoryError::from)?;

        Ok(permissions)
    }

    async fn get_role_by_id(&self, id: Uuid) -> Result<Role, RepositoryError> {
        let role: Role = sqlx::query_as(
            r#"SELECT id, name, description, created_at, updated_at FROM roles WHERE id = $1"#,
        )
        .bind(id)
        .fetch_one(self.database.pool())
        .await
        .map_err(RepositoryError::from)?;

        Ok(role)
    }
}
