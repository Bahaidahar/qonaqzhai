# Qonaqzhai — Diploma Brief

Event services marketplace for Kazakhstan. Customers plan events (weddings, toi, beshik-toi, corporate, birthdays) with an AI assistant, browse local vendors (venues, catering, photo/video, decor, music, traditional services), book and pay online, leave reviews. Vendors manage their profile and bookings. Admins moderate the platform.

**Tech stack:** Go (backend, Clean Architecture) · Next.js 16 (frontend, Feature-Sliced Design) · Flutter (mobile app) · SQLite / PostgreSQL · Gemini AI · Gmail SMTP · Firebase Cloud Messaging · Freedom Pay / PayBox (test mode) · Docker · GitHub Actions.

---

## Register / Login

- Three user roles: customer, vendor, admin
- JWT access token (short TTL) + refresh token rotation
- Password reset via email (one-time token)
- bcrypt (cost 12) password hashing
- Rate limiting on auth endpoints

## Create Event Request

- Customer fills event parameters: date, guest count, budget, city, style
- Request feeds the AI planner and vendor search
- Supports KZ / RU / EN languages

## Use AI Planner

- Chat with Gemini AI
- Prompt engineering: zero-shot prompt + injection of real vendors from the database
- AI recommends matching vendors based on budget and event description
- Per-user rate limit on AI calls

## Browse Vendors

- Unified Vendor entity with categories: Venue, Catering, Photo, Video, Decor, Music, Traditional Services, etc.
- Full-text search on name and description (SQLite FTS5)
- Filters: category, city, price range, minimum rating
- Sort: price asc/desc, rating desc, newest
- Pagination with metadata
- URL-synced filters on the frontend

## Book Vendor

- Customer creates a booking for a selected vendor
- Lifecycle: pending → accepted / declined → completed / cancelled
- Vendor receives the booking in their inbox
- Customer tracks status from their dashboard

## Place Order (Booking)

- Customer confirms event details, dates, options
- Idempotent order creation (no duplicate bookings on retry)
- Server-side validation of vendor availability

## Payment

- **Freedom Pay (PayBox.money)** test mode — local Kazakhstani payment aggregator
- Aggregates Kaspi QR + Halyk Bank + Visa/Mastercard cards through one gateway
- Endpoint creates a payment intent → customer is redirected to Freedom Pay checkout
- Signed callback / webhook → booking transitions to "paid"
- Idempotent webhook handling by transaction id
- Stripe was considered and rejected: Stripe does not support merchants in Kazakhstan
- Fallback: a stub checkout for UI demo until a PSP contract is signed

## Reviews

- Review can be left only by a customer with a completed booking (one review per booking)
- 1–5 star rating + text comment
- Vendor average rating cached on the vendor row, recomputed on review insert / delete
- Admin can delete abusive reviews

## Moderation

- Admin approves / rejects / suspends vendors
- Admin suspends / restores users
- Admin moderates reviews
- Audit log of admin actions

---

## Component Diagram

Architecture includes:

- **Frontend Web** — Next.js 16, App Router, Feature-Sliced Design
- **Mobile App** — Flutter, shared OpenAPI client
- **Backend Monolith (Go)** — Clean Architecture: HTTP / Use Cases / Domain / Adapters layers
- **Database** — SQLite (MVP) / PostgreSQL (production)
- **AI Planner** — Gemini AI integration
- **Notification Service** — in-process worker (goroutine + channel queue), email + push
- **Payment Service** — Freedom Pay / PayBox adapter
- **External Integrations** — Gemini AI, Gmail SMTP, Firebase Cloud Messaging, Freedom Pay

**Justification (stated in thesis):** classical microservice components (API Gateway, separate Auth Server, message broker) were evaluated and rejected at the MVP scale. A Clean Architecture monolith provides equivalent code-level modularity with O(1) operational complexity. The architecture preserves a clear migration path: each use case package is a candidate microservice.

## Package Diagram

Packages:

- **UI** — `src/app`, `src/pages`, `src/widgets`, `src/features`, `src/entities`, `src/shared/ui` (Next.js + FSD)
- **API Layer** — `internal/adapter/http` (handlers, middleware, DTOs, router) + `src/shared/api` on the frontend
- **Services** — `internal/usecase/*` (application logic, port interfaces)
- **Data Access** — `internal/adapter/repo` + `internal/infra/db` (repositories, migrations)
- **External Integrations** — `internal/adapter/{ai,mail,push,pay}` (Gemini, SMTP, FCM, Freedom Pay)
- **Monitoring** — `internal/infra/logger` (slog structured logging), `/healthz` endpoint, optional Uptime Kuma
- **Security** — `internal/infra/token` (JWT), middleware (auth, rate limit, CORS), validation, bcrypt

