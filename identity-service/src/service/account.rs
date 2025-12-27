use crate::proto::account::account_service_server::AccountService as AccountServiceTrait;
use crate::proto::account::{
    CreateRequest, CreateResponse, DeleteRequest, DeleteResponse, ResetPasswordRequest,
    ResetPasswordResponse,
};
use crate::utils::hash_string;
use tonic::{Request, Response, Status};

pub struct AccountService {
    db_pool: sqlx::PgPool,
}

impl AccountService {
    pub fn new(db_pool: sqlx::PgPool) -> Self {
        Self { db_pool }
    }
}

#[tonic::async_trait]
impl AccountServiceTrait for AccountService {
    async fn create(
        &self,
        request: Request<CreateRequest>,
    ) -> Result<Response<CreateResponse>, Status> {
        let create_req = request.into_inner();

        let id: uuid::Uuid = sqlx::query_scalar(
            r#"
                INSERT INTO accounts (username, email, password) VALUES ($1, $2, $3) RETURNING id
            "#,
        )
        .bind(create_req.username)
        .bind(create_req.email)
        .bind(hash_string(&create_req.password).unwrap())
        .fetch_one(&self.db_pool)
        .await
        .map_err(|_| Status::internal("Failed to create account"))?;

        let response = CreateResponse {
            success: true,
            user_id: id.to_string(),
        };

        Ok(Response::new(response))
    }

    async fn delete(
        &self,
        request: Request<DeleteRequest>,
    ) -> Result<Response<DeleteResponse>, Status> {
        let delete_req = request.into_inner();

        sqlx::query(r#"DELETE FROM accounts WHERE id = $1"#)
            .bind(uuid::Uuid::parse_str(&delete_req.user_id).unwrap())
            .execute(&self.db_pool)
            .await
            .map_err(|_| Status::internal("Failed to delete account"))?;

        let response = DeleteResponse { success: true };

        Ok(Response::new(response))
    }

    async fn reset_password(
        &self,
        request: Request<ResetPasswordRequest>,
    ) -> Result<Response<ResetPasswordResponse>, Status> {
        todo!()
    }
}
