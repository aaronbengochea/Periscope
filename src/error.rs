use thiserror::Error;

#[derive(Error, Debug)]
pub enum PeriscopeError {
    #[error("Configuration error: {0}")]
    Config(String),

    #[error("API request failed: {0}")]
    Api(#[from] reqwest::Error),

    #[error("JSON parsing error: {0}")]
    Json(#[from] serde_json::Error),

    #[error("Environment variable error: {0}")]
    Env(#[from] std::env::VarError),

    #[error("Invalid input: {0}")]
    InvalidInput(String),

    #[error("Rate limit exceeded")]
    RateLimited,

    #[error("Authentication failed")]
    Unauthorized,

    #[error("Resource not found: {0}")]
    NotFound(String),
}

pub type Result<T> = std::result::Result<T, PeriscopeError>;
