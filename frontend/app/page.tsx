"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { fetchOptionsChain } from "@/lib/api";
import { OptionsChain } from "@/components/OptionsChain";
import { ExpirationDropdown } from "@/components/ExpirationDropdown";
import { extractUniqueExpirations } from "@/lib/dateUtils";

export default function Home() {
  const [ticker, setTicker] = useState("AAPL");
  const [searchInput, setSearchInput] = useState("AAPL");
  const [selectedExpiration, setSelectedExpiration] = useState<string | null>(null);

  const { data, isLoading, error } = useQuery({
    queryKey: ["options", ticker],
    queryFn: () => fetchOptionsChain(ticker),
    enabled: !!ticker,
  });

  // Extract current price and ticker from the first contract that has it
  const currentPrice = data?.results?.find(
    (contract) => contract.underlying_asset?.price
  )?.underlying_asset?.price || 0;

  const underlyingTicker = data?.results?.find(
    (contract) => contract.underlying_asset?.ticker
  )?.underlying_asset?.ticker || ticker;

  // Extract unique expiration dates
  const expirations = data ? extractUniqueExpirations(data.results) : [];

  // Auto-select first expiration when data loads
  if (data && expirations.length > 0 && !selectedExpiration) {
    setSelectedExpiration(expirations[0].date);
  }

  // Filter contracts by selected expiration
  const filteredData = data && selectedExpiration
    ? {
        ...data,
        results: data.results.filter(
          contract => contract.details?.expiration_date === selectedExpiration
        ),
      }
    : data;

  // Debug: log current price extraction
  if (data && data.results.length > 0) {
    console.log("[Frontend Page] Current price extraction:");
    console.log("  - Searching through", data.results.length, "contracts");
    console.log("  - First contract underlying_asset:", data.results[0]?.underlying_asset);
    console.log("  - Extracted current price:", currentPrice);
    console.log("  - Total expirations:", expirations.length);
    console.log("  - Selected expiration:", selectedExpiration);
    console.log("  - Filtered contracts:", filteredData?.results.length);

    if (currentPrice === 0) {
      console.warn("[Frontend Page] ⚠ Current price is 0! Checking all contracts...");
      data.results.slice(0, 5).forEach((contract, i) => {
        console.log(`  - Contract ${i} underlying_asset:`, contract.underlying_asset);
      });
    } else {
      console.log(`[Frontend Page] ✓ Current price found: $${currentPrice}`);
    }
  }

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchInput.trim()) {
      setTicker(searchInput.toUpperCase().trim());
    }
  };

  return (
    <main className="min-h-screen bg-[#0a0a0a] text-white p-6">
      <div className="max-w-[95%] mx-auto">
        {/* Header */}
        <div className="mb-6">
          <h1 className="text-3xl font-bold mb-4">Periscope</h1>

          {/* Search Bar */}
          <form onSubmit={handleSearch} className="flex gap-2 mb-4">
            <input
              type="text"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              placeholder="Enter ticker symbol (e.g., AAPL, SPY)"
              className="flex-1 max-w-md px-4 py-2 bg-[#1a1a1a] border border-gray-700 rounded-lg focus:outline-none focus:border-blue-500 text-white"
            />
            <button
              type="submit"
              className="px-6 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg font-medium transition-colors"
            >
              Search
            </button>
          </form>

          {/* Expiration Dropdown */}
          {expirations.length > 0 && (
            <div className="flex items-center gap-3">
              <label className="text-sm text-gray-400">Expiration:</label>
              <ExpirationDropdown
                expirations={expirations}
                selectedExpiration={selectedExpiration}
                onSelect={setSelectedExpiration}
              />
            </div>
          )}
        </div>

        {/* Loading State */}
        {isLoading && (
          <div className="text-center py-12">
            <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-blue-500 border-r-transparent"></div>
            <p className="mt-4 text-gray-400">Loading options chain for {ticker}...</p>
          </div>
        )}

        {/* Error State */}
        {error && (
          <div className="bg-red-900/20 border border-red-500/50 rounded-lg p-4">
            <p className="text-red-400">Failed to load options chain: {error.message}</p>
          </div>
        )}

        {/* Options Chain */}
        {filteredData && !isLoading && (
          <OptionsChain
            data={filteredData}
            currentPrice={currentPrice}
          />
        )}
      </div>
    </main>
  );
}
