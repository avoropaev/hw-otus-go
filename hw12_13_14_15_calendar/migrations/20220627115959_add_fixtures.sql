-- +goose Up

INSERT INTO events(guid, title, start_at, end_at, description, user_guid, notify_before, notified)
VALUES
       ('a83c340e-fcf2-4a93-84ff-3fbbed8da001', 'will be deleted', '2021-01-01 04:39:49.178000', '2021-01-01 06:39:49.178000', '', 'a83c340e-fcf2-4a93-84ff-3fbbed8da003', '0 years 0 mons 0 days 0 hours 16 mins 40.0 secs', true),
       ('a83c340e-fcf2-4a93-84ff-3fbbed8da002', 'will be deleted 2', '2021-01-01 04:39:49.178000', '2021-01-01 06:39:49.178000', '', 'a83c340e-fcf2-4a93-84ff-3fbbed8da003', '0 years 0 mons 0 days 0 hours 16 mins 40.0 secs', true),

       ('a83c340e-fcf2-4a93-84ff-3fbbed8da003', 'will be notified', '2022-06-01 04:39:49.178000', '2022-06-01 06:39:49.178000', '', 'a83c340e-fcf2-4a93-84ff-3fbbed8da003', '0 years 0 mons 0 days 0 hours 16 mins 40.0 secs', false),
       ('a83c340e-fcf2-4a93-84ff-3fbbed8da004', 'will be notified 2', '2022-06-01 04:39:49.178000', '2022-06-01 06:39:49.178000', '', 'a83c340e-fcf2-4a93-84ff-3fbbed8da003', '0 years 0 mons 0 days 0 hours 16 mins 40.0 secs', false),
       ('a83c340e-fcf2-4a93-84ff-3fbbed8da005', 'will be notified later', '2022-12-01 04:39:49.178000', '2022-12-01 06:39:49.178000', '', 'a83c340e-fcf2-4a93-84ff-3fbbed8da003', '0 years 0 mons 0 days 0 hours 16 mins 40.0 secs', false);

-- +goose Down
truncate events;
