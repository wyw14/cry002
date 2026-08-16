CREATE INDEX IF NOT EXISTS idx_cases_search ON case_files(year,status,classification_id,updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_cases_department ON case_files(responsible_department,handler_id);
CREATE INDEX IF NOT EXISTS idx_material_case_pages ON materials(case_id,page_from);
CREATE INDEX IF NOT EXISTS idx_versions_case_version ON case_versions(case_id,version DESC);
CREATE INDEX IF NOT EXISTS idx_borrow_status_due ON borrow_requests(status,due_at);
CREATE INDEX IF NOT EXISTS idx_refresh_family ON refresh_tokens(family_id) WHERE revoked_at IS NULL;
