use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct Greeks {
    pub delta: Option<f64>,
    pub gamma: Option<f64>,
    pub theta: Option<f64>,
    pub vega: Option<f64>,
}

impl Greeks {
    pub fn is_complete(&self) -> bool {
        self.delta.is_some() && self.gamma.is_some() && self.theta.is_some() && self.vega.is_some()
    }
}

impl std::fmt::Display for Greeks {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "Delta: {:.6}, Gamma: {:.6}, Theta: {:.6}, Vega: {:.6}",
            self.delta.unwrap_or(0.0),
            self.gamma.unwrap_or(0.0),
            self.theta.unwrap_or(0.0),
            self.vega.unwrap_or(0.0)
        )
    }
}
