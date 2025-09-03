-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create type preference_enum as enum (
    'CARDIO',
    'WEIGHT'
);

create type weight_unit_enum as enum (
    'KG',
    'LBS'
);

create type height_unit_enum as enum (
    'CM',
    'INCH'
);

create table if not exists users (
    id serial primary key,
    email text unique not null,
    password_hash text not null,
    name varchar(60),
    preference preference_enum,
    weight_unit weight_unit_enum,
    height_unit height_unit_enum,
    weight int,
    height int,
    image_uri text,
    created_at timestamptz default current_timestamp,
    updated_at timestamptz default current_timestamp
);

CREATE INDEX IF NOT EXISTS users_email_idx ON users USING HASH (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table if exists users;
-- +goose StatementEnd
