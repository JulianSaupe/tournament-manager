mod app;
mod config;
mod db;
mod interceptor;
mod models;
mod proto;
mod service;
mod utils;

use crate::app::App;
use crate::config::Config;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let config = Config::from_env()?;
    let app = App::new(config);
    app.run().await?;
    Ok(())
}
