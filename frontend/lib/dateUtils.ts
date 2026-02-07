// Utility functions for formatting expiration dates

export interface ExpirationInfo {
  date: string; // ISO format: YYYY-MM-DD
  daysToExpiry: number;
  dayOfWeek: string;
  formattedDate: string; // "Feb 20" or "Jan 15, 2027"
  displayText: string; // "14D Fri Feb 20"
}

export function calculateDaysToExpiry(expirationDate: string): number {
  // Parse date as local time to avoid timezone issues
  // Format: YYYY-MM-DD
  const [year, month, day] = expirationDate.split('-').map(Number);

  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const expiry = new Date(year, month - 1, day); // month is 0-indexed
  expiry.setHours(0, 0, 0, 0);

  const diffTime = expiry.getTime() - today.getTime();
  const diffDays = Math.round(diffTime / (1000 * 60 * 60 * 24));

  return diffDays;
}

export function getDayOfWeek(date: string): string {
  const days = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

  // Parse date as local time to avoid timezone issues
  const [year, month, day] = date.split('-').map(Number);
  const d = new Date(year, month - 1, day); // month is 0-indexed

  return days[d.getDay()];
}

export function formatExpirationDate(date: string): string {
  // Parse date as local time to avoid timezone issues
  const [year, month, day] = date.split('-').map(Number);
  const d = new Date(year, month - 1, day); // month is 0-indexed

  const currentYear = new Date().getFullYear();
  const expiryYear = d.getFullYear();

  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
                  'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

  const monthName = months[d.getMonth()];
  const dayOfMonth = d.getDate();

  if (expiryYear === currentYear) {
    return `${monthName} ${dayOfMonth}`;
  } else {
    return `${monthName} ${dayOfMonth}, ${expiryYear}`;
  }
}

export function formatExpirationDisplay(date: string): ExpirationInfo {
  const daysToExpiry = calculateDaysToExpiry(date);
  const dayOfWeek = getDayOfWeek(date);
  const formattedDate = formatExpirationDate(date);

  const displayText = `${daysToExpiry}D ${dayOfWeek} ${formattedDate}`;

  return {
    date,
    daysToExpiry,
    dayOfWeek,
    formattedDate,
    displayText,
  };
}

interface ContractWithExpiration {
  details?: {
    expiration_date?: string;
  };
}

export function extractUniqueExpirations(contracts: ContractWithExpiration[]): ExpirationInfo[] {
  const expirationSet = new Set<string>();

  contracts.forEach(contract => {
    const expiry = contract.details?.expiration_date;
    if (expiry) {
      expirationSet.add(expiry);
    }
  });

  const expirations = Array.from(expirationSet)
    .map(date => formatExpirationDisplay(date))
    .filter(exp => exp.daysToExpiry >= 0) // Only include today and future dates
    .sort((a, b) => a.daysToExpiry - b.daysToExpiry);

  return expirations;
}
