import { z } from "zod";
import { gatewayRequest } from "../client.js";

export async function list() {
  return gatewayRequest("/api/cards");
}

export const addInput = z
  .object({
    number: z.string().regex(/^\d{13,19}$/, "PAN must be 13-19 digits"),
    expMonth: z.number().int().min(1).max(12),
    expYear: z.number().int().min(2024).max(2100),
    holder: z.string().min(1),
    makeDefault: z.boolean().default(false),
  })
  .strict();

export async function add(input: z.infer<typeof addInput>) {
  return gatewayRequest("/api/cards", { method: "POST", body: input });
}

export const idInput = z.object({ id: z.string().uuid() }).strict();

export async function setDefault(input: z.infer<typeof idInput>) {
  return gatewayRequest(`/api/cards/${input.id}/default`, { method: "POST" });
}

export async function remove(input: z.infer<typeof idInput>) {
  return gatewayRequest(`/api/cards/${input.id}`, { method: "DELETE" });
}
