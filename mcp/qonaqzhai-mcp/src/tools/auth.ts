import { z } from "zod";
import { gatewayRequest } from "../client.js";

export const loginInput = z.object({
  email: z.string().email(),
  password: z.string().min(1),
});

export async function login(input: z.infer<typeof loginInput>) {
  const body = await gatewayRequest<{
    token: string;
    refreshToken?: string;
    user: { id: string; email: string; role: string; status: string };
  }>("/api/login", { method: "POST", auth: false, body: input });
  process.env.QONAQZHAI_TOKEN = body.token;
  return body;
}

export async function me() {
  return gatewayRequest("/api/me");
}
