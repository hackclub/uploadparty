# UploadParty — One‑Pager

Date: 2025-09-17
Version: 0.1.0

What is UploadParty?
- A gamified, web-based music creation community where teens upload beats, get feedback, and compete in challenges.
- Backed by Native Instruments (Komplete licenses) and Hack Club support.
- Platform application: designed as a scalable platform (web app + APIs) rather than a one-off site.
- Capacity baseline: we expect to hold audio files for approximately 100 daily active users (DAU) initially, with room to scale.

Problem
- Teens want a fun, safe place to create and share music but lack engaging, beginner‑friendly platforms with community and incentives.

Solution
- A simple upload → share → compete loop with badges, leaderboards (showcasing hours spent on projects/songs), and challenges to keep creators motivated.
- Frictionless large‑file uploads, fast playback, and social actions (likes, comments, follows, collaborations).
- No buying/selling or licensing transactions — the focus is uploading beats, tracking time, building community, and earning prizes.

Who it’s for
- High‑school creators (musicians, producers) starting or leveling up their beat‑making journey.
- Educators and community organizers running music/tech clubs.

Key Features (MLP)
- Secure account creation and profiles
- Beat uploads (MP3/WAV), metadata, and playback
- Likes, comments, and basic feed
- Time‑boxed challenges with submissions and winners
- Badges and scoring to celebrate progress

Value Proposition
- Learn by doing: low friction uploads + instant feedback
- Stay engaged: gamification and challenges
- Grow reputation: badges, leaderboards that showcase hours spent on projects/songs.

How it works (Tech at a glance)
- Frontend: Next.js 15 (App Router), React 19, Tailwind CSS 4
- Backend: Go (Gin), GORM, PostgreSQL, Redis, JWT auth
- Storage: AWS S3 for audio; optional Nginx reverse proxy
- DevOps: Docker Compose for local parity; health checks and rate limits

Milestones
- M0 RSVP Launchpad (Launch: Friday, September 26, 2025): Build a simple website that allows people to RSVP and redirects them to join the community launching on September 26.
- M1 Platform Foundation : Work begins immediately upon the RSVP Launchpad going live (M0) and is planned to be entirely completed before Oct 10. Establish core platform (auth/profiles, backend/frontend wiring, security/CORS) and prepare for subsequent feature increments.
- M2 Uploads: S3 integration, size/MIME validation
- M3 Social: Likes/comments, basic feed
- M4 Challenges: Create/list/submit, scoring, winners
- M5 Beta Hardening: Rate limits, timeouts, telemetry, docs (Fully complete before: Friday, October 10, 2025)
- M6 Hackatime Extension: Build an extension/integration with https://hackatime.hackclub.com that links time-tracked activity to the user’s VST workflow (e.g., deep-link from activity to open the relevant VST/project).