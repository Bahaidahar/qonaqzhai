export type BlockType = "plan" | "budget" | "vendors";

export interface PlanBlock {
  type: "plan";
  title: string;
  eventType: string;
  date: string;
  city: string;
  guests: number;
  budget: number;
}

export interface BudgetBlock {
  type: "budget";
  total: number;
  categories: { name: string; amount: number; pct: number }[];
}

export interface VendorsBlock {
  type: "vendors";
  query: string;
  items: {
    id: string;
    name: string;
    category: string;
    rating: number;
    priceFrom: number;
    city: string;
    tags?: string[];
  }[];
}

export type Block = PlanBlock | BudgetBlock | VendorsBlock;

export interface ChatMessage {
  id: string;
  role: "user" | "ai";
  text?: string;
  blocks?: Block[];
  chips?: string[];
  /** When true, the message text animates char-by-char on mount. Stripped on persist. */
  streaming?: boolean;
}
