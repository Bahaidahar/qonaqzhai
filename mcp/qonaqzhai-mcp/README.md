# qonaqzhai MCP server

Exposes the qonaqzhai gateway (`http://localhost:8080`) as
[Model Context Protocol](https://modelcontextprotocol.io) tools. Lets
Claude Desktop / Claude Code call the backend with strongly-typed JSON
arguments instead of free-form HTTP.

## Layout

```
src/
  client.ts            HTTP client wrapping the gateway + bearer token
  index.ts             MCP server entry — registers tools, routes calls
  tools/
    auth.ts            login + me
    vendors.ts         public catalog: search, get, reviews, services
    bookings.ts        create/list/transition/pay
    vendor_self.ts     vendor profile + services CRUD
    messaging.ts       booking-bound chat threads
    notifications.ts   inbox
    ai_chat.ts         qonaqzhai AI planner
    reviews.ts         submit
    cards.ts           saved payment cards
```

## Tools

| Tool | Description |
| --- | --- |
| `auth_login` | email + password → bearer token (stashed in env) |
| `auth_me` | current user |
| `vendors_search` | catalog with filters + sort |
| `vendors_get` | one vendor |
| `vendors_reviews` | reviews for a vendor |
| `vendors_services` | service packages for a vendor |
| `bookings_create` | request booking (customer) |
| `bookings_list` | my bookings (role-aware) |
| `bookings_get` | one booking |
| `bookings_transition` | accept / decline / cancel |
| `bookings_pay` / `bookings_pay_mock` | real PayBox / dev mock |
| `vendor_profile_get` / `vendor_profile_upsert` | vendor's own listing |
| `vendor_services_*` | vendor service CRUD |
| `threads_list` / `threads_get` / `threads_send` | booking chats |
| `notifications_list` | inbox |
| `ai_chat_send` / `ai_chat_list` / `ai_chat_get` | AI planner |
| `reviews_submit` | star + text review |
| `cards_*` | saved card list / add / default / delete |

## Build + run

```bash
cd mcp/qonaqzhai-mcp
npm install            # or pnpm install
npm run build
```

The server speaks **stdio** — point an MCP client at the `dist/index.js`
binary. Two ways to authenticate:

1. `auth_login` tool — call it first; the bearer is kept in `QONAQZHAI_TOKEN`
   for the rest of the session.
2. Pre-set `QONAQZHAI_TOKEN` in the environment (handy for CI / scripts).

## Claude Desktop config

`~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "qonaqzhai": {
      "command": "node",
      "args": ["/absolute/path/to/diploma/mcp/qonaqzhai-mcp/dist/index.js"],
      "env": {
        "QONAQZHAI_API": "http://localhost:8080"
      }
    }
  }
}
```

## Claude Code config

`.mcp.json` at the repo root or per-project:

```json
{
  "mcpServers": {
    "qonaqzhai": {
      "type": "stdio",
      "command": "node",
      "args": ["./mcp/qonaqzhai-mcp/dist/index.js"],
      "env": {
        "QONAQZHAI_API": "http://localhost:8080"
      }
    }
  }
}
```

## Smoke test

```bash
# from the repo root, with the backend stack running on :8080
cd mcp/qonaqzhai-mcp
npm run dev <<'EOF'
{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}
EOF
```

You should see the tool catalog printed on stdout.

## Local example session

```text
> auth_login { "email": "customer1@demo.kz", "password": "demo12345" }
{ token: "…", user: { … } }

> vendors_search { "category": "Venue", "city": "Almaty", "maxPrice": 600000 }
{ items: [ { id: "…", name: "Rixos Almaty Ballroom", priceFrom: 500000, … }, … ] }

> bookings_create { "vendorId": "…", "eventDate": "2026-07-12", "guestCount": 80, "amount": 500000 }
{ id: "…", status: "pending", … }

> ai_chat_send { "message": "Plan a wedding for 100 guests in Almaty, 5M ₸ budget" }
{ chatId: "…", text: "…", blocks: [ { type: "plan", … } ] }
```
