use crate::db::Database;
use crate::db::repository_error::RepositoryError;
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
    async fn assign_role(&self, user_id: Uuid, role_id: Uuid) -> Result<(), RepositoryError>;
    async fn revoke_role(&self, user_id: Uuid, role_id: Uuid) -> Result<(), RepositoryError>;
    async fn get_roles_for_user(&self, user_id: Uuid) -> Result<Vec<Role>, RepositoryError>;
    async fn get_user_permissions(&self, user_id: Uuid) -> Result<Vec<String>, RepositoryError>;
    async fn assign_permission_to_role(
        &self,
        role_id: Uuid,
        permission_id: Uuid,
    ) -> Result<(), RepositoryError>;
}

#[tonic::async_trait]
impl AuthorizationRepositoryTrait for AuthorizationRepository {
    async fn assign_role(&self, user_id: Uuid, role_id: Uuid) -> Result<(), RepositoryError> {
        sqlx::query(r#"INSERT INTO users_roles (user_id, role_id) VALUES ($1, $2)"#)
            .bind(user_id)
            .bind(role_id)
            .execute(self.database.pool())
            .await
            .map_err(RepositoryError::from)?;

        Ok(())
    }

    async fn revoke_role(&self, user_id: Uuid, role_id: Uuid) -> Result<(), RepositoryError> {
        let result = sqlx::query(r#"DELETE FROM users_roles WHERE user_id = $1 AND role_id = $2"#)
            .bind(user_id)
            .bind(role_id)
            .execute(self.database.pool())
            .await
            .map_err(RepositoryError::from)?;

        if result.rows_affected() == 0 {
            return Err(RepositoryError::NotFound);
        }

        Ok(())
    }

    async fn get_roles_for_user(&self, user_id: Uuid) -> Result<Vec<Role>, RepositoryError> {
        let roles: Vec<Role> = sqlx::query_as(
            r#"SELECT r.id, r.name, r.description, r.created_at, r.updated_at
                   FROM users_roles ur INNER JOIN roles r ON ur.role_id = r.id
                   WHERE ur.user_id = $1"#,
        )
        .bind(user_id)
        .fetch_all(self.database.pool())
        .await
        .map_err(RepositoryError::from)?;

        Ok(roles)
    }

    async fn get_user_permissions(&self, user_id: Uuid) -> Result<Vec<String>, RepositoryError> {
        let permissions: Vec<String> = sqlx::query_scalar(
            r#"SELECT p.name FROM permissions p
                    INNER JOIN role_permissions rp ON p.id = rp.permission_id
                    INNER JOIN users_roles ur ON rp.role_id = ur.role_id
                    WHERE ur.user_id = $1"#,
        )
        .bind(user_id)
        .fetch_all(self.database.pool())
        .await
        .map_err(RepositoryError::from)?;

        Ok(permissions)
    }

    async fn assign_permission_to_role(
        &self,
        role_id: Uuid,
        permission_id: Uuid,
    ) -> Result<(), RepositoryError> {
        sqlx::query(r#"INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)"#)
            .bind(role_id)
            .bind(permission_id)
            .execute(self.database.pool())
            .await
            .map_err(RepositoryError::from)?;

        Ok(())
    }
}
