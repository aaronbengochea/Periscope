import axios from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
});

// Types matching Go backend models
export interface Greeks {
  delta?: number;
  gamma?: number;
  theta?: number;
  vega?: number;
  rho?: number;
}

export interface LastQuote {
  bid?: number;
  ask?: number;
  bid_size?: number;
  ask_size?: number;
}

export interface LastTrade {
  price?: number;
  size?: number;
}

export interface DayBar {
  open?: number;
  high?: number;
  low?: number;
  close?: number;
  volume?: number;
}

export interface Session {
  change?: number;
  change_percent?: number;
  close?: number;
  high?: number;
  low?: number;
  open?: number;
  previous_close?: number;
  volume?: number;
}

export interface ContractDetails {
  ticker?: string;
  contract_type?: 'call' | 'put';
  strike_price?: number;
  expiration_date?: string;
  exercise_style?: string;
  shares_per_contract?: number;
}

export interface UnderlyingAsset {
  ticker?: string;
  price?: number;
}

export interface OptionContract {
  details?: ContractDetails;
  greeks?: Greeks;
  implied_volatility?: number;
  open_interest?: number;
  last_quote?: LastQuote;
  last_trade?: LastTrade;
  day?: DayBar;
  session?: Session;
  underlying_asset?: UnderlyingAsset;
}

export interface OptionsChainResponse {
  status: string;
  request_id: string;
  results: OptionContract[];
  next_url?: string;
}

// API functions
export async function fetchOptionsChain(ticker: string): Promise<OptionsChainResponse> {
  console.log(`[Frontend API] Fetching options chain for ticker: ${ticker}`);

  const { data } = await apiClient.get<OptionsChainResponse>(`/options/${ticker}`);

  console.log(`[Frontend API] Received response:`, {
    status: data.status,
    request_id: data.request_id,
    total_contracts: data.results?.length || 0,
    first_contract: data.results?.[0] || null,
  });

  // Log full first contract for debugging
  if (data.results && data.results.length > 0) {
    console.log(`[Frontend API] First contract (full):`, JSON.stringify(data.results[0], null, 2));
  }

  return data;
}

export async function fetchContractDetails(contractTickers: string[]): Promise<OptionsChainResponse> {
  console.log(`[Frontend API] Fetching details for ${contractTickers.length} contracts`);

  const { data } = await apiClient.post<OptionsChainResponse>('/options/details', {
    contract_tickers: contractTickers,
  });

  console.log(`[Frontend API] Received contract details:`, {
    status: data.status,
    total_contracts: data.results?.length || 0,
  });

  return data;
}
