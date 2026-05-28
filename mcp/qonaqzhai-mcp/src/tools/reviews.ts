import { z } from "zod";
import { gatewayRequest } from "../client.js";

export const submitInput = z
  .object({
    bookingId: z.string().uuid(),
    rating: z.number().int().min(1).max(5),
    text: z.string().optional(),
  })
  .strict();

export async function submit(input: z.infer<typeof submitInput>) {
  return gatewayRequest("/api/reviews", { method: "POST", body: input });
}
