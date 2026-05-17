export type CardBrand = "visa" | "mastercard" | "amex" | "discover" | "unknown";

export interface PaymentCard {
  id: string;
  userId: string;
  brand: CardBrand;
  last4: string;
  expMonth: number;
  expYear: number;
  holder: string;
  isDefault: boolean;
  createdAt: string;
}
