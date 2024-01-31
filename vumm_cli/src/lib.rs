use clap::{Parser, Subcommand};

mod commands;

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
#[command(propagate_version = true)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
#[clap(version = env!("CARGO_PKG_VERSION"), author = env!("CARGO_PKG_AUTHORS"), about = env!("CARGO_PKG_DESCRIPTION"))]
enum Commands {
    Info(commands::get_info::ModInfo)
}

#[tokio::main]
pub async fn run_cli() {
    let cli: Cli = Cli::parse();

    match &cli.command {
        Commands::Info(cmd) => cmd.run().await,
    }
}
