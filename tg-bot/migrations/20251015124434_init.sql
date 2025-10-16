-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reminder(
	id SERIAL PRIMARY KEY,
	chat_id BIGINT NOT NULL,
	text TEXT NOT NULL,
	type_reminder VARCHAR(32) NOT NULL,
	week_day TEXT,
	group_send_id BIGINT NOT NULL,
	time VARCHAR(128) NOT NULL,
	full_time VARCHAR(312) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_group(
	id SERIAL PRIMARY KEY,
	chat_id_group BIGINT NOT NULL,
	user_id BIGINT NOT NULL,
	title_group TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM reminder;
DELETE FROM user_group;
-- +goose StatementEnd
