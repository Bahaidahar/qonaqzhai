export type ServiceUnit = "fixed" | "hour" | "item" | "person" | "day";

export const SERVICE_UNITS: { value: ServiceUnit; label: string }[] = [
  { value: "fixed", label: "Fixed" },
  { value: "hour", label: "Per hour" },
  { value: "item", label: "Per item" },
  { value: "person", label: "Per person" },
  { value: "day", label: "Per day" },
];

export interface Service {
  id: string;
  vendorId: string;
  name: string;
  description: string;
  price: number;
  unit: ServiceUnit;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface ServiceInput {
  name: string;
  description?: string;
  price: number;
  unit: ServiceUnit;
  isActive?: boolean;
}
