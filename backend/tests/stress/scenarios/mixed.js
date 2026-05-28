// Workload: realistic mixed traffic. ~80 rps total — under the gateway's
// 100 rps per-IP cap. 50 % catalog (anon), 20 % bookings list, 20 % /me,
// 10 % booking create.

import { loginCustomer, get, post, expect2xx } from "./_lib.js";

export const options = {
  scenarios: {
    realistic: {
      executor: "constant-arrival-rate",
      rate: 80,
      timeUnit: "1s",
      duration: "60s",
      preAllocatedVUs: 60,
      maxVUs: 120,
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.02"],
    http_req_duration: ["p(95)<300", "p(99)<800"],
  },
};

const CATEGORIES = ["Venue", "Catering", "Music & DJ", "Photo & Video"];

export function setup() {
  const token = loginCustomer();
  const list = get("/api/vendors?limit=1", token).json("items");
  return { token, vendorId: list && list.length > 0 ? list[0].id : null };
}

function isoFuture(d) {
  return new Date(Date.now() + d * 86400_000).toISOString().slice(0, 10);
}

export default function (data) {
  const roll = Math.random();
  if (roll < 0.5) {
    const cat = CATEGORIES[Math.floor(Math.random() * CATEGORIES.length)];
    expect2xx(
      get(`/api/vendors?category=${encodeURIComponent(cat)}&limit=20`),
      "catalog",
    );
  } else if (roll < 0.7) {
    expect2xx(get("/api/bookings", data.token), "bookings list");
  } else if (roll < 0.9) {
    expect2xx(get("/api/me", data.token), "me");
  } else if (data.vendorId) {
    expect2xx(
      post(
        "/api/bookings",
        {
          vendorId: data.vendorId,
          eventDate: isoFuture(30),
          guestCount: 80,
          amount: 350000,
        },
        data.token,
      ),
      "create booking",
    );
  }
}
