use std::collections::HashMap;
use std::fs::{self, File};
use std::io::Read;
use std::path::Path;

use anyhow::{Context, Result};
use clap::Args;
use flate2::write::GzEncoder;
use flate2::Compression;
use serde::{Deserialize, Serialize};
use tar::Builder;
use vumm_api::auth::ModMetadata;

use crate::config::Config;

const MOD_JSON_FILENAME: &str = "mod.json";
const DEFAULT_IGNORES: &[&str] = &[
    ".git",
    ".gitignore",
    ".github/",

    "node_modules/",
    "package.json",
    "package-lock.json",
    "yarn.lock",
    ".vscode/",
    ".idea/",
    "*.iml",
    ".vummignore",
];

fn parse_mod_string(input: &str) -> (String, String) {
    if let Some((name, version)) = input.split_once('@') {
        (name.to_string(), version.to_string())
    } else {
        (input.to_string(), String::new())
    }
}

#[derive(Serialize, Deserialize)]
pub struct ModJson {
    pub name: String,
    pub version: String,
    #[serde(default)]
    pub dependencies: HashMap<String, String>,
    #[serde(default)]
    pub description: Option<String>,
    #[serde(default)]
    pub author: Option<String>,
}

impl From<ModJson> for ModMetadata {
    fn from(m: ModJson) -> Self {
        ModMetadata {
            name: m.name.to_lowercase(),
            version: m.version,
            description: m.description,
            author: m.author,
        }
    }
}

fn load_mod_json(project_path: &Path) -> Result<ModJson> {
    let mod_json_path = project_path.join(MOD_JSON_FILENAME);
    let contents = fs::read_to_string(&mod_json_path)
        .with_context(|| format!("Failed to read {}", MOD_JSON_FILENAME))?;
    let contents = contents.trim_start_matches('\u{feff}');
    let mod_json: ModJson = serde_json::from_str(contents)
        .with_context(|| format!("Failed to parse {}", MOD_JSON_FILENAME))?;
    Ok(mod_json)
}

fn should_ignore(path: &Path, ignores: &[String]) -> bool {
    let path_str = path.to_string_lossy();
    for ignore in ignores {
        if ignore.starts_with('^') {
            let prefix = &ignore[1..];
            if path_str.starts_with(prefix) {
                return true;
            }
        } else if path_str.contains(ignore) {
            return true;
        }
    }
    false
}

fn create_tarball(project_path: &Path, output_path: &Path) -> Result<()> {
    let _mod_json = load_mod_json(project_path)?;
    
    let ignores = load_ignores(project_path)?;
    
    let file = File::create(output_path)?;
    let encoder = GzEncoder::new(file, Compression::default());
    let mut builder = Builder::new(encoder);
    
    let mut entries = Vec::new();
    for entry in fs::read_dir(project_path)? {
        let entry = entry?;
        let path = entry.path();
        if path.is_file() || path.is_dir() {
            entries.push(path);
        }
    }
    
    for path in entries {
        let relative_path = path.strip_prefix(project_path).unwrap();
        
        if should_ignore(relative_path, &ignores) {
            continue;
        }
        
        let mut file = File::open(&path)?;
        let mut buffer = Vec::new();
        file.read_to_end(&mut buffer)?;
        
        let mut header = tar::Header::new_old();
        header.set_path(relative_path)?;
        if path.is_dir() {
            header.set_entry_type(tar::EntryType::Directory);
            header.set_mode(0o755);
        } else {
            header.set_size(buffer.len() as u64);
            header.set_entry_type(tar::EntryType::Regular);
            header.set_mode(0o644);
        }
        
        builder.append(&header, buffer.as_slice())?;
    }
    
    builder.finish()?;
    
    Ok(())
}

