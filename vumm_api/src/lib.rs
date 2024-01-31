pub mod mods;

use serde::de::DeserializeOwned;
use serde_json::Value;

use crate::mods::ModsEndpoint;

#[derive(thiserror::Error, Debug)]
pub enum Error {
    #[error("bad request")]
    BadRequest(reqwest::Response),

    #[error("unauthorized")]
    Unauthorized(reqwest::Response),

    #[error("forbidden")]
    Forbidden(reqwest::Response),

    #[error("not found")]
    NotFound(reqwest::Response),

    #[error("status code {}", reqwest::Response::status(.0))]
    Unknown(reqwest::Response),

    #[error("failed to process request: {0}")]
    Internal(#[from] reqwest::Error),
}

pub type ClientResult<T> = Result<T, Error>;

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
            bearer_token: None,
        }
    }

    pub fn mods(&self) -> ModsEndpoint {
        ModsEndpoint { client: self }
    }

    pub fn set_bearer_token(&mut self, token: String) {
        self.bearer_token = Some(token);
    }

    pub async fn get(&self, path: String) -> ClientResult<reqwest::Response> {
        self.request(reqwest::Method::GET, path, |req| req).await
    }

    pub async fn post(&self, path: String, body: &Value) -> ClientResult<reqwest::Response> {
        self.request(reqwest::Method::POST, path, |req| req.json(body))
            .await
    }

    pub async fn put(&self, path: String, body: &Value) -> ClientResult<reqwest::Response> {
        self.request(reqwest::Method::PUT, path, |req| req.json(body))
            .await
    }

    pub async fn delete(&self, path: String, body: &Value) -> ClientResult<reqwest::Response> {
        self.request(reqwest::Method::DELETE, path, |req| req.json(body))
            .await
    }

    async fn request<B>(
        &self,
        method: reqwest::Method,
        path: String,
        request_builder: B,
    ) -> ClientResult<reqwest::Response>
    where
        B: FnOnce(reqwest::RequestBuilder) -> reqwest::RequestBuilder,
    {
        let url = format!("{}{}", self.base_url, path);
        let mut request = self.http_client.request(method.clone(), url);

        if let Some(token) = &self.bearer_token {
            request = request.header("Authorization", token);
        }

        request = request_builder(request);

        let response = request.send().await?;

        if response.status().is_success() {
            Ok(response)
        } else {
            Err(match response.status() {
                // TODO: handle validation errors
                reqwest::StatusCode::BAD_REQUEST => Error::BadRequest(response),
                reqwest::StatusCode::UNAUTHORIZED => Error::Unauthorized(response),
                reqwest::StatusCode::FORBIDDEN => Error::Forbidden(response),
                reqwest::StatusCode::NOT_FOUND => Error::NotFound(response),
                _ => Error::Unknown(response),
            })
        }
    }

    async fn parse_json_response<T: DeserializeOwned>(
        &self,
        response: reqwest::Response,
    ) -> ClientResult<T> {
        response.json::<T>().await.map_err(Into::into)
    }
}