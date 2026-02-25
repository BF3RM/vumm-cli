use anyhow::Result;
use clap::Args;

use crate::config::Config;

#[derive(Args)]
#[command(about = "Check for updates", long_about = None)]
pub struct Update;

impl Update {
    pub async fn run(&self, _config: &Config) -> Result<()> {
        let current_version = env!("CARGO_PKG_VERSION");
        
        println!("Current version: v{}", current_version);
        println!("Please check GitHub releases for updates:");
        println!("https://github.com/BF3RM/vumm/releases");

        Ok(())
    }
}
