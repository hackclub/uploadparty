-- UploadParty initial database schema migration (moved to backend/migrations)
-- This file is safe to re-run thanks to IF NOT EXISTS guards.

-- 1) Enum for project status
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'project_status') THEN
        CREATE TYPE project_status AS ENUM ('in_progress', 'complete');
    END IF;
END$$;

-- 2) Users
CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    email         VARCHAR(255) NOT NULL,
    username      VARCHAR(50)  NOT NULL,
    password_hash TEXT         NOT NULL,
    display_name  VARCHAR(100),
    bio           VARCHAR(280),
    public        BOOLEAN      NOT NULL DEFAULT TRUE
);

-- Unique constraints per model tags
CREATE UNIQUE INDEX IF NOT EXISTS ux_users_email ON users (LOWER(email));
CREATE UNIQUE INDEX IF NOT EXISTS ux_users_username ON users (LOWER(username));

-- 3) Projects
CREATE TABLE IF NOT EXISTS projects (
    id                BIGSERIAL PRIMARY KEY,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    user_id           BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    title             VARCHAR(200) NOT NULL,
    daw               VARCHAR(100),
    plugin_version    VARCHAR(50),
    duration_seconds  INTEGER      NOT NULL DEFAULT 0,
    metadata          JSONB,
    status            project_status NOT NULL DEFAULT 'in_progress',
    completed_at      TIMESTAMPTZ,
    public            BOOLEAN      NOT NULL DEFAULT FALSE
);

-- Helpful uniqueness: upsert by (user_id, title)
CREATE UNIQUE INDEX IF NOT EXISTS ux_projects_user_title ON projects (user_id, LOWER(title));
CREATE INDEX IF NOT EXISTS ix_projects_user ON projects (user_id);
CREATE INDEX IF NOT EXISTS ix_projects_public ON projects (public);

-- 4) Plugins
CREATE TABLE IF NOT EXISTS plugins (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,

    name       VARCHAR(120) NOT NULL,
    vendor     VARCHAR(120),
    version    VARCHAR(50),
    format     VARCHAR(20), -- e.g., VST3, AU, AAX
    metadata   JSONB
);

-- Uniqueness per project: name unique within a project
CREATE UNIQUE INDEX IF NOT EXISTS idx_project_name ON plugins (project_id, LOWER(name));
CREATE INDEX IF NOT EXISTS ix_plugins_project ON plugins (project_id);

-- 5) updated_at maintenance trigger (optional)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'set_updated_at') THEN
        CREATE OR REPLACE FUNCTION set_updated_at()
        RETURNS TRIGGER AS $$
        BEGIN
            NEW.updated_at := NOW();
            RETURN NEW;
        END;
        $$ LANGUAGE plpgsql;
    END IF;
END$$;

-- Attach trigger to tables
DO $$
BEGIN
    -- users
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'tr_users_set_updated_at'
    ) THEN
        CREATE TRIGGER tr_users_set_updated_at
        BEFORE UPDATE ON users
        FOR EACH ROW EXECUTE FUNCTION set_updated_at();
    END IF;

    -- projects
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'tr_projects_set_updated_at'
    ) THEN
        CREATE TRIGGER tr_projects_set_updated_at
        BEFORE UPDATE ON projects
        FOR EACH ROW EXECUTE FUNCTION set_updated_at();
    END IF;

    -- plugins
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'tr_plugins_set_updated_at'
    ) THEN
        CREATE TRIGGER tr_plugins_set_updated_at
        BEFORE UPDATE ON plugins
        FOR EACH ROW EXECUTE FUNCTION set_updated_at();
    END IF;
END$$;
