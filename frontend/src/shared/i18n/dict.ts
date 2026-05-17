export type Locale = "kz" | "ru" | "en";

export const LOCALES: { code: Locale; label: string; short: string }[] = [
  { code: "kz", label: "Қазақша", short: "KZ" },
  { code: "ru", label: "Русский", short: "RU" },
  { code: "en", label: "English", short: "EN" },
];

export const DEFAULT_LOCALE: Locale = "ru";

type Dict = {
  // chrome
  status_online: string;
  hero_title_a: string;
  hero_title_b: string;
  input_placeholder_main: string;
  input_placeholder_reply: string;
  input_hint: string;
  capability_plan: string;
  capability_vendors: string;
  capability_budget: string;
  capability_guests: string;
  suggestion_toi: string;
  suggestion_corporate: string;
  suggestion_photographer: string;
  suggestion_invitation: string;
  sidebar_new_chat: string;
  sidebar_recent: string;
  sidebar_plan: string;
  thinking: string;
  // settings
  settings_title: string;
  settings_back: string;
  settings_section_account: string;
  settings_section_pref: string;
  settings_section_about: string;
  settings_language: string;
  settings_language_hint: string;
  settings_appearance: string;
  settings_appearance_hint: string;
  settings_email: string;
  settings_phone: string;
  settings_plan_label: string;
  settings_plan_value: string;
  settings_version: string;
  settings_signout: string;
  settings_save: string;
  settings_name: string;
  settings_model: string;
  settings_built_in: string;
  theme_light: string;
  theme_dark: string;
  theme_system: string;
  // chat block headers
  block_event_plan: string;
  block_draft: string;
  block_budget: string;
  block_vendors: string;
  // plan stat
  stat_date: string;
  stat_city: string;
  stat_guests: string;
  stat_budget: string;
  // vendor block
  vendor_view_all: string;
  vendor_matches: string;
  // auth
  auth_welcome: string;
  auth_signin_tab: string;
  auth_signup_tab: string;
  auth_email: string;
  auth_password: string;
  auth_name: string;
  auth_signin_btn: string;
  auth_signup_btn: string;
  auth_no_account: string;
  auth_has_account: string;
  auth_loading: string;
  auth_role_question: string;
  auth_role_customer: string;
  auth_role_customer_hint: string;
  auth_role_vendor: string;
  auth_role_vendor_hint: string;
  auth_demo_title: string;
  // nav
  nav_chat: string;
  nav_vendors: string;
  nav_bookings: string;
  nav_vendor_profile: string;
  nav_admin_users: string;
  nav_notifications: string;
  nav_messages: string;
  nav_admin_vendors: string;
  // common
  common_loading: string;
  common_save: string;
  common_cancel: string;
  common_search: string;
  common_back: string;
  common_delete: string;
  threads_title: string;
  threads_hint: string;
  threads_empty: string;
  threads_unknown_vendor: string;
  threads_unknown_user: string;
  notif_signup_welcome_title: string;
  notif_signup_welcome_body: string;
  notif_password_reset_title: string;
  notif_password_reset_body: string;
  notif_booking_created_title: string;
  notif_booking_created_body: string;
  notif_booking_accepted_title: string;
  notif_booking_accepted_body: string;
  notif_booking_declined_title: string;
  notif_booking_declined_body: string;
  notif_booking_paid_title: string;
  notif_booking_paid_body: string;
  notif_vendor_approved_title: string;
  notif_vendor_approved_body: string;
  notif_vendor_rejected_title: string;
  notif_vendor_rejected_body: string;
  notif_thread_message_title: string;
  notif_channel_email: string;
  notif_channel_push: string;
  notif_status_pending: string;
  notif_status_sent: string;
  notif_status_failed: string;
  booking_status_pending: string;
  booking_status_accepted: string;
  booking_status_declined: string;
  booking_status_cancelled: string;
  booking_status_completed: string;
  booking_status_paid: string;
  // cards (saved payment instruments — mock)
  cards_title: string;
  cards_hint: string;
  cards_empty: string;
  cards_add: string;
  cards_save: string;
  cards_default: string;
  cards_make_default: string;
  cards_field_number: string;
  cards_field_exp: string;
  cards_field_holder: string;
  cards_err_exp: string;
  cards_disclaimer: string;
  nav_cards: string;
  bookings_btn_pay: string;
  bookings_btn_add_card: string;
  // admin users
  users_title: string;
  users_hint: string;
  users_filter_all: string;
  users_col_name: string;
  users_col_email: string;
  users_col_role: string;
  users_col_status: string;
  users_col_actions: string;
  users_btn_suspend: string;
  users_btn_activate: string;
  // admin vendors
  admin_title: string;
  admin_hint: string;
  admin_stat_users: string;
  admin_stat_pending: string;
  admin_stat_approved: string;
  admin_stat_bookings: string;
  admin_group_pending: string;
  admin_group_approved: string;
  admin_group_rejected: string;
  admin_empty: string;
  admin_btn_approve: string;
  admin_btn_reject: string;
  admin_btn_suspend: string;
  admin_btn_reapprove: string;
  funnel_submitted: string;
  funnel_pending: string;
  funnel_approved: string;
  funnel_rejected: string;
  admin_btn_preview: string;
  admin_btn_hide: string;
  admin_preview_no_photos: string;
  admin_preview_no_description: string;
  admin_preview_no_services: string;
  admin_preview_services: string;
  admin_preview_description: string;
  // vendors catalog (customer)
  vendors_title: string;
  vendors_hint: string;
  vendors_search_ph: string;
  vendors_filter_all: string;
  vendors_empty: string;
  vendors_from: string;
  vendors_price_max_ph: string;
  vendors_sort_aria: string;
  vendors_sort_newest: string;
  vendors_sort_price_asc: string;
  vendors_sort_price_desc: string;
  vendors_sort_rating_desc: string;
  vendors_total_suffix: string;
  // vendor detail
  vendor_detail_back: string;
  vendor_detail_starting: string;
  vendor_detail_date: string;
  vendor_detail_guests: string;
  vendor_detail_note: string;
  vendor_detail_note_ph: string;
  vendor_detail_btn: string;
  vendor_detail_sent_title: string;
  vendor_detail_sent_hint: string;
  vendor_detail_view_bookings: string;
  vendor_detail_not_found: string;
  // bookings list (customer)
  bookings_title: string;
  bookings_hint: string;
  bookings_empty: string;
  bookings_btn_cancel: string;
  bookings_status_pending: string;
  bookings_status_accepted: string;
  bookings_status_declined: string;
  bookings_status_cancelled: string;
  // vendor profile
  vendor_profile_title: string;
  vendor_profile_hint: string;
  vendor_profile_section_basics: string;
  vendor_profile_section_photos: string;
  vendor_profile_field_name: string;
  vendor_profile_field_name_ph: string;
  vendor_profile_field_category: string;
  vendor_profile_field_city: string;
  vendor_profile_field_price: string;
  vendor_profile_field_desc: string;
  vendor_profile_field_desc_ph: string;
  vendor_profile_listed_from: string;
  vendor_profile_photos_hint: string;
  vendor_profile_no_photos: string;
  vendor_profile_btn_upload: string;
  vendor_profile_btn_delete: string;
  vendor_profile_status_pending: string;
  vendor_profile_status_approved: string;
  vendor_profile_status_rejected: string;
  // vendor bookings inbox
  vendor_bookings_title: string;
  vendor_bookings_hint: string;
  vendor_bookings_empty: string;
  vendor_bookings_btn_accept: string;
  vendor_bookings_btn_decline: string;
  vendor_bookings_guests: string;
  // category labels (value stays English, label translates)
  category_venue: string;
  category_catering: string;
  category_music: string;
  category_photo: string;
  category_decor: string;
  category_cakes: string;
  category_other: string;
  // city labels
  city_almaty: string;
  city_astana: string;
  city_shymkent: string;
  city_karaganda: string;
  city_aktobe: string;
  city_atyrau: string;
  // booking status (vendor inbox & bookings list)
  status_pending: string;
  status_accepted: string;
  status_declined: string;
  status_cancelled: string;
  // user role labels
  role_customer: string;
  role_vendor: string;
  role_admin: string;
  user_status_active: string;
  user_status_suspended: string;
  // services
  services_title: string;
  services_hint: string;
  services_btn_new: string;
  services_empty: string;
  services_inactive: string;
  services_disable: string;
  services_enable: string;
  services_confirm_delete: string;
  services_dialog_new: string;
  services_dialog_edit: string;
  services_field_name: string;
  services_field_name_ph: string;
  services_field_description: string;
  services_field_price: string;
  services_field_unit: string;
  services_unit_fixed: string;
  services_unit_hour: string;
  services_unit_item: string;
  services_unit_person: string;
  services_unit_day: string;
  services_save_failed: string;
  services_load_failed: string;
  services_loading: string;
  services_empty_public: string;
  services_pick: string;
  services_per_hour: string;
  services_per_item: string;
  services_per_person: string;
  services_per_day: string;
  // reviews
  reviews_title: string;
  reviews_loading: string;
  reviews_empty: string;
  reviews_submit: string;
  reviews_share_ph: string;
  reviews_thanks: string;
  // notifications page
  notifications_title: string;
  notifications_hint: string;
  notifications_empty: string;
  // admin analytics
  charts_bookings_per_day: string;
  charts_top_categories: string;
  charts_funnel: string;
  charts_no_data: string;
  // auth forgot/reset
  auth_forgot_title: string;
  auth_forgot_hint: string;
  auth_forgot_btn_send: string;
  auth_forgot_done_title: string;
  auth_forgot_done_hint: string;
  auth_forgot_enter_token: string;
  auth_forgot_back: string;
  auth_forgot_link: string;
  auth_reset_title: string;
  auth_reset_token: string;
  auth_reset_token_ph: string;
  auth_reset_new_password: string;
  auth_reset_btn: string;
  auth_reset_done_title: string;
  auth_reset_done_hint: string;
  auth_network_error: string;
  // chat header
  chat_header_label: string;
  // booking pay button
  booking_pay: string;
  // vendor detail action labels
  vendor_detail_book_now: string;
};

