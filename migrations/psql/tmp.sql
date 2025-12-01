-- Journal table and seed data

CREATE TABLE IF NOT EXISTS journals (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO journals (id, name) VALUES
    ('7f8a9e3a-1c2b-4d5e-8f90-1234567890ab', 'Personal'),
    ('1b2c3d4e-5f60-7181-92a3-b4c5d6e7f809', 'Memory')
ON CONFLICT DO NOTHING;

-- Entry table with foreign key constraint and seed data

CREATE TABLE IF NOT EXISTS entries (
    id          TEXT PRIMARY KEY,
    journal_id  TEXT NOT NULL,
    title       TEXT NOT NULL,
    content     TEXT NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE entries ADD CONSTRAINT fk_journal_id FOREIGN KEY (journal_id) REFERENCES journals(id); 

INSERT INTO entries (id, journal_id, title, content) VALUES
    ('18dd943a-400b-4075-9490-daef64aefade', '7f8a9e3a-1c2b-4d5e-8f90-1234567890ab', 'Entry 1', 'Content 1'),
    ('29cc193a-400b-4075-9490-daef64aefade', '7f8a9e3a-1c2b-4d5e-8f90-1234567890ab', 'Entry 2', 'Content 2'),
    ('3aee3d3a-400b-4075-9490-daef64aefade', '7f8a9e3a-1c2b-4d5e-8f90-1234567890ab', 'Entry unknown', 'Unknown content'),
    ('4b00613a-400b-4075-9490-daef64aefade', '7f8a9e3a-1c2b-4d5e-8f90-1234567890ab', 'Entry 4', 'Content 4')
ON CONFLICT DO NOTHING;

-- Tags table and seed data

CREATE TABLE IF NOT EXISTS tags (
    id          TEXT PRIMARY KEY,
    label       TEXT NOT NULL UNIQUE
);

INSERT INTO tags (id, label) VALUES
    ('a1b2c3d4-5e6f-7890-abcd-ef1234567890', 'important'),
    ('b2c3d4e5-6f70-8901-bcde-f12345678901', 'memo'),
    ('c3d4e5f6-7081-9012-cdef-123456789012', 'personal'),
    ('d4e5f6g7-8192-0123-def0-234567890123', 'todo')
ON CONFLICT DO NOTHING;

-- Tag assignments table

CREATE TABLE IF NOT EXISTS tag_assignments (
    tag_id      TEXT NOT NULL,
    entry_id    TEXT NOT NULL,
    PRIMARY KEY (tag_id, entry_id),
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    FOREIGN KEY (entry_id) REFERENCES entries(id) ON DELETE CASCADE
);

-- Assign tags to entries:
-- Entry 1: all tags
-- Entry 2: few tags (important, memo)
-- Entry 3: none
-- Entry 4: only one tag (personal)

INSERT INTO tag_assignments (tag_id, entry_id) VALUES
    -- Entry 1: all tags
    ('a1b2c3d4-5e6f-7890-abcd-ef1234567890', '18dd943a-400b-4075-9490-daef64aefade'),
    ('b2c3d4e5-6f70-8901-bcde-f12345678901', '18dd943a-400b-4075-9490-daef64aefade'),
    ('c3d4e5f6-7081-9012-cdef-123456789012', '18dd943a-400b-4075-9490-daef64aefade'),
    ('d4e5f6g7-8192-0123-def0-234567890123', '18dd943a-400b-4075-9490-daef64aefade'),
    -- Entry 2: few tags (important, memo)
    ('a1b2c3d4-5e6f-7890-abcd-ef1234567890', '29cc193a-400b-4075-9490-daef64aefade'),
    ('b2c3d4e5-6f70-8901-bcde-f12345678901', '29cc193a-400b-4075-9490-daef64aefade'),
    -- Entry 3: none (no assignments)
    -- Entry 4: only one tag (personal)
    ('c3d4e5f6-7081-9012-cdef-123456789012', '4b00613a-400b-4075-9490-daef64aefade')
ON CONFLICT DO NOTHING;