fn load_ignores(project_path: &Path) -> Result<Vec<String>> {
    let mut ignores: Vec<String> = DEFAULT_IGNORES.iter().map(|s| s.to_string()).collect();
    
    let vummignore_path = project_path.join(".vummignore");
    if vummignore_path.exists() {
        let contents = fs::read_to_string(&vummignore_path)?;
        for line in contents.lines() {
            let line = line.trim();
            if !line.is_empty() && !line.starts_with('#') {
                ignores.push(line.to_string());
            }
        }
    }
    
    Ok(ignores)
}

#[derive(Args)]
#[command(about = "Pack current directory into a tarball", long_about = None)]
pub struct Pack {
    #[arg(long, default_value = "30m")]
    timeout: String,
}

impl Pack {
    pub async fn run(&self, _config: &Config) -> Result<()> {
        let project_path = std::env::current_dir()
            .context("Failed to get current directory")?;
        
        println!("Packing mod in {}", project_path.display());
        
        let mod_json = load_mod_json(&project_path)?;
        
        println!("Mod: {} v{}", mod_json.name, mod_json.version);
        
        let output_filename = format!("{}-{}.tar.gz", mod_json.name, mod_json.version);
        let output_path = project_path.join(&output_filename);
        
        create_tarball(&project_path, &output_path)?;
        
        println!("Packed to {}", output_filename);
        
        Ok(())
    }
}

#[derive(Args)]
#[command(about = "Publish a mod to the registry", long_about = None)]
pub struct Publish {
    #[arg(long, default_value = "latest")]
    tag: String,

    #[arg(long, default_value = "30m")]
    timeout: String,
}

impl Publish {
    pub async fn run(&self, config: &Config) -> Result<()> {
        let token = config.get_token()
            .context("Not logged in. Run 'vumm login' first")?;
        
        let project_path = std::env::current_dir()
            .context("Failed to get current directory")?;
        
        println!("Publishing mod from {}", project_path.display());
        
        let mod_json = load_mod_json(&project_path)?;
        
        println!("Publishing: {} v{}", mod_json.name, mod_json.version);
        
        let output_filename = format!("{}-{}.tar.gz", mod_json.name, mod_json.version);
        let output_path = project_path.join(&output_filename);
        
        create_tarball(&project_path, &output_path)?;
        
        let mut file = File::open(&output_path)?;
        let mut buffer = Vec::new();
        file.read_to_end(&mut buffer)?;
        
        let mut client = vumm_api::Client::new();
        if let Some(registry) = &config.registry {
            client.set_base_url(registry.clone());
        }
        client.set_bearer_token(token.clone());
        
        let metadata: ModMetadata = mod_json.into();
        let metadata_name = metadata.name.clone();
        
        let response = client.mods_publish()
            .publish(metadata, &self.tag, buffer.as_slice())
            .await?;
        
        if response.status().is_success() {
            println!("Successfully published {} as {}", metadata_name, self.tag);
        } else {
            eprintln!("Failed to publish: {}", response.status());
        }
        
        fs::remove_file(&output_path)?;
        
        Ok(())
    }
}

#[derive(Args)]
#[command(about = "Unpublish a mod version from the registry", long_about = None)]
pub struct Unpublish {
    #[arg(help = "Mod name (can include version like mod@1.0.0)", required = true)]
    mod_name: String,
}

impl Unpublish {
    pub async fn run(&self, config: &Config) -> Result<()> {
        let token = config.get_token()
            .context("Not logged in. Run 'vumm login' first")?;
        
        let (mod_name, version) = parse_mod_string(&self.mod_name);
        
        if version.is_empty() {
            anyhow::bail!("Version is required for unpublish");
        }
        
        let mut client = vumm_api::Client::new();
        if let Some(registry) = &config.registry {
            client.set_base_url(registry.clone());
        }
        client.set_bearer_token(token.clone());
        
        println!("Unpublishing {}@{}", mod_name, version);
        
        let response = client.mods_publish()
            .unpublish(&mod_name, &version)
            .await?;
        
        if response.status().is_success() {
            println!("Successfully unpublished {}@{}", mod_name, version);
        } else {
            eprintln!("Failed to unpublish: {}", response.status());
        }
        
        Ok(())
    }
}
