CREATE TABLE sections (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name         TEXT NOT NULL,
    icon         TEXT NOT NULL DEFAULT '',
    cols         INTEGER NOT NULL DEFAULT 3,
    collapsed    BOOLEAN NOT NULL DEFAULT false,
    sort_order   INTEGER NOT NULL DEFAULT 0,
    section_type TEXT NOT NULL DEFAULT 'services',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
