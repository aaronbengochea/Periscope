import { OptionsChainResponse, OptionContract } from "@/lib/api";

interface OptionsChainProps {
  data: OptionsChainResponse;
  currentPrice: number;
}

interface StrikeRow {
  strike: number;
  call?: OptionContract;
  put?: OptionContract;
}

export function OptionsChain({ data, currentPrice }: OptionsChainProps) {
  // Group contracts by strike price
  const strikeMap = new Map<number, StrikeRow>();

  data.results.forEach((contract) => {
    const strike = contract.details?.strike_price;
    const type = contract.details?.contract_type;

    if (!strike || !type) return;

    if (!strikeMap.has(strike)) {
      strikeMap.set(strike, { strike });
    }

    const row = strikeMap.get(strike)!;
    if (type === "call") {
      row.call = contract;
    } else if (type === "put") {
      row.put = contract;
    }
  });

  // Sort by strike price
  const strikes = Array.from(strikeMap.values()).sort((a, b) => a.strike - b.strike);

  // Calculate % change
  const calculateChange = (contract: OptionContract) => {
    const open = contract.day?.open;
    const current = contract.last_trade?.price;
    if (open && current && open > 0) {
      return ((current - open) / open) * 100;
    }
    return null;
  };

  // Calculate net change ($)
  const calculateNetChange = (contract: OptionContract) => {
    const open = contract.day?.open;
    const current = contract.last_trade?.price;
    if (open && current) {
      return current - open;
    }
    return null;
  };

  // Format number with decimals
  const fmt = (num?: number, decimals = 2) => {
    if (num === undefined || num === null) return "-";
    return num.toFixed(decimals);
  };

  // Format large numbers (volume, OI)
  const fmtLarge = (num?: number) => {
    if (num === undefined || num === null) return "-";
    if (num >= 1000) {
      return (num / 1000).toFixed(1) + "K";
    }
    return num.toString();
  };

  return (
    <div className="overflow-x-auto">
      <div className="inline-block min-w-full">
        {/* Header */}
        <div className="grid grid-cols-[1fr_auto_1fr] gap-0 bg-[#1a1a1a] border-b border-gray-700">
          {/* CALLS Header */}
          <div className="grid grid-cols-11 gap-px bg-gray-700 p-px">
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Bid</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Ask</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">IV</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Δ</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Γ</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Θ</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">V</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">OI</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Δ%</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Vol</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Δ$</div>
          </div>

          {/* STRIKE Header */}
          <div className="bg-[#1a1a1a] px-4 py-3 text-xs font-semibold text-center min-w-[80px]">
            STRIKE
          </div>

          {/* PUTS Header */}
          <div className="grid grid-cols-11 gap-px bg-gray-700 p-px">
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Δ$</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Vol</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Δ%</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">OI</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">V</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Θ</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Γ</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Δ</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">IV</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Ask</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Bid</div>
          </div>
        </div>

        {/* Data Rows */}
        {strikes.map((row) => {
          const isATM = Math.abs(row.strike - currentPrice) < 5; // Within $5 of current price
          const callChange = row.call ? calculateChange(row.call) : null;
          const putChange = row.put ? calculateChange(row.put) : null;
          const callNetChange = row.call ? calculateNetChange(row.call) : null;
          const putNetChange = row.put ? calculateNetChange(row.put) : null;

          return (
            <div
              key={row.strike}
              className={`grid grid-cols-[1fr_auto_1fr] gap-0 border-b border-gray-800 hover:bg-[#1a1a1a] ${
                isATM ? "bg-yellow-900/20" : ""
              }`}
            >
              {/* CALLS Data */}
              <div className="grid grid-cols-11 gap-px bg-gray-800 p-px">
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmt(row.call?.last_quote?.bid)}</div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmt(row.call?.last_quote?.ask)}</div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">
                  {row.call?.implied_volatility ? `${fmt(row.call.implied_volatility * 100, 1)}%` : "-"}
                </div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmt(row.call?.greeks?.delta)}</div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmt(row.call?.greeks?.gamma, 3)}</div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmt(row.call?.greeks?.theta)}</div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmt(row.call?.greeks?.vega)}</div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmtLarge(row.call?.open_interest)}</div>
                <div className={`bg-[#0a0a0a] px-2 py-2 text-xs text-center ${
                  callChange !== null ? (callChange >= 0 ? "text-green-400" : "text-red-400") : ""
                }`}>
                  {callChange !== null ? `${fmt(callChange, 1)}%` : "-"}
                </div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmtLarge(row.call?.day?.volume)}</div>
                <div className={`bg-[#0a0a0a] px-2 py-2 text-xs text-center ${
                  callNetChange !== null ? (callNetChange >= 0 ? "text-green-400" : "text-red-400") : ""
                }`}>
                  {callNetChange !== null ? fmt(callNetChange) : "-"}
                </div>
              </div>

              {/* STRIKE */}
              <div className={`bg-[#0a0a0a] px-4 py-2 text-sm font-semibold text-center min-w-[80px] ${
                isATM ? "text-yellow-400" : ""
              }`}>
                {fmt(row.strike)}
              </div>

              {/* PUTS Data */}
              <div className="grid grid-cols-11 gap-px bg-gray-800 p-px">
                <div className={`bg-[#0a0a0a] px-2 py-2 text-xs text-center ${
                  putNetChange !== null ? (putNetChange >= 0 ? "text-green-400" : "text-red-400") : ""
                }`}>
                  {putNetChange !== null ? fmt(putNetChange) : "-"}
                </div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmtLarge(row.put?.day?.volume)}</div>
                <div className={`bg-[#0a0a0a] px-2 py-2 text-xs text-center ${
                  putChange !== null ? (putChange >= 0 ? "text-green-400" : "text-red-400") : ""
                }`}>
                  {putChange !== null ? `${fmt(putChange, 1)}%` : "-"}
                </div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmtLarge(row.put?.open_interest)}</div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmt(row.put?.greeks?.vega)}</div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmt(row.put?.greeks?.theta)}</div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmt(row.put?.greeks?.gamma, 3)}</div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmt(row.put?.greeks?.delta)}</div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">
                  {row.put?.implied_volatility ? `${fmt(row.put.implied_volatility * 100, 1)}%` : "-"}
                </div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmt(row.put?.last_quote?.ask)}</div>
                <div className="bg-[#0a0a0a] px-2 py-2 text-xs text-center">{fmt(row.put?.last_quote?.bid)}</div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Current Price Indicator */}
      <div className="mt-4 text-sm text-gray-400">
        Current Price: <span className="text-white font-semibold">${fmt(currentPrice)}</span>
      </div>
    </div>
  );
}
