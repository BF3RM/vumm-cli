pub mod mods;

use crate::mods::ModsEndpoint;

pub struct Client {
    http_client: reqwest::Client,
    base_url: String,
    bearer_token: Option<String>,
}

impl Client {
    pub fn new() -> Self {
        Self {
            http_client: reqwest::Client::new(),
            base_url: String::from("https://vumm.bf3reality.com/api/v1"),
            bearer_token: None
        }
    }

    pub fn mods(&self) -> ModsEndpoint {
        ModsEndpoint { client: self }
    }

    pub fn set_bearer_token(&mut self, token: String) {
        self.bearer_token = Some(token);
    }

    pub async fn get(&self, path: String) -> Result<reqwest::Response, reqwest::Error> {
        let res = self.create_request(reqwest::Method::GET, path)
            .send()
            .await?;

        return Ok(res);
    }

    pub async fn post(&self, path: String, body: String) -> Result<reqwest::Response, reqwest::Error> {
        let res = self.create_request(reqwest::Method::POST, path)
            .body(body)
            .send()
            .await
            .expect("Failed send post request");

        return Ok(res);
    }

    pub async fn put(&self, path: String, body: String) -> Result<reqwest::Response, reqwest::Error> {
        let res = self.create_request(reqwest::Method::PUT, path)
            .body(body)
            .send()
            .await
            .expect("Failed send put request");

        return Ok(res);
    }

    pub async fn delete(&self, path: String) -> Result<reqwest::Response, reqwest::Error> {
        let res = self.create_request(reqwest::Method::DELETE, path)
            .send()
            .await
            .expect("Failed send delete request");

        return Ok(res);
    }

    fn create_request(&self, method: reqwest::Method, path: String) -> reqwest::RequestBuilder {
        let url = format!("{}{}", self.base_url, path);
        let mut builder = self.http_client.request(method, url);

        if let Some(token) = &self.bearer_token {
            builder = builder.header("Authorization", token);
        }

        return builder;
    }
}

pub fn add(left: usize, right: usize) -> usize {
    left + right
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let result = add(2, 2);
        assert_eq!(result, 4);
    }
}
