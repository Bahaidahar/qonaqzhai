import { z } from "zod";
import { gatewayRequest } from "../client.js";

export const createInput = z
  .object({
    vendorId: z.string().uuid(),
    eventDate: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, "YYYY-MM-DD"),
    guestCount: z.number().int().min(1),
    note: z.string().optional(),
    amount: z.number().int().min(0).optional(),
    serviceId: z.string().uuid().optional(),
  })
  .strict();

export async function create(input: z.infer<typeof createInput>) {
  return gatewayRequest("/api/bookings", {
    method: "POST",
    body: input,
  });
}

export async function listMine() {
  return gatewayRequest<{ items: unknown[] | null }>("/api/bookings");
}

export const idInput = z.object({ id: z.string().uuid() }).strict();

export async function getOne(input: z.infer<typeof idInput>) {
  return gatewayRequest(`/api/bookings/${input.id}`);
}

export const transitionInput = z
  .object({
    id: z.string().uuid(),
    status: z.enum(["accepted", "declined", "cancelled", "completed"]),
  })
  .strict();

export async function transition(input: z.infer<typeof transitionInput>) {
  return gatewayRequest(`/api/bookings/${input.id}`, {
    method: "PATCH",
    body: { status: input.status },
  });
}

export async function pay(input: z.infer<typeof idInput>) {
  return gatewayRequest(`/api/bookings/${input.id}/pay`, { method: "POST" });
}

export async function payMock(input: z.infer<typeof idInput>) {
  return gatewayRequest(`/api/bookings/${input.id}/pay/mock`, { method: "POST" });
}
