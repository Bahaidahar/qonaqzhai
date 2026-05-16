"use client";

import { useI18n } from "@/shared/i18n/context";
import type { DictKey } from "@/shared/i18n/dict";

const CATEGORY_KEYS: Record<string, DictKey> = {
  Venue: "category_venue",
  Catering: "category_catering",
  "Music & DJ": "category_music",
  "Photo & Video": "category_photo",
  "Decor & Florists": "category_decor",
  Cakes: "category_cakes",
  Other: "category_other",
};

const CITY_KEYS: Record<string, DictKey> = {
  Almaty: "city_almaty",
  Astana: "city_astana",
  Shymkent: "city_shymkent",
  Karaganda: "city_karaganda",
  Aktobe: "city_aktobe",
  Atyrau: "city_atyrau",
};

const ROLE_KEYS: Record<string, DictKey> = {
  customer: "role_customer",
  vendor: "role_vendor",
  admin: "role_admin",
};

const USER_STATUS_KEYS: Record<string, DictKey> = {
  active: "user_status_active",
  suspended: "user_status_suspended",
};

const BOOKING_STATUS_KEYS: Record<string, DictKey> = {
  pending: "status_pending",
  accepted: "status_accepted",
  declined: "status_declined",
  cancelled: "status_cancelled",
};

export function useLabels() {
  const { t } = useI18n();
  return {
    category: (value: string) =>
      CATEGORY_KEYS[value] ? t(CATEGORY_KEYS[value]) : value,
    city: (value: string) =>
      CITY_KEYS[value] ? t(CITY_KEYS[value]) : value,
    role: (value: string) =>
      ROLE_KEYS[value] ? t(ROLE_KEYS[value]) : value,
    userStatus: (value: string) =>
      USER_STATUS_KEYS[value] ? t(USER_STATUS_KEYS[value]) : value,
    bookingStatus: (value: string) =>
      BOOKING_STATUS_KEYS[value] ? t(BOOKING_STATUS_KEYS[value]) : value,
  };
}
