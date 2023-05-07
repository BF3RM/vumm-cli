use std::collections::HashMap;

use semver::{Version, VersionReq};
use vumm_api::{mods::ModVersion, Client};

pub struct DependencyResolver<'a> {
    api_client: Client,

    to_resolve: HashMap<&'a str, ModVersionConstraint<'a>>,
    lock: HashMap<String, LockedMod>,
}

struct ModVersionConstraint<'a> {
    name: &'a str,
    tag: Option<&'a str>,
    version: Option<VersionReq>,
}

struct LockedMod {
    version: Version,
    dependencies: HashMap<String, String>,
}

fn resolve_mod_version_constraint<'a>(
    mod_name: &'a str,
    mod_version: &'a str,
) -> ModVersionConstraint<'a> {
    let req = VersionReq::parse(mod_version);

    match req {
        Ok(version_req) => ModVersionConstraint {
            name: mod_name,
            tag: None,
            version: Some(version_req),
        },
        Err(err) => {
            println!("Error parsing version constraint: {}", err);
            ModVersionConstraint {
                name: mod_name,
                tag: Some(mod_version),
                version: None,
            }
        }
    }
}

impl<'a> DependencyResolver<'a> {
    fn new(api_client: Client) -> DependencyResolver<'a> {
        DependencyResolver {
            api_client: api_client,
            to_resolve: HashMap::new(),
            lock: HashMap::new(),
        }
    }

    async fn resolve_dependencies(&mut self, mod_name: &'a str, mod_version: &'a str) {
        self.to_resolve.insert(
            mod_name,
            resolve_mod_version_constraint(mod_name, mod_version),
        );

        for (name, constraint) in self.to_resolve.iter() {
            // Check if installed mod matches version constraint
            if let Some(installed) = self.lock.get(*name) {
                if let Some(constraint_version) = &constraint.version {
                    if !constraint_version.matches(&installed.version) {
                        println!("Version mismatch");
                        continue;
                    }
                }
            }

            let mod_version = self.resolve_dependency_version(constraint).await;
            println!("Resolved version: {:?}", mod_version)
        }
    }

    async fn resolve_dependency_version(
        &self,
        mod_constraint: &ModVersionConstraint<'a>,
    ) -> Option<ModVersion> {
        let mod_response = self
            .api_client
            .mods()
            .get(mod_constraint.name.to_string())
            .await;

        match mod_response {
            Ok(mod_info) => {
                if let Some(tag) = mod_constraint.tag {
                    return mod_info.get_version_by_tag(tag);
                }

                if let Some(version) = &mod_constraint.version {
                    return mod_info.get_version_by_constraint(version);
                }

                return None;
            }
            Err(err) => {
                println!("Error getting mod info: {}", err);
                return None;
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_resolve_version_constraint() {
        let mod_name = "test";
        let mod_version = "1.0.0";

        let constraint = resolve_mod_version_constraint(mod_name, mod_version);

        assert_eq!(constraint.name, mod_name);
        assert_eq!(constraint.tag, None);
        assert_eq!(
            constraint.version,
            Some(VersionReq::parse(mod_version).unwrap())
        );
    }

    #[test]
    fn test_resolve_version_constraint_tag() {
        let mod_name = "test";
        let mod_version = "latest";

        let constraint = resolve_mod_version_constraint(mod_name, mod_version);

        assert_eq!(constraint.name, mod_name);
        assert_eq!(constraint.tag, Some(mod_version));
        assert_eq!(constraint.version, None);
    }

    #[tokio::test]
    async fn test_resolve_dependency_version() {
        let mut api_client = Client::new();
        api_client.set_bearer_token("007e9635-7a92-4796-adb9-26d468de6b74".to_string());
        let resolver = DependencyResolver::new(api_client);

        let mod_name = "realitymod";
        let mod_version = "latest";

        let constraint = resolve_mod_version_constraint(mod_name, mod_version);

        let mod_version = resolver.resolve_dependency_version(&constraint).await;
        println!("Resolved version: {:?}", mod_version);
        assert_eq!(mod_version.is_some(), true);
    }
}
