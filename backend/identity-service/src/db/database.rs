use crate::db::DbPool;
use sqlx::postgres::PgPoolOptions;

pub struct Database {
    pool: DbPool,
}

impl Database {
    pub async fn new(database_url: &str) -> Result<Self, sqlx::Error> {
        let pool = PgPoolOptions::new()
            .max_connections(5)
            .connect(&database_url)
            .await?;

        println!("Database connection pool established");
        println!("Running migrations...");

        sqlx::migrate!("./migrations").run(&pool).await?;

        println!("Migrations completed");

        Ok(Self { pool })
    }

    pub fn pool(&self) -> &DbPool {
        &self.pool
    }
}
