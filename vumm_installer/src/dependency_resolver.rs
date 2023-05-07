use std::collections::{HashMap, VecDeque};

use vumm_api::{mods::ModVersion, Client};

use super::lock::{LockFile, LockedMod};
use super::ModVersionConstraint;

#[derive(thiserror::Error, Debug)]
pub enum DependencyResolverError {
    #[error("Mod not found: {0}")]
    ModNotFound(String),

    #[error("Insufficient permissions to access mod: {0}")]
    ModNoAccess(String),

    #[error("Mod version not found: {0}@{1}")]
    ModVersionNotFound(String, String),

    #[error(transparent)]
    Other(#[from] anyhow::Error),
}

pub struct DependencyResolver {
    api_client: Client,
}

impl DependencyResolver {
    pub fn new(api_client: Client) -> DependencyResolver {
        DependencyResolver {
            api_client: api_client,
        }
    }

    pub async fn resolve_dependencies(
        &self,
        mod_name: String,
        mod_version: String,
    ) -> Result<LockFile, DependencyResolverError> {
        let mut unresolved = VecDeque::new();
        let mut lock_file = LockFile {
            mods: HashMap::new(),
        };

        unresolved.push_back(ModVersionConstraint::new(mod_name, mod_version));

        while let Some(constraint) = unresolved.pop_front() {
            // Check if installed mod matches version constraint
            // if let Some(installed) = self.lock.get(*name) {
            //     if let Some(constraint_version) = &constraint.version {
            //         if !constraint_version.matches(&installed.version) {
            //             println!("Version mismatch");
            //             continue;
            //         }
            //     }
            // }

            println!("Resolving dependency constraint {}", constraint);

            let resolved_version = self.resolve_dependency_version(&constraint).await?;

            println!("Resolved version: {:?}", resolved_version);

            // TODO: Should I check if the mod was already added?

            lock_file.mods.insert(
                constraint.name.to_string(),
                LockedMod {
                    name: constraint.name.to_string(),
                    version: resolved_version.version.clone(),
                    dependencies: resolved_version.dependencies.clone().unwrap_or_default(),
                },
            );

            if let Some(dependencies) = resolved_version.dependencies {
                for (dep_name, dep_version) in dependencies.iter() {
                    // Skip internal veniceext mod
                    if dep_name == "veniceext" {
                        continue;
                    }

                    // TODO: Check if dependency is already part of the lock file

                    // TODO: Check if not already resolving this dependency and whether the version matches

                    unresolved.push_back(ModVersionConstraint {
                        name: dep_name.clone(),
                        tag: None,
                        version: Some(dep_version.clone()),
                    });
                }
            }
        }

        return Ok(lock_file);
    }

    pub async fn resolve_dependency_version(
        &self,
        mod_constraint: &ModVersionConstraint,
    ) -> Result<ModVersion, DependencyResolverError> {
        let mod_response = self
            .api_client
            .mods()
            .get(mod_constraint.name.to_string())
            .await;

        match mod_response {
            Ok(mod_info) => {
                if let Some(tag) = &mod_constraint.tag {
                    mod_info.get_version_by_tag(tag).ok_or(
                        DependencyResolverError::ModVersionNotFound(
                            mod_constraint.name.to_string(),
                            tag.to_string(),
                        ),
                    )
                } else if let Some(version) = &mod_constraint.version {
                    mod_info.get_version_by_constraint(version).ok_or(
                        DependencyResolverError::ModVersionNotFound(
                            mod_constraint.name.to_string(),
                            version.to_string(),
                        ),
                    )
                } else {
                    Err(DependencyResolverError::Other(anyhow::anyhow!(
                        "Invalid mod version constraint"
                    )))
                }
            }
            Err(err) => Err(match err {
                vumm_api::Error::NotFound(_) => {
                    DependencyResolverError::ModNotFound(mod_constraint.name.to_string())
                }
                vumm_api::Error::Forbidden(_) => {
                    DependencyResolverError::ModNoAccess(mod_constraint.name.to_string())
                }
                _ => DependencyResolverError::Other(err.into()),
            }),
        }
    }
}
