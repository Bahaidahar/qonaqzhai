// Workload: bcrypt verify cost. Stays just under the auth-svc per-IP limit
// (20 rps) so the bottleneck is bcrypt, not the limiter.

import { post, expect2xx } from "./_lib.js";

export const options = {
  scenarios: {
    login_burst: {
      executor: "constant-arrival-rate",
      rate: 15,
      timeUnit: "1s",
      duration: "30s",
      preAllocatedVUs: 20,
      maxVUs: 50,
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.02"],
    "http_req_duration{name:/api/login}": ["p(95)<150", "p(99)<300"],
  },
};

const ACCOUNTS = [
  { email: "customer1@demo.kz", password: "demo12345" },
  { email: "vendor1@demo.kz", password: "demo12345" },
];

export default function () {
  const a = ACCOUNTS[Math.floor(Math.random() * ACCOUNTS.length)];
  const r = post("/api/login", { email: a.email, password: a.password });
  expect2xx(r, "login");
}
