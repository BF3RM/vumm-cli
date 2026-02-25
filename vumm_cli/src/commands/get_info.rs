use clap::Args;
use semver::VersionReq;
use vumm_api::{Client, Error};

use crate::config::Config;

const SHOW_LAST_VERSION_COUNT: usize = 3;

fn parse_mod_string(input: &str) -> (String, Option<String>) {
    if let Some((name, version)) = input.split_once('@') {
        (name.to_string(), Some(version.to_string()))
    } else {
        (input.to_string(), None)
    }
}

#[derive(Args)]
#[command(
    about = "Get information about a specific mod",
    arg_required_else_help = true
)]
pub struct ModInfo {
    #[arg(help = "Name of the mod (can include tag like mod@latest or mod@1.0.0)", required = true)]
    mod_name: String,
    #[arg(help = "Version/Tag of the mod OR Version comparator [e.g. '<0.5'] (optional)")]
    mod_version: Option<String>,
}

impl ModInfo {
    pub async fn run(&mut self, config: &Config) -> anyhow::Result<()> {
        let mut client = Client::new();

        if let Some(registry) = &config.registry {
            client.set_base_url(registry.clone());
        }
        if let Some(token) = config.get_token() {
            client.set_bearer_token(token.clone());
        }

        let (mod_name, tag_or_version) = parse_mod_string(&self.mod_name);
        self.mod_name = mod_name;
        if let Some(v) = tag_or_version {
            self.mod_version = Some(v);
        }

        let mod_ = match client.mods().get(self.mod_name.clone()).await {
            Ok(returned_mod) => returned_mod,
            Err(Error::NotFound(response)) => {
                eprintln!("Mod not found: {}", response.status());
                return Ok(());
            }
            Err(e) => {
                eprintln!("Error retrieving mod: {}", e);
                return Ok(());
            }
        };

        if let Some(mod_version) = &self.mod_version {
            match VersionReq::parse(mod_version) {
                Ok(version_req) => {
                    if let Some(mod_version) = mod_.get_version_by_constraint(&version_req) {
                        println!("Found matching version: {}", mod_version);
                    } else {
                        println!("No version found matching the provided constraint.");
                    }
                }
                Err(_) => {
                    if let Some(mod_version) = mod_.get_version_by_tag(mod_version) {
                        println!("Found version by tag: {}", mod_version);
                    } else {
                        println!("No version found with the provided tag.");
                    }
                }
            }
        } else {
            println!(
                "Mod found: {}",
                mod_.get_last_versions(SHOW_LAST_VERSION_COUNT)
            );
        }
        
        Ok(())
    }
}
