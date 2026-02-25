use clap::{Parser, Subcommand};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

mod commands;
mod config;
mod context;

use config::Config;

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
#[command(propagate_version = true)]
struct Cli {
    #[command(subcommand)]
    command: Commands,

    #[arg(short, long, global = true)]
    verbose: bool,

    #[arg(long, global = true)]
    registry: Option<String>,
}

#[derive(Subcommand)]
#[clap(version = env!("CARGO_PKG_VERSION"), author = env!("CARGO_PKG_AUTHORS"), about = env!("CARGO_PKG_DESCRIPTION"))]
enum Commands {
    Info(commands::get_info::ModInfo),
    Login(commands::auth::Login),
    Register(commands::auth::Register),
    Logout(commands::auth::Logout),
    Install(commands::install::Install),
    Uninstall(commands::install::Uninstall),
    Pack(commands::publish::Pack),
    Publish(commands::publish::Publish),
    Unpublish(commands::publish::Unpublish),
    Grant(commands::permissions::Grant),
    Revoke(commands::permissions::Revoke),
    Update(commands::update::Update),
}

#[tokio::main]
pub async fn run_cli() -> anyhow::Result<()> {
    let cli: Cli = Cli::parse();

    if cli.verbose {
        tracing_subscriber::registry()
            .with(
                tracing_subscriber::EnvFilter::try_from_default_env()
                    .unwrap_or_else(|_| "debug".into()),
            )
            .with(tracing_subscriber::fmt::layer())
            .init();
    }

    let config = Config::load().unwrap_or_default();

    match cli.command {
        Commands::Info(mut cmd) => cmd.run(&config).await?,
        Commands::Login(cmd) => cmd.run(&config).await?,
        Commands::Register(cmd) => cmd.run(&config).await?,
        Commands::Logout(cmd) => cmd.run(&config).await?,
        Commands::Install(cmd) => cmd.run(&config).await?,
        Commands::Uninstall(cmd) => cmd.run(&config).await?,
        Commands::Pack(cmd) => cmd.run(&config).await?,
        Commands::Publish(cmd) => cmd.run(&config).await?,
        Commands::Unpublish(cmd) => cmd.run(&config).await?,
        Commands::Grant(cmd) => cmd.run(&config).await?,
        Commands::Revoke(cmd) => cmd.run(&config).await?,
        Commands::Update(cmd) => cmd.run(&config).await?,
    }
    
    Ok(())
}
