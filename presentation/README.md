# Qonaqzhai — Diploma Defense Package

All artifacts for the diploma defense, generated 2026-05-28.

## Files

```
presentation/
├── README.md                      this file
├── diploma.md                     22-slide deck (English)
├── diploma.html                   HTML render of diploma.md
├── diploma.pdf                    print-ready PDF (22 slides ≈ 25 pp)
├── style.css                      print styling (A4, 18mm margins)
└── research/
    ├── competitors.md             comparative analysis (~1,050 words)
    ├── stack-justification.md     stack rationale w/ 2025–26 sources (~1,400 words)
    ├── code-map.md                full architecture map (~7,200 words)
    └── features-roadmap.md        feature ideas + demo flow
```

## Regenerate PDF

```bash
cd presentation

# 1. MD → HTML
pandoc diploma.md -o diploma.html \
  --standalone \
  --metadata title="Qonaqzhai — Diploma Defense" \
  --toc --toc-depth=1 \
  -c style.css

# 2. HTML → PDF (Chrome headless, no LaTeX needed)
"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" \
  --headless --disable-gpu --no-pdf-header-footer \
  --print-to-pdf=diploma.pdf \
  "file://$PWD/diploma.html"
```

Alternative if LaTeX is installed (sharper typography):

```bash
brew install --cask basictex
pandoc diploma.md -o diploma.pdf --pdf-engine=xelatex -V geometry:margin=2cm
```

## Slide outline (22)

| # | Title |
|---|-------|
| 1 | Title |
| 2 | Problem Statement |
| 3 | Solution at a Glance |
| 4 | Domain Model |
| 5 | Stack Overview |
| 6 | Backend Architecture |
| 7 | Service-to-Service Communication |
| 8 | Clean Architecture per Service |
| 9 | Web Frontend (FSD) |
| 10 | Mobile Architecture |
| 11 | Data Model |
| 12 | Real-time Chat (WebSocket) |
| 13 | Stack Justification (Comparative) |
| 14 | Comparative Analysis (Competitors) |
| 15 | AI Event Planner (Signature Feature) |
| 16 | Security Model |
| 17 | Testing & Quality |
| 18 | Agentic Engineering & MCP |
| 19 | Live Demo Script (8 min) |
| 20 | Roadmap |
| 21 | Lessons Learned |
| 22 | Q&A |
+ Appendix A (numbers), Appendix B (repo layout)

## What teacher asked for — coverage

| Ask | Where it lives |
|---|---|
| Comparative analysis | Slide 14 + `research/competitors.md` |
| Stack explanation | Slide 5 + Slide 13 + `research/stack-justification.md` |
| Architecture | Slides 4, 6–12 + `research/code-map.md` |
| Research part (web-sourced) | Slide 13 (inline source links), `competitors.md`, `stack-justification.md` |
| Stack | Slide 5 |
| Architecture diagrams | Slides 3, 4, 6, 7, 8, 9, 10, 12, 15 |
| MCP / skills | Slide 18 |
| Web-search-backed | All research files cite 2025–2026 sources |
| Presentation | `diploma.md` → `diploma.pdf` |
| Feature ideas | Slide 20 + `research/features-roadmap.md` |

## Customize before printing

- Set author institution on the YAML front-matter of `diploma.md`.
- Drop logos or screenshots into `presentation/assets/` and reference them in slide bodies.
- Adjust colors in `style.css` (`#0a3d62` is the primary).

## Speaker time budget

22 slides @ ~45 s/slide ≈ 16 min. Demo (Slide 19) replaces walking through 14–17 if time short.
