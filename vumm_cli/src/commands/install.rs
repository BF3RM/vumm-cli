use anyhow::Result;
use clap::Args;
use vumm_api::Client;
use std::fs;
use std::path::Path;

use crate::config::Config;
use vumm_installer::DependencyResolver;

fn copy_dir_all(src: &Path, dst: &Path) -> Result<()> {
    fs::create_dir_all(dst)?;
    for entry in fs::read_dir(src)? {
        let entry = entry?;
        let ty = entry.file_type()?;
        if ty.is_dir() {
            copy_dir_all(&entry.path(), &dst.join(entry.file_name()))?;
        } else {
            fs::copy(entry.path(), dst.join(entry.file_name()))?;
        }
    }
    Ok(())
}

#[derive(Args)]
#[command(about = "Install a mod", long_about = None)]
pub struct Install {
    #[arg(help = "Name of the mod OR name@version of mod", required = true)]
    mod_name: String,

    #[arg(long, default_value = "30m")]
    timeout: String,

    #[arg(long, short = 'd', help = "Target directory for mod symlinks")]
    target_dir: Option<String>,
}

impl Install {
    pub async fn run(&self, config: &Config) -> Result<()> {
        let _timeout = parse_duration::parse(&self.timeout)
            .map_err(|e| anyhow::anyhow!("Invalid timeout: {}", e))?;

        let mut client = Client::new();
        if let Some(registry) = &config.registry {
            client.set_base_url(registry.clone());
        }
        if let Some(token) = config.get_token() {
            client.set_bearer_token(token.clone());
        }

        let (mod_name, mod_version) = parse_mod_string(&self.mod_name);

        println!("Installing {}...", self.mod_name);

        let resolver = DependencyResolver::new(client.clone());
        
        let lock_file = resolver
            .resolve_dependencies(mod_name.clone(), mod_version.clone())
            .await?;

        let target_dir = match &self.target_dir {
            Some(d) => std::path::PathBuf::from(d),
            None => std::env::current_dir()?.join("mods"),
        };

        let cache_dir = get_cache_directory()?;
        
        print!("Installing to: {} ... ", target_dir.display());
        std::fs::create_dir_all(&target_dir)?;
        println!("OK");
        
        std::fs::create_dir_all(&cache_dir)?;

        println!("Installing {} mods...", lock_file.mods.len());

        for (name, locked_mod) in &lock_file.mods {
            print!("  {} v{} ... ", name, locked_mod.version);
            
            let mod_cache_dir = cache_dir.join(name).join(locked_mod.version.to_string());
            
            if !mod_cache_dir.exists() {
                let archive = client.mods()
                    .download_version(name.clone(), locked_mod.version.to_string())
                    .await;
                
                match archive {
                    Ok(mut archive) => {
                        let temp_dir = cache_dir.join(".tmp").join(name).join(locked_mod.version.to_string());
                        std::fs::create_dir_all(&temp_dir)?;
                        
                        for entry in archive.entries()? {
                            let mut entry = entry?;
                            let path = entry.path()?.into_owned();
                            let full_path = temp_dir.join(&path);
                            
                            if let Some(parent) = full_path.parent() {
                                std::fs::create_dir_all(parent)?;
                            }
                            
                            entry.unpack(&full_path)?;
                        }
                        
                        if let Err(e) = std::fs::rename(&temp_dir, &mod_cache_dir) {
                            let _ = std::fs::remove_dir_all(&temp_dir);
                            return Err(e).map_err(|e| anyhow::anyhow!("Failed to move temp dir: {}", e));
                        }
                    }
                    Err(e) => {
                        println!("FAILED");
                        eprintln!("    Download error: {}", e);
                        continue;
                    }
                }
            }

            let target_link = target_dir.join(name);
            
            #[cfg(windows)]
            {
                if target_link.exists() {
                    if target_link.is_dir() {
                        let _ = std::fs::remove_dir_all(&target_link);
                    } else {
                        let _ = std::fs::remove_file(&target_link);
                    }
                }
                copy_dir_all(&mod_cache_dir, &target_link)?;
            }
            
            #[cfg(not(windows))]
            {
                if target_link.exists() {
                    std::fs::remove_file(&target_link).ok();
                }
                std::os::unix::fs::symlink(&mod_cache_dir, &target_link)?;
            }
            
            println!("OK");
        }

        let lock_file_path = target_dir.join("vumm.lock");
        lock_file.save(lock_file_path.to_str().unwrap())?;
        
        println!("\nDone! Installed {} mods.", lock_file.mods.len());
        
        Ok(())
    }
}

fn get_cache_directory() -> Result<std::path::PathBuf> {
    if let Some(proj_dirs) = directories::ProjectDirs::from("com", "bf3rm", "vumm") {
        Ok(proj_dirs.cache_dir().to_path_buf())
    } else {
        anyhow::bail!("Could not determine cache directory")
    }
}

fn parse_mod_string(input: &str) -> (String, String) {
    if let Some((name, version)) = input.split_once('@') {
        (name.to_string(), version.to_string())
    } else {
        (input.to_string(), "latest".to_string())
    }
}

#[derive(Args)]
#[command(about = "Uninstall a mod", long_about = None)]
pub struct Uninstall {
    #[arg(help = "Name of the mod to uninstall (can include tag like mod@latest)", required = true)]
    mod_name: String,
}

impl Uninstall {
    pub async fn run(&self, _config: &Config) -> Result<()> {
        let (mod_name, _tag) = parse_mod_string(&self.mod_name);
        let target_dir = std::env::current_dir()?.join("mods");
        let target_link = target_dir.join(&mod_name);
        
        if target_link.exists() {
            std::fs::remove_dir_all(&target_link)?;
            println!("Uninstalled {}", mod_name);
        } else {
            println!("{} is not installed", mod_name);
        }
        
        Ok(())
    }
}
