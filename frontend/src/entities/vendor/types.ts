export type VendorStatus = "pending" | "approved" | "rejected";

export interface Vendor {
  id: string;
  userId: string;
  name: string;
  category: string;
  city: string;
  description: string;
  priceFrom: number;
  status: VendorStatus;
  ratingAvg: number;
  ratingCount: number;
  photoIds: string[];
  createdAt: string;
  updatedAt: string;
}

export interface Photo {
  id: string;
  vendorId: string;
  mime: string;
  size: number;
  createdAt: string;
}

export interface VendorSearchParams {
  q?: string;
  category?: string;
  city?: string;
  priceMin?: number;
  priceMax?: number;
  ratingMin?: number;
  sort?: "price_asc" | "price_desc" | "rating_desc" | "newest";
  page?: number;
  limit?: number;
}

export interface VendorSearchResult {
  items: Vendor[] | null;
  total: number;
  page: number;
  limit: number;
}
