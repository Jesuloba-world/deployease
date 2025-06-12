-- +goose Up
-- +goose StatementBegin
CREATE TYPE deployment_status AS ENUM ('pending', 'in_progress', 'success', 'failed', 'cancelled');

CREATE TABLE deployments (
    id VARCHAR(32) PRIMARY KEY,
    project_id VARCHAR(32) NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    status deployment_status NOT NULL DEFAULT 'pending',
    commit_hash VARCHAR(255),
    deployed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_deployments_project_id ON deployments (project_id);

CREATE INDEX idx_deployments_status ON deployments (status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS deployments;

DROP TYPE IF EXISTS deployment_status;
-- +goose StatementEnd