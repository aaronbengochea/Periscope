"use client";

import { useState, useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { fetchOptionsChain, fetchContractDetails, OptionContract } from "@/lib/api";
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

  // Get contracts for selected expiration
  const contractsForExpiration = useMemo(() => {
    if (!data || !selectedExpiration) return [];
    return data.results.filter(
      contract => contract.details?.expiration_date === selectedExpiration
    );
  }, [data, selectedExpiration]);

  // Extract contract tickers for the selected expiration
  const contractTickers = useMemo(() => {
    return contractsForExpiration
      .map(c => c.details?.ticker)
      .filter((ticker): ticker is string => !!ticker);
  }, [contractsForExpiration]);

  // Fetch detailed data for selected expiration contracts
  const { data: detailedData } = useQuery({
    queryKey: ["contract-details", ticker, selectedExpiration],
    queryFn: () => fetchContractDetails(contractTickers),
    enabled: contractTickers.length > 0,
    staleTime: 30000, // Cache for 30 seconds
  });

  // Merge detailed data with basic contract data
  const enrichedContracts = useMemo(() => {
    if (!detailedData) return contractsForExpiration;

    // Create a map of detailed contracts by ticker for quick lookup
    const detailsMap = new Map<string, OptionContract>();
    detailedData.results.forEach(contract => {
      if (contract.details?.ticker) {
        detailsMap.set(contract.details.ticker, contract);
      }
    });

    // Merge detailed data into basic contracts
    return contractsForExpiration.map(contract => {
      const ticker = contract.details?.ticker;
      if (ticker && detailsMap.has(ticker)) {
        const detailed = detailsMap.get(ticker)!;
        // Merge: detailed data takes precedence
        return {
          ...contract,
          ...detailed,
          // Ensure details aren't overwritten with undefined
          details: detailed.details || contract.details,
        };
      }
      return contract;
    });
  }, [contractsForExpiration, detailedData]);

  // Filter contracts by selected expiration
  const filteredData = data && selectedExpiration
    ? {
        ...data,
        results: enrichedContracts,
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

          {/* Search Bar and Current Price Row */}
          <div className="flex gap-4 mb-4 items-center w-full">
            <form onSubmit={handleSearch} className="flex gap-2 items-center">
              <input
                type="text"
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
                placeholder="Enter ticker symbol (e.g., AAPL, SPY)"
                className="w-80 px-4 py-2 bg-[#1a1a1a] border border-gray-700 rounded-lg focus:outline-none focus:border-blue-500 text-white"
              />
              <button
                type="submit"
                className="px-6 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg font-medium transition-colors"
              >
                Search
              </button>
            </form>

            {/* Current Price Display - Right Side */}
            {currentPrice > 0 && (
              <div className="ml-auto flex items-center gap-3 px-4 py-2 bg-[#1a1a1a] border border-gray-700 rounded-lg">
                <span className="text-sm text-gray-400">Current Price:</span>
                <span className="text-xl font-bold text-yellow-400">${currentPrice.toFixed(2)}</span>
                <span className="text-sm text-gray-500">({underlyingTicker})</span>
              </div>
            )}
          </div>

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
