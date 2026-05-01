-- Workflow Definitions
CREATE TABLE workflow_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    version INT NOT NULL DEFAULT 1,
    definition JSONB NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(name, version)
);

-- Workflow Instances
CREATE TABLE workflow_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    definition_id UUID REFERENCES workflow_definitions(id),
    status VARCHAR(50) NOT NULL,
    current_step VARCHAR(255),
    input JSONB,
    output JSONB,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Workflow Events (Event Sourcing)
CREATE TABLE workflow_events (
    id BIGSERIAL PRIMARY KEY,
    instance_id UUID REFERENCES workflow_instances(id),
    event_type VARCHAR(100) NOT NULL,
    step_id VARCHAR(255),
    event_data JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Step Executions
CREATE TABLE step_executions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instance_id UUID REFERENCES workflow_instances(id),
    step_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    attempt INT DEFAULT 1,
    input JSONB,
    output JSONB,
    error TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Human Approvals
CREATE TABLE approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instance_id UUID REFERENCES workflow_instances(id),
    step_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'PENDING',
    approver_email VARCHAR(255),
    approved_by VARCHAR(255),
    approved_at TIMESTAMP,
    comments TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Compensations
CREATE TABLE compensations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instance_id UUID REFERENCES workflow_instances(id),
    step_id VARCHAR(255) NOT NULL,
    compensation_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'PENDING',
    result JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    executed_at TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_instances_status ON workflow_instances(status);
CREATE INDEX idx_instances_created ON workflow_instances(created_at);
CREATE INDEX idx_events_instance ON workflow_events(instance_id);
CREATE INDEX idx_events_created ON workflow_events(created_at);
CREATE INDEX idx_steps_instance ON step_executions(instance_id);
CREATE INDEX idx_approvals_status ON approvals(status);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_workflow_definitions_updated_at BEFORE UPDATE ON workflow_definitions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_workflow_instances_updated_at BEFORE UPDATE ON workflow_instances
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
