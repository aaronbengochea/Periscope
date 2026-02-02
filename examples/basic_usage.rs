use periscope::{client::OptionsChainParams, Config, MassiveClient};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // Load configuration from environment
    let config = Config::from_env()?;

    // Create the API client
    let client = MassiveClient::new(&config);

    // Fetch options chain for AAPL
    let params = OptionsChainParams {
        limit: Some(5),
        contract_type: Some("call".to_string()),
        ..Default::default()
    };

    let response = client.get_options_chain("AAPL", Some(params)).await?;

    println!("Fetched {} contracts", response.results.len());

    for contract in &response.results {
        if let Some(details) = &contract.details {
            println!(
                "{}: ${} strike, expires {}",
                details.ticker.as_deref().unwrap_or("N/A"),
                details.strike_price.unwrap_or(0.0),
                details.expiration_date.as_deref().unwrap_or("N/A")
            );

            if let Some(greeks) = &contract.greeks {
                println!("  {}", greeks);
            }
        }
    }

    Ok(())
}
