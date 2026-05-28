/**
 * Thin HTTP client for the qonaqzhai gateway. Resolves the base URL and bearer
 * token from environment variables on each call so the MCP server can be
 * pointed at staging or rotated tokens without a restart.
 */
const BASE_URL = process.env.QONAQZHAI_API ?? "http://localhost:8080";

function token(): string | undefined {
  return process.env.QONAQZHAI_TOKEN || undefined;
}

interface RequestOptions {
  method?: string;
  body?: unknown;
  auth?: boolean;
  query?: Record<string, string | number | undefined>;
}

function buildUrl(path: string, query?: RequestOptions["query"]): string {
  const url = new URL(BASE_URL + path);
  if (query) {
    for (const [k, v] of Object.entries(query)) {
      if (v === undefined || v === null || v === "") continue;
      url.searchParams.set(k, String(v));
    }
  }
  return url.toString();
}

export class GatewayError extends Error {
  constructor(public status: number, message: string, public body: unknown) {
    super(`${status}: ${message}`);
    this.name = "GatewayError";
  }
}

export async function gatewayRequest<T = unknown>(
  path: string,
  opts: RequestOptions = {},
): Promise<T> {
  const headers: Record<string, string> = {
    "content-type": "application/json",
  };
  if (opts.auth ?? true) {
    const tok = token();
    if (!tok) {
      throw new Error(
        "QONAQZHAI_TOKEN is not set — call the `auth_login` tool first or export the token.",
      );
    }
    headers.authorization = `Bearer ${tok}`;
  }

  const res = await fetch(buildUrl(path, opts.query), {
    method: opts.method ?? "GET",
    headers,
    body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
  });

  const text = await res.text();
  let parsed: unknown = null;
  try {
    parsed = text.length > 0 ? JSON.parse(text) : null;
  } catch {
    // gateway sometimes returns a plain string for non-JSON paths
    parsed = text;
  }

  if (!res.ok) {
    const message =
      parsed && typeof parsed === "object" && parsed !== null && "error" in parsed
        ? String((parsed as { error: unknown }).error)
        : `request failed`;
    throw new GatewayError(res.status, message, parsed);
  }

  return parsed as T;
}
