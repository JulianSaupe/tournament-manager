mod database;
pub mod repository;

pub use database::Database;
pub use repository::*;

use sqlx::postgres::PgPool;

pub type DbPool = PgPool;
