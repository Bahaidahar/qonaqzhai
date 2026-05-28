// Workload: AI chat endpoint. Today the handler is a stub, so this
// measures the gateway + auth verify path more than any real LLM cost.
// When Gemini is wired in, expect p95 to jump 800 ms+ and concurrency
// thresholds to come from the upstream's quota, not us.

import { loginCustomer, post, expect2xx } from "./_lib.js";

export const options = {
  scenarios: {
    chat: {
      executor: "constant-arrival-rate",
      rate: 20,
      timeUnit: "1s",
      duration: "30s",
      preAllocatedVUs: 25,
      maxVUs: 60,
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.02"],
    "http_req_duration{name:/api/chat}": ["p(95)<400"],
  },
};

const PROMPTS = [
  "plan a toi for 120 in Almaty, 5M tenge",
  "birthday for 30 kids with animators",
  "corporate offsite for 80",
  "wedding photographer in Astana",
];

export function setup() {
  return { token: loginCustomer() };
}

export default function (data) {
  const r = post(
    "/api/chat",
    { message: PROMPTS[Math.floor(Math.random() * PROMPTS.length)] },
    data.token,
  );
  expect2xx(r, "chat");
}
