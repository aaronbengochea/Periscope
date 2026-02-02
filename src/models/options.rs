use super::Greeks;
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct OptionsChainResponse {
    pub status: String,
    pub request_id: String,
    pub results: Vec<OptionContract>,
    pub next_url: Option<String>,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct OptionContract {
    pub details: Option<ContractDetails>,
    pub greeks: Option<Greeks>,
    pub implied_volatility: Option<f64>,
    pub open_interest: Option<i64>,
    pub last_quote: Option<LastQuote>,
    pub last_trade: Option<LastTrade>,
    pub day: Option<DayBar>,
    pub underlying_asset: Option<UnderlyingAsset>,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct ContractDetails {
    pub ticker: Option<String>,
    pub contract_type: Option<String>,
    pub strike_price: Option<f64>,
    pub expiration_date: Option<String>,
    pub exercise_style: Option<String>,
    pub shares_per_contract: Option<i32>,
}

impl ContractDetails {
    pub fn contract_type_enum(&self) -> Option<ContractType> {
        self.contract_type.as_ref().and_then(|t| match t.as_str() {
            "call" => Some(ContractType::Call),
            "put" => Some(ContractType::Put),
            _ => None,
        })
    }

    pub fn exercise_style_enum(&self) -> Option<ExerciseStyle> {
        self.exercise_style
            .as_ref()
            .and_then(|s| match s.as_str() {
                "american" => Some(ExerciseStyle::American),
                "european" => Some(ExerciseStyle::European),
                "bermudan" => Some(ExerciseStyle::Bermudan),
                _ => None,
            })
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ContractType {
    Call,
    Put,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ExerciseStyle {
    American,
    European,
    Bermudan,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct LastQuote {
    pub bid: Option<f64>,
    pub ask: Option<f64>,
    pub bid_size: Option<i64>,
    pub ask_size: Option<i64>,
}

impl LastQuote {
    pub fn mid_price(&self) -> Option<f64> {
        match (self.bid, self.ask) {
            (Some(bid), Some(ask)) => Some((bid + ask) / 2.0),
            _ => None,
        }
    }

    pub fn spread(&self) -> Option<f64> {
        match (self.bid, self.ask) {
            (Some(bid), Some(ask)) => Some(ask - bid),
            _ => None,
        }
    }
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct LastTrade {
    pub price: Option<f64>,
    pub size: Option<i64>,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct DayBar {
    pub open: Option<f64>,
    pub high: Option<f64>,
    pub low: Option<f64>,
    pub close: Option<f64>,
    pub volume: Option<i64>,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct UnderlyingAsset {
    pub ticker: Option<String>,
    pub price: Option<f64>,
}
