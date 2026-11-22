CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,

    title VARCHAR(120) NOT NULL,

    description TEXT NOT NULL DEFAULT '',

    status VARCHAR(20) NOT NULL DEFAULT 'new'
        CHECK (status IN ('new', 'in_progress', 'done')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
