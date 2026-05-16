# Qonaqzhai — Дипломная работа (бриф для письменной части)

Платформа для организации событий в Казахстане. Маркетплейс event-сервисов (площадки, кейтеринг, фото/видео, декор, музыка и т.д.) с AI-планировщиком, бронированием, оплатой и системой отзывов.

**Технологии:** Go (backend, Clean Architecture) · Next.js 16 (frontend, Feature-Sliced Design) · Flutter (мобильное приложение) · SQLite / PostgreSQL · Gemini AI · Gmail SMTP · Firebase (FCM push) · Freedom Pay / PayBox (test mode) · Docker · GitHub Actions · nginx + Let's Encrypt.

---

## Register / Login

- Регистрация и вход для трёх ролей: customer, vendor, admin
- JWT access token (короткий TTL) + refresh token с ротацией
- Восстановление пароля по email (one-time token)
- Хеширование паролей bcrypt (cost 12)
- Rate limiting на auth endpoints

## Create Event Request

- Customer создаёт запрос с параметрами: дата, количество гостей, бюджет, город, стиль
- Запрос используется AI-планировщиком и поиском vendors
- Поддержка KZ / RU / EN языков

## Use AI Planner

- Чат с Gemini AI
- Prompt engineering: zero-shot prompt + injection реальных vendors из БД
- AI рекомендует подходящих исполнителей под бюджет и описание события
- Ограничение запросов на пользователя (rate limit)

## Browse Vendors (search + filter)

- Единая сущность Vendor с категориями: Venue, Catering, Photo, Video, Decor, Music, Traditional Services и др.
- Поиск по названию и описанию (SQLite FTS5)
- Фильтры: категория, город, диапазон цены, минимальный рейтинг
- Сортировка: по цене, по рейтингу, по новизне
- Пагинация с метаданными
- URL-синхронизация фильтров на фронте

## Book Vendor

- Customer создаёт бронирование выбранного vendor
- Lifecycle: pending → accepted / declined → completed / cancelled
- Vendor получает уведомление в inbox
- Customer видит статус в своём кабинете

## Payment

- **Freedom Pay (PayBox.money)** test mode — локальный казахстанский платёжный агрегатор
- Поддерживает Kaspi QR, Halyk, банковские карты Visa/Mastercard через единый шлюз
- Endpoint создания платёжного intent → редирект customer на PayBox checkout
- Callback / webhook с проверкой подписи → перевод booking в статус paid
- Идемпотентность webhook-ов по transaction id
- Stripe рассмотрен и отклонён: Stripe не поддерживает мерчантов в Казахстане
- Резервный fallback: stub-checkout для демонстрации UI-флоу до подписания договора с PSP

## Reviews & Ratings

- Отзыв доступен только customer'у с completed booking (один отзыв на бронирование)
- Шкала 1–5 звёзд + текстовый комментарий
- Средний рейтинг кэшируется в записи vendor, пересчитывается при insert/delete
- Модерация отзывов админом (удаление abuse)

## Moderation

- Admin одобряет/отклоняет/блокирует vendors
- Admin блокирует/разблокирует users
- Admin модерирует отзывы
- Журнал действий admin для аудита

## Notifications

- **Email (Gmail SMTP):** подтверждение регистрации, восстановление пароля, статусы бронирований, одобрение vendor
- **Push (Firebase Cloud Messaging):** статусы бронирований и платежей для мобильного приложения (Flutter)
- In-process очередь (Go channels + worker goroutine), асинхронная отправка
- Таблица notifications для inbox внутри приложения

## Admin Dashboard & Analytics

- KPI-плитки: всего users, vendors, bookings, revenue, DAU
- Графики: bookings по дням (line chart), топ категорий (bar chart), воронка одобрения vendors
- Time series endpoint с фильтрами from / to / metric
- Библиотека графиков: Recharts (frontend)

---

## Component Diagram

Архитектура включает:

