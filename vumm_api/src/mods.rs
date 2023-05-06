use std::{collections::HashMap, io::Cursor};

use flate2::read::GzDecoder;
use serde::Deserialize;
use tar::Archive;

use crate::{Client, ClientResult};

#[derive(Deserialize, Debug)]
pub struct Mod {
    pub name: String,
    pub description: Option<String>,
    pub author: Option<String>,
    pub tags: HashMap<String, String>,
    pub versions: HashMap<String, ModVersion>,
}

#[derive(Deserialize, Debug)]
pub struct ModVersion {
    pub name: String,
    pub description: Option<String>,
    pub author: Option<String>,
    pub version: String,
    pub dependencies: HashMap<String, String>,
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
