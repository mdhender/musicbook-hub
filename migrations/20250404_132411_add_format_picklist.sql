-- migrations/20250404_132411_add_format_picklist.sql
-- Copyright (c) 2025 Michael D Henderson. All rights reserved.

create table if not exists format_picklist
(
    id          integer primary key autoincrement,
    format      text not null unique,
    description text not null
);

insert into format_picklist(format, description)
values ('sheet music', 'Single piece or short folio, typically softcover and intended for performance'),
       ('music book', 'Bound book of sheet music or exercises, such as method books, anthologies, or collections'),
       ('method book', 'Instructional book for learning an instrument, often organized by level'),
       ('score', 'Full musical score for ensembles, orchestras, or chamber music'),
       ('lead sheet', 'Single-line melody with chords, often used in jazz or pop music'),
       ('fake book', 'Gig-style book with hundreds of lead sheets for performance'),
       ('manuscript', 'Blank staff paper or note-taking formats for composers and students'),
       ('programming book', 'Technical or instructional book focused on coding, software, or computer science topics'),
       ('textbook', 'Educational book intended for academic study, often includes theory and exercises'),
       ('reference book',
        'Non-fiction book used for lookups or guidance, such as dictionaries, style guides, or API references');

create table books_new
(
    id          integer primary key autoincrement,
    title       text     not null,
    author      text     not null default '',
    condition   text     not null default '',
    format      text     not null default '',
    description text     not null default '',
    instrument  text     not null default '',
    public      integer  not null default 0 check (public in (0, 1)),
    created_at  datetime not null default current_timestamp,
    updated_at  datetime not null default current_timestamp
);

-- Copy data from old table
insert into books_new (id, title, author, instrument, condition, format, description, public)
select id,
       title,
       author,
       instrument,
       condition,
       format,
       description,
       public
from books;

-- drop old table and rename
drop table books;
alter table books_new
    rename to books;

create trigger if not exists set_updated_at
    after update
    on books
    for each row
begin
    update books set updated_at = CURRENT_TIMESTAMP where id = old.id;
end;

create index if not exists idx_books_format on books (format);

