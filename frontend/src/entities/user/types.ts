export type Role = "customer" | "vendor" | "admin";
export type UserStatus = "active" | "suspended";

export interface User {
  id: string;
  email: string;
  name: string;
  role: Role;
  status: UserStatus;
  createdAt: string;
}
