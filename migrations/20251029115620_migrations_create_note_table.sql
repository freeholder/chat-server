-- +goose Up
create table request (
    Id SERIAL PRIMARY KEY,
    usernames TEXT[]
);
-- +goose Down
SELECT 'down SQL query';
