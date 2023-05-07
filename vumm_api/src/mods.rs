use std::{collections::HashMap, io::Cursor};

use flate2::read::GzDecoder;
use semver::{Version, VersionReq};
use serde::Deserialize;
use tar::Archive;

use crate::{Client, ClientResult};

#[derive(Deserialize, Debug, Clone)]
pub struct Mod {
    pub name: String,
    pub description: Option<String>,
    pub author: Option<String>,
    pub tags: HashMap<String, String>,
    pub versions: HashMap<String, ModVersion>,
}

#[derive(Deserialize, Debug, Clone)]
pub struct ModVersion {
    pub name: String,
    pub description: Option<String>,
    pub author: Option<String>,
    pub version: Version,
    pub dependencies: Option<HashMap<String, VersionReq>>,
}

impl Mod {
    pub fn get_version_by_tag(&self, tag: &str) -> Option<ModVersion> {
        let tag_version = self.tags.get(tag)?;

        self.versions.get(tag_version).cloned()
    }

    pub fn get_version_by_constraint(&self, constraint: &VersionReq) -> Option<ModVersion> {
        let mut versions = self.versions.values().collect::<Vec<&ModVersion>>();

        versions.sort_by(|a, b| b.version.cmp(&a.version));

        for version in versions {
            if constraint.matches(&version.version) {
                return Some(version.clone());
            }
        }

        return None;
    }
}

pub struct ModsEndpoint<'a> {
    pub client: &'a Client,
}

impl ModsEndpoint<'_> {
    /// Get a mod by it's name
    /// # Example
    ///
    /// ```no_run
    /// use vumm_api::Client;
    ///
    /// let client = Client::new();
    /// let mod_name = String::from("mapeditor");
    ///
    /// let mod_response = client.mods().get(mod_name);
    /// ```
    ///
    /// # Arguments
    ///
    /// `mod_name` - The name of the mod to get
    pub async fn get(&self, mod_name: String) -> ClientResult<Mod> {
        let path = format!("/mods/{}", mod_name);

        return self
            .client
            .parse_json_response::<Mod>(self.client.get(path).await?)
            .await;
    }

    pub async fn get_version(
        &self,
        mod_name: String,
        mod_version: String,
    ) -> ClientResult<ModVersion> {
        let path = format!("/mods/{}/{}", mod_name, mod_version);

        return self
            .client
            .parse_json_response::<ModVersion>(self.client.get(path).await?)
            .await;
    }

    pub async fn download_version(
        &self,
        mod_name: String,
        mod_version: String,
    ) -> ClientResult<Archive<GzDecoder<Cursor<Vec<u8>>>>> {
        let path = format!("/mods/{}/{}/download", mod_name, mod_version);

        let res = self.client.get(path).await?;

        let bytes = res.bytes().await?.to_vec();

        // Create a Cursor to read the downloaded bytes
        let cursor = Cursor::new(bytes);

        // Open the gzipped tar archive
        let decoder = GzDecoder::new(cursor);
        let archive = Archive::new(decoder);

        return Ok(archive);
    }
}

#[cfg(test)]

mod tests {
    use super::*;

    #[test]
    fn test_mod_get_version_by_tag() {
        let mod_json = r#"
        {
            "name": "mapeditor",
            "description": "A map editor for Voxel Universe",
            "author": "Voxel Universe",
            "tags": {
                "latest": "0.1.0"
            },
            "versions": {
                "0.1.0": {
                    "name": "mapeditor",
                    "description": "A map editor for Voxel Universe",
                    "author": "Voxel Universe",
                    "version": "0.1.0",
                    "dependencies": {}
                }
            }
        }
    "#;

        let mod_obj: Mod = serde_json::from_str(mod_json).unwrap();

        let version = mod_obj.get_version_by_tag("latest").unwrap();

        assert_eq!(version.version, Version::parse("0.1.0").unwrap());
    }

    #[test]
    fn test_mod_get_version_by_constraint() {
        let mod_json = r#"
        {
            "name": "mapeditor",
            "description": "A map editor for Voxel Universe",
            "author": "Voxel Universe",
            "tags": {
                "latest": "0.1.0"
            },
            "versions": {
                "0.1.0": {
                    "name": "mapeditor",
                    "description": "A map editor for Voxel Universe",
                    "author": "Voxel Universe",
                    "version": "0.1.0",
                    "dependencies": {}
                },
                "0.2.0": {
                    "name": "mapeditor",
                    "description": "A map editor for Voxel Universe",
                    "author": "Voxel Universe",
                    "version": "0.2.0",
                    "dependencies": {}
                }
            }
        }
    "#;

        let mod_obj: Mod = serde_json::from_str(mod_json).unwrap();

        let version = mod_obj
            .get_version_by_constraint(&VersionReq::parse(">= 0.2.0").unwrap())
            .unwrap();

        assert_eq!(version.version, Version::parse("0.2.0").unwrap());
    }
}
