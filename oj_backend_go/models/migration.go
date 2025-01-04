package models

import (
	"fmt"

	"github.com/rudrakshsattabhayya/oj_backend_go/config"

	"github.com/rudrakshsattabhayya/oj_backend_go/models/auth"
	"github.com/rudrakshsattabhayya/oj_backend_go/models/oj"
)

func MigrateDB() {
	fmt.Println("Migrating Database...")

	db := config.GetDB()

	models := []interface{}{
		&auth.User{},
		&oj.ProblemId{},
		&oj.Tag{},
		&oj.Problem{},
		&oj.Submission{},
		&oj.ProblemTag{},
	}

	err := db.Migrator().AutoMigrate(models...)

	if err != nil {
		fmt.Printf("Migration failed: %v\n", err)
		return
	}

	fmt.Println("Migration Complete!")
}
