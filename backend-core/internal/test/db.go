package test

import (
	"path"
	"runtime"
	"testing"

	"github.com/Wave-95/boards/backend-core/db"
	"github.com/Wave-95/boards/backend-core/internal/config"
)

func DB(t *testing.T) *db.DB {
	// load env vars into config
	dir := getSourcePath()
	cfg, err := config.Load(dir + "/../../.env")
	if err != nil {
		t.Errorf("Issue loading config:%v", err)
		t.FailNow()
	}
	// connect to db
	db, err := db.Connect(cfg.DatabaseConfig)
	if err != nil {
		t.Errorf("Issue connecting db:%v", err)
		t.FailNow()
	}
	return db
}

// getSourcePath returns the directory containing the source code that is calling this function.
// thanks go-rest-api
func getSourcePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
