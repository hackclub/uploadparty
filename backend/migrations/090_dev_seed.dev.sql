-- DEV-ONLY seed data. This file is applied only when running migrate with -env=dev or MIGRATIONS_ENV=dev.
-- Creates a demo user and a sample project with one plugin.
-- DO NOT run in production.

-- Upsert demo user
INSERT INTO users (id, email, username, password_hash, display_name, bio, public, created_at, updated_at)
VALUES (
  1001,
  'demo@example.com',
  'demo',
  -- bcrypt hash for password: demo123
  '$2a$10$0qD1Yw0b1v2CzvQWb2y7sOa7k2GZ0c0o7Z1j7r9qz7E1y7h6iNzu6',
  'Demo User',
  'This is a demo account for local testing.',
  TRUE,
  NOW(), NOW()
)
ON CONFLICT (id) DO UPDATE SET
  email = EXCLUDED.email,
  username = EXCLUDED.username,
  display_name = EXCLUDED.display_name,
  bio = EXCLUDED.bio,
  public = EXCLUDED.public,
  updated_at = NOW();

-- Create a sample project for demo user if not exists
INSERT INTO projects (id, user_id, title, daw, plugin_version, duration_seconds, metadata, status, completed_at, public, created_at, updated_at)
VALUES (
  2001,
  1001,
  'My First Beat',
  'FL Studio',
  '1.0.0',
  120,
  '{}',
  'in_progress',
  NULL,
  TRUE,
  NOW(), NOW()
)
ON CONFLICT (id) DO NOTHING;

-- Attach a sample plugin
INSERT INTO plugins (id, project_id, name, vendor, version, format, metadata, created_at, updated_at)
VALUES (
  3001,
  2001,
  'SuperSynth',
  'Acme Audio',
  '2.3.1',
  'VST3',
  '{"preset":"Init"}',
  NOW(), NOW()
)
ON CONFLICT (id) DO NOTHING;
