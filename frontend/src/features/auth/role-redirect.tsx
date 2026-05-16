"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/features/auth/context";
import type { Role } from "@/shared/api";

const ROLE_HOME: Record<Role, string> = {
  customer: "/",
  vendor: "/vendor",
  admin: "/admin",
};

/**
 * Redirects user to their role's home if not already there.
 * Use for the root chat page so vendors/admins land on their dashboards.
 */
export function RedirectIfWrongRole({
  expected,
}: {
  expected: Role | Role[];
}) {
  const { user } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!user) return;
    const ok = Array.isArray(expected)
      ? expected.includes(user.role)
      : user.role === expected;
    if (!ok) router.replace(ROLE_HOME[user.role]);
  }, [user, expected, router]);

  return null;
}