export type DictKey = keyof Dict;

export const DICT: Record<Locale, Dict> = {
  en: {
    status_online: "Online · Gemini 2.5 Flash",
    hero_title_a: "What's on your mind,",
    hero_title_b: "Aigerim?",
    input_placeholder_main: "Describe your event or ask anything...",
    input_placeholder_reply: "Reply...",
    input_hint:
      "qonaqzhai can make mistakes. Verify important info with vendors.",
    capability_plan: "Plan events",
    capability_vendors: "Match vendors",
    capability_budget: "Budget split",
    capability_guests: "Guest mgmt",
    suggestion_toi: "Plan a toi for 150 guests in Almaty",
    suggestion_corporate: "Suggest budget for 80-person corporate event",
    suggestion_photographer: "Find me a photographer in Astana",
    suggestion_invitation: "Generate invitation for August wedding",
    sidebar_new_chat: "New chat",
    sidebar_recent: "Recent",
    sidebar_plan: "Free plan",
    thinking: "Thinking",
    settings_title: "Settings",
    settings_back: "Back to chat",
    settings_section_account: "Account",
    settings_section_pref: "Preferences",
    settings_section_about: "About",
    settings_language: "Language",
    settings_language_hint: "Interface language for menus and prompts.",
    settings_appearance: "Appearance",
    settings_appearance_hint: "Choose how qonaqzhai looks to you.",
    settings_email: "Email",
    settings_phone: "Phone",
    settings_plan_label: "Plan",
    settings_plan_value: "Free · MVP",
    settings_version: "Version",
    settings_signout: "Sign out",
    settings_save: "Save changes",
    settings_name: "Name",
    settings_model: "Model",
    settings_built_in: "Built in",
    theme_light: "Light",
    theme_dark: "Dark",
    theme_system: "System",
    block_event_plan: "event_plan",
    block_draft: "Draft",
    block_budget: "budget_breakdown",
    block_vendors: "vendors",
    stat_date: "Date",
    stat_city: "City",
    stat_guests: "Guests",
    stat_budget: "Budget",
    vendor_view_all: "View all matches",
    vendor_matches: "matches",
    auth_welcome: "Welcome to qonaqzhai",
    auth_signin_tab: "Sign in",
    auth_signup_tab: "Sign up",
    auth_email: "Email",
    auth_password: "Password",
    auth_name: "Name",
    auth_signin_btn: "Sign in",
    auth_signup_btn: "Create account",
    auth_no_account: "No account yet?",
    auth_has_account: "Already have an account?",
    auth_loading: "Loading...",
    auth_role_question: "I want to...",
    auth_role_customer: "Plan an event",
    auth_role_customer_hint: "I'm organizing my own toi/wedding/event",
    auth_role_vendor: "Offer services",
    auth_role_vendor_hint: "I'm a vendor (venue, catering, photo, etc.)",
    auth_demo_title: "Quick demo login",
    nav_chat: "Chat",
    nav_vendors: "Vendors",
    nav_bookings: "Bookings",
    nav_vendor_profile: "My profile",
    nav_admin_users: "Users",
    nav_notifications: "Notifications",
    nav_messages: "Messages",
    nav_admin_vendors: "Vendor moderation",
    common_loading: "Loading...",
    threads_title: "Messages",
    threads_hint: "Direct chats with vendors / customers, scoped to each accepted booking.",
    threads_empty: "No active conversations. Threads open automatically when a vendor accepts a booking.",
    threads_unknown_vendor: "Vendor",
    threads_unknown_user: "User",
    notif_signup_welcome_title: "Welcome",
    notif_signup_welcome_body: "Your account is ready.",
    notif_password_reset_title: "Password reset",
    notif_password_reset_body: "We sent a reset link to your inbox.",
    notif_booking_created_title: "New booking request",
    notif_booking_created_body: "You have a new booking request — please review it in your inbox.",
    notif_booking_accepted_title: "Booking accepted",
    notif_booking_accepted_body: "The vendor accepted your booking.",
    notif_booking_declined_title: "Booking declined",
    notif_booking_declined_body: "The vendor declined your booking.",
    notif_booking_paid_title: "Booking paid",
    notif_booking_paid_body: "Your booking has been paid.",
    notif_vendor_approved_title: "Vendor approved",
    notif_vendor_approved_body: "Your vendor profile is now visible in the catalog.",
    notif_vendor_rejected_title: "Vendor rejected",
    notif_vendor_rejected_body: "Your vendor profile was rejected by moderation.",
    notif_thread_message_title: "New message",
    notif_channel_email: "email",
    notif_channel_push: "push",
    notif_status_pending: "pending",
    notif_status_sent: "sent",
    notif_status_failed: "failed",
    booking_status_pending: "Pending",
    booking_status_accepted: "Accepted",
    booking_status_declined: "Declined",
    booking_status_cancelled: "Cancelled",
    booking_status_completed: "Completed",
    booking_status_paid: "Paid",
    common_save: "Save",
    common_cancel: "Cancel",
    common_search: "Search...",
    common_back: "Back",
    common_delete: "Delete",
    cards_title: "Saved cards",
    cards_hint: "Add a card to pay vendors after they accept your booking. Mock — no real payment is taken.",
    cards_empty: "No cards yet.",
    cards_add: "Add card",
    cards_save: "Save card",
    cards_default: "Default",
    cards_make_default: "Set default",
    cards_field_number: "Card number",
    cards_field_exp: "Expiry (MM/YY)",
    cards_field_holder: "Cardholder",
    cards_err_exp: "Invalid expiry date",
    cards_disclaimer: "Demo only — we never call a real payment processor and the PAN is discarded immediately.",
    nav_cards: "Cards",
    bookings_btn_pay: "Pay",
    bookings_btn_add_card: "Add card to pay",
    users_title: "Users",
    users_hint: "Manage all platform users.",
    users_filter_all: "all",
    users_col_name: "Name",
    users_col_email: "Email",
    users_col_role: "Role",
    users_col_status: "Status",
    users_col_actions: "Actions",
    users_btn_suspend: "Suspend",
    users_btn_activate: "Activate",
    admin_title: "Admin · Vendors",
    admin_hint: "Approve or reject vendor profiles before they appear in the catalog.",
    admin_stat_users: "Users",
    admin_stat_pending: "Vendors pending",
    admin_stat_approved: "Vendors approved",
    admin_stat_bookings: "Bookings",
    admin_group_pending: "Pending",
    admin_group_approved: "Approved",
    admin_group_rejected: "Rejected",
    admin_empty: "Nothing here",
    admin_btn_approve: "Approve",
    admin_btn_reject: "Reject",
    admin_btn_suspend: "Suspend",
    admin_btn_reapprove: "Re-approve",
    funnel_submitted: "Submitted",
    funnel_pending: "Pending",
    funnel_approved: "Approved",
    funnel_rejected: "Rejected",
    admin_btn_preview: "Preview",
    admin_btn_hide: "Hide",
    admin_preview_no_photos: "No photos uploaded.",
    admin_preview_no_description: "No description provided.",
    admin_preview_no_services: "No services listed.",
    admin_preview_services: "Services",
    admin_preview_description: "Description",
    vendors_title: "Vendors",
    vendors_hint: "Verified partners for your next event.",
    vendors_search_ph: "Search...",
    vendors_filter_all: "All",
    vendors_price_max_ph: "Max price (₸)",
    vendors_sort_aria: "Sort",
    vendors_sort_newest: "Newest",
    vendors_sort_price_asc: "Price ↑",
    vendors_sort_price_desc: "Price ↓",
    vendors_sort_rating_desc: "Top rated",
    vendors_total_suffix: "total",
    vendors_empty: "No vendors found",
    vendors_from: "from",
    vendor_detail_back: "Back to catalog",
    vendor_detail_starting: "Starting from",
    vendor_detail_date: "Event date",
    vendor_detail_guests: "Guests",
    vendor_detail_note: "Note",
    vendor_detail_note_ph: "Anything we should know...",
    vendor_detail_btn: "Request booking",
    vendor_detail_sent_title: "Request sent",
    vendor_detail_sent_hint: "The vendor will respond shortly. Check your bookings.",
    vendor_detail_view_bookings: "View bookings",
    vendor_detail_not_found: "Vendor not found",
    bookings_title: "My bookings",
    bookings_hint: "Requests you sent to vendors.",
    bookings_empty: "No bookings yet",
    bookings_btn_cancel: "Cancel",
    bookings_status_pending: "pending",
    bookings_status_accepted: "accepted",
    bookings_status_declined: "declined",
    bookings_status_cancelled: "cancelled",
    vendor_profile_title: "My profile",
    vendor_profile_hint: "Your public listing in the qonaqzhai catalog.",
    vendor_profile_section_basics: "Basics",
    vendor_profile_section_photos: "Photos",
    vendor_profile_field_name: "Business name",
    vendor_profile_field_name_ph: "Rixos Almaty Ballroom",
    vendor_profile_field_category: "Category",
    vendor_profile_field_city: "City",
    vendor_profile_field_price: "Price from (₸)",
    vendor_profile_field_desc: "Description",
    vendor_profile_field_desc_ph: "Premier venue in the heart of Almaty...",
    vendor_profile_listed_from: "Listed from",
    vendor_profile_photos_hint: "JPG/PNG up to 5MB. First photo is your cover.",
    vendor_profile_no_photos: "No photos yet",
    vendor_profile_btn_upload: "Upload",
    vendor_profile_btn_delete: "Delete",
    vendor_profile_status_pending: "Pending admin approval — your profile is hidden from customers until reviewed.",
    vendor_profile_status_approved: "Approved — your profile is live in the catalog.",
    vendor_profile_status_rejected: "Rejected — contact support to discuss.",
    vendor_bookings_title: "Bookings",
    vendor_bookings_hint: "Incoming requests from customers.",
    vendor_bookings_empty: "No bookings yet",
    vendor_bookings_btn_accept: "Accept",
    vendor_bookings_btn_decline: "Decline",
    vendor_bookings_guests: "guests",
    category_venue: "Venue",
    category_catering: "Catering",
    category_music: "Music & DJ",
    category_photo: "Photo & Video",
    category_decor: "Decor & Florists",
    category_cakes: "Cakes",
    category_other: "Other",
    city_almaty: "Almaty",
    city_astana: "Astana",
    city_shymkent: "Shymkent",
    city_karaganda: "Karaganda",
    city_aktobe: "Aktobe",
    city_atyrau: "Atyrau",
    status_pending: "pending",
    status_accepted: "accepted",
    status_declined: "declined",
    status_cancelled: "cancelled",
    role_customer: "customer",
    role_vendor: "vendor",
    role_admin: "admin",
    user_status_active: "active",
    user_status_suspended: "suspended",
    services_title: "Services",
    services_hint: "Offerings shown to customers — they pick one when booking.",
    services_btn_new: "New service",
    services_empty: "No services yet — add your first one.",
    services_inactive: "inactive",
    services_disable: "Disable",
    services_enable: "Enable",
    services_confirm_delete: "Delete this service?",
    services_dialog_new: "New service",
    services_dialog_edit: "Edit service",
    services_field_name: "Name",
    services_field_name_ph: "Wedding photography",
    services_field_description: "Description",
    services_field_price: "Price (₸)",
    services_field_unit: "Unit",
    services_unit_fixed: "Fixed",
    services_unit_hour: "Per hour",
    services_unit_item: "Per item",
    services_unit_person: "Per person",
    services_unit_day: "Per day",
    services_save_failed: "Save failed",
    services_load_failed: "Load failed",
    services_loading: "Loading services…",
    services_empty_public: "Vendor has not published any services.",
    services_pick: "Pick a service",
    services_per_hour: " / hr",
    services_per_item: " / item",
    services_per_person: " / person",
    services_per_day: " / day",
    reviews_title: "Reviews",
    reviews_loading: "Loading reviews…",
    reviews_empty: "No reviews yet.",
    reviews_submit: "Submit review",
    reviews_share_ph: "Share your experience…",
    reviews_thanks: "Thanks for the review!",
    notifications_title: "Notifications",
    notifications_hint: "Booking updates, vendor approvals, payment receipts.",
    notifications_empty: "No notifications yet.",
    charts_bookings_per_day: "Bookings per day",
    charts_top_categories: "Top categories",
    charts_funnel: "Vendor approval funnel",
    charts_no_data: "No data yet",
    auth_forgot_title: "Forgot password",
    auth_forgot_hint: "We'll email you a reset link.",
    auth_forgot_btn_send: "Send reset link",
    auth_forgot_done_title: "Check your inbox",
    auth_forgot_done_hint: "If an account exists, a password reset link has been sent. Follow the email, or paste the token below.",
    auth_forgot_enter_token: "Enter reset token →",
    auth_forgot_back: "Back to sign in",
    auth_forgot_link: "Forgot password?",
    auth_reset_title: "Reset password",
    auth_reset_token: "Reset token",
    auth_reset_token_ph: "Paste from email",
    auth_reset_new_password: "New password",
    auth_reset_btn: "Reset password",
    auth_reset_done_title: "Password updated",
    auth_reset_done_hint: "Redirecting to sign in…",
    auth_network_error: "Network error",
    chat_header_label: "Chat",
    booking_pay: "Pay",
    vendor_detail_book_now: "Book now",
  },
  ru: {
    status_online: "Онлайн · Gemini 2.5 Flash",
    hero_title_a: "Что планируешь,",
    hero_title_b: "Айгерим?",
    input_placeholder_main: "Опиши событие или задай вопрос...",
    input_placeholder_reply: "Ответить...",
    input_hint:
      "qonaqzhai может ошибаться. Проверяйте важную информацию с поставщиками.",
    capability_plan: "Планирование",
    capability_vendors: "Подбор подрядчиков",
    capability_budget: "Бюджет",
    capability_guests: "Гости",
    suggestion_toi: "Спланируй той на 150 гостей в Алматы",
    suggestion_corporate: "Предложи бюджет на корпоратив на 80 человек",
    suggestion_photographer: "Найди фотографа в Астане",
    suggestion_invitation: "Сгенерируй приглашение на августовскую свадьбу",
    sidebar_new_chat: "Новый чат",
    sidebar_recent: "Недавние",
    sidebar_plan: "Бесплатный тариф",
    thinking: "Думаю",
    settings_title: "Настройки",
    settings_back: "К чату",
    settings_section_account: "Аккаунт",
    settings_section_pref: "Предпочтения",
    settings_section_about: "О приложении",
    settings_language: "Язык",
    settings_language_hint: "Язык интерфейса меню и подсказок.",
    settings_appearance: "Внешний вид",
    settings_appearance_hint: "Выберите тему оформления.",
    settings_email: "Email",
    settings_phone: "Телефон",
    settings_plan_label: "Тариф",
    settings_plan_value: "Бесплатный · MVP",
    settings_version: "Версия",
    settings_signout: "Выйти",
    settings_save: "Сохранить",
    settings_name: "Имя",
    settings_model: "Модель",
    settings_built_in: "Сделано в",
    theme_light: "Светлая",
    theme_dark: "Тёмная",
    theme_system: "Системная",
    block_event_plan: "план_события",
    block_draft: "Черновик",
    block_budget: "бюджет",
    block_vendors: "подрядчики",
    stat_date: "Дата",
    stat_city: "Город",
    stat_guests: "Гостей",
    stat_budget: "Бюджет",
    vendor_view_all: "Все варианты",
    vendor_matches: "совпадений",
    auth_welcome: "Добро пожаловать в qonaqzhai",
    auth_signin_tab: "Вход",
    auth_signup_tab: "Регистрация",
    auth_email: "Email",
    auth_password: "Пароль",
    auth_name: "Имя",
    auth_signin_btn: "Войти",
    auth_signup_btn: "Создать аккаунт",
    auth_no_account: "Нет аккаунта?",
    auth_has_account: "Уже есть аккаунт?",
    auth_loading: "Загрузка...",
    auth_role_question: "Я хочу...",
    auth_role_customer: "Планировать событие",
    auth_role_customer_hint: "Я организую свой той/свадьбу/мероприятие",
    auth_role_vendor: "Оказывать услуги",
    auth_role_vendor_hint: "Я подрядчик (зал, кейтеринг, фото и т.д.)",
    auth_demo_title: "Быстрый демо-вход",
    nav_chat: "Чат",
    nav_vendors: "Подрядчики",
    nav_bookings: "Бронирования",
    nav_vendor_profile: "Мой профиль",
    nav_admin_users: "Пользователи",
    nav_notifications: "Уведомления",
    nav_messages: "Сообщения",
    nav_admin_vendors: "Модерация",
    common_loading: "Загрузка...",
    threads_title: "Сообщения",
    threads_hint: "Прямые чаты с подрядчиками / клиентами по принятым бронированиям.",
    threads_empty: "Нет активных переписок. Чаты открываются автоматически, когда подрядчик принимает бронь.",
    threads_unknown_vendor: "Подрядчик",
    threads_unknown_user: "Пользователь",
    notif_signup_welcome_title: "Добро пожаловать",
    notif_signup_welcome_body: "Аккаунт готов.",
    notif_password_reset_title: "Сброс пароля",
    notif_password_reset_body: "Мы отправили ссылку для сброса на вашу почту.",
    notif_booking_created_title: "Новый запрос на бронь",
    notif_booking_created_body: "Поступил новый запрос на бронирование — проверьте во вкладке броней.",
    notif_booking_accepted_title: "Бронь принята",
    notif_booking_accepted_body: "Подрядчик принял вашу бронь.",
    notif_booking_declined_title: "Бронь отклонена",
    notif_booking_declined_body: "Подрядчик отклонил вашу бронь.",
    notif_booking_paid_title: "Оплата прошла",
    notif_booking_paid_body: "Бронь оплачена.",
    notif_vendor_approved_title: "Профиль одобрен",
    notif_vendor_approved_body: "Ваш профиль виден в каталоге.",
    notif_vendor_rejected_title: "Профиль отклонён",
    notif_vendor_rejected_body: "Модерация отклонила ваш профиль.",
    notif_thread_message_title: "Новое сообщение",
    notif_channel_email: "email",
    notif_channel_push: "push",
    notif_status_pending: "ожидает",
    notif_status_sent: "отправлено",
    notif_status_failed: "ошибка",
    booking_status_pending: "Ожидает",
    booking_status_accepted: "Принято",
    booking_status_declined: "Отклонено",
    booking_status_cancelled: "Отменено",
    booking_status_completed: "Завершено",
    booking_status_paid: "Оплачено",
    common_save: "Сохранить",
    common_cancel: "Отмена",
    common_search: "Поиск...",
    common_back: "Назад",
    common_delete: "Удалить",
    cards_title: "Карты",
    cards_hint: "Добавьте карту, чтобы оплатить бронь после её подтверждения подрядчиком. Мок — реальная оплата не производится.",
    cards_empty: "Нет сохранённых карт.",
    cards_add: "Добавить карту",
    cards_save: "Сохранить карту",
    cards_default: "По умолчанию",
    cards_make_default: "Сделать основной",
    cards_field_number: "Номер карты",
    cards_field_exp: "Срок (ММ/ГГ)",
    cards_field_holder: "Держатель",
    cards_err_exp: "Неверный срок действия",
    cards_disclaimer: "Демо — мы не обращаемся к реальному платёжному провайдеру и не сохраняем номер карты.",
    nav_cards: "Карты",
    bookings_btn_pay: "Оплатить",
    bookings_btn_add_card: "Добавить карту",
    users_title: "Пользователи",
    users_hint: "Управление всеми пользователями платформы.",
    users_filter_all: "все",
    users_col_name: "Имя",
    users_col_email: "Email",
    users_col_role: "Роль",
    users_col_status: "Статус",
    users_col_actions: "Действия",
    users_btn_suspend: "Заблокировать",
    users_btn_activate: "Активировать",
    admin_title: "Админ · Подрядчики",
    admin_hint: "Одобряй или отклоняй профили подрядчиков перед публикацией в каталоге.",
    admin_stat_users: "Пользователи",
    admin_stat_pending: "На модерации",
    admin_stat_approved: "Одобрено",
    admin_stat_bookings: "Бронирования",
    admin_group_pending: "На модерации",
    admin_group_approved: "Одобрено",
    admin_group_rejected: "Отклонено",
    admin_empty: "Пусто",
    admin_btn_approve: "Одобрить",
    admin_btn_reject: "Отклонить",
    admin_btn_suspend: "Снять",
    admin_btn_reapprove: "Восстановить",
    funnel_submitted: "Поданы",
    funnel_pending: "На модерации",
    funnel_approved: "Одобрены",
    funnel_rejected: "Отклонены",
    admin_btn_preview: "Предпросмотр",
    admin_btn_hide: "Свернуть",
    admin_preview_no_photos: "Нет загруженных фото.",
    admin_preview_no_description: "Описание не заполнено.",
    admin_preview_no_services: "Услуги не добавлены.",
    admin_preview_services: "Услуги",
    admin_preview_description: "Описание",
    vendors_title: "Подрядчики",
    vendors_hint: "Проверенные партнёры для твоего события.",
    vendors_search_ph: "Поиск...",
    vendors_filter_all: "Все",
    vendors_price_max_ph: "Макс. цена (₸)",
    vendors_sort_aria: "Сортировка",
    vendors_sort_newest: "Сначала новые",
    vendors_sort_price_asc: "Цена ↑",
    vendors_sort_price_desc: "Цена ↓",
    vendors_sort_rating_desc: "Топ рейтинг",
    vendors_total_suffix: "всего",
    vendors_empty: "Подрядчиков не найдено",
    vendors_from: "от",
    vendor_detail_back: "К каталогу",
    vendor_detail_starting: "Стартовая цена",
    vendor_detail_date: "Дата события",
    vendor_detail_guests: "Гостей",
    vendor_detail_note: "Комментарий",
    vendor_detail_note_ph: "Что-нибудь добавить...",
    vendor_detail_btn: "Запросить бронь",
    vendor_detail_sent_title: "Запрос отправлен",
    vendor_detail_sent_hint: "Подрядчик скоро ответит. Проверь свои бронирования.",
    vendor_detail_view_bookings: "Мои бронирования",
    vendor_detail_not_found: "Подрядчик не найден",
    bookings_title: "Мои бронирования",
    bookings_hint: "Запросы которые ты отправил подрядчикам.",
    bookings_empty: "Пока нет бронирований",
    bookings_btn_cancel: "Отменить",
    bookings_status_pending: "ожидает",
    bookings_status_accepted: "принят",
    bookings_status_declined: "отклонён",
    bookings_status_cancelled: "отменён",
    vendor_profile_title: "Мой профиль",
    vendor_profile_hint: "Твоя публичная карточка в каталоге qonaqzhai.",
    vendor_profile_section_basics: "Основное",
    vendor_profile_section_photos: "Фото",
    vendor_profile_field_name: "Название",
    vendor_profile_field_name_ph: "Rixos Almaty Ballroom",
    vendor_profile_field_category: "Категория",
    vendor_profile_field_city: "Город",
    vendor_profile_field_price: "Цена от (₸)",
    vendor_profile_field_desc: "Описание",
    vendor_profile_field_desc_ph: "Премиум зал в центре Алматы...",
    vendor_profile_listed_from: "В каталоге от",
    vendor_profile_photos_hint: "JPG/PNG до 5MB. Первое фото — обложка.",
    vendor_profile_no_photos: "Фото пока нет",
    vendor_profile_btn_upload: "Загрузить",
    vendor_profile_btn_delete: "Удалить",
    vendor_profile_status_pending: "На модерации — твой профиль не виден заказчикам пока админ не одобрит.",
    vendor_profile_status_approved: "Одобрено — твой профиль виден в каталоге.",
    vendor_profile_status_rejected: "Отклонено — свяжись с поддержкой.",
    vendor_bookings_title: "Бронирования",
    vendor_bookings_hint: "Входящие запросы от заказчиков.",
    vendor_bookings_empty: "Пока нет бронирований",
    vendor_bookings_btn_accept: "Принять",
    vendor_bookings_btn_decline: "Отклонить",
    vendor_bookings_guests: "гостей",
    category_venue: "Зал",
    category_catering: "Кейтеринг",
    category_music: "Музыка и DJ",
    category_photo: "Фото и видео",
    category_decor: "Декор и флористы",
    category_cakes: "Торты",
    category_other: "Другое",
    city_almaty: "Алматы",
    city_astana: "Астана",
    city_shymkent: "Шымкент",
    city_karaganda: "Караганда",
    city_aktobe: "Актобе",
    city_atyrau: "Атырау",
    status_pending: "ожидает",
    status_accepted: "принят",
    status_declined: "отклонён",
    status_cancelled: "отменён",
    role_customer: "заказчик",
    role_vendor: "подрядчик",
    role_admin: "админ",
    user_status_active: "активен",
    user_status_suspended: "заблокирован",
    services_title: "Услуги",
    services_hint: "Услуги, которые видит заказчик и выбирает при бронировании.",
    services_btn_new: "Новая услуга",
    services_empty: "Услуг пока нет — добавь первую.",
    services_inactive: "скрыта",
    services_disable: "Скрыть",
    services_enable: "Показать",
    services_confirm_delete: "Удалить эту услугу?",
    services_dialog_new: "Новая услуга",
    services_dialog_edit: "Редактировать услугу",
    services_field_name: "Название",
    services_field_name_ph: "Свадебная фотосъёмка",
    services_field_description: "Описание",
    services_field_price: "Цена (₸)",
    services_field_unit: "Единица",
    services_unit_fixed: "Фикс",
    services_unit_hour: "За час",
    services_unit_item: "За штуку",
    services_unit_person: "За гостя",
    services_unit_day: "За день",
    services_save_failed: "Не удалось сохранить",
    services_load_failed: "Не удалось загрузить",
    services_loading: "Загрузка услуг…",
    services_empty_public: "У подрядчика пока нет опубликованных услуг.",
    services_pick: "Выбери услугу",
    services_per_hour: " / час",
    services_per_item: " / шт",
    services_per_person: " / гость",
    services_per_day: " / день",
    reviews_title: "Отзывы",
    reviews_loading: "Загрузка отзывов…",
    reviews_empty: "Пока нет отзывов.",
    reviews_submit: "Отправить отзыв",
    reviews_share_ph: "Поделись впечатлениями…",
    reviews_thanks: "Спасибо за отзыв!",
    notifications_title: "Уведомления",
    notifications_hint: "Обновления бронирований, модерации и платежей.",
    notifications_empty: "Уведомлений пока нет.",
    charts_bookings_per_day: "Брони по дням",
    charts_top_categories: "Топ категорий",
    charts_funnel: "Воронка одобрения подрядчиков",
    charts_no_data: "Пока нет данных",
    auth_forgot_title: "Сброс пароля",
    auth_forgot_hint: "Отправим письмо со ссылкой для сброса.",
    auth_forgot_btn_send: "Отправить ссылку",
    auth_forgot_done_title: "Проверь почту",
    auth_forgot_done_hint: "Если аккаунт существует — отправили ссылку. Следуй ей или вставь токен ниже.",
    auth_forgot_enter_token: "Ввести токен →",
    auth_forgot_back: "Назад ко входу",
    auth_forgot_link: "Забыли пароль?",
    auth_reset_title: "Новый пароль",
    auth_reset_token: "Токен сброса",
    auth_reset_token_ph: "Вставь из письма",
    auth_reset_new_password: "Новый пароль",
    auth_reset_btn: "Сменить пароль",
    auth_reset_done_title: "Пароль обновлён",
    auth_reset_done_hint: "Перенаправляем на вход…",
    auth_network_error: "Сетевая ошибка",
    chat_header_label: "Чат",
    booking_pay: "Оплатить",
    vendor_detail_book_now: "Забронировать",
  },
  kz: {
    status_online: "Онлайн · Gemini 2.5 Flash",
    hero_title_a: "Не жоспарлайсың,",
    hero_title_b: "Айгерім?",
    input_placeholder_main: "Іс-шараңды сипатта немесе сұра...",
    input_placeholder_reply: "Жауап беру...",
    input_hint:
      "qonaqzhai қателесуі мүмкін. Маңызды ақпаратты мердігермен тексеріңіз.",
    capability_plan: "Жоспарлау",
    capability_vendors: "Мердігер таңдау",
    capability_budget: "Бюджет",
    capability_guests: "Қонақтар",
    suggestion_toi: "Алматыда 150 қонаққа той жоспарла",
    suggestion_corporate: "80 адамға корпоративке бюджет ұсын",
    suggestion_photographer: "Астанадан фотограф тап",
    suggestion_invitation: "Тамыздағы үйленуге шақыру жаса",
    sidebar_new_chat: "Жаңа чат",
    sidebar_recent: "Соңғы",
    sidebar_plan: "Тегін тариф",
    thinking: "Ойланып жатырмын",
    settings_title: "Баптаулар",
    settings_back: "Чатқа оралу",
    settings_section_account: "Аккаунт",
    settings_section_pref: "Қалаулар",
    settings_section_about: "Қосымша туралы",
    settings_language: "Тіл",
    settings_language_hint: "Мәзір мен ұсыныс тілі.",
    settings_appearance: "Сыртқы көрініс",
    settings_appearance_hint: "Қалаған тақырыпты таңдаңыз.",
    settings_email: "Email",
    settings_phone: "Телефон",
    settings_plan_label: "Тариф",
    settings_plan_value: "Тегін · MVP",
    settings_version: "Нұсқа",
    settings_signout: "Шығу",
    settings_save: "Сақтау",
    settings_name: "Аты",
    settings_model: "Модель",
    settings_built_in: "Жасалған жер",
    theme_light: "Ашық",
    theme_dark: "Қараңғы",
    theme_system: "Жүйелік",
    block_event_plan: "той_жоспары",
    block_draft: "Жоба",
    block_budget: "бюджет",
    block_vendors: "мердігерлер",
    stat_date: "Күні",
    stat_city: "Қала",
    stat_guests: "Қонақ",
    stat_budget: "Бюджет",
    vendor_view_all: "Барлық нұсқалар",
    vendor_matches: "нұсқа",
    auth_welcome: "qonaqzhai-ге қош келдіңіз",
    auth_signin_tab: "Кіру",
    auth_signup_tab: "Тіркелу",
    auth_email: "Email",
    auth_password: "Құпиясөз",
    auth_name: "Аты",
    auth_signin_btn: "Кіру",
    auth_signup_btn: "Аккаунт жасау",
    auth_no_account: "Аккаунт жоқ па?",
    auth_has_account: "Аккаунт бар ма?",
    auth_loading: "Жүктеу...",
    auth_role_question: "Мен қалаймын...",
    auth_role_customer: "Іс-шара жоспарлау",
    auth_role_customer_hint: "Өзімнің тойымды/үйленуімді ұйымдастырамын",
    auth_role_vendor: "Қызмет көрсету",
    auth_role_vendor_hint: "Мен мердігермін (зал, кейтеринг, фото т.б.)",
    auth_demo_title: "Жылдам демо-кіру",
    nav_chat: "Чат",
    nav_vendors: "Мердігерлер",
    nav_bookings: "Брондар",
    nav_vendor_profile: "Профилім",
    nav_admin_users: "Қолданушылар",
    nav_notifications: "Хабарландырулар",
    nav_messages: "Хабарламалар",
    nav_admin_vendors: "Модерация",
    common_loading: "Жүктеу...",
    threads_title: "Хабарламалар",
    threads_hint: "Қабылданған брондар бойынша мердігер / клиентпен тікелей чат.",
    threads_empty: "Белсенді чат жоқ. Мердігер бронды қабылдағанда чат автоматты ашылады.",
    threads_unknown_vendor: "Мердігер",
    threads_unknown_user: "Қолданушы",
    notif_signup_welcome_title: "Қош келдіңіз",
    notif_signup_welcome_body: "Аккаунт дайын.",
    notif_password_reset_title: "Құпиясөзді қалпына келтіру",
    notif_password_reset_body: "Поштаңызға қалпына келтіру сілтемесін жібердік.",
    notif_booking_created_title: "Жаңа брондау сұранысы",
    notif_booking_created_body: "Жаңа сұраныс келді — брон қойындысын қараңыз.",
    notif_booking_accepted_title: "Брон қабылданды",
    notif_booking_accepted_body: "Мердігер броныңызды қабылдады.",
    notif_booking_declined_title: "Брон бас тартылды",
    notif_booking_declined_body: "Мердігер броннан бас тартты.",
    notif_booking_paid_title: "Төлем расталды",
    notif_booking_paid_body: "Брон төленді.",
    notif_vendor_approved_title: "Профиль мақұлданды",
    notif_vendor_approved_body: "Профиліңіз каталогта көрінеді.",
    notif_vendor_rejected_title: "Профиль қабылданбады",
    notif_vendor_rejected_body: "Модерация профильді қабылдамады.",
    notif_thread_message_title: "Жаңа хабарлама",
    notif_channel_email: "email",
    notif_channel_push: "push",
    notif_status_pending: "күтілуде",
    notif_status_sent: "жіберілді",
    notif_status_failed: "қате",
    booking_status_pending: "Күтілуде",
    booking_status_accepted: "Қабылданды",
    booking_status_declined: "Бас тартылды",
    booking_status_cancelled: "Болдырылмады",
    booking_status_completed: "Аяқталды",
    booking_status_paid: "Төленді",
    common_save: "Сақтау",
    common_cancel: "Бас тарту",
    common_search: "Іздеу...",
    common_back: "Артқа",
    common_delete: "Жою",
    cards_title: "Карталар",
    cards_hint: "Орындаушы броньды растағаннан кейін төлеу үшін карта қосыңыз. Мок — нақты төлем жасалмайды.",
    cards_empty: "Сақталған карта жоқ.",
    cards_add: "Карта қосу",
    cards_save: "Картаны сақтау",
    cards_default: "Әдепкі",
    cards_make_default: "Әдепкі ету",
    cards_field_number: "Карта нөмірі",
    cards_field_exp: "Мерзімі (АА/ЖЖ)",
    cards_field_holder: "Иесі",
    cards_err_exp: "Жарамсыз мерзім",
    cards_disclaimer: "Демо — нақты төлем провайдеріне жүгінбейміз, картаның толық нөмірін сақтамаймыз.",
    nav_cards: "Карталар",
    bookings_btn_pay: "Төлеу",
    bookings_btn_add_card: "Карта қосу",
    users_title: "Қолданушылар",
    users_hint: "Платформа қолданушыларын басқару.",
    users_filter_all: "барлығы",
    users_col_name: "Аты",
    users_col_email: "Email",
    users_col_role: "Рөлі",
    users_col_status: "Күйі",
    users_col_actions: "Әрекеттер",
    users_btn_suspend: "Бұғаттау",
    users_btn_activate: "Қосу",
    admin_title: "Админ · Мердігерлер",
    admin_hint: "Каталогта жариялау алдында мердігер профильдерін мақұлда не қабылдама.",
    admin_stat_users: "Қолданушылар",
    admin_stat_pending: "Күтілуде",
    admin_stat_approved: "Мақұлданды",
    admin_stat_bookings: "Брондар",
    admin_group_pending: "Күтілуде",
    admin_group_approved: "Мақұлданды",
    admin_group_rejected: "Қабылданбады",
    admin_empty: "Бос",
    admin_btn_approve: "Мақұлдау",
    admin_btn_reject: "Қабылдамау",
    admin_btn_suspend: "Алып тастау",
    admin_btn_reapprove: "Қайтару",
    funnel_submitted: "Жіберілді",
    funnel_pending: "Модерацияда",
    funnel_approved: "Мақұлданды",
    funnel_rejected: "Қабылданбады",
    admin_btn_preview: "Алдын ала қарау",
    admin_btn_hide: "Жасыру",
    admin_preview_no_photos: "Фото жоқ.",
    admin_preview_no_description: "Сипаттама жоқ.",
    admin_preview_no_services: "Қызметтер жоқ.",
    admin_preview_services: "Қызметтер",
    admin_preview_description: "Сипаттама",
    vendors_title: "Мердігерлер",
    vendors_hint: "Сенімді әріптестер сенің іс-шараңа.",
    vendors_search_ph: "Іздеу...",
    vendors_filter_all: "Барлығы",
    vendors_price_max_ph: "Макс. баға (₸)",
    vendors_sort_aria: "Сұрыптау",
    vendors_sort_newest: "Жаңалары",
    vendors_sort_price_asc: "Баға ↑",
    vendors_sort_price_desc: "Баға ↓",
    vendors_sort_rating_desc: "Үздік рейтинг",
    vendors_total_suffix: "барлығы",
    vendors_empty: "Мердігер табылмады",
    vendors_from: "бастап",
    vendor_detail_back: "Каталогқа",
    vendor_detail_starting: "Бастапқы баға",
    vendor_detail_date: "Іс-шара күні",
    vendor_detail_guests: "Қонақтар",
    vendor_detail_note: "Ескертпе",
    vendor_detail_note_ph: "Бір нәрсе қосу қажет пе...",
    vendor_detail_btn: "Бронға өтінім",
    vendor_detail_sent_title: "Өтінім жіберілді",
    vendor_detail_sent_hint: "Мердігер жақын арада жауап береді. Брондарыңды тексер.",
    vendor_detail_view_bookings: "Менің брондарым",
    vendor_detail_not_found: "Мердігер табылмады",
    bookings_title: "Менің брондарым",
    bookings_hint: "Мердігерлерге жіберген өтінімдерің.",
    bookings_empty: "Әзірге брон жоқ",
    bookings_btn_cancel: "Бас тарту",
    bookings_status_pending: "күтілуде",
    bookings_status_accepted: "қабылданды",
    bookings_status_declined: "қабылданбады",
    bookings_status_cancelled: "бас тартылды",
    vendor_profile_title: "Менің профилім",
    vendor_profile_hint: "qonaqzhai каталогындағы көпшілік картаң.",
    vendor_profile_section_basics: "Негізгі",
    vendor_profile_section_photos: "Фото",
    vendor_profile_field_name: "Атауы",
    vendor_profile_field_name_ph: "Rixos Almaty Ballroom",
    vendor_profile_field_category: "Санат",
    vendor_profile_field_city: "Қала",
    vendor_profile_field_price: "Баға (₸)",
    vendor_profile_field_desc: "Сипаттама",
    vendor_profile_field_desc_ph: "Алматы орталығындағы премиум зал...",
    vendor_profile_listed_from: "Каталогта бағасы",
    vendor_profile_photos_hint: "JPG/PNG 5MB дейін. Бірінші фото — мұқаба.",
    vendor_profile_no_photos: "Фото әлі жоқ",
    vendor_profile_btn_upload: "Жүктеу",
    vendor_profile_btn_delete: "Жою",
    vendor_profile_status_pending: "Модерацияда — әкімші мақұлдағанша профилің қонақтарға көрінбейді.",
    vendor_profile_status_approved: "Мақұлданды — профилің каталогта көрінеді.",
    vendor_profile_status_rejected: "Қабылданбады — қолдау қызметіне хабарлас.",
    vendor_bookings_title: "Брондар",
    vendor_bookings_hint: "Қонақтардан кіріс өтінімдер.",
    vendor_bookings_empty: "Әзірге брон жоқ",
    vendor_bookings_btn_accept: "Қабылдау",
    vendor_bookings_btn_decline: "Қабылдамау",
    vendor_bookings_guests: "қонақ",
    category_venue: "Зал",
    category_catering: "Кейтеринг",
    category_music: "Музыка және DJ",
    category_photo: "Фото және видео",
    category_decor: "Декор және флористер",
    category_cakes: "Торттар",
    category_other: "Басқа",
    city_almaty: "Алматы",
    city_astana: "Астана",
    city_shymkent: "Шымкент",
    city_karaganda: "Қарағанды",
    city_aktobe: "Ақтөбе",
    city_atyrau: "Атырау",
    status_pending: "күтілуде",
    status_accepted: "қабылданды",
    status_declined: "қабылданбады",
    status_cancelled: "бас тартылды",
    role_customer: "клиент",
    role_vendor: "мердігер",
    role_admin: "әкімші",
    user_status_active: "белсенді",
    user_status_suspended: "бұғатталған",
    services_title: "Қызметтер",
    services_hint: "Тапсырыс берушілерге көрсетілетін қызметтер — олар бронь кезінде біреуін таңдайды.",
    services_btn_new: "Жаңа қызмет",
    services_empty: "Қызметтер әлі жоқ — бірінші қызметті қос.",
    services_inactive: "жасырын",
    services_disable: "Жасыру",
    services_enable: "Көрсету",
    services_confirm_delete: "Бұл қызметті жоюды растайсыз ба?",
    services_dialog_new: "Жаңа қызмет",
    services_dialog_edit: "Қызметті өзгерту",
    services_field_name: "Атауы",
    services_field_name_ph: "Үйленуге фотосессия",
    services_field_description: "Сипаттама",
    services_field_price: "Баға (₸)",
    services_field_unit: "Өлшем",
    services_unit_fixed: "Бекітілген",
    services_unit_hour: "Сағатына",
    services_unit_item: "Дана үшін",
    services_unit_person: "Қонаққа",
    services_unit_day: "Күніне",
    services_save_failed: "Сақтау сәтсіз",
    services_load_failed: "Жүктеу сәтсіз",
    services_loading: "Қызметтер жүктелуде…",
    services_empty_public: "Мердігер әлі қызмет жарияламаған.",
    services_pick: "Қызметті таңда",
    services_per_hour: " / сағ",
    services_per_item: " / дана",
    services_per_person: " / қонақ",
    services_per_day: " / күн",
    reviews_title: "Пікірлер",
    reviews_loading: "Пікірлер жүктелуде…",
    reviews_empty: "Әзірге пікір жоқ.",
    reviews_submit: "Пікір қалдыру",
    reviews_share_ph: "Әсеріңізбен бөлісіңіз…",
    reviews_thanks: "Пікіріңіз үшін рахмет!",
    notifications_title: "Хабарландырулар",
    notifications_hint: "Бронь, мердігер модерациясы мен төлемдер бойынша жаңарулар.",
    notifications_empty: "Хабарландырулар әлі жоқ.",
    charts_bookings_per_day: "Күнделікті брондар",
    charts_top_categories: "Топ санаттар",
    charts_funnel: "Мердігер мақұлдау воронкасы",
    charts_no_data: "Деректер әлі жоқ",
    auth_forgot_title: "Құпиясөзді ұмыттыңыз ба",
    auth_forgot_hint: "Сізге қалпына келтіру сілтемесін жібереміз.",
    auth_forgot_btn_send: "Сілтеме жіберу",
    auth_forgot_done_title: "Пошта жәшігіңізді тексеріңіз",
    auth_forgot_done_hint: "Аккаунт бар болса — сілтеме жіберілді. Хатпен жалғастырыңыз немесе төменге токенді қойыңыз.",
    auth_forgot_enter_token: "Токенді енгізу →",
    auth_forgot_back: "Кіруге оралу",
    auth_forgot_link: "Құпиясөзді ұмыттыңыз ба?",
    auth_reset_title: "Жаңа құпиясөз",
    auth_reset_token: "Қалпына келтіру токені",
    auth_reset_token_ph: "Хаттан көшіріп қой",
    auth_reset_new_password: "Жаңа құпиясөз",
    auth_reset_btn: "Құпиясөзді ауыстыру",
    auth_reset_done_title: "Құпиясөз жаңартылды",
    auth_reset_done_hint: "Кіру бетіне бағыттаудамыз…",
    auth_network_error: "Желі қатесі",
    chat_header_label: "Чат",
    booking_pay: "Төлеу",
    vendor_detail_book_now: "Бронь",
  },
};
