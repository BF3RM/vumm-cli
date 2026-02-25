use std::{collections::HashMap, fs};

use semver::{Version, VersionReq};
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct LockedMod {
    pub name: String,
    pub version: Version,
    pub dependencies: HashMap<String, VersionReq>,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct LockFile {
    pub mods: HashMap<String, LockedMod>,
}

impl LockFile {
    pub fn new() -> Self {
        LockFile {
            mods: HashMap::new(),
        }
    }

    pub fn load(path: &str) -> Result<Self, anyhow::Error> {
        let json = fs::read_to_string(path)?;
        serde_json::from_str(json.as_str()).map_err(|e| e.into())
    }

    pub fn save(&self, path: &str) -> Result<(), anyhow::Error> {
        let json = serde_json::to_string_pretty(self)?;
        fs::write(path, json.as_bytes()).map_err(|e| e.into())
    }
}
