// Workload: public catalog browse. Constant arrival rate well below the
// per-IP gateway limit (100 rps) so we measure backend latency, not the
// rate-limiter. To find the saturation point, see saturation.js.

import { BASE, expect2xx } from "./_lib.js";
import http from "k6/http";

export const options = {
  scenarios: {
    catalog: {
      executor: "constant-arrival-rate",
      rate: 60,
      timeUnit: "1s",
      duration: "30s",
      preAllocatedVUs: 40,
      maxVUs: 80,
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.01"],
    "http_req_duration{name:/api/vendors}": ["p(95)<200", "p(99)<400"],
  },
};

const CATEGORIES = ["Venue", "Catering", "Music & DJ", "Photo & Video", "Decor & Florists"];
const SORTS = ["newest", "price_asc", "price_desc", "rating_desc"];

export default function () {
  const cat = CATEGORIES[Math.floor(Math.random() * CATEGORIES.length)];
  const sort = SORTS[Math.floor(Math.random() * SORTS.length)];
  const r = http.get(
    `${BASE}/api/vendors?category=${encodeURIComponent(cat)}&sort=${sort}&limit=20`,
    { tags: { name: "/api/vendors" } },
  );
  expect2xx(r, "search");
}
