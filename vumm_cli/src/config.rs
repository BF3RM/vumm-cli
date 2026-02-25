use anyhow::Result;
use directories::ProjectDirs;
use keyring::Entry;
use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;

const SERVICE_NAME: &str = "vumm";
const TOKEN_KEY: &str = "token";

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct Config {
    pub registry: Option<String>,
}

impl Config {
    pub fn load() -> Result<Self> {
        let config_path = Self::config_path()?;

        if !config_path.exists() {
            return Ok(Self::default());
        }

        let contents = fs::read_to_string(&config_path)?;
        let config: Config = serde_json::from_str(&contents)?;
        Ok(config)
    }

    #[allow(dead_code)]
    pub fn save(&self) -> Result<()> {
        let config_path = Self::config_path()?;

        if let Some(parent) = config_path.parent() {
            fs::create_dir_all(parent)?;
        }

        let contents = serde_json::to_string_pretty(self)?;
        fs::write(&config_path, contents)?;
        Ok(())
    }

    pub fn get_token(&self) -> Option<String> {
        let entry = Entry::new(SERVICE_NAME, TOKEN_KEY).ok()?;
        entry.get_password().ok()
    }

    pub fn has_token(&self) -> bool {
        self.get_token().is_some()
    }

    pub fn set_token(&self, token: &str) -> Result<()> {
        let entry = Entry::new(SERVICE_NAME, TOKEN_KEY)?;
        entry.set_password(token)?;
        Ok(())
    }

    pub fn clear_token(&self) -> Result<()> {
        let entry = Entry::new(SERVICE_NAME, TOKEN_KEY)?;
        match entry.delete_credential() {
            Ok(()) => Ok(()),
            Err(e)
                if e.to_string().contains("not found")
                    || e.to_string().contains("Element not found") =>
            {
                Err(anyhow::anyhow!("not found"))
            }
            Err(e) => Err(anyhow::anyhow!(e)),
        }
    }

    fn config_path() -> Result<PathBuf> {
        if let Some(proj_dirs) = ProjectDirs::from("com", "bf3rm", "vumm") {
            return Ok(proj_dirs.config_dir().join("config.json"));
        }
        let local_app_data = std::env::var("LOCALAPPDATA")
            .map_err(|_| anyhow::anyhow!("Could not determine config directory"))?;
        Ok(PathBuf::from(local_app_data)
            .join("bf3rm")
            .join("vumm")
            .join("config.json"))
    }
}
