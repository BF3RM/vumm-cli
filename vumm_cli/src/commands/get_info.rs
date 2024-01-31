use clap::Args;
use vumm_api::{Client, Error};

#[derive(Args)]
#[command(about = "Get information about a specific mod")]
pub struct ModInfo {
    #[arg(help = "Name of the mod")]
    mod_name: String,
    #[arg(help = "Version of the mod (optional)")]
    mod_version: Option<String>,
}

impl ModInfo {
    pub async fn run(&self) {
        let client = Client::new();

        match &self.mod_version {
            Some(mod_version) => {
                // Fetch specific mod version information
                match client.mods().get_version(self.mod_name.clone(), mod_version.to_owned()).await {
                    Ok(mod_version) => println!("{}", mod_version),
                    Err(Error::NotFound(response)) => println!("Mod version not found: {}", response.status()),
                    Err(e) => println!("Error: {}", e),
                }
            }
            None => {
                // Fetch all versions of the mod
                match client.mods().get(self.mod_name.clone()).await {
                    Ok(mod_) => println!("{}", mod_),
                    Err(Error::NotFound(response)) => println!("Mod not found: {}", response.status()),
                    Err(e) => println!("Error: {}", e),
                }
            }
        }
    }
}
