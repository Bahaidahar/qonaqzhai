# Qonaqzhai Mobile

Flutter client for the Qonaqzhai event services marketplace.

## Architecture

Simplified Clean Architecture, feature-first, **Riverpod + MVVM** (no BLoC).

```
lib/
├── core/          # network, theme, router, di
├── features/
│   ├── auth/
│   │   ├── data/         # DTOs, datasources, repository impl
│   │   ├── domain/       # entities + abstract repo + use cases
│   │   └── presentation/ # screens + ViewModels (Riverpod)
│   ├── vendor_catalog/
│   ├── booking/
│   ├── ai_chat/
│   ├── payment/
│   ├── reviews/
│   └── notifications/
└── main.dart
```

## Layer rules

- `domain` is pure Dart — no Flutter, no JSON, no Dio.
- `data` implements the `domain` repositories using Dio + DTOs.
- `presentation` exposes Riverpod `StateNotifier`s (ViewModels) consumed by `ConsumerWidget` screens.

## Setup

```
flutter pub get
flutter run --dart-define=API_BASE_URL=http://localhost:8080
```

The default base URL is `http://localhost:8080`. Override via `--dart-define=API_BASE_URL=https://qonaqzhai.kz`.

## Dependencies

- `flutter_riverpod` — state management + DI
- `dio` — HTTP with auth + refresh interceptor
- `go_router` — navigation
- `firebase_messaging` — push (token registered with `/api/notifications/tokens`)
- `flutter_secure_storage` — access / refresh tokens

## Next steps

- Generate Dart client from `backend/docs/openapi.yaml` via `openapi-generator-cli` to eliminate hand-rolled DTOs.
- Wire `firebase_messaging` token-on-login flow.
- Add review submission + vendor detail screens.
