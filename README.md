# Mindcare — Project NLN

**NLN** is a Duolingo-style experience for mental health education and check-ins, with copy and framing tailored for **Nepal**. Users answer a short questionnaire (~2 minutes), see a **non-diagnostic** “where you stand” view, and access **learning paths** that match their support zone.

> **This is a self-check for emotional wellbeing, not a diagnosis.** Screening should always be paired with appropriate follow-up and referral paths where risk is elevated.

---

## What the product does (MVP goals)

- **Quick assessment** — ~90 seconds to 2 minutes; 12 core questions plus a **conditional** safety item.
- **Support zones** — Plain-language bands (green / yellow / orange / red) based on a simple risk model, not clinical labels.
- **Safety first** — If certain answers indicate elevated concern, the app **must** route users to **human support** and crisis-appropriate content—not gamified celebration.
- **Learn in context** — Educational modules grounded in validated constructs (e.g. PHQ-style mood/anxiety items, GAD-2–style worry, WHO-5–inspired wellbeing, functioning), rewritten in **simpler English** and **Nepal-friendly** language (e.g. “overthinking” where it helps honesty).

External references used in the product design for **green-zone** public-health content include [WHO — Doing What Matters in Times of Stress](https://www.who.int/publications/i/item/9789240003927), [CDC — How Right Now (stress)](https://www.cdc.gov/howrightnow/emotion/stress/index.html), and [CDC — Sleep and health](https://www.cdc.gov/sleep/about/index.html).

---

## Assessment structure (spec)

| Section | Focus | Items (example) |
|--------|--------|-------------------|
| **Mood & anxiety** | Low mood, anxiety, worry (0–3 frequency scale) | Q1–Q4 |
| **Wellbeing** | WHO-5–style positives (0–5 scale) | Q5–Q8 |
| **Daily function** | Impact on work, focus, sleep (0–4) | Q9–Q11 |
| **Support** | “Someone I can talk to honestly” (Likert) | Q12 |
| **Safety (conditional)** | Heaviness / not worth continuing | Q13 |

**When to show Q13:** Recommended trigger when mood/anxiety subtotal ≥ **7** *or* function subtotal ≥ **8**. Wording should allow **skip**; answers like **Sometimes** / **Often** trigger an **immediate safety / human support** flow (overrides score-based zones).

**Risk model (product scoring, not a clinical score):**

- Mood/anxiety: sum Q1–Q4 (0–12).
- Wellbeing risk: **20 −** sum Q5–Q8 (wellbeing 0–20 inverted).
- Function: sum Q9–Q11 (0–12).
- Support risk: **4 −** Q12 (0–4).

**Total risk** = mood/anxiety + wellbeing risk + function + support (max **48**).

**Zones (support-oriented copy, not diagnoses):**

| Zone | Score | Intent |
|------|--------|--------|
| Green | 0–12 | Doing okay — prevention, habits, basic psychoeducation |
| Yellow | 13–22 | Under pressure — stress, overthinking, sleep, regulation |
| Orange | 23–34 | Struggling — deeper self-help, journaling, strong nudge to talk to someone |
| Red | 35–48 | Needs human support — resources, helplines, no “confetti” tone |

**Hard rule:** Any **Q13** response indicating **Sometimes** or **Often** → **safety support screen** regardless of total score.

---

## Learning paths by zone (spec)

Each zone gets **tracks** (e.g. emotional fitness, stress prevention, self-awareness, healthy habits) with short lessons, micro-activities, quizzes, and light gamification—**toned down** for orange/red and **disabled where inappropriate** after safety escalation.

Safety override flow emphasizes Nepal-relevant and general helplines, trusted contacts, and optional grounding-only steps—not heavy reflection.

---

## Repository layout

```
Mindcare/
├── frontend/          # React 19 + TypeScript + Vite SPA
├── backend/           # Go module: DB helpers, user model (API server stub)
│   ├── internal/
│   └── migrations/   # PostgreSQL (golang-migrate)
└── Makefile          # Database migration targets
```

**Current implementation status**

- **Frontend:** Login screen (demo **email-only** auth stored in `sessionStorage`), protected dashboard placeholder, React Router. Branding in the UI may still say **G7** in places; product direction is **NLN** as above.
- **Backend:** `internal/db` (Postgres via `lib/pq`), `internal/user` (save/get by email with nickname, age, `daily_ember`, streak, avatar). `internal/cmd/main.go` is a **stub** (no HTTP server yet).
- **Database:** Migrations define `users` and gamification-related columns; see `backend/migrations/` and `backend/migrations/README.md`.

---

## Prerequisites

- **Node.js** (for the frontend).
- **Go** (for the backend module and tooling).
- **PostgreSQL** (for schema/migrations when you wire the API).

---

## Frontend

From `frontend/`:

```bash
npm install
npm run dev
```

Build and preview:

```bash
npm run build
npm run preview
```

---

## Database migrations

Install the [golang-migrate](https://github.com/golang-migrate/migrate) CLI, create the database (default name in the Makefile is `g7`), set `DATABASE_URL`, then from the **repo root**:

```bash
make migrate-up
```

Details, rollbacks, and creating new migrations: **`backend/migrations/README.md`**.

---

## Contributing & ethics

- Treat all user-facing language as **supportive and non-judgmental**; avoid implying diagnosis.
- **Never** ship screening without vetted escalation and resource content for your deployment region (e.g. Nepal helplines, institutional counseling, crisis lines).
- Align UX with medical/legal guidance for your jurisdiction; this README is **not** clinical or legal advice.

---

## License

See repository license if present; otherwise add one when you open-source or distribute.
