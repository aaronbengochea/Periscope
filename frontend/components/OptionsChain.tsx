import React from "react";
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

  // Find the index where we should insert the current price row
  // Insert it between the strike just below and just above the current price
  const currentPriceIndex = strikes.findIndex(row => row.strike >= currentPrice);

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
        {/* CALLS / PUTS Label Row */}
        <div className="grid grid-cols-[1fr_auto_1fr] gap-0 bg-[#0a0a0a] border-b-2 border-blue-500">
          {/* CALLS Label */}
          <div className="px-4 py-2 text-center">
            <span className="text-lg font-bold text-green-400">CALLS</span>
          </div>

          {/* Empty center */}
          <div className="px-4 py-2 min-w-[80px]"></div>

          {/* PUTS Label */}
          <div className="px-4 py-2 text-center">
            <span className="text-lg font-bold text-red-400">PUTS</span>
          </div>
        </div>

        {/* Column Headers */}
        <div className="grid grid-cols-[1fr_auto_1fr] gap-0 bg-[#1a1a1a] border-b border-gray-700">
          {/* CALLS Column Headers */}
          <div className="grid grid-cols-11 gap-px bg-gray-700 p-px">
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Bid</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Ask</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">IV</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Delta</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Gamma</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Theta</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Vega</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">OI</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Chg%</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Vol</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Chg$</div>
          </div>

          {/* STRIKE Header */}
          <div className="bg-[#1a1a1a] px-4 py-3 text-xs font-semibold text-center min-w-[80px]">
            STRIKE
          </div>

          {/* PUTS Column Headers */}
          <div className="grid grid-cols-11 gap-px bg-gray-700 p-px">
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Chg$</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Vol</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Chg%</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">OI</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Vega</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Theta</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Gamma</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Delta</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">IV</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Ask</div>
            <div className="bg-[#1a1a1a] px-2 py-3 text-xs font-semibold text-center">Bid</div>
          </div>
        </div>

        {/* Data Rows */}
        {strikes.map((row, index) => {
          const isATM = Math.abs(row.strike - currentPrice) < 5; // Within $5 of current price
          const callChange = row.call ? calculateChange(row.call) : null;
          const putChange = row.put ? calculateChange(row.put) : null;
          const callNetChange = row.call ? calculateNetChange(row.call) : null;
          const putNetChange = row.put ? calculateNetChange(row.put) : null;

          // Check if we should insert current price row before this strike
          const showCurrentPriceRow = index === currentPriceIndex && currentPrice > 0 && currentPrice < row.strike;

          return (
            <React.Fragment key={`strike-${row.strike}`}>
              {/* Current Price Row - inserted at the correct position */}
              {showCurrentPriceRow && (
                <div
                  key={`current-price-${currentPrice}`}
                  className="grid grid-cols-[1fr_auto_1fr] gap-0 border-y-2 border-yellow-400 bg-yellow-900/30"
                >
                  {/* Empty left side (Calls) */}
                  <div className="bg-[#0a0a0a] px-2 py-3"></div>

                  {/* Current Price - Center */}
                  <div className="bg-[#0a0a0a] px-4 py-3 text-base font-bold text-center min-w-[80px] text-yellow-400">
                    ${fmt(currentPrice)} ← Current
                  </div>

                  {/* Empty right side (Puts) */}
                  <div className="bg-[#0a0a0a] px-2 py-3"></div>
                </div>
              )}

              {/* Regular strike row */}
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
            </React.Fragment>
          );
        })}
      </div>

      {/* Current Price Indicator */}
      <div className="mt-4 flex items-center gap-3">
        <span className="text-sm text-gray-400">Current Price:</span>
        {currentPrice > 0 ? (
          <>
            <span className="text-xl font-bold text-yellow-400">${fmt(currentPrice)}</span>
            {data.results.find(c => c.underlying_asset?.ticker)?.underlying_asset?.ticker && (
              <span className="text-sm text-gray-500">
                ({data.results.find(c => c.underlying_asset?.ticker)?.underlying_asset?.ticker})
              </span>
            )}
          </>
        ) : (
          <span className="text-sm text-red-400">Price not available in API response</span>
        )}
      </div>
    </div>
  );
}
