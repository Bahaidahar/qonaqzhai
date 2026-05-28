import { gatewayRequest } from "../client.js";

export async function list() {
  return gatewayRequest("/api/notifications");
}
