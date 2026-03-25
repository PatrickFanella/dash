CREATE TABLE service_section_mappings (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    section_id UUID NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(service_id, section_id)
);

CREATE INDEX idx_ssm_section_id ON service_section_mappings(section_id);
CREATE INDEX idx_ssm_service_id ON service_section_mappings(service_id);
