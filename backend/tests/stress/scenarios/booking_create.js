// Workload: write-path saga. Constant arrival rate under the core per-IP
// limit (30 rps) so saturation comes from Postgres inserts + gRPC fan-out,
// not the limiter.

import { loginCustomer, get, post, expect2xx } from "./_lib.js";

export const options = {
  scenarios: {
    booking: {
      executor: "constant-arrival-rate",
      rate: 20,
      timeUnit: "1s",
      duration: "30s",
      preAllocatedVUs: 30,
      maxVUs: 60,
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.02"],
    "http_req_duration{name:/api/bookings}": ["p(95)<300", "p(99)<600"],
  },
};

export function setup() {
  const token = loginCustomer();
  const list = get("/api/vendors?limit=1", token).json("items");
  if (!list || list.length === 0) throw new Error("no vendors seeded");
  return { token, vendorId: list[0].id };
}

function isoFuture(daysOffset) {
  const d = new Date(Date.now() + daysOffset * 86400_000);
  return d.toISOString().slice(0, 10);
}

export default function (data) {
  const r = post(
    "/api/bookings",
    {
      vendorId: data.vendorId,
      eventDate: isoFuture(30 + Math.floor(Math.random() * 90)),
      guestCount: 40 + Math.floor(Math.random() * 200),
      amount: 200000 + Math.floor(Math.random() * 800000),
    },
    data.token,
  );
  expect2xx(r, "create booking");
}
