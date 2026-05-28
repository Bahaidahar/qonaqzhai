#!/usr/bin/env node
/**
 * qonaqzhai MCP server.
 *
 * Exposes the gateway as a curated set of MCP tools so a Claude-compatible
 * client can plan events, browse vendors, book, message, and manage a vendor
 * profile through structured calls instead of free-form HTTP.
 *
 * Transport: stdio (default for Claude Desktop / Claude Code).
 * Auth: bearer token from `QONAQZHAI_TOKEN`; tools accept it on `auth_login`
 *       and stash the result in the same env var for the rest of the session.
 */
import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
  type Tool,
} from "@modelcontextprotocol/sdk/types.js";
import { z } from "zod";

import { GatewayError } from "./client.js";
import * as auth from "./tools/auth.js";
import * as vendors from "./tools/vendors.js";
import * as bookings from "./tools/bookings.js";
import * as vendorSelf from "./tools/vendor_self.js";
import * as messaging from "./tools/messaging.js";
import * as notifications from "./tools/notifications.js";
import * as aiChat from "./tools/ai_chat.js";
import * as reviews from "./tools/reviews.js";
import * as cards from "./tools/cards.js";

interface ToolDef {
  name: string;
  description: string;
  schema?: z.ZodTypeAny;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  handler: (input: any) => Promise<unknown>;
}

// Helper to register tools with a single signature regardless of input shape.
function tool(def: ToolDef): ToolDef {
  return def;
}

const TOOLS: ToolDef[] = [
  // --- auth ---
  tool({
    name: "auth_login",
    description:
      "Log in with email + password. Stores the bearer token in QONAQZHAI_TOKEN so subsequent tools work.",
    schema: auth.loginInput,
    handler: auth.login,
  }),
  tool({
    name: "auth_me",
    description: "Return the current user (requires QONAQZHAI_TOKEN).",
    handler: () => auth.me(),
  }),

  // --- public vendor catalog ---
  tool({
    name: "vendors_search",
    description:
      "Search the vendor catalog. Filters by category, city (defaults to Almaty), price, rating, sort.",
    schema: vendors.searchInput,
    handler: vendors.search,
  }),
  tool({
    name: "vendors_get",
    description: "Fetch a single vendor by id.",
    schema: vendors.getInput,
    handler: vendors.get,
  }),
  tool({
    name: "vendors_reviews",
    description: "List reviews for a vendor.",
    schema: vendors.getInput,
    handler: vendors.listReviews,
  }),
  tool({
    name: "vendors_services",
    description: "List the service packages (price + unit) the vendor offers.",
    schema: vendors.getInput,
    handler: vendors.listServices,
  }),

  // --- bookings (customer + vendor) ---
  tool({
    name: "bookings_create",
    description:
      "Create a booking request against a vendor. The vendor reviews and accepts/declines before any charge.",
    schema: bookings.createInput,
    handler: bookings.create,
  }),
  tool({
    name: "bookings_list",
    description:
      "List my bookings. Customer sees their requests; vendor sees their incoming inbox.",
    handler: () => bookings.listMine(),
  }),
  tool({
    name: "bookings_get",
    description: "Fetch a single booking by id.",
    schema: bookings.idInput,
    handler: bookings.getOne,
  }),
  tool({
    name: "bookings_transition",
    description:
      "Update a booking status. Vendor uses accepted/declined; customer uses cancelled.",
    schema: bookings.transitionInput,
    handler: bookings.transition,
  }),
  tool({
    name: "bookings_pay",
    description:
      "Start a real payment for a booking. Returns a redirectUrl to PayBox; opens an external session.",
    schema: bookings.idInput,
    handler: bookings.pay,
  }),
  tool({
    name: "bookings_pay_mock",
    description: "Mock-pay a booking (development only).",
    schema: bookings.idInput,
    handler: bookings.payMock,
  }),

  // --- vendor profile (vendor role) ---
  tool({
    name: "vendor_profile_get",
    description: "Return the calling vendor's profile.",
    handler: () => vendorSelf.get(),
  }),
  tool({
    name: "vendor_profile_upsert",
    description: "Create or replace the calling vendor's public listing.",
    schema: vendorSelf.upsertInput,
    handler: vendorSelf.upsert,
  }),
  tool({
    name: "vendor_services_list",
    description: "List the calling vendor's service packages.",
    handler: () => vendorSelf.listServices(),
  }),
  tool({
    name: "vendor_services_add",
    description: "Add a service package to the calling vendor's listing.",
    schema: vendorSelf.serviceInput,
    handler: vendorSelf.addService,
  }),
  tool({
    name: "vendor_services_update",
    description: "Update a service package on the calling vendor's listing.",
    schema: vendorSelf.updateServiceInput,
    handler: vendorSelf.updateService,
  }),
  tool({
    name: "vendor_services_delete",
    description: "Delete a service package from the calling vendor's listing.",
    schema: vendorSelf.idInput,
    handler: vendorSelf.deleteService,
  }),

  // --- messaging ---
  tool({
    name: "threads_list",
    description:
      "List my booking-bound chat threads (opened when the vendor accepts a booking).",
    handler: () => messaging.listThreads(),
  }),
  tool({
    name: "threads_get",
    description: "Fetch one thread and its message history.",
    schema: messaging.idInput,
    handler: messaging.getThread,
  }),
  tool({
    name: "threads_send",
    description: "Send a text message into a thread.",
    schema: messaging.sendInput,
    handler: messaging.send,
  }),

  // --- notifications ---
  tool({
    name: "notifications_list",
    description: "Return the calling user's notification inbox.",
    handler: () => notifications.list(),
  }),

  // --- AI chat ---
  tool({
    name: "ai_chat_send",
    description:
      "Send a message to qonaqzhai's planning AI. Returns the assistant message + structured blocks (plan / budget / vendors). Pass `chatId` to continue an existing session, omit for a new one.",
    schema: aiChat.chatInput,
    handler: aiChat.chat,
  }),
  tool({
    name: "ai_chat_list",
    description: "List my saved AI chats.",
    handler: () => aiChat.listChats(),
  }),
  tool({
    name: "ai_chat_get",
    description: "Fetch one saved AI chat with its full message history.",
    schema: aiChat.idInput,
    handler: aiChat.getChat,
  }),

  // --- reviews ---
  tool({
    name: "reviews_submit",
    description: "Submit a 1-5 star review against a completed booking.",
    schema: reviews.submitInput,
    handler: reviews.submit,
  }),

  // --- payment cards ---
  tool({
    name: "cards_list",
    description: "List my saved payment cards.",
    handler: () => cards.list(),
  }),
  tool({
    name: "cards_add",
    description:
      "Save a payment card. The PAN is tokenised by the payment service — only the last 4 digits are stored.",
    schema: cards.addInput,
    handler: cards.add,
  }),
  tool({
    name: "cards_set_default",
    description: "Promote a saved card to default.",
    schema: cards.idInput,
    handler: cards.setDefault,
  }),
  tool({
    name: "cards_delete",
    description: "Delete a saved card.",
    schema: cards.idInput,
    handler: cards.remove,
  }),
];

