CREATE TABLE IF NOT EXISTS expenses (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    date timestamp(0) with time zone NOT NULL,
    recipient text NOT NULL,
    description text NOT NULL,
    amount bigint NOT NULL,
    currency CHAR(3) NOT NULL
);