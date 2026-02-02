use crate::config::Config;
use crate::error::Result;
use crate::models::OptionsChainResponse;
use reqwest::Client;
use tracing::{debug, instrument};

pub struct MassiveClient {
    client: Client,
    base_url: String,
    api_key: String,
}

#[derive(Debug, Default)]
pub struct OptionsChainParams {
    pub strike_price: Option<f64>,
    pub expiration_date: Option<String>,
    pub contract_type: Option<String>,
    pub limit: Option<i32>,
}

impl MassiveClient {
    pub fn new(config: &Config) -> Self {
        Self {
            client: Client::new(),
            base_url: config.massive_base_url.clone(),
            api_key: config.massive_api_key.clone(),
        }
    }

    #[instrument(skip(self))]
    pub async fn get_options_chain(
        &self,
        underlying_ticker: &str,
        params: Option<OptionsChainParams>,
    ) -> Result<OptionsChainResponse> {
        let url = format!("{}/snapshot/options/{}", self.base_url, underlying_ticker);
        let params = params.unwrap_or_default();

        debug!("Fetching options chain for {}", underlying_ticker);

        let mut request = self
            .client
            .get(&url)
            .query(&[("apiKey", &self.api_key)]);

        if let Some(limit) = params.limit {
            request = request.query(&[("limit", limit.to_string())]);
        }

        if let Some(strike) = params.strike_price {
            request = request.query(&[("strike_price", strike.to_string())]);
        }

        if let Some(expiration) = &params.expiration_date {
            request = request.query(&[("expiration_date", expiration)]);
        }

        if let Some(contract_type) = &params.contract_type {
            request = request.query(&[("contract_type", contract_type)]);
        }

        let response = request.send().await?.json::<OptionsChainResponse>().await?;

        debug!(
            "Received {} contracts for {}",
            response.results.len(),
            underlying_ticker
        );

        Ok(response)
    }
}