## Object Diagram

Objects included (sample runtime state during a booking flow):

- **user** — Customer "Aigerim", role = customer, status = active
- **booking order** — bookingId, customerId, vendorId, eventDate, status = accepted, amount = 1 800 000 ₸, paymentId
- **vendor** — "Rixos Almaty Ballroom", category = Venue, city = Almaty, priceFrom = 1 500 000 ₸, ratingAvg = 4.7, status = approved
- **vendor service item** — service offered by a vendor (e.g. "Ballroom rental, 200 guests", price, photos)
- **AI plan** — array of ChatMessage exchanged between customer and Gemini, including budget breakdown and vendor recommendations
- **review** — linked to a completed booking, rating = 5, text
- **notification** — type = booking.created, channel = email + push, status = sent
- **payment** — Freedom Pay transactionId, amount, status = paid

## Sequence Diagram

Main process: **event creation and booking workflow**

1. Customer opens the AI chat → backend forwards the prompt to Gemini with injected vendor catalog → AI returns suggestions
2. Customer applies filters in the catalog → backend runs an FTS5 + filtered SQL query → returns a paginated vendor list
3. Customer creates a booking → backend INSERTs booking(pending) → notification worker emails and pushes the vendor
4. Customer initiates payment → backend creates a Freedom Pay payment intent → redirect to Freedom Pay checkout
5. After payment Freedom Pay calls the signed webhook → backend UPDATEs booking(paid) → notifies both parties
6. Vendor clicks Accept → backend UPDATEs booking(accepted) → notification to the customer
7. After the event, customer leaves a review → backend INSERTs the review → recomputes vendor ratingAvg

## Activity Diagram

Main flow:

1. **User planning event** — login, open AI planner, describe event (date, guests, budget, style, city)
2. **AI recommendation** — Gemini returns vendor suggestions enriched with the real catalog
3. Refine results through filters (category / price / city / rating)
4. Open vendor page and read reviews
5. **Booking** — create booking record
6. **Payment** — pay now via Freedom Pay OR defer payment (pending)
7. Vendor accepts or declines → email + push notification to customer
8. **Confirmation** — booking transitions to accepted, then to completed after the event
9. Customer leaves a rating and review

## High-Level Architecture Diagram

Layers:

- **Client Layer** — Next.js Web (FSD) + Flutter Mobile
- **API Layer** — single Go backend exposing REST + OpenAPI 3
- **Application Layer** — Clean Architecture use cases, in-process notification worker
- **Integration Layer** — Gemini AI, Gmail SMTP, Firebase Cloud Messaging, Freedom Pay / PayBox
- **Data Layer** — SQLite / PostgreSQL with golang-migrate migrations
- **Monitoring & Logging** — slog structured JSON logs, `/healthz` endpoint, request correlation IDs, optional Uptime Kuma dashboard
- **External Services** — Gemini AI, Google SMTP, Firebase, Freedom Pay

Includes:

- Authentication handled inside the monolith (JWT issuer in `internal/infra/token`) — no separate Auth Server required
- Go monolith with Clean Architecture layers — modular at code level instead of process level
- PostgreSQL ready (SQLite as MVP default, single config switch)
- In-process channel queue for notifications (replaces Kafka / RabbitMQ at MVP scale)
- slog structured logging (replaces ELK at MVP scale)
- `/healthz` + Uptime Kuma (lightweight substitute for Prometheus / Grafana)
- Redis intentionally omitted — rating caching is column-level in the database; rate-limit state is in-memory

**Rejected with justification:** API Gateway as a separate service, dedicated Auth Server, microservices split, Kafka / RabbitMQ, ELK, Prometheus / Grafana, Redis. All are documented in the thesis as future migration steps once scale demands them.

---

## Academic Structure

The diploma structure includes:

### 1. Theoretical Foundations of the Project

Topics:

- Event management industry in Kazakhstan
- Digital transformation of small and medium enterprises (SME)
- AI recommendation systems and large language models
- Event-services marketplaces: business models and monetization
- Existing solutions and gaps (GoSwana, Wezoom, Marry.kz, international platforms The Bash, WeddingWire)
- Methodologies: Clean Architecture (R. Martin) and Feature-Sliced Design

### 2. System Analysis and Design Modelling

Includes:

- Functional and non-functional requirements
- System analysis and use case modelling
- ER diagram of the database
- UML diagrams: Use Case, Class, Sequence, Activity, Component, Package, Deployment
- High-level architecture diagram
- API contract design (OpenAPI 3)

### 3. Implementation, Integration and Testing

Includes:

