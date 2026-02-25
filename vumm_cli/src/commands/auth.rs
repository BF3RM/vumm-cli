use std::io::{self, Write};

use anyhow::Result;
use clap::Args;
use vumm_api::auth::PermissionType;

use crate::config::Config;

#[derive(Args)]
#[command(about = "Login with username and password", long_about = None)]
pub struct Login {
    #[arg(long, default_value = "publish")]
    token_type: String,
}

impl Login {
    pub async fn run(&self, config: &Config) -> Result<()> {
        let permission_type = PermissionType::from_str(&self.token_type);

        let mut client = vumm_api::Client::new();
        if let Some(registry) = &config.registry {
            client.set_base_url(registry.clone());
        }

        print!("Username: ");
        io::stdout().flush()?;
        let mut username = String::new();
        io::stdin().read_line(&mut username)?;
        let username = username.trim().to_string();

        print!("Password: ");
        io::stdout().flush()?;
        let password = rpassword::read_password()?;
        let password = password.trim().to_string();

        print!("Logging in... ");
        io::stdout().flush()?;

        let result = client
            .auth()
            .login(username, password, permission_type)
            .await?;

        config.set_token(&result.token)?;
        println!("OK");

        Ok(())
    }
}

#[derive(Args)]
#[command(about = "Register a new account", long_about = None)]
pub struct Register {
    #[arg(long, default_value = "publish")]
    token_type: String,
}

impl Register {
    pub async fn run(&self, config: &Config) -> Result<()> {
        let permission_type = PermissionType::from_str(&self.token_type);

        let mut client = vumm_api::Client::new();
        if let Some(registry) = &config.registry {
            client.set_base_url(registry.clone());
        }

        print!("Username: ");
        io::stdout().flush()?;
        let mut username = String::new();
        io::stdin().read_line(&mut username)?;
        let username = username.trim().to_string();

        print!("Password: ");
        io::stdout().flush()?;
        let password = rpassword::read_password()?;
        let password = password.trim().to_string();

        print!("Confirm Password: ");
        io::stdout().flush()?;
        let password_confirm = rpassword::read_password()?;
        let password_confirm = password_confirm.trim().to_string();

        if password != password_confirm {
            anyhow::bail!("Passwords do not match");
        }

        print!("Registering... ");
        io::stdout().flush()?;

        let result = client
            .auth()
            .register(username, password, permission_type)
            .await?;

        config.set_token(&result.token)?;
        println!("OK");

        Ok(())
    }
}

#[derive(Args)]
#[command(about = "Logout and delete stored credentials", long_about = None)]
pub struct Logout;

impl Logout {
    pub async fn run(&self, config: &Config) -> Result<()> {
        if !config.has_token() {
            println!("No active session (not logged in)");
            return Ok(());
        }

        match config.clear_token() {
            Ok(()) => println!("Logged out successfully"),
            Err(e) => eprintln!("Error clearing token: {}", e),
        }
        Ok(())
    }
}
