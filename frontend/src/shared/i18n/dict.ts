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
  nav_admin_vendors: string;
  // common
  common_loading: string;
  common_save: string;
  common_cancel: string;
  common_search: string;
  common_back: string;
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
  // vendors catalog (customer)
  vendors_title: string;
  vendors_hint: string;
  vendors_search_ph: string;
  vendors_filter_all: string;
  vendors_empty: string;
  vendors_from: string;
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
    nav_admin_vendors: "Vendor moderation",
    common_loading: "Loading...",
    common_save: "Save",
    common_cancel: "Cancel",
    common_search: "Search...",
    common_back: "Back",
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
    vendors_title: "Vendors",
    vendors_hint: "Verified partners for your next event.",
    vendors_search_ph: "Search...",
    vendors_filter_all: "All",
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
    nav_admin_vendors: "Модерация",
    common_loading: "Загрузка...",
    common_save: "Сохранить",
    common_cancel: "Отмена",
    common_search: "Поиск...",
    common_back: "Назад",
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
    vendors_title: "Подрядчики",
    vendors_hint: "Проверенные партнёры для твоего события.",
    vendors_search_ph: "Поиск...",
    vendors_filter_all: "Все",
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
    nav_admin_vendors: "Модерация",
    common_loading: "Жүктеу...",
    common_save: "Сақтау",
    common_cancel: "Бас тарту",
    common_search: "Іздеу...",
    common_back: "Артқа",
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
    vendors_title: "Мердігерлер",
    vendors_hint: "Сенімді әріптестер сенің іс-шараңа.",
    vendors_search_ph: "Іздеу...",
    vendors_filter_all: "Барлығы",
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
  },
};
