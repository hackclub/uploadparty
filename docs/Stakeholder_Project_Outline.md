# Project Outline for Stakeholders

Title: UploadParty Project Outline
Date: 2025-09-17
Version: 0.1.0
Prepared by: Dhamari (PM: Dhamari)

--------------------------------------------------------------------------------

1. Introduction and Project Overview
   1.1 Project Goals and Objectives (The "Why")
   - Deliver a music-focused web platform that enables creators to upload, share, and collaborate on beats while engaging through challenges and community features.
   - Provide a streamlined, secure experience for uploading large audio files and managing content with scalable storage and a responsive UI.
   - Increase creator engagement and retention via gamified elements (badges, challenges, leaderboards used to showcase hours spent on projects/songs) and social interactions (likes, comments, collaborations).
   - Intergrate features for partner integrations (e.g., Native Instruments), with prize/recognition incentives.

   1.2 Project Deliverables (The "What")
   - Backend API (Go/Gin) with JWT auth, media upload endpoints, and core domain services (users, beats, comments, challenges, scoring).
   - PostgreSQL schema and GORM-based models with auto-migrations.
   - Redis-backed caching/session support (as applicable).
   - Frontend Next.js app (App Router) for marketing pages and authenticated app (dashboard, upload flow, challenge participation).
   - S3 integration for storage of audio files; optional Nginx reverse proxy for production.
   - Infrastructure-as-code via Docker Compose for local parity.
   - Documentation: Developer Guidelines (.junie/guidelines.md), this Stakeholder Project Outline, API reference (basic), environment setup instructions.
   - Operational runbooks: health checks, basic monitoring hooks, and deployment steps (initial).

   1.3 Project Scope and Context
   - In scope: User auth, profile, beat uploads, basic content browsing, likes/comments, challenges and submissions, scoring, badges, and minimal moderation tooling.
   - Program focus: Uploading beats, tracking time spent creating, community participation, and earning prizes/recognition — not facilitating sales or licensing transactions.
   - Out of scope (initial release): Any buying/selling or payments, full e-commerce marketplace, advanced licensing workflows, mobile native apps, full-featured analytics dashboards, content fingerprinting, and advanced moderation/AI content filters.
   - Context: Fits within a creator-tools portfolio, integrates with external providers (S3; optionally Native Instruments). Must interoperate with standard web clients and optionally sit behind Nginx in production. CORS configured to a single frontend origin.
   - Platform application: This is a platform application (web app + APIs) intended to serve a growing community, not a one-off site.
   - Capacity baseline: We expect to hold audio files for approximately 100 daily active users (DAU) initially, with headroom to scale as adoption grows.

2. Project Approach and Strategy
   2.1 Chosen Process Model
   - Approach: Agile (Scrum-ish) with two-week iterations and continuous integration. Rationale: Rapidly evolving product needs, frequent feedback cycles, and a small cross-functional team.
   - Fit: Aligns with customer preference for iterative visibility and the team’s experience with rapid web delivery and DevOps practices.

   2.2 Key Principles and Practices
   - Mission-first: Inspire and enable as many teens as possible to get into music production; lower barriers to entry through accessible tools, education, and community.

3. Project Organization and Responsibilities (The "Who" and "Where")
   3.1 Organizational Structure
   - Executive Sponsor → Product Manager → Engineering Lead → Backend/Frontend Engineers → QA/UAT → DevOps.
   - Supporting roles: Designer, Security Advisor, and Data/Analytics as-needed.

   3.2 Key Personnel and Responsibilities
   - Current staffing model: Single-owner/operator. Dhamari is presently responsible for all roles listed below (PM, Engineering Lead, Backend, Frontend, QA/UAT, and DevOps).
   - Native Instruments support: Sponsor-only role. NI sponsors the program with products/prizes/licenses and is not responsible for operations. Any use of NI brand names, logos, or assets requires prior written approval. NI will review/approve all public website changes that reference NI prior to publication.
   - Product Manager: Owns roadmap and prioritization; interfaces with stakeholders; defines success metrics; ensures discovery and usability.
   - Engineering Lead: Owns architecture and technical quality; coordinates releases; ensures security and performance baselines.
   - Backend Engineer(s): Implement Go/Gin services, DB schema, and integrations; maintain API contracts and performance.
   - Frontend Engineer(s): Implement Next.js UI, accessibility, state management, and integration with API.
   - QA/UAT: Define acceptance criteria; execute test plans; validate releases.
   - DevOps: Own CI/CD, infrastructure (Docker, Nginx), observability, and incident response runbooks.

   3.3 Communication and Interaction
   - Cadence: Weekly stakeholder status (written), bi-weekly sprint review/demo (live), daily internal standup, and ad-hoc design/tech sessions.
   - Channels: Slack/Teams for async, project management tool (Jira/Linear) for tracking, documentation in repo (docs/ and .junie/).
   - Status format: Risks, blockers, velocity, burndown, release health, and upcoming milestones.

4. Project Plan Summary (The "When" and "How Much")
   4.1 High-Level Schedule and Milestones
   - M0: RSVP Launchpad – Build a simple website that allows people to RSVP and redirects them to join the community launching on September 26. (Launch date: Friday, September 26, 2025)
   - M1: Platform application foundation – Work begins immediately upon the RSVP Launchpad going live (M0) and is planned to be entirely completed before the target public launch date. Stand up the core platform across backend and frontend: initialize Go/Gin API and Next.js app, implement JWT auth and basic profile, wire Postgres/Redis, set up Docker Compose/CI, and ensure scalable architecture and CORS/security baselines. (Fully complete before: Friday, October 10, 2025)
   - M2: Uploads + S3 integration – file size limits, MIME validation, and happy-path upload/download.
   - M3: Social features – likes/comments; basic feed.
   - M4: Challenges – create/list/submit; scoring; winners.
   - M5: Hardening + Beta – rate limits, timeouts, telemetry, documentation.
   - M6: Hackatime extension – build https://hackatime.hackclub.com integration that links time-tracked activity to the user's VST workflow (e.g., deep links to the relevant VST or associated project), ship MVP browser extension.
   - Tracking: Roadmap in PM tool; Gantt-lite or milestone board; burndown reports; earned value tracking for larger contracts.