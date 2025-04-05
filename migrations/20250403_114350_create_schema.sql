-- migrations/20250403_114350_create_schema.sql
-- Copyright (c) 2025 Michael D Henderson. All rights reserved.

create table if not exists books
(
    id          integer primary key autoincrement,
    title       text    not null,
    author      text    not null default '',
    instrument  text    not null default '',
    condition   text    not null default '',
    description text    not null default '',
    public      integer not null default 0 check (public in (0, 1))
);
