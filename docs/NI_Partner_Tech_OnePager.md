# UploadParty × Native Instruments — Technical One‑Pager (For Partner Review)

Date: September 17, 2025  
Prepared by: Dhamari  
Audience: Jack Tarricone (Native Instruments) + Partner Team  
Document Type: Partner one‑pager (Google‑Doc style headings)

---

## 1) Executive Summary
- UploadParty is an online, gamified challenge/community for teen music producers (target: 500 participants) with active collaboration from Native Instruments.  
- Educational mission: “learning by doing” — teens upload beats, track time spent creating, participate in challenges, and earn recognition/prizes. No marketplace or buying/selling.  
- Launch target: September 26, 2025. Immediate deliverable is an RSVP landing page and basic onboarding; full platform features roll out post‑launch.  
- Request: Short extension for the RSVP website delivery to Monday (Sept 22, 2025) to reflect our now‑finalized roadmap and ensure quality.

## 2) Timeline, Commitments, and Extension Request
- Partnership status: Confirmed; NI to provide 500 Komplete 15 license codes (request in progress).  
- Launch date: Friday, September 26, 2025 — aligned on both sides.  
- Website for review: Initially promised EOD Friday; requesting extension to EOD Monday, September 22, 2025. Rationale: When we began, the roadmap was not fully defined. We now have a clear plan and want to ensure the RSVP experience (copy, branding, and data capture) is polished for NI review and approval.  
- Co‑branding & approvals: All public materials (site, copy, graphics) will be held for NI approval before publication. We will provide a preview link and asset pack.

## 3) What Ships First (M0): RSVP Launchpad
- Goal: A simple, fast site for teens to RSVP and get redirected to join the community that launches on September 26.  
- Key elements:  
  - Hero section with UploadParty × Native Instruments brand presence  
  - Clear program value and date (Sept 26)  
  - RSVP form (name, email, age confirmation, region/time zone)  
  - Consent language and privacy notice  
  - Post‑RSVP redirect to community join flow (email confirmation + link)  
  - Basic analytics (visit, form submit)  
- Review package for NI: Staging URL, copy deck, hero image options, and tracking plan for sign‑offs.

## 4) Technical Architecture (At a Glance)
- Frontend: Next.js 15 (App Router), React 19, Tailwind CSS 4.  
  - Edge‑friendly static/server rendering for speed and reliability.  
  - Accessibility and mobile readiness prioritized.  
- Backend: Go (Gin), PostgreSQL (GORM), Redis, JWT auth.  
  - API base path /api/v1 with security middlewares (CORS, size limits, timeouts, rate limits).  
  - Stateless services suitable for horizontal scaling.  
- Storage & Media: AWS S3 for audio (multipart uploads, lifecycle policies).  
- Infra & Security: Docker Compose for local parity; optional Nginx reverse proxy; non‑root containers, read‑only filesystems; strict CORS, JWT, and content‑type/size validation.  
- Capacity target: Designed as a platform application; initial baseline ~100 DAU storage with headroom to scale.  
- Observability: Health endpoint, structured logging; adding basic telemetry during beta.

## 5) Roadmap & Milestones (Condensed)
- M0: RSVP Launchpad — ship RSVP site, redirect to community join flow (launch day: Sept 26).  
- M1: Platform foundation — auth, profiles, DB/cache wiring, CI, security/CORS baselines.  
- M2: Uploads + S3 — large‑file uploads (MP3/WAV), size/MIME validation, durable storage, playback path.  
- M3: Social — likes/comments, basic feed, community interactions.  
- M4: Challenges — declare challenges, submit entries, scoring, winners; badges/leaderboard (hours spent).  
- M5: Hardening + Beta — rate limits, timeouts, telemetry, docs; scale testing.  
- M6: Hackatime extension — integrate https://hackatime.hackclub.com for time‑tracked activity linked to VST workflows.  

## 6) NI Co‑Branding, Approvals, and Integration
- Role: Sponsor-only. NI provides prizes/licenses and brand support; any use of NI brand names, logos, or assets requires prior written approval.
- Approval workflow:
  1) Share staging link, copy, and assets by Monday EOD (Sept 22).  
  2) Collect NI feedback within 24–48 hours; turn around revisions quickly.  
  3) Lock content no later than Thursday for Friday launch (Sept 26).  
- Product integration: NI brand presence on RSVP and throughout platform; Komplete Start participation path defined in copy; codes distributed via approved mechanism post‑challenge completion (to be finalized together).  
- Cross‑promotion: Provide partner asset kit (images/copy) for social amplification across NI and UploadParty channels.

## 7) Risks & Mitigations (Launch‑adjacent)
- Copy/brand approvals run late → Mitigation: Reserve design/dev bandwidth for same‑day revisions; share early previews.  
- High traffic spike on launch day → Mitigation: CDN/static pre‑render, API rate limiting, S3 offload, basic autoscale plan.  
- Email deliverability for RSVP confirmations → Mitigation: Use reputable provider, SPF/DKIM, test seed accounts.  
- Code distribution/eligibility verification → Mitigation: Gate on challenge completion via verifiable submission workflow; design together with NI.

## 8) What We Need From NI This Week
- Review and feedback on the RSVP copy, hero image, and sign‑up fields.  
- Confirmation of brand usage rules (logos, tone, co‑marks).  
- Guidance on voucher/code distribution mechanism and messaging.  
- Any legal/privacy language NI requires on the RSVP page.

## 9) Why the Short Extension Helps
- We began before the roadmap was fully clear; we now have a concrete, agreed technical plan and milestone sequence.  
- The additional 1–2 business days ensures the RSVP experience meets quality expectations for brand, clarity, and stability — leading to a smoother launch week and fewer last‑minute changes.

---

## 11) Engagement Plan & Permission Request
- To maintain momentum between RSVP day and platform launch, we would like to proactively reach out to a small set of teen‑friendly producers/mentors to share prompts, sample packs, and short videos. This keeps participants engaged and preparing for launch week content.
- Request: Is it okay if I (Dhamari) begin coordinating outreach to a few producers this week, with NI looped in for visibility/approval on names and materials before anything is published?
- Safeguards: Light‑touch content, no co‑marketing without NI approval, alignment with UploadParty educational mission and brand guidelines.

## 12) Updated Timeline Notes
- RSVP Launch (M0): Friday, September 26, 2025 — RSVP page live and community opening.
- Platform Work Start: Begins immediately upon the RSVP Launchpad going live (M0) to maintain momentum and ensure delivery quality.
- Platform Launch (M1): Fully complete before Friday, October 10, 2025 (two weeks after Sept 26). Core platform functionality goes live (auth/profiles foundation) with subsequent feature increments following the roadmap.

Appendix (Optional)  
- Health endpoint (dev): http://localhost:8080/health  
- API base: /api/v1  
- Frontend (dev): http://localhost:3000