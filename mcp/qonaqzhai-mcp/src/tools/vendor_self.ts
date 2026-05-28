import { z } from "zod";
import { gatewayRequest } from "../client.js";

export async function get() {
  return gatewayRequest("/api/me/vendor");
}

export const upsertInput = z
  .object({
    name: z.string().min(1),
    category: z.enum(["Venue", "Catering", "Music & DJ", "Photo & Video", "Decor & Florists", "Cakes", "Other"]),
    city: z.string().default("Almaty"),
    description: z.string().default(""),
    priceFrom: z.number().int().min(0).default(0),
  })
  .strict();

export async function upsert(input: z.infer<typeof upsertInput>) {
  return gatewayRequest("/api/me/vendor", { method: "POST", body: input });
}

export async function listServices() {
  return gatewayRequest("/api/me/vendor/services");
}

export const serviceInput = z
  .object({
    name: z.string().min(1),
    description: z.string().default(""),
    price: z.number().int().min(0),
    unit: z.enum(["fixed", "hour", "item", "person", "day"]).default("fixed"),
  })
  .strict();

export async function addService(input: z.infer<typeof serviceInput>) {
  return gatewayRequest("/api/me/vendor/services", {
    method: "POST",
    body: input,
  });
}

export const updateServiceInput = serviceInput.extend({ id: z.string().uuid() });

export async function updateService(input: z.infer<typeof updateServiceInput>) {
  const { id, ...rest } = input;
  return gatewayRequest(`/api/me/vendor/services/${id}`, {
    method: "PATCH",
    body: rest,
  });
}

export const idInput = z.object({ id: z.string().uuid() }).strict();

export async function deleteService(input: z.infer<typeof idInput>) {
  return gatewayRequest(`/api/me/vendor/services/${input.id}`, { method: "DELETE" });
}
