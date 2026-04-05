use identity_service::app::App;
use identity_service::config::Config;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let config = Config::from_env()?;
    let app = App::new(config);
    app.run().await?;
    Ok(())
}