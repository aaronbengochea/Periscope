import os
import requests
import json
from dotenv import load_dotenv

load_dotenv()

API_KEY = os.getenv("MASSIVE_API_KEY")
BASE_URL = os.getenv("MASSIVE_BASE_URL")


def get_options_chain_snapshot(underlying_ticker, **kwargs):
    """
    Fetch options chain snapshot for a given underlying asset.

    Args:
        underlying_ticker: The underlying stock ticker (e.g., 'AAPL')
        **kwargs: Optional filters (strike_price, expiration_date, contract_type, limit)

    Returns:
        dict: API response containing options chain data with greeks
    """
    url = f"{BASE_URL}/snapshot/options/{underlying_ticker}"

    params = {"apiKey": API_KEY}
    params.update(kwargs)

    response = requests.get(url, params=params)
    response.raise_for_status()

    return response.json()


def print_options_chain(data):
    """Pretty print the options chain data with greeks."""
    print(f"Status: {data.get('status')}")
    print(f"Request ID: {data.get('request_id')}")
    print("-" * 80)

    results = data.get("results", [])
    print(f"Total contracts returned: {len(results)}\n")

    for option in results:
        details = option.get("details", {})
        greeks = option.get("greeks", {})
        last_quote = option.get("last_quote", {})
        last_trade = option.get("last_trade", {})

        print(f"Contract: {details.get('ticker', 'N/A')}")
        print(f"  Type: {details.get('contract_type', 'N/A').upper()}")
        print(f"  Strike: ${details.get('strike_price', 'N/A')}")
        print(f"  Expiration: {details.get('expiration_date', 'N/A')}")
        print(f"  Exercise Style: {details.get('exercise_style', 'N/A')}")

        print(f"  Greeks:")
        print(f"    Delta: {greeks.get('delta', 'N/A')}")
        print(f"    Gamma: {greeks.get('gamma', 'N/A')}")
        print(f"    Theta: {greeks.get('theta', 'N/A')}")
        print(f"    Vega:  {greeks.get('vega', 'N/A')}")

        print(f"  Implied Volatility: {option.get('implied_volatility', 'N/A')}")
        print(f"  Open Interest: {option.get('open_interest', 'N/A')}")

        if last_quote:
            print(f"  Last Quote: Bid ${last_quote.get('bid', 'N/A')} / Ask ${last_quote.get('ask', 'N/A')}")

        if last_trade:
            print(f"  Last Trade: ${last_trade.get('price', 'N/A')} ({last_trade.get('size', 'N/A')} contracts)")

        print()


if __name__ == "__main__":
    # Example: Get AAPL options chain snapshot
    ticker = "AAPL"

    print(f"Fetching options chain snapshot for {ticker}...")
    print("=" * 80)

    try:
        data = get_options_chain_snapshot(
            ticker,
            limit=10  # Limit to 10 contracts for readability
        )

        # Print formatted output
        print_options_chain(data)

        # Also print raw JSON for detailed analysis
        print("=" * 80)
        print("Raw JSON Response:")
        print("=" * 80)
        print(json.dumps(data, indent=2))

    except requests.exceptions.HTTPError as e:
        print(f"HTTP Error: {e}")
        print(f"Response: {e.response.text}")
    except Exception as e:
        print(f"Error: {e}")
