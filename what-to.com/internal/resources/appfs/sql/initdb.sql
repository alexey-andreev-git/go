CREATE TABLE IF NOT EXISTS entities (
    id SERIAL PRIMARY KEY COMMENT 'entity''s id',
    reference INT NOT NULL COMMENT 'entity''s reference',
    comment VARCHAR(256) NOT NULL COMMENT 'entity''s reference'
);
