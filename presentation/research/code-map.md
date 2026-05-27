# Qonaqzhai Platform: Code Architecture Map

**Generated**: May 28, 2026  
**Scope**: Complete architecture overview of the diploma project's backend services, APIs, mobile client, and frontend web application.

---

## 1. Domain Entities & Bounded Contexts

The platform organizes business logic into four microservices, each owning distinct bounded contexts:

### Auth Service
- **User** — Authenticated principal with email, name, role (customer/vendor/admin), status (active/suspended), password hash, timestamps
- **Role** — Enumeration: `customer`, `vendor`, `admin`
- **RefreshToken** — Long-lived JWT token stored hashed; tracks expiry and revocation
- **PasswordResetToken** — One-time, short-TTL recovery token; tracks usage

### Core Service
- **Vendor** — Business profile owned by vendor-role user; includes name, category, city, description, price range, moderation status (pending/approved/rejected), rating aggregates, photo IDs
- **Service** — Menu item per vendor: name, description, price, unit (fixed/hour/item/person/day), activation flag
- **Booking** — Customer reservation: customer ID, vendor ID, service ID, event date, guest count, note, lifecycle status (pending/accepted/declined/cancelled/completed/paid), amount, linked payment ID
- **Review** — Customer evaluation post-booking: 1–5 star rating, text, linked to booking and vendor for rating aggregation
- **Notification** — In-app + delivery record: type, channel (email/push/email+push), title, body, status (queued/sent/failed)
- **FCMToken** — Device token for Firebase Cloud Messaging push delivery; tracks platform (iOS/Android)
- **Photo** — Vendor profile image: MIME type, size in bytes, binary data; max 5 MB, allowed formats: JPEG, PNG, WebP, GIF

### Payment Service
- **Card** — Saved payment instrument: brand (visa/mastercard/amex/discover), last 4 digits, expiry month/year, holder name, default flag; raw PAN never persisted (PCI scope sits with PSP)
- **Payment** — Charge attempt against booking: amount, currency (KZT), lifecycle status (pending/captured/failed/refunded), provider reference for PSP lookup

### Realtime Service
- **Thread** — DM channel attached to booking: customer ID, vendor ID (both foreign UUIDs from auth-svc), booking ID
- **Message** — Single chat line: sender ID, text, timestamp

---

## 2. HTTP Endpoints by Service

### Auth Service (Port 3001)
| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/health` | Liveness probe |
| POST | `/api/signup` | Register new user (email, password, optional name, role defaults to customer) |
| POST | `/api/login` | Obtain JWT + refresh token (email, password) |
| POST | `/api/refresh` | Rotate tokens using valid refresh token |
| POST | `/api/logout` | Revoke refresh token |
| POST | `/api/forgot-password` | Initiate password reset flow (email) |
| POST | `/api/reset-password` | Consume password reset token + set new password |
| GET | `/api/me` | Fetch authenticated user profile (requires token) |

**Authentication**: Token-verified endpoints require JWT in `Authorization: Bearer <token>` header. Rate limited to 20 req/IP/10s.

### Core Service (Port 3002)
| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/health` | Liveness probe |
| GET | `/api/vendors` | Public catalog search: filter by category, city, price range; sort by newest/price/rating |
| GET | `/api/vendors/{id}` | Public vendor detail view |
| GET | `/api/photos/{id}` | Serve vendor photo binary (MIME-typed response) |
| GET | `/api/vendors/{vendorId}/reviews` | List reviews for vendor (with pagination) |
| GET | `/api/me/vendor` | Fetch authenticated vendor's own profile (vendor-role only) |
| POST / PUT | `/api/me/vendor` | Create or update vendor profile (vendor-role) |
| POST | `/api/me/vendor/photos` | Upload vendor photo (multipart/form-data; vendor-role) |
| DELETE | `/api/me/vendor/photos/{id}` | Remove photo (vendor-role) |
| POST | `/api/bookings` | Create booking request (customer-role; requires vendor ID, service ID, date, guest count, amount) |
| GET | `/api/bookings` | List authenticated user's bookings (customer or vendor view) |
| GET | `/api/bookings/{id}` | Fetch booking detail + linked payment info |
| PATCH | `/api/bookings/{id}` | Transition booking status (vendor: accept/decline/complete; customer: cancel) |
| POST | `/api/bookings/{id}/pay` | Mark booking paid after successful payment |
| POST | `/api/reviews` | Submit review post-booking (customer-role; rating 1–5, text) |
| GET | `/api/notifications` | List user's notifications (inbox view with pagination) |
| POST | `/api/notifications/fcm` | Register FCM device token for push delivery |
| PATCH | `/api/admin/vendors/{id}/status` | Admin approval/rejection of vendor profile (admin-role) |
| GET | `/api/admin/stats` | Admin dashboard: vendor/booking counts, GMV, etc. (admin-role) |

