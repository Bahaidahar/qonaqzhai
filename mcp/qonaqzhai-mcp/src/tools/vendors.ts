import { z } from "zod";
import { gatewayRequest } from "../client.js";

export const searchInput = z
  .object({
    q: z.string().optional().describe("Free-text query."),
    category: z
      .enum(["Venue", "Catering", "Music & DJ", "Photo & Video", "Decor & Florists", "Cakes", "Other"])
      .optional(),
    city: z.string().default("Almaty"),
    minPrice: z.number().int().min(0).optional(),
    maxPrice: z.number().int().min(0).optional(),
    minRating: z.number().min(0).max(5).optional(),
    sort: z.enum(["newest", "price_asc", "price_desc", "rating_desc"]).default("newest"),
    page: z.number().int().min(1).default(1),
    limit: z.number().int().min(1).max(50).default(20),
  })
  .strict();

export async function search(input: z.infer<typeof searchInput>) {
  return gatewayRequest<{ items: unknown[]; total: number }>("/api/vendors", {
    auth: false,
    query: {
      q: input.q,
      category: input.category,
      city: input.city,
      min_price: input.minPrice,
      max_price: input.maxPrice,
      min_rating: input.minRating,
      sort: input.sort,
      page: input.page,
      limit: input.limit,
    },
  });
}

export const getInput = z.object({ id: z.string().uuid() }).strict();

export async function get(input: z.infer<typeof getInput>) {
  return gatewayRequest(`/api/vendors/${input.id}`, { auth: false });
}

export async function listReviews(input: z.infer<typeof getInput>) {
  return gatewayRequest(`/api/vendors/${input.id}/reviews`, { auth: false });
}

export async function listServices(input: z.infer<typeof getInput>) {
  return gatewayRequest(`/api/vendors/${input.id}/services`, { auth: false });
}
