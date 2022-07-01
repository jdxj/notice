package db

import (
	"context"
	"testing"
)

func TestDB(t *testing.T) {
	gormDB.WithContext(context.Background())
	sqlDB, err := gormDB.DB()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}
