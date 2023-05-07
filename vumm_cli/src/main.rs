use vumm_api::ClientError;
use clap::{Subcommand, Parser};

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
    SayHello(commands::talk::SayHello),
    SayHelloNicely(commands::talk::SayHelloNicely),
}

#[tokio::main]
async fn main() {
    let cli: Cli = Cli::parse();

    match &cli.command {
        Commands::SayHello(cmd) => cmd.run().await,
        Commands::SayHelloNicely(cmd) => cmd.run().await,
    }
}