const TOOLS_BY_NAME = new Map(TOOLS.map((t) => [t.name, t]));

function toolDescriptor(def: ToolDef): Tool {
  const schemaJson = def.schema
    ? jsonSchemaFromZod(def.schema as z.ZodTypeAny)
    : { type: "object" as const, properties: {} };
  return {
    name: def.name,
    description: def.description,
    inputSchema: schemaJson,
  };
}

// Minimal Zod → JSON Schema generator covering the shapes the tools use. The
// official @modelcontextprotocol/sdk does not ship a transformer for Zod, and
// pulling zod-to-json-schema for a handful of shapes is overkill.
function jsonSchemaFromZod(schema: z.ZodTypeAny): Tool["inputSchema"] {
  return walk(schema) as Tool["inputSchema"];
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function walk(s: z.ZodTypeAny): any {
  // Unwrap default + optional + nullable so the inner type drives the schema.
  if (s instanceof z.ZodDefault) return walk(s._def.innerType);
  if (s instanceof z.ZodOptional) return walk(s._def.innerType);
  if (s instanceof z.ZodNullable) return { ...walk(s._def.innerType), nullable: true };

  if (s instanceof z.ZodObject) {
    const shape = s.shape as Record<string, z.ZodTypeAny>;
    const properties: Record<string, unknown> = {};
    const required: string[] = [];
    for (const [k, v] of Object.entries(shape)) {
      properties[k] = walk(v);
      if (!(v instanceof z.ZodOptional) && !(v instanceof z.ZodDefault)) {
        required.push(k);
      }
    }
    const out: Record<string, unknown> = {
      type: "object",
      properties,
      additionalProperties: false,
    };
    if (required.length > 0) out.required = required;
    return out;
  }
  if (s instanceof z.ZodString) {
    const out: Record<string, unknown> = { type: "string" };
    if (s.description) out.description = s.description;
    return out;
  }
  if (s instanceof z.ZodNumber) return { type: "number" };
  if (s instanceof z.ZodBoolean) return { type: "boolean" };
  if (s instanceof z.ZodEnum) {
    return { type: "string", enum: s.options as unknown as string[] };
  }
  if (s instanceof z.ZodLiteral) {
    return { const: s.value };
  }
  if (s instanceof z.ZodArray) {
    return { type: "array", items: walk(s.element) };
  }
  return {};
}

async function dispatch(
  name: string,
  rawArgs: unknown,
): Promise<{ content: { type: "text"; text: string }[]; isError?: boolean }> {
  const def = TOOLS_BY_NAME.get(name);
  if (!def) {
    return {
      content: [{ type: "text", text: `unknown tool: ${name}` }],
      isError: true,
    };
  }
  try {
    const input = def.schema ? def.schema.parse(rawArgs ?? {}) : (rawArgs ?? {});
    const result = await def.handler(input);
    return {
      content: [
        { type: "text", text: JSON.stringify(result, null, 2) },
      ],
    };
  } catch (err) {
    if (err instanceof GatewayError) {
      return {
        content: [
          {
            type: "text",
            text: `gateway error ${err.status}: ${err.message}\n${JSON.stringify(err.body, null, 2)}`,
          },
        ],
        isError: true,
      };
    }
    if (err instanceof z.ZodError) {
      return {
        content: [{ type: "text", text: `invalid arguments: ${err.message}` }],
        isError: true,
      };
    }
    const message = err instanceof Error ? err.message : String(err);
    return {
      content: [{ type: "text", text: `error: ${message}` }],
      isError: true,
    };
  }
}

async function main(): Promise<void> {
  const server = new Server(
    {
      name: "qonaqzhai-mcp",
      version: "0.1.0",
    },
    {
      capabilities: { tools: {} },
    },
  );

  server.setRequestHandler(ListToolsRequestSchema, async () => ({
    tools: TOOLS.map((t) => toolDescriptor(t as ToolDef)),
  }));

  server.setRequestHandler(CallToolRequestSchema, async (req) => {
    return dispatch(req.params.name, req.params.arguments);
  });

  const transport = new StdioServerTransport();
  await server.connect(transport);
}

main().catch((err) => {
  // Log to stderr so stdio transport stays clean for MCP framing on stdout.
  // eslint-disable-next-line no-console
  console.error("qonaqzhai-mcp fatal:", err);
  process.exit(1);
});
