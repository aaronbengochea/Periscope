use crate::error::{PeriscopeError, Result};
use std::env;

#[derive(Debug, Clone)]
pub struct Config {
    pub massive_api_key: String,
    pub massive_base_url: String,
}

impl Config {
    pub fn from_env() -> Result<Self> {
        dotenvy::dotenv().ok();

        let massive_api_key = env::var("MASSIVE_API_KEY")
            .map_err(|_| PeriscopeError::Config("MASSIVE_API_KEY must be set".into()))?;

        let massive_base_url = env::var("MASSIVE_BASE_URL")
            .unwrap_or_else(|_| "https://api.massive.com/v3".into());

        Ok(Self {
            massive_api_key,
            massive_base_url,
        })
    }
}
