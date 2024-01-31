use clap::Args;
use vumm_api::{Client, Error};
use semver::VersionReq;

#[derive(Args)]
#[command(about = "Get information about a specific mod", arg_required_else_help = true)]
pub struct ModInfo {
    #[arg(help = "Name of the mod", required = true)]
    mod_name: String,
    #[arg(help = "Version/Tag of the mod (optional)")]
    mod_version: Option<String>,
}

impl ModInfo {
    pub async fn run(&self) {
        let client = Client::new();

        let mod_ = match client.mods().get(self.mod_name.clone()).await {
            Ok(returned_mod) => returned_mod,
            Err(Error::NotFound(response)) => {
                println!("Mod not found: {}", response.status());
                return;
            }
            Err(e) => {
                println!("Error retrieving mod: {}", e);
                return;
            }
        };

        if let Some(mod_version) = &self.mod_version {
            match VersionReq::parse(mod_version) {
                Ok(version_req) => {
                    // VersionConstraint provided
                    if let Some(mod_version) = mod_.get_version_by_constraint(&version_req) {
                        println!("Found matching version: {}", mod_version);
                    } else {
                        println!("No version found matching the provided constraint.");
                    }
                }
                Err(_) => {
                    // Invalid VersionConstraint treat as a tag
                    if let Some(mod_version) = mod_.get_version_by_tag(mod_version) {
                        println!("Found version by tag: {}", mod_version);
                    } else {
                        println!("No version found with the provided tag.");
                    }
                }
            }
        } else {
            // If no specific version requirement or tag is provided, simply print the mod
            println!("Mod found: {}", mod_);
        }
    }
}


