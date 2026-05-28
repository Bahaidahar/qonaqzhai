import { z } from "zod";
import { gatewayRequest } from "../client.js";

export const chatInput = z
  .object({
    message: z.string().min(1),
    chatId: z.string().uuid().optional(),
  })
  .strict();

export async function chat(input: z.infer<typeof chatInput>) {
  return gatewayRequest("/api/chat", { method: "POST", body: input });
}

export async function listChats() {
  return gatewayRequest("/api/chats");
}

export const idInput = z.object({ id: z.string().uuid() }).strict();

export async function getChat(input: z.infer<typeof idInput>) {
  return gatewayRequest(`/api/chats/${input.id}`);
}
