-- +goose Up
create index events_start_at_and_end_at
    on events (start_at, end_at);

create index events_start_at
    on events (start_at);

-- +goose Down
drop index events_start_at;
drop index events_start_at_and_end_at;
