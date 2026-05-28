import { z } from "zod";
import { gatewayRequest } from "../client.js";

export async function listThreads() {
  return gatewayRequest("/api/threads");
}

export const idInput = z.object({ id: z.string().uuid() }).strict();

export async function getThread(input: z.infer<typeof idInput>) {
  return gatewayRequest(`/api/threads/${input.id}`);
}

export const sendInput = z
  .object({
    threadId: z.string().uuid(),
    text: z.string().min(1).max(4000),
  })
  .strict();

export async function send(input: z.infer<typeof sendInput>) {
  return gatewayRequest(`/api/threads/${input.threadId}/messages`, {
    method: "POST",
    body: { text: input.text },
  });
}
