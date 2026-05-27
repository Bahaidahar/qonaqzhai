# Stack Justification

Research-backed justification for the technology choices made in this diploma project, grounded in 2025–2026 industry data, benchmarks, and adoption trends.

## 1. Why Go Microservices?

Go was chosen for the backend services because it sits at the intersection of three properties that matter most for a distributed booking platform: predictable low-latency performance, first-class concurrency, and operational simplicity. In 2025 HTTP-server benchmarks, Go-based stacks (e.g. Fiber) processed roughly **4.5M requests in 30 seconds versus ~2M for Node.js/Express**, and Go outperformed Node.js on CPU-bound workloads by an average of **~2.6×** ([Netguru — Golang vs Node 2025](https://www.netguru.com/blog/golang-vs-node), [OceanSoft — Node.js vs Go for Microservices](https://oceansoftsol.com/blog-nodejs-vs-go-microservices.html)). Compared to Java, Go 1.24's Swiss-Tables map delivered up to **60% faster map operations** and Go 1.25's experimental Green Tea GC directly targets the memory-wall problem in containerized deployments ([Java Code Geeks — Go 1.24 vs Java 25 for Microservices 2026](https://www.javacodegeeks.com/2026/05/go-1-24-vs-java-25-for-microservices-an-updated-honest-benchmark-in-2026.html)).

Adoption signals confirm the choice: Uber runs Go across thousands of microservices and Cloudflare saved ~97 cores by applying Profile-Guided Optimization to its Go binaries ([Netguru — 17 Major Companies Using Go](https://www.netguru.com/blog/companies-that-use-golang), [Uber — Domain-Oriented Microservice Architecture](https://www.uber.com/us/en/blog/microservice-architecture/)).

### gRPC + Protocol Buffers

For service-to-service communication we use **gRPC over Protocol Buffers**, while keeping HTTP/JSON for the public gateway. Industry benchmarks in 2025 put gRPC at **5–10× faster than REST**, with Protobuf payloads **70–90% smaller than equivalent JSON** ([Markaicode — gRPC vs REST Benchmarks 2025](https://markaicode.com/grpc-vs-rest-benchmarks-2025/), [Toptal — gRPC vs REST](https://www.toptal.com/developers/grpc/grpc-vs-rest-api)). One reported case: **90 000 RPS over gRPC versus 66 000 RPS over REST** for the same workload — a meaningful margin when auth ↔ core ↔ payment chains call each other thousands of times per booking flow.

### Microservices at diploma scale — an honest answer

Splitting into auth / core / payment / realtime / gateway is admittedly **operational overkill for diploma traffic**. The justification is pedagogical: the project demonstrates *deployable* knowledge of service boundaries, gRPC contracts, per-service databases, and gateway aggregation — competencies expected in modern backend roles. The trade-off accepted is higher Docker Compose complexity and longer cold-start chains in exchange for a portfolio-grade architecture.

## 2. Why Next.js 16 + React 19?

Next.js 16 (October 2025) was selected because it ships the first stable implementation of **Cache Components + Partial Pre-Rendering (PPR)** — a model that combines a static HTML shell with streamed dynamic holes in a single route, eliminating the historical "all static *or* all dynamic" choice ([Next.js 16 release notes](https://nextjs.org/blog/next-16), [Ashish Gogula — Practical PPR Guide in Next.js 16](https://www.ashishgogula.in/blogs/a-practical-guide-to-partial-prerendering-in-next-js-16)). The App Router in 16 runs on **React 19.2** features (View Transitions, useEffectEvent, Activity), and 16.1 (Dec 2025) added stable Turbopack file-system caching for `next dev` ([Next.js blog](https://nextjs.org/blog)).

Against alternatives:
- **Remix 3** dropped React entirely in 2025, which is a high-risk bet for a thesis project ([Merge — Remix vs Next.js 2025](https://merge.rocks/blog/remix-vs-nextjs-2025-comparison)).
- **Nuxt** is the natural Vue companion but we have a React-skilled team.
- **SvelteKit** and **Astro** lead satisfaction in **State of JS 2024/2025**, but Next.js still dominates raw usage at **~60–70% adoption** among meta-framework users ([Strapi — State of JS 2025 Takeaways](https://strapi.io/blog/state-of-javascript-2025-key-takeaways), [Leapcell — 2025 Frontend Showdown](https://leapcell.io/blog/the-2025-frontend-framework-showdown-next-js-nuxt-js-sveltekit-and-astro)).

We layered **Feature-Sliced Design (FSD)** on top to enforce strict horizontal isolation between features, and **Playwright** for E2E — the de facto modern choice over Cypress for parallel cross-browser runs.

## 3. Why Flutter Over React Native / Native?

Flutter wins on three measurable dimensions in 2025:

- **Performance**: Flutter's Impeller engine renders at **60–120 fps** with up to **30% better frame rendering on modern hardware** than the legacy Skia backend, and consistently beats React Native on heavy-animation and nested-list benchmarks ([Synergy Boat — Flutter vs RN vs Native Benchmark 2025](https://www.synergyboat.com/blog/flutter-vs-react-native-vs-native-performance-benchmark-2025), [Nomtek — Flutter vs React Native 2025](https://www.nomtek.com/blog/flutter-vs-react-native)).
- **Adoption**: Stack Overflow 2024 ranked Flutter as the **most-used cross-platform framework at 46% vs RN's 35%**, and Statista's 2025 data confirms the same gap ([Droids on Roids — Flutter vs RN 2025](https://www.thedroidsonroids.com/blog/flutter-vs-react-native-comparison)). Stack Overflow 2025 keeps Flutter ahead overall (9.12% vs 8.43%).
- **Single codebase**: One Dart codebase ships to iOS and Android with native-compiled binaries — critical for a solo diploma builder.

Architecturally, **Riverpod** is now the canonical state-management + DI solution (no `BuildContext` dependency, compile-safe), and we apply **MVVM + Clean Architecture** (data / domain / presentation) to keep the app testable. `go_router` handles declarative navigation, `Dio` handles HTTP with interceptors for JWT refresh, and **Firebase Cloud Messaging** delivers push notifications cross-platform.

## 4. Why Per-Service PostgreSQL + gRPC?

We implement the canonical **Database-per-Service** pattern as defined by Chris Richardson: *each microservice owns its schema exclusively, and no other service touches it directly* ([microservices.io — Database per Service](https://microservices.io/patterns/data/database-per-service.html)). This means **no cross-service foreign keys**, no shared `JOIN`s between `auth.users` and `core.bookings` — services communicate only over gRPC contracts.

The accepted cost is **eventual consistency**: a booking and its payment record will not commit atomically. We mitigate this through the **Saga pattern** with compensating actions (cancel-booking on payment-failure) and event publication, again following Richardson's playbook ([microservices.io — Event-Driven Architecture](https://microservices.io/patterns/data/event-driven-architecture.html)). The CAP theorem makes distributed ACID transactions impractical anyway — eventual consistency is the honest, scalable answer.

PostgreSQL is the per-service engine of choice: open-source, JSONB-capable, with mature replication and Docker tooling. It earned the largest "admired + desired" gap in the Stack Overflow 2025 survey, and Docker itself jumped **+17 points** in adoption — the biggest single-year move of any tool tracked ([Stack Overflow 2025 Survey](https://survey.stackoverflow.co/2025/technology)).

## 5. Why MCP + Claude API in 2026?

The Anthropic **Claude API** powers our AI event-planning chat because Claude 4.5/Opus leads coding and reasoning benchmarks, and **Claude Code** itself grew **10× in three months** post-launch with enterprise subscriptions quadrupling in 2025 ([CIO — How Agentic AI Will Reshape Engineering 2026](https://www.cio.com/article/4134741/how-agentic-ai-will-reshape-engineering-workflows-in-2026.html)).

During development we lean on **Model Context Protocol (MCP)** servers — the open standard Anthropic introduced in Nov 2024, **adopted by OpenAI in March 2025**, by Google DeepMind, and previewed in Windows 11 at Microsoft Build 2025 ([The New Stack — Why MCP Won](https://thenewstack.io/why-the-model-context-protocol-won/), [Anthropic — Donating MCP to the Linux Foundation](https://www.anthropic.com/news/donating-the-model-context-protocol-and-establishing-of-the-agentic-ai-foundation)). By December 2025 Anthropic reported **97M+ monthly SDK downloads** and **10 000+ active MCP servers** in production. Anthropic donated MCP governance to the **Agentic AI Foundation under the Linux Foundation**, cementing it as the cross-vendor agent integration standard.

Choosing MCP-based agentic workflows aligns the diploma with **Stack Overflow 2025's finding that 84% of developers now use AI tools**, and IDC's projection of a **$36.2B AI developer-tool market by end of 2026** — this is not an experimental bet, it is the mainstream 2026 developer workflow.

## Comparative Stack Table

| Layer / Choice | Main Alternative | Why We Picked Ours | Trade-off Accepted |
|---|---|---|---|
| **Go microservices** | Node.js / Java Spring | ~2.6× faster CPU-bound; lower memory; static binaries; mature gRPC tooling | Smaller talent pool than Node; less ORM richness than JPA |
| **gRPC + Protobuf (internal)** | REST/JSON everywhere | 5–10× throughput, 70–90% smaller payloads, typed contracts via `.proto` | Harder to debug with curl; not browser-native (we keep REST at the gateway) |
| **Next.js 16 + React 19** | Remix 3 / SvelteKit / Nuxt | Stable PPR + Cache Components, biggest ecosystem, React 19.2 features | Higher learning curve than SvelteKit; satisfaction has slipped vs Astro |
| **Feature-Sliced Design** | Atomic Design / classic `src/components` | Hard horizontal isolation between features; clear dependency rules | More boilerplate for small features |
| **Flutter 3.24+ (Impeller)** | React Native (Fabric) / Native iOS+Android | 46% vs 35% cross-platform share; 60–120 fps Impeller; single Dart codebase | Larger APK size; Dart is niche outside Flutter |
| **Riverpod + MVVM + Clean Arch** | Provider / Bloc | Compile-safe DI, no `BuildContext`, testable layers | Steeper learning curve than Provider |
| **Per-service PostgreSQL** | Shared monolithic DB | Service autonomy, independent scaling, no cross-team schema coupling | No cross-service `JOIN`s; eventual consistency via Sagas |
| **JWT auth** | Session cookies + Redis | Stateless, scales horizontally across services without sticky sessions | Revocation is harder; requires careful refresh-token rotation |
| **PayBox gateway** | Stripe / direct card processing | PCI compliance offloaded; supports local KZ payment methods | Vendor lock-in to PayBox SDK and webhook format |
| **Claude API + MCP** | OpenAI GPT-4 / self-hosted LLM | Best-in-class reasoning; MCP is the cross-vendor agent standard (Linux Foundation, 97M+ SDK downloads) | API cost vs self-hosted; rate limits on free tier |
| **Playwright E2E** | Cypress / Selenium | Parallel cross-browser; better CI integration; first-class TypeScript | Slower startup than Cypress for tiny suites |
| **Docker Compose deploy** | Kubernetes / bare metal | Simple, reproducible, fits diploma scope; +17pt Docker adoption in 2025 | Not production-grade for >1 host; manual scaling |

---

### Summary

Every layer of this stack is anchored in a **measured 2025–2026 industry trend**: Go for performance-critical microservices, gRPC for internal RPC, Next.js 16 for the modern React frontier, Flutter for cross-platform mobile leadership, per-service Postgres for proper service autonomy, and MCP + Claude for an agentic AI workflow that reflects how engineering will actually be done in 2026.