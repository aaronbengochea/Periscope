"use client";

import { useState, useRef, useEffect } from "react";
import { ExpirationInfo } from "@/lib/dateUtils";

interface ExpirationDropdownProps {
  expirations: ExpirationInfo[];
  selectedExpiration: string | null;
  onSelect: (date: string) => void;
}

export function ExpirationDropdown({
  expirations,
  selectedExpiration,
  onSelect,
}: ExpirationDropdownProps) {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }

    if (isOpen) {
      document.addEventListener("mousedown", handleClickOutside);
      return () => document.removeEventListener("mousedown", handleClickOutside);
    }
  }, [isOpen]);

  const selectedInfo = expirations.find(e => e.date === selectedExpiration);

  return (
    <div className="relative" ref={dropdownRef}>
      {/* Dropdown Button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-2 px-4 py-2 bg-[#1a1a1a] border border-gray-700 rounded-lg hover:border-gray-600 transition-colors min-w-[200px]"
      >
        <span className="flex-1 text-left text-sm">
          {selectedInfo ? (
            <span className="font-mono">{selectedInfo.displayText}</span>
          ) : (
            <span className="text-gray-400">Select expiration</span>
          )}
        </span>
        <svg
          className={`w-4 h-4 transition-transform ${isOpen ? "rotate-180" : ""}`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      {/* Dropdown Menu */}
      {isOpen && (
        <div className="absolute z-50 mt-2 w-full min-w-[280px] bg-[#1a1a1a] border border-gray-700 rounded-lg shadow-xl max-h-[400px] overflow-y-auto">
          {expirations.length === 0 ? (
            <div className="px-4 py-3 text-sm text-gray-400">No expirations available</div>
          ) : (
            expirations.map((expiration) => (
              <button
                key={expiration.date}
                onClick={() => {
                  onSelect(expiration.date);
                  setIsOpen(false);
                }}
                className={`w-full px-4 py-3 text-left text-sm font-mono hover:bg-[#2a2a2a] transition-colors border-b border-gray-800 last:border-b-0 ${
                  expiration.date === selectedExpiration
                    ? "bg-[#2a2a2a] text-blue-400"
                    : "text-white"
                }`}
              >
                <div className="flex items-center gap-2">
                  <span className="text-yellow-400 font-semibold min-w-[40px]">
                    {expiration.daysToExpiry}D
                  </span>
                  <span className="text-gray-400">{expiration.dayOfWeek}</span>
                  <span>{expiration.formattedDate}</span>
                </div>
              </button>
            ))
          )}
        </div>
      )}
    </div>
  );
}
