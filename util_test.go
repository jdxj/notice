package notice

import (
	"context"
	"testing"
)

func TestDB(t *testing.T) {
	db.WithContext(context.Background())
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}
