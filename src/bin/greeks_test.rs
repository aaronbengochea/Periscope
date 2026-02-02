use anyhow::Result;
use clap::Parser;
use periscope::{client::MassiveClient, Config, OptionsChainResponse};
use tracing::info;
use tracing_subscriber::EnvFilter;

#[derive(Parser, Debug)]
#[command(name = "greeks_test")]
#[command(about = "Fetch and display options chain data with Greeks")]
struct Args {
    /// Underlying ticker symbol
    #[arg(short, long, default_value = "AAPL")]
    ticker: String,

    /// Number of contracts to fetch
    #[arg(short, long, default_value = "10")]
    limit: i32,

    /// Output raw JSON
    #[arg(long)]
    json: bool,
}

fn print_options_chain(data: &OptionsChainResponse) {
    println!("Status: {}", data.status);
    println!("Request ID: {}", data.request_id);
    println!("{}", "-".repeat(80));
    println!("Total contracts returned: {}\n", data.results.len());

    for option in &data.results {
        if let Some(details) = &option.details {
            println!(
                "Contract: {}",
                details.ticker.as_deref().unwrap_or("N/A")
            );
            println!(
                "  Type: {}",
                details
                    .contract_type
                    .as_deref()
                    .unwrap_or("N/A")
                    .to_uppercase()
            );
            println!(
                "  Strike: ${}",
                details
                    .strike_price
                    .map(|p| p.to_string())
                    .unwrap_or_else(|| "N/A".to_string())
            );
            println!(
                "  Expiration: {}",
                details.expiration_date.as_deref().unwrap_or("N/A")
            );
            println!(
                "  Exercise Style: {}",
                details.exercise_style.as_deref().unwrap_or("N/A")
            );
        }

        println!("  Greeks:");
        if let Some(greeks) = &option.greeks {
            println!(
                "    Delta: {}",
                greeks
                    .delta
                    .map(|v| format!("{:.6}", v))
                    .unwrap_or_else(|| "N/A".to_string())
            );
            println!(
                "    Gamma: {}",
                greeks
                    .gamma
                    .map(|v| format!("{:.6}", v))
                    .unwrap_or_else(|| "N/A".to_string())
            );
            println!(
                "    Theta: {}",
                greeks
                    .theta
                    .map(|v| format!("{:.6}", v))
                    .unwrap_or_else(|| "N/A".to_string())
            );
            println!(
                "    Vega:  {}",
                greeks
                    .vega
                    .map(|v| format!("{:.6}", v))
                    .unwrap_or_else(|| "N/A".to_string())
            );
        } else {
            println!("    N/A");
        }

        println!(
            "  Implied Volatility: {}",
            option
                .implied_volatility
                .map(|v| format!("{:.4}", v))
                .unwrap_or_else(|| "N/A".to_string())
        );
        println!(
            "  Open Interest: {}",
            option
                .open_interest
                .map(|v| v.to_string())
                .unwrap_or_else(|| "N/A".to_string())
        );

        if let Some(quote) = &option.last_quote {
            println!(
                "  Last Quote: Bid ${} / Ask ${}",
                quote
                    .bid
                    .map(|v| format!("{:.2}", v))
                    .unwrap_or_else(|| "N/A".to_string()),
                quote
                    .ask
                    .map(|v| format!("{:.2}", v))
                    .unwrap_or_else(|| "N/A".to_string())
            );
        }

        if let Some(trade) = &option.last_trade {
            println!(
                "  Last Trade: ${} ({} contracts)",
                trade
                    .price
                    .map(|v| format!("{:.2}", v))
                    .unwrap_or_else(|| "N/A".to_string()),
                trade
                    .size
                    .map(|v| v.to_string())
                    .unwrap_or_else(|| "N/A".to_string())
            );
        }

        println!();
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::from_default_env())
        .init();

    let args = Args::parse();
    let config = Config::from_env()?;
    let client = MassiveClient::new(&config);

    info!("Fetching options chain for {}", args.ticker);
    println!("Fetching options chain snapshot for {}...", args.ticker);
    println!("{}", "=".repeat(80));

    let params = periscope::client::OptionsChainParams {
        limit: Some(args.limit),
        ..Default::default()
    };

    let data = client.get_options_chain(&args.ticker, Some(params)).await?;

    if args.json {
        println!("{}", serde_json::to_string_pretty(&data)?);
    } else {
        print_options_chain(&data);

        println!("{}", "=".repeat(80));
        println!("Raw JSON Response:");
        println!("{}", "=".repeat(80));
        println!("{}", serde_json::to_string_pretty(&data)?);
    }

    Ok(())
}
