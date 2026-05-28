// Shared helpers reused across every scenario. Keeps individual scripts
// focused on workload shape, not boilerplate.

import http from "k6/http";
import { check, fail } from "k6";

export const BASE = __ENV.BASE_URL || "http://localhost:8080";

export function post(path, body, token) {
  const headers = { "Content-Type": "application/json" };
  if (token) headers["Authorization"] = `Bearer ${token}`;
  return http.post(`${BASE}${path}`, JSON.stringify(body), {
    headers,
    tags: { name: path },
  });
}

export function get(path, token) {
  const headers = {};
  if (token) headers["Authorization"] = `Bearer ${token}`;
  return http.get(`${BASE}${path}`, { headers, tags: { name: path } });
}

export function patch(path, body, token) {
  const headers = { "Content-Type": "application/json" };
  if (token) headers["Authorization"] = `Bearer ${token}`;
  return http.patch(`${BASE}${path}`, JSON.stringify(body), {
    headers,
    tags: { name: path },
  });
}

export function loginCustomer() {
  const r = post("/api/login", {
    email: "customer1@demo.kz",
    password: "demo12345",
  });
  if (r.status !== 200) fail(`login failed: ${r.status}`);
  return r.json("token");
}

export function loginVendor() {
  const r = post("/api/login", {
    email: "vendor1@demo.kz",
    password: "demo12345",
  });
  if (r.status !== 200) fail(`login failed: ${r.status}`);
  return r.json("token");
}

export function expect2xx(r, name) {
  return check(r, { [`${name} 2xx`]: (x) => x.status >= 200 && x.status < 300 });
}
