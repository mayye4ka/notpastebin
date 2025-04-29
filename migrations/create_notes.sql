-- +migrate Up
CREATE TABLE IF NOT EXISTS notes (
    admin_hash varchar(32) NOT NULL CONSTRAINT unique_admin_hash UNIQUE,
    reader_hash varchar(32) NOT NULL CONSTRAINT unique_reader_hash UNIQUE,
    note text
);

-- +migrate Down
DROP TABLE IF EXISTS notes;