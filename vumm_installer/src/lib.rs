mod dependency_resolver;
mod lock;

use core::fmt;
use semver::VersionReq;

pub use dependency_resolver::{DependencyResolver, DependencyResolverError};
pub use lock::{LockFile, LockedMod};

#[derive(Clone, Debug)]
pub struct ModVersionConstraint {
    name: String,
    tag: Option<String>,
    version: Option<VersionReq>,
}

impl ModVersionConstraint {
    pub fn new(mod_name: String, mod_version: String) -> Self {
        let req = VersionReq::parse(mod_version.as_str());

        match req {
            Ok(version_req) => ModVersionConstraint {
                name: mod_name,
                tag: None,
                version: Some(version_req),
            },
            Err(_) => ModVersionConstraint {
                name: mod_name,
                tag: Some(mod_version),
                version: None,
            },
        }
    }

    pub fn get_name(&self) -> &String {
        &self.name
    }

    pub fn get_tag(&self) -> &Option<String> {
        &self.tag
    }

    pub fn get_version(&self) -> &Option<VersionReq> {
        &self.version
    }
}

impl fmt::Display for ModVersionConstraint {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        if let Some(version) = &self.version {
            write!(f, "{}@{}", self.name, version)
        } else {
            write!(f, "{}@{}", self.name, self.tag.as_ref().unwrap())
        }
    }
}

#[cfg(test)]
mod tests {
    use vumm_api::Client;

    use super::*;

    // #[test]
    // fn test_resolve_version_constraint() {
    //     let mod_name = String::from("test");
    //     let mod_version = String::from("1.0.0");

    //     let constraint = resolve_mod_version_constraint(mod_name, mod_version);

    //     assert_eq!(constraint.name, mod_name);
    //     assert_eq!(constraint.tag, None);
    //     assert_eq!(
    //         constraint.version,
    //         Some(VersionReq::parse(mod_version.as_str()).unwrap())
    //     );
    // }

    // #[test]
    // fn test_resolve_version_constraint_tag() {
    //     let mod_name = String::from("test");
    //     let mod_version = String::from("latest");

    //     let constraint = resolve_mod_version_constraint(mod_name, mod_version);

    //     assert_eq!(constraint.name, mod_name);
    //     assert_eq!(constraint.tag, Some(mod_version));
    //     assert_eq!(constraint.version, None);
    // }

    #[tokio::test]
    async fn test_resolve_dependencies() {
        let api_client = Client::new();
        let resolver = DependencyResolver::new(api_client);

        let lock_file = resolver
            .resolve_dependencies("mapeditor".to_string(), "preview".to_string())
            .await;

        match lock_file {
            Ok(lock_file) => {
                println!("Lock file: {:#?}", lock_file);
            }
            Err(err) => {
                println!("Error: {}", err);
            }
        }
    }
}
