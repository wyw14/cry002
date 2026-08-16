package postgres

func (userRow) TableName() string           { return "users" }
func (departmentRow) TableName() string     { return "departments" }
func (classificationRow) TableName() string { return "classifications" }
func (caseRow) TableName() string           { return "case_files" }
func (materialRow) TableName() string       { return "materials" }
func (versionRow) TableName() string        { return "case_versions" }
func (reviewRow) TableName() string         { return "reviews" }
func (borrowRow) TableName() string         { return "borrow_requests" }
func (attachmentRow) TableName() string     { return "attachments" }
func (auditRow) TableName() string          { return "audit_events" }
func (tokenRow) TableName() string          { return "refresh_tokens" }