- **Frontend Web** — Next.js 16, App Router, Feature-Sliced Design
- **Mobile App** — Flutter, общий OpenAPI-клиент
- **Reverse Proxy** — nginx с TLS (Let's Encrypt)
- **Backend Monolith (Go)** — Clean Architecture с разделением на слои HTTP / Use Cases / Domain / Adapters
- **Database** — SQLite (MVP) / PostgreSQL (production-ready)
- **External Integrations:**
  - Gemini AI (планировщик)
  - Gmail SMTP (email)
  - Firebase Cloud Messaging (push)
  - Freedom Pay / PayBox (платежи)
- **In-process Notification Worker** — goroutine + channel queue

**Примечание:** микросервисы, Kafka, ELK, отдельный Auth Server рассмотрены и отклонены на этапе MVP. Решение обосновано в записке: монолит с Clean Architecture обеспечивает эквивалентную модульность на уровне кода без операционных издержек распределённой системы.

## Package Diagram

**Backend (Go, Clean Architecture):**

- `cmd.qonaqzhai` — composition root
- `internal.domain` — сущности, value objects, бизнес-ошибки
- `internal.usecase` — application services + port interfaces
- `internal.adapter.http` — handlers, middleware, DTO, router
- `internal.adapter.repo` — реализации репозиториев (SQLite / Postgres)
- `internal.adapter.ai` — Gemini client
- `internal.adapter.mail` — SMTP client
- `internal.adapter.push` — Firebase Cloud Messaging client
- `internal.adapter.pay` — Freedom Pay / PayBox client
- `internal.infra.db` — соединение, миграции
- `internal.infra.token` — JWT issuer
- `internal.infra.config` — env loading
- `internal.infra.logger` — slog
- `internal.app` — DI wiring

**Frontend (Next.js, FSD):**

- `app` — Next.js App Router
- `pages` — композиции страниц
- `widgets` — сложные UI-блоки
- `features` — пользовательские сценарии
- `entities` — доменные сущности
- `shared.ui` — переиспользуемые компоненты
- `shared.api` — HTTP-клиент
- `shared.lib` — утилиты и хуки
- `shared.config` — env, константы
- `shared.i18n` — переводы

## Object Diagram

Примеры объектов в runtime во время booking flow:

- **user** — Customer "Aigerim", active, role = customer
- **vendor** — "Rixos Almaty Ballroom", category = Venue, city = Almaty, priceFrom = 1 500 000 ₸, status = approved, ratingAvg = 4.7
- **booking** — связь user ↔ vendor, дата мероприятия, статус = accepted, amount = 1 800 000 ₸, paymentId
- **aiPlan** — массив ChatMessage из диалога customer ↔ Gemini
- **review** — связан с booking, rating = 5, text
- **notification** — type = booking.created, channel = email+push, status = sent

## Sequence Diagram

**Главный процесс: создание события и бронирование**

1. Customer открывает AI-чат → backend пересылает запрос в Gemini с инжектом vendors из БД → AI возвращает рекомендации
2. Customer применяет фильтры в каталоге → backend выполняет SQL-запрос с FTS5 и фильтрами → возвращает страницу vendors
3. Customer создаёт booking → backend INSERT booking(pending), notification worker отправляет email + push vendor'у
4. Customer инициирует оплату → backend создаёт PayBox payment → редирект на Freedom Pay checkout
5. После оплаты Freedom Pay вызывает callback с подписью → backend UPDATE booking(paid), отправляет уведомления обеим сторонам
6. Vendor нажимает Accept → backend UPDATE booking(accepted), notification customer'у
7. После события customer оставляет отзыв → backend INSERT review, пересчёт ratingAvg у vendor

## Activity Diagram

**Главный поток: пользователь планирует событие**

1. Login / Signup
2. Открытие AI-планировщика
3. Описание события (дата, гости, бюджет, стиль, город)
4. AI возвращает рекомендации vendors (Gemini + инжект реальных vendors)
5. Уточнение через фильтры (категория / цена / город / рейтинг)
6. Открытие страницы vendor + чтение отзывов
7. Создание бронирования
8. Ветка: оплата сейчас (Freedom Pay / PayBox) ИЛИ оплата позже (pending payment)
9. Vendor принимает / отклоняет → уведомление customer'у (email + push)
10. Событие проходит → booking переходит в completed
11. Customer оставляет отзыв с рейтингом
12. Конец

## High-Level Architecture Diagram

**Слои:**

- **Client Layer** — Next.js Web (FSD) + Flutter Mobile
- **Edge Layer** — nginx reverse proxy + TLS
- **Application Layer** — Go monolith (Clean Architecture) с in-process notification worker
- **Data Layer** — SQLite / PostgreSQL + миграции (golang-migrate)
- **Integration Layer** — Gemini AI, Gmail SMTP, Firebase Cloud Messaging, Freedom Pay / PayBox
- **Cross-cutting** — slog structured logging, /healthz, correlation ID middleware
- **DevOps** — Docker Compose, GitHub Actions CI/CD, опционально Uptime Kuma для мониторинга

**Не используется (обоснованно):**

- API Gateway (один backend, nginx достаточен)
- Микросервисы (монолит с Clean Arch)
- Kafka / RabbitMQ (in-process channel queue)
- ELK / Prometheus / Grafana (slog + Uptime Kuma)
- Отдельный Auth Server (JWT внутри монолита)
- Redis (не нужен на масштабе MVP)

---

## Academic Structure (структура записки)

### 1. Теоретические основы проекта

- Индустрия event-сервисов в Казахстане
- Цифровизация малого и среднего бизнеса
- AI-рекомендательные системы (LLM-based)
- Маркетплейс-платформы: бизнес-модели и монетизация
- Анализ существующих решений и выявленных пробелов (GoSwana, Wezoom, Marry.kz, международные The Bash, WeddingWire)
- Методологии: Clean Architecture (R. Martin) и Feature-Sliced Design

### 2. Системный анализ и проектирование

- Функциональные и нефункциональные требования
- Архитектура решения
- ER-диаграмма базы данных
- UML-диаграммы (Use Case, Class, Sequence, Activity, Component, Package, Deployment)
- High-level архитектурная схема

### 3. Реализация, интеграция и тестирование

- Реализация backend (Go, Clean Architecture)
- Реализация frontend (Next.js + FSD)
- Реализация мобильного приложения (Flutter)
- Схема БД и система миграций
- Дизайн REST API + OpenAPI 3 (Swagger UI)
- Интеграция Gemini AI (prompt engineering)
- Интеграция Gmail SMTP, Firebase Cloud Messaging, Freedom Pay / PayBox
- Безопасность: JWT + refresh, password reset, rate limiting, bcrypt, OWASP Top 10
- Тестирование: unit, integration, E2E (Playwright); coverage 80%+

### 4. Развёртывание, поддержка и влияние проекта

- Контейнеризация (multi-stage Dockerfiles)
- Оркестрация (Docker Compose)
- Reverse proxy + HTTPS (nginx + certbot)
- CI/CD (GitHub Actions: lint → test → build → deploy)
- Логирование (slog structured JSON)
- Мониторинг (healthcheck + Uptime Kuma)
- Экономический эффект
- Социальный эффект
- Ограничения и направления развития
- Результаты исследовательской части

---

## Functional Requirements

**Customer может:**

- зарегистрироваться / войти / выйти / сбросить пароль
- создать запрос на событие
- использовать AI-планировщик
- искать и фильтровать vendors
- бронировать услуги vendor
- оплачивать онлайн (Freedom Pay / PayBox: Kaspi QR, Halyk, карты)
- оставлять отзывы и рейтинг
- получать уведомления (email + push)
- переключать язык (KZ / RU / EN)

**Vendor может:**

- зарегистрироваться как vendor
- управлять профилем (название, категория, город, цена, описание, фото)
- получать запросы на бронирование
- принимать / отклонять / завершать бронирования
- видеть свои отзывы
- получать уведомления

**Admin может:**

- модерировать vendors (одобрить / отклонить / заблокировать)
- управлять users (заблокировать / разблокировать)
- модерировать отзывы
- видеть статистику и аналитику (KPI + графики)

## Non-Functional Requirements

- **Безопасность:** OWASP Top 10, JWT + refresh rotation, bcrypt cost 12, rate limiting, HTTPS, параметризованные SQL-запросы, валидация входных данных
- **Производительность:** p95 латентность API < 300 ms (без AI), поисковый запрос < 100 ms с индексами
- **Надёжность:** stateless backend, ACID-транзакции, идемпотентные webhook-и
- **Масштабируемость:** stateless монолит, готовый к горизонтальному масштабированию; БД легко мигрируется на Postgres
- **Наблюдаемость:** structured logs (slog), `/healthz`, correlation ID
- **Расширяемость:** границы Clean Architecture позволяют менять адаптеры без правок в usecases
- **Кроссплатформенность:** Web (Next.js) + Mobile (Flutter, iOS + Android)
- **Интернационализация:** runtime-переключение KZ / RU / EN
- **Developer Experience:** OpenAPI спецификация, Swagger UI, code coverage 80%+

---

## Economic and Social Impact

**Экономический:**

- Цифровизация SME-сегмента event-индустрии Казахстана
- Снижение барьера выхода в онлайн для vendors без своего сайта
- Комиссионная модель монетизации (5–10% с booking)
- Снижение издержек customer'а: одна платформа вместо обзвона десятков vendors (экономия ~15 ч на событие)
- Поддержка локальной экономики через интеграцию с региональными платёжными системами

**Социальный:**

- Упрощение организации культурных событий (той, бесік-той, қыз ұзату)
- Поддержка традиций через категории национальных услуг (домбра-ансамбли, баурсак-кейтеринг, шанырак-декор)
- Доступность через mobile-first подход
- Инклюзивность благодаря трёхъязычному интерфейсу (KZ / RU / EN)
- Прозрачность рынка через систему отзывов и рейтингов
- Повышение качества услуг за счёт конкуренции по рейтингу

---

## Future Research Directions

- Generative AI: LLM-генерация полного сценария события с таймлайном
- Predictive analytics: прогнозирование спроса по категориям и сезонам
- Анализ объяснимости AI (explainability): обоснование рекомендаций
- Обнаружение мошенничества: выявление аномалий в отзывах и платежах
- Поведенческая аналитика: оптимизация воронки customer journey
- Federated reputation: переносимость рейтингов между платформами
- A/B-тестирование AI-промптов с измерением relevance и cost

**Тема дипломного исследования:** AI-ассистированная рекомендация vendors — сравнение zero-shot prompt и prompt с инжектом каталога. Метрики: релевантность (оценка 1–5), latency, стоимость токенов. Корпус: 50 тестовых запросов на KZ/RU. Результат: количественное обоснование выбранной prompt-стратегии.

---

## Что отсутствует и обосновано в записке

- Микросервисная архитектура → монолит (Clean Architecture обеспечивает модульность)
- Очереди сообщений (Kafka / RabbitMQ) → in-process channel queue
- ELK / Prometheus / Grafana → slog structured + Uptime Kuma
- Redis → не требуется на масштабе MVP
- Отдельный Auth Server → JWT в монолите
- Restaurant / Rental items как разные сущности → унифицированная сущность Vendor с полем category
- Customer-support chat → AI-планировщик закрывает основной сценарий общения
