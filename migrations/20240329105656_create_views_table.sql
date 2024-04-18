-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS views
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resume_id CHAR(24) NOT NULL,
    company_id UUID NOT NULL,
    viewed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE views;
-- +goose StatementEnd
