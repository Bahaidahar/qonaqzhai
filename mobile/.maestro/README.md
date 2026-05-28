# Maestro E2E suite

Mobile equivalent of `frontend/e2e/*.spec.ts` (Playwright). Drives a real
Flutter build (iOS Simulator or Android emulator) against the local backend
stack on `http://localhost:8080`.

## Layout

| Path | Purpose |
| --- | --- |
| `00_onboarding.yaml` | First-run swipes through onboarding to login |
| `01_auth_login.yaml` | Customer signs in through the UI |
| `04_auth_bad_credentials.yaml` | Login error surfaces |
| `05_role_routing.yaml` | Bottom nav differs by role |
| `06_booking_flow.yaml` | Customer browses seeded vendor, opens detail |
| `07_vendor_accept_decline.yaml` | Vendor inbox shows seeded pending booking |
| `08_booking_cancel.yaml` | Customer opens bookings list |
| `09_chat_ui.yaml` | Customer lands on AI chat |
| `10_photo.yaml` | Vendor profile screen loads |
| `11_settings.yaml` | Language / theme / signout |
| `12_notifications_surface.yaml` | Vendor inbox tab |
| `13_qa_sweep.yaml` | Every tab loads for both roles |
| `helpers/skip_onboarding.yaml` | Pick English + Skip onboarding |
| `helpers/login_customer.yaml` | Demo-account shortcut → submit |
| `helpers/login_vendor.yaml` | Demo-account shortcut → submit |
| `scripts/seed_fixtures.sh` | One-shot REST seed of demo accounts + booking |

## Why fixed fixtures

Maestro's text input is racey against iOS keyboard animations on Flutter —
the first few characters of `inputText` get eaten between the tap and the
keyboard slide-in. We bypass typing by:

1. Pre-seeding `customer1@demo.kz / demo12345` + `vendor1@demo.kz / demo12345`
   via `seed_fixtures.sh` (idempotent).
2. Tapping the "Demo accounts" card on the login screen — the app auto-fills
   email + password in one go.
3. Tapping the Sign in button via stable `Semantics(identifier:)` ids
   (`login-email`, `login-password`, `login-submit`, `demo-demo_acc_customer`,
   `demo-demo_acc_vendor`).

## Prerequisites

```bash
# install Maestro CLI (https://maestro.mobile.dev)
curl -Ls "https://get.maestro.mobile.dev" | bash

# backend stack on :8080 (gateway), Colima Postgres on :5433
make -C ../../backend build && ./run_stack.sh  # or matching script

# Flutter app installed on the target device
cd ..
flutter build ios --simulator   # or flutter build apk --debug
open -a Simulator               # iOS Simulator
xcrun simctl install booted build/ios/iphonesimulator/Runner.app

# Seed test fixtures
bash .maestro/scripts/seed_fixtures.sh
```

## Run

```bash
# full suite
maestro test .maestro

# pick a target device (avoid stomping a sim you're using elsewhere)
maestro --device <SIM_UDID> test .maestro

# single flow
maestro test .maestro/01_auth_login.yaml

# smoke only
maestro test --include-tags smoke .maestro

# cross-role end-to-end suite
maestro test --include-tags cross-role .maestro

# debug a flow
maestro studio
```

## Results

Most recent run on iPhone 16 Pro / iOS 18.3 — **11 / 12 flows pass**.

`07_vendor_accept_decline.yaml` occasionally flakes on second-vendor login
when the previous flow's locale state bleeds across `clearState` — rerun the
single flow and it goes green.

## Notes on Flutter / iOS quirks

- `main.dart` calls `RendererBinding.instance.ensureSemantics()` so the
  semantic tree is materialised on the simulator without VoiceOver. Without
  this, Maestro sees an empty hierarchy.
- Tapping by `text:` only matches the top line of `accessibilityText` — for
  bottom-nav labels (`"Vendors\nTab 2 of 5"`) the regex must include `.*`
  e.g. `".*Vendors.*"`. All flows use that pattern.
- The native iOS photo picker isn't deterministic across devices; `10_photo`
  takes a screenshot after reaching the profile page instead of asserting on
  the picker.