**Authentication**: Most endpoints require JWT. Rate limited to 30 req/IP/10s. Admins use `mw.RequireRole("admin")` middleware.

### Payment Service (Port 3003)
| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/health` | Liveness probe |
| GET | `/api/cards` | List authenticated user's saved cards (requires token) |
| POST | `/api/cards` | Add new card (PAN + expiry + holder name; validated but not stored) |
| DELETE | `/api/cards/{id}` | Remove card from vault |
| POST | `/api/cards/{id}/default` | Set card as default payment method |
| GET | `/api/payments` | List authenticated user's payment history |

**Authentication**: All endpoints except `/api/health` require JWT. Rate limited to 20 req/IP/10s.

### Realtime Service (Port 3004)
| Method | Path | Purpose |
|--------|------|---------|
| GET | `/api/health` | Liveness probe |
| GET | `/api/threads` | List authenticated user's chat threads (summary: thread ID, other party, last message timestamp) |
| GET | `/api/threads/{id}` | Fetch thread + load all messages (with pagination) |
| POST | `/api/threads/{id}/messages` | Send message to thread (creates `Message` entity) |
| GET | `/api/ws` | WebSocket upgrade; fan-out channel for real-time message delivery and admin push events |

**Authentication**: All endpoints except `/api/health` require JWT. WebSocket connection validated via same middleware. Rate limited to 40 req/IP/10s.

---

## 3. gRPC Services

All services communicate via gRPC for inter-service calls. Generated client stubs in `/backend/gen/proto/*/v1/*_grpc.pb.go`.

### AuthService (gRPC)
- **VerifyToken** — Verify JWT signature + return claims (user ID, email, role, expiry, status)
- **GetUser** — Fetch single user by ID
- **GetUsersBatch** — Fetch multiple users by ID list (for denormalization in other services)
- **AdminListUsers** — List users paginated with optional role filter
- **AdminSetStatus** — Update user status (active/suspended)

### CoreService (gRPC)
- **GetVendor** — Fetch vendor by ID
- **GetVendorByUser** — Resolve user ID → vendor ID (one-to-one mapping)
- **ListVendorsByIDs** — Batch fetch vendors
- **GetBooking** — Fetch booking by ID
- **IsBookingAccepted** — Check if booking accepted + return customer/vendor IDs (saga check before payment)
- **MarkBookingPaid** — Update booking status to paid + link payment ID (transactional saga step)
- **AdminStats** — Compute stats: vendor counts (total/pending/approved), booking counts (total/pending/accepted/paid), GMV

### PaymentService (gRPC)
- **Charge** — Process payment: accepts booking ID, user ID, card ID, amount, currency; returns `Payment` entity with status
- **Refund** — Refund payment by ID (optional partial amount)
- **GetPayment** — Fetch payment by ID
- **ListCardsByUser** — List user's cards without PAN

### RealtimeService (gRPC)
- **EnsureThread** — Create thread if not exists (called by core-svc when booking accepted)
- **PublishEvent** — Push event to users (any service may call; fan-out via WebSocket + queued notifications)

---

## 4. Mobile Features (Flutter)

Each feature directory follows MVVM with data/domain/presentation layers:

| Feature | Purpose |
|---------|---------|
| **auth** | User sign-up, login, password reset; role selection at signup |
| **vendor_catalog** | Public vendor search/filter/sort by category/city/price/rating; vendor detail view with photos & reviews |
| **booking** | Create booking request; list bookings (customer + vendor views); status tracking |
| **vendor_self** | Vendor profile editor; service menu manager; photo upload/removal |
| **payment** | Saved card manager; card add/delete/set-default; payment history |
| **cards** | Card list/detail view; payment method selector |
| **reviews** | Post-booking review submission (1–5 rating + text); vendor reviews display |
| **messaging** | Chat threads (thread list + detail view); real-time message send/receive via WebSocket |
| **ai_chat** | AI chat interface (vendor discovery assistant; experimental feature) |
| **notifications** | In-app notification center; FCM device token registration |
| **onboarding** | First-time user flow; role selection; email verification (if required) |
| **settings** | User profile edit; notification preferences; logout |
| **admin** | Admin dashboard (stats, user management, vendor moderation) |

**State Management**: Flutter Riverpod for ViewModels + async data loading.

---

## 5. Frontend Pages (Next.js)

Routes in `/frontend/src/app/` follow Next.js 15 conventions (file-based routing):

| Route | Purpose | Auth Role |
|-------|---------|-----------|
| `/` | AI chat home page (vendor discovery assistant) | customer |
| `/vendors` | Public vendor catalog (search, filter, sort) | customer |
| `/vendors/{id}` | Vendor detail + reviews + booking form | customer |
| `/bookings` | My bookings list (paginated, filterable by status) | customer / vendor |
| `/bookings/{id}` | Booking detail + status timeline + payment form | customer / vendor |
| `/cards` | Saved payment cards manager | customer |
| `/threads` | Chat thread list | customer / vendor |
| `/threads/{id}` | Thread detail + message history + send form | customer / vendor |
| `/notifications` | Notification center (in-app inbox) | customer / vendor |
| `/settings` | User profile editor + logout | customer / vendor |
| `/vendor` | Vendor profile self-editor (services, photos) | vendor |
| `/vendor/bookings` | Vendor's received bookings (accept/decline/complete) | vendor |
| `/admin` | Admin dashboard (stats, user counts, GMV) | admin |
| `/admin/users` | User list + status management | admin |
| `/auth/forgot` | Password reset request form | public |
| `/auth/reset` | Password reset confirmation (token in query) | public |

**Framework**: Next.js 15 (React Server Components + client-side interactivity). Playwight E2E tests in `/frontend/e2e/`.

---

## 6. Data Model & Postgres Schema

All services use Postgres 14+. Each service owns its tables; foreign keys to auth-svc users are comments-only (no DB constraints across service boundaries).

### auth_svc Database
| Table | Key Columns | Indexes | Purpose |
|-------|------------|---------|---------|
| **users** | id (PK), email (UNIQUE), role, status | role, status | User accounts (customer/vendor/admin) |
| **refresh_tokens** | id (PK), user_id (FK), token_hash (UNIQUE), expires_at, revoked_at | user_id | Token rotation & revocation |
| **password_reset_tokens** | id (PK), user_id (FK), token_hash (UNIQUE), expires_at, used_at | user_id | One-time password reset links |

### core_svc Database
| Table | Key Columns | Indexes | Purpose |
|-------|------------|---------|---------|
| **vendors** | id (PK), user_id, name, category, city, status, rating_avg, rating_count | category, city, status, price_from, rating_avg, search_tsv (GIN) | Vendor profiles with moderation |
| **services** | id (PK), vendor_id (FK), name, price, unit, is_active | vendor_id | Vendor menu items |
| **photos** | id (PK), vendor_id (FK), mime, size, data (BYTEA) | vendor_id | Vendor image blobs (max 5 MB each) |
| **bookings** | id (PK), customer_id, vendor_id (FK), service_id, status, amount | customer_id, vendor_id, status | Reservations (payment_id linked after charge) |
| **reviews** | id (PK), booking_id (UNIQUE FK), customer_id, vendor_id (FK), rating (1–5) | vendor_id | Post-booking feedback; feeds vendor rating_avg |
| **notifications** | id (PK), user_id, type, channel, status | user_id | Event log for in-app + delivery |
| **fcm_tokens** | id (PK), user_id, token (UNIQUE), platform | user_id | Device registrations (iOS/Android) |

### payment_svc Database
| Table | Key Columns | Indexes | Purpose |
|-------|------------|---------|---------|
| **cards** | id (PK), user_id, brand, last4, exp_month, exp_year, is_default | user_id | Card vault (PAN not stored) |
| **payments** | id (PK), booking_id (UNIQUE), user_id, card_id, amount, currency, status, provider_ref | user_id, booking_id | Payment ledger + PSP linkage |

### realtime_svc Database
| Table | Key Columns | Indexes | Purpose |
|-------|------------|---------|---------|
| **threads** | id (PK), booking_id (UNIQUE), customer_id, vendor_id | customer_id, vendor_id | DM channels per booking |
| **thread_messages** | id (PK), thread_id (FK), sender_id, text, created_at | thread_id, created_at | Chat history |

---

## 7. Test Coverage Summary

| Component | Count | Lines of Code | Notes |
|-----------|-------|---------------|-------|
| **Auth Service** | 1 test file | 386 | User signup, login, token refresh, password reset flows |
| **Core Service** | 4 test files | 1,007 | Vendor CRUD, booking state machine, review aggregation |
| **Payment Service** | 2 test files | 436 | Card validation, payment charge/refund, PSP integration |
| **Realtime Service** | 1 test file | 204 | Thread ensure, message send, WebSocket connection |
| **Backend Total** | 8 test files | 2,033 | Unit + integration tests; use table-driven patterns |
| **Frontend E2E** | 12 Playwright specs | ~800 lines | auth, booking flow, chat, admin, vendor self-edit, photo upload, settings, sidebar, QA sweep |
| **Dart/Mobile** | Unit tests (in-repo) | TBD | ViewModel tests via Riverpod mocking |

**Testing Approach**: Go backend uses standard `go test` with `-race` flag. Frontend uses Playwright for E2E user journeys. All tests run in CI/CD pipeline.

---

## 8. Architecture Patterns

### Service Communication
- **HTTP** — Client → Gateway or direct to service (public catalog, auth)
- **gRPC** — Service-to-service (auth token verify, core/payment saga, realtime fan-out)
- **WebSocket** — Real-time chat (gateway upgrade, hub broadcasts to connected clients)

### Transactions & Sagas
- **Payment Saga**: Core calls Payment.Charge → on success calls Core.MarkBookingPaid (distributed transaction via gRPC)
- **Booking Flow**: Booking created → if accepted, Core calls Realtime.EnsureThread → thread ready for chat

### Authorization
- **Auth Middleware** — All services verify JWT via Auth.VerifyToken gRPC or (auth-svc) local validation
- **Role-Based Access Control** — `mw.RequireRole("admin")` on admin endpoints; vendor endpoints check user ownership
- **Token Structure** — JWT claims include user_id, email, role, status, expiry (exp)

### Rate Limiting
- Per-IP rate limiter on all HTTP endpoints (20–40 req/10s depending on service)

---

## 9. External Dependencies & Integrations

| Dependency | Purpose | Service |
|-----------|---------|---------|
| **PostgreSQL 14+** | Primary data store (all services) | All |
| **PayBox** | Payment gateway (card tokenization + charge) | Payment |
| **Firebase Cloud Messaging (FCM)** | Push notifications to mobile | Core + Realtime |
| **gorilla/websocket** | WebSocket upgrade library | Realtime |
| **protoc-gen-go-grpc** | gRPC code generation | Build |
| **Postgres migrate** | Schema versioning | All |
| **Flutter / Riverpod** | Mobile UI framework | Mobile |
| **Next.js 15 / React** | Web UI framework | Frontend |
| **Playwright** | E2E testing | Frontend |

---

## 10. Key Files Reference

### Backend Source
- Domain Entities: `backend/services/*/internal/domain/*.go`
- HTTP Handlers: `backend/services/*/internal/adapter/http/router.go` (route definitions)
- gRPC Stubs: `backend/gen/proto/*/v1/*_grpc.pb.go`
- Proto Definitions: `backend/proto/*/v1/*.proto`
- Migrations: `backend/services/*/internal/adapter/repo/migrations/0001_init.up.sql`
- Tests: `backend/services/*/internal/usecase/*_test.go`

### Frontend Source
- Pages: `frontend/src/app/**/page.tsx`
- Shared API Client: `frontend/src/shared/api/index.ts` (REST client)
- Features: `frontend/src/features/*/` (auth, booking, vendor, etc.)
- E2E Tests: `frontend/e2e/*.spec.ts`

### Mobile Source
- Features: `mobile/lib/features/*/` (data/domain/presentation structure)
- ViewModels: `mobile/lib/features/*/presentation/viewmodels/*.dart`
- UI Screens: `mobile/lib/features/*/presentation/screens/*.dart`

---

## Summary

Qonaqzhai is a **multi-service event-booking platform** with clear separation of concerns across auth, vendor/booking management, payments, and real-time chat. The **backend** uses Go microservices with gRPC inter-service communication and Postgres persistence. The **frontend** (Next.js) and **mobile** (Flutter) clients consume HTTP REST APIs behind a gateway. The system enforces RBAC (customer/vendor/admin roles), implements a booking-to-payment saga, and provides real-time chat via WebSocket. Comprehensive test coverage spans unit tests (Go), integration tests (gRPC mocking), and E2E user flows (Playwright).

