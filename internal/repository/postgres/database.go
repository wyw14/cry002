package postgres

import (
	"context"
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct{ DB *gorm.DB }

func Open(ctx context.Context, dsn string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}
	return &Database{DB: db}, nil
}
func (d *Database) Ready(ctx context.Context) error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
func (d *Database) Migrate(ctx context.Context) error {
	return d.DB.WithContext(ctx).AutoMigrate(&userRow{}, &departmentRow{}, &classificationRow{}, &caseRow{}, &materialRow{}, &versionRow{}, &reviewRow{}, &borrowRow{}, &attachmentRow{}, &auditRow{}, &tokenRow{})
}

var _ *sql.DB
