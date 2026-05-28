// Workload: find the saturation point. Ramp until the gateway's per-IP
// limit kicks in and the error rate explodes. The breakpoint is the
// inflection where 5xx + 429 start outpacing 2xx — that's the published
// "safe operating limit" we hand vendors.

import { BASE, expect2xx } from "./_lib.js";
import http from "k6/http";

export const options = {
  scenarios: {
    ramp: {
      executor: "ramping-arrival-rate",
      startRate: 10,
      timeUnit: "1s",
      preAllocatedVUs: 100,
      maxVUs: 400,
      stages: [
        { duration: "20s", target: 50 },
        { duration: "20s", target: 100 },
        { duration: "20s", target: 200 },
        { duration: "20s", target: 400 },
        { duration: "20s", target: 0 },
      ],
    },
  },
};

export default function () {
  const r = http.get(`${BASE}/api/vendors?limit=20`, {
    tags: { name: "/api/vendors" },
  });
  expect2xx(r, "search");
}
