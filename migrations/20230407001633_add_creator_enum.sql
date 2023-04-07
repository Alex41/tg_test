-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
create type message_creator as enum ('operator', 'user');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop type message_creator;
-- +goose StatementEnd
