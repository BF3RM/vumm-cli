use clap::Args;
use vumm_api::{Client, Error};
use semver::VersionReq;

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

        match client.mods().get(self.mod_name.clone()).await {
            Ok(mod_) => {
                if let Some(mod_version_req) = &self.mod_version {
                    match VersionReq::parse(mod_version_req) {
                        Ok(version_req) => {
                            match mod_.get_version_by_constraint(&version_req) {
                                Some(mod_version) => println!("Found matching version: {}", mod_version),
                                None => println!("No version found matching the provided constraint.")
                            }
                        },
                        Err(e) => println!("Failed to parse version requirement: {}", e),
                    }
                } else {
                    // If no specific version requirement is provided, simply print the mod
                    println!("{}", mod_);
                }
            },
            Err(Error::NotFound(response)) => println!("{}", response.status()),
            Err(e) => println!("Error retrieving mod: {}", e),
        }
    }
}
