"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { fetchOptionsChain } from "@/lib/api";
import { OptionsChain } from "@/components/OptionsChain";

export default function Home() {
  const [ticker, setTicker] = useState("AAPL");
  const [searchInput, setSearchInput] = useState("AAPL");

  const { data, isLoading, error } = useQuery({
    queryKey: ["options", ticker],
    queryFn: () => fetchOptionsChain(ticker),
    enabled: !!ticker,
  });

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
          <form onSubmit={handleSearch} className="flex gap-2">
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
        {data && !isLoading && (
          <OptionsChain
            data={data}
            currentPrice={data.results[0]?.underlying_asset?.price || 0}
          />
        )}
      </div>
    </main>
  );
}
