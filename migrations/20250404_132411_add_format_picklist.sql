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

alter table books
    add column format text not null default '';
alter table books
    add column created_at datetime not null default CURRENT_TIMESTAMP;
alter table books
    add column updated_at datetime not null default CURRENT_TIMESTAMP;

create trigger if not exists set_updated_at
    after update on books
    for each row
begin
    update books set updated_at = CURRENT_TIMESTAMP where id = old.id;
end;

create index if not exists idx_books_format on books(format);

