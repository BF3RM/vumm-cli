use serde::{Deserialize, Serialize};

use crate::{Client, ClientResult};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum PermissionType {
    #[serde(rename = "Readonly")]
    Readonly,
    #[serde(rename = "Publish")]
    Publish,
}

impl PermissionType {
    pub fn from_str(s: &str) -> Self {
        match s.to_lowercase().as_str() {
            "publish" => PermissionType::Publish,
            _ => PermissionType::Readonly,
        }
    }
}

pub fn deserialize_permission<'de, D>(deserializer: D) -> Result<PermissionType, D::Error>
where
    D: serde::Deserializer<'de>,
{
    let s = String::deserialize(deserializer)?;
    Ok(PermissionType::from_str(&s))
}

impl Default for PermissionType {
    fn default() -> Self {
        Self::Readonly
    }
}

impl std::str::FromStr for PermissionType {
    type Err = String;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        Ok(PermissionType::from_str(s))
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AuthResult {
    pub token: String,
    #[serde(rename = "type", deserialize_with = "deserialize_permission")]
    pub permission_type: PermissionType,
}

#[derive(Debug, Serialize)]
struct CredentialsDto {
    username: String,
    password: String,
    #[serde(rename = "type")]
    permission_type: PermissionType,
}

pub struct AuthEndpoint<'a> {
    pub client: &'a Client,
}

impl AuthEndpoint<'_> {
    pub async fn login(
        &self,
        username: String,
        password: String,
        permission_type: PermissionType,
    ) -> ClientResult<AuthResult> {
        self.auth_request("login", username, password, permission_type)
            .await
    }

    pub async fn register(
        &self,
        username: String,
        password: String,
        permission_type: PermissionType,
    ) -> ClientResult<AuthResult> {
        self.auth_request("register", username, password, permission_type)
            .await
    }

    async fn auth_request(
        &self,
        endpoint: &str,
        username: String,
        password: String,
        permission_type: PermissionType,
    ) -> ClientResult<AuthResult> {
        let path = format!("/auth/{}", endpoint);
        let body = serde_json::json!(CredentialsDto {
            username,
            password,
            permission_type,
        });

        self.client
            .parse_json_response::<AuthResult>(self.client.post(path, &body).await?)
            .await
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ModMetadata {
    pub name: String,
    pub version: String,
    pub description: Option<String>,
    pub author: Option<String>,
}

#[derive(Debug, Serialize)]
struct PublishModDto {
    #[serde(flatten)]
    metadata: ModMetadata,
    tag: String,
    archive: ModArchiveDto,
}

#[derive(Debug, Serialize)]
struct ModArchiveDto {
    data: String,
    length: usize,
    content_type: String,
}

#[derive(Debug, Serialize)]
struct GrantPermissionDto {
    #[serde(rename = "username")]
    username: String,
    #[serde(rename = "permission")]
    permission_type: Option<PermissionType>,
    #[serde(rename = "tag")]
    tag: String,
}

pub struct ModsPublishEndpoint<'a> {
    pub client: &'a Client,
}

impl ModsPublishEndpoint<'_> {
    pub async fn publish<R: tokio::io::AsyncRead + Unpin>(
        &self,
        metadata: ModMetadata,
        tag: &str,
        mut archive: R,
    ) -> ClientResult<reqwest::Response> {
        use tokio::io::AsyncReadExt;

        let mut buffer = Vec::new();
        archive.read_to_end(&mut buffer).await?;

        let length = buffer.len();
        let content_type = detect_content_type(&buffer);
        let data = base64_encode(&buffer);

        let path = format!("/mods/{}/{}", metadata.name, metadata.version);
        let body = serde_json::json!(PublishModDto {
            metadata,
            tag: tag.to_string(),
            archive: ModArchiveDto {
                data,
                length,
                content_type,
            },
        });

        self.client.put(path, &body).await
    }

    pub async fn unpublish(&self, mod_name: &str, mod_version: &str) -> ClientResult<reqwest::Response> {
        let path = format!("/mods/{}/{}", mod_name, mod_version);
        self.client.delete(path, &serde_json::Value::Null).await
    }

    pub async fn grant(
        &self,
        mod_name: &str,
        mod_tag: &str,
        username: &str,
        permission_type: Option<PermissionType>,
    ) -> ClientResult<reqwest::Response> {
        let path = format!("/mods/{}/grant", mod_name);
        let body = serde_json::json!(GrantPermissionDto {
            username: username.to_string(),
            permission_type,
            tag: mod_tag.to_string(),
        });

        self.client.post(path, &body).await
    }

    pub async fn revoke(
        &self,
        mod_name: &str,
        mod_tag: &str,
        username: &str,
    ) -> ClientResult<reqwest::Response> {
        let path = format!("/mods/{}/revoke", mod_name);
        let body = serde_json::json!(GrantPermissionDto {
            username: username.to_string(),
            permission_type: None,
            tag: mod_tag.to_string(),
        });

        self.client.post(path, &body).await
    }
}

fn detect_content_type(data: &[u8]) -> String {
    if data.starts_with(&[0x1f, 0x8b]) {
        return "application/gzip".to_string();
    }
    if data.starts_with(b"PK") {
        return "application/zip".to_string();
    }
    "application/octet-stream".to_string()
}

fn base64_encode(data: &[u8]) -> String {
    const CHARS: &[u8] = b"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    let mut result = String::new();
    let mut i = 0;

    while i < data.len() {
        let b0 = data[i] as usize;
        let b1 = if i + 1 < data.len() { data[i + 1] as usize } else { 0 };
        let b2 = if i + 2 < data.len() { data[i + 2] as usize } else { 0 };

        result.push(CHARS[b0 >> 2] as char);
        result.push(CHARS[((b0 & 0x03) << 4) | (b1 >> 4)] as char);

        if i + 1 < data.len() {
            result.push(CHARS[((b1 & 0x0f) << 2) | (b2 >> 6)] as char);
        } else {
            result.push('=');
        }

        if i + 2 < data.len() {
            result.push(CHARS[b2 & 0x3f] as char);
        } else {
            result.push('=');
        }

        i += 3;
    }

    result
}
