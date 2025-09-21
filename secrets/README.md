This folder is for local/dev secrets only. Do NOT commit real secrets.

What goes here (examples):
- gcp-service-account.json — Google Cloud service account credentials for local development.
- licenses.token — Generic bearer token for the external license directory (if enabled). Example provider: Airtable.
- licenses.dsn — Opaque DSN string for the external license directory (e.g., "base=appXXXXXXXXXXXX;table=Licenses").

Git rules:
- Everything in this folder is ignored by default. Only README.md and *.example files are tracked.

How to use:
1) Copy the example files provided and fill your local values:
   - cp secrets/gcp-service-account.json.example secrets/gcp-service-account.json
   - cp secrets/licenses.token.example secrets/licenses.token
   - cp secrets/licenses.dsn.example secrets/licenses.dsn
2) Ensure your env points to these files where applicable:
   - GOOGLE_APPLICATION_CREDENTIALS=./secrets/gcp-service-account.json (when running locally)
   - LICENSES_TOKEN: put the content of secrets/licenses.token into your .env (or inject as env in prod)
   - LICENSES_DSN: put the content of secrets/licenses.dsn into your .env (or inject as env in prod)

Docker Compose usage (api service):
- The compose file mounts ./secrets at /app/secrets (read-only) inside the container.
- The default GOOGLE_APPLICATION_CREDENTIALS is set to /app/secrets/gcp-service-account.json.

Security reminders:
- Never commit real secrets. Use .env for local/dev and platform-injected env for prod.
- The license provider is abstracted and not disclosed in code. Use LICENSES_PROVIDER=airtable to enable if desired, otherwise leave as none.