- Backend implementation in Go with Clean Architecture (layer-by-layer walkthrough)
- Frontend implementation in Next.js with Feature-Sliced Design
- Mobile implementation in Flutter (shared OpenAPI client)
- Database schema and migration system (golang-migrate)
- API integration (REST + Swagger UI)
- AI integration: Gemini prompt engineering, vendor catalog injection
- External integrations: Gmail SMTP, Firebase Cloud Messaging, Freedom Pay / PayBox
- Security implementation: JWT + refresh tokens, password reset, rate limiting, bcrypt, OWASP Top 10
- Testing: unit, integration, end-to-end (Playwright); target coverage 80 %+

### 4. Deployment, Support and Project Impact

Includes:

- Containerization (multi-stage Dockerfiles)
- Orchestration (Docker Compose)
- CI / CD (GitHub Actions: lint → test → build → deploy)
- Monitoring (healthcheck endpoint, optional Uptime Kuma)
- Logging strategy (slog structured JSON)
- Economic impact assessment
- Social impact assessment
- Project limitations and future improvements
- Results of the research component

---

## Functional Requirements

**Users (customers) can:**

- register / login / logout / reset password
- create event requests
- use the AI planner
- browse vendors (search and filter)
- book vendor services
- pay online (Freedom Pay: Kaspi QR, Halyk, cards)
- leave reviews and ratings
- receive notifications (email + push)
- switch interface language (KZ / RU / EN)

**Admins can:**

- moderate vendors (approve / reject / suspend)
- manage users (suspend / restore)
- monitor transactions and platform KPIs
- moderate reviews
- view analytics dashboard with charts

**Business partners (vendors) can:**

- register and manage their profile (name, category, city, price, description, photos)
- manage their service offerings
- confirm or decline booking requests
- mark bookings as completed
- view their reviews and rating
- receive notifications

## Non-Functional Requirements

System must provide:

- **Scalability** — stateless backend ready for horizontal scaling behind a load balancer; database extractable to PostgreSQL with one config switch
- **Reliability** — ACID database transactions, idempotent webhooks, graceful shutdown, restart-safe stateless services
- **Security** — OWASP Top 10 mitigations, JWT short TTL + refresh rotation, bcrypt cost 12, rate limiting, HTTPS only in production, parameterized SQL, input validation
- **High performance** — p95 API latency < 300 ms excluding AI calls, search query < 100 ms with proper indexes
- **Monitoring / logging** — slog structured JSON logs, `/healthz` endpoint, request correlation IDs, optional Uptime Kuma dashboard
- **Extensibility** — Clean Architecture boundaries allow swapping adapters (DB engine, payment provider, push provider) without touching business logic
- **Cross-platform support** — Web (Next.js, all modern browsers) + Mobile (Flutter for iOS and Android)
- **Internationalization** — runtime language switching across KZ / RU / EN
- **Developer experience** — OpenAPI 3 contract, Swagger UI, automated tests, CI gates

---

## Economic and Social Impact Ideas

**Economic:**

- **Marketplace growth** — digitalization of a fragmented event-services market in Kazakhstan
- **SME support** — vendors gain an online sales channel without building their own website
- **Automation** — automated booking, payment and notification workflows replace manual phone-based coordination
- **Reduced planning costs** — customers save approximately 15 hours per event by using one platform instead of contacting many vendors individually
- Commission-based revenue model (5–10 % per booking) generates marketplace income
- Local payment integration (Freedom Pay aggregating Kaspi QR and Halyk) keeps transaction fees in-region

**Social:**

- **Easier event organization** — single interface from planning through payment to feedback
- **Support of local traditions** — dedicated categories for toi, beshik-toi, kyz uzatu, traditional music ensembles, national cuisine catering, national decor
- **Digital transformation** — onboards event-services SMEs into the digital economy
- **Accessibility** — mobile-first design reaches users without desktops; trilingual interface (KZ / RU / EN) lowers language barriers
- **Trust and transparency** — public reviews and ratings improve service quality through competition
- **Time recovery** — automating coordination frees families for the human side of celebrations

---

## Future Research Directions

Potential research areas:

- **AI recommendation systems** — benchmarking prompt strategies (zero-shot vs catalog-injected) on a Kazakh / Russian query corpus
- **Predictive analytics** — demand forecasting by category and season
- **Generative AI event planning** — full event scenario generation with timelines, budget breakdowns, vendor sequences
- **Distributed systems** — migration path from monolith to microservices when scale demands it
- **AI explainability** — surfacing why each vendor was recommended, increasing user trust
- **Fraud detection** — anomaly scoring on review patterns and payment flows
- **Behavioral analytics** — customer journey funnel optimization and conversion analysis
- **Federated reputation** — portability of vendor ratings across platforms
- **Multilingual NLP** — improving Gemini prompt quality for Kazakh-language event requests
