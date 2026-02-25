use anyhow::{Context, Result};
use clap::Args;
use vumm_api::auth::PermissionType;

use crate::config::Config;

#[derive(Args)]
#[command(about = "Grant permissions to a user for a mod", long_about = None)]
pub struct Grant {
    #[arg(help = "Mod name", required = true)]
    mod_name: String,

    #[arg(help = "Username to grant permission", required = true)]
    username: String,

    #[arg(help = "Permission type: publish or readonly", required = true)]
    permission: String,
}

impl Grant {
    pub async fn run(&self, config: &Config) -> Result<()> {
        let token = config.get_token()
            .context("Not logged in. Run 'vumm login' first")?;

        let permission_type = PermissionType::from_str(&self.permission);

        let mut client = vumm_api::Client::new();
        if let Some(registry) = &config.registry {
            client.set_base_url(registry.clone());
        }
        client.set_bearer_token(token.clone());

        println!("Granting {} permission to {} for {}",
            self.permission, self.username, self.mod_name);

        let response = client.mods_publish()
            .grant(&self.mod_name, "", &self.username, Some(permission_type))
            .await?;

        if response.status().is_success() {
            println!("Successfully granted permission");
        } else {
            eprintln!("Failed to grant permission: {}", response.status());
        }

        Ok(())
    }
}

#[derive(Args)]
#[command(about = "Revoke permissions from a user for a mod", long_about = None)]
pub struct Revoke {
    #[arg(help = "Mod name", required = true)]
    mod_name: String,

    #[arg(help = "Username to revoke permission", required = true)]
    username: String,
}

impl Revoke {
    pub async fn run(&self, config: &Config) -> Result<()> {
        let token = config.get_token()
            .context("Not logged in. Run 'vumm login' first")?;

        let mut client = vumm_api::Client::new();
        if let Some(registry) = &config.registry {
            client.set_base_url(registry.clone());
        }
        client.set_bearer_token(token.clone());

        println!("Revoking permission from {} for {}",
            self.username, self.mod_name);

        let response = client.mods_publish()
            .revoke(&self.mod_name, "", &self.username)
            .await?;

        if response.status().is_success() {
            println!("Successfully revoked permission");
        } else {
            eprintln!("Failed to revoke permission: {}", response.status());
        }

        Ok(())
    }
}
