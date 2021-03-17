package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/goAuth"
	"github.com/goAuth/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	// command line argument for db path
	// dbPath := flag.String("path", "/tmp/go-auth.db", "Path to sqlite database")
	dbPath := flag.String("path", "/home/mj/Projects/goAuth/go-auth.db", "Path to sqlite database")

	// parse the command line arguments
	flag.Parse()

	// *dbPath = "/home/mj/Projects/goAuth/go-auth.db"

	// open or create sqlite database
	db, err := gorm.Open(sqlite.Open(*dbPath), &gorm.Config{})

	// check if no successful
	if err != nil {
		fmt.Printf("Can not open %v due to error: %v", *dbPath, err)
	}

	// migrate db
	db.AutoMigrate(&model.User{}, &model.Group{}, &model.Policy{})

	// insert test data
	adminUser, normalUser := insertTestData(db)

	fmt.Printf("Two users have prepared for tests:\n\t1- %v (administrator)\n\t2- %v (normal user)\n",
		adminUser.Username, normalUser.Username,
	)

	/**
	Above steps is an example of preparing datatbase
	model.User{}, model.Group{}, model.Policy{} are mandatory
	fields in your db model, you must first have this models &
	then you can use this package, OR use the above sample
	to make your own!
	**/

	// load user with groups
	db.Set("gorm:auto_preload", true).Where("id = ?", adminUser.ID).First(&adminUser)
	log.Println(adminUser)

	ga := goAuth.Init([]goAuth.GoAuthPolicy{
		{
			Section: "admin.dashboard.*",
			UGO:     14,
		},
		{
			Section: "app.user.*",
			UGO:     14,
		},
	})

	for _, p := range ga.Policies {
		fmt.Println(p.UGO.Bools())
		fmt.Println(p.UGO)
	}

}

// insert some test data into database
func insertTestData(db *gorm.DB) (model.User, model.User) {

	// admin group
	adminGroup := model.Group{
		Name: "administrator",
		Desc: "Unlimited group which has access to any thing",
	}
	// create admin group if not exists
	db.Where("name = ?", adminGroup.Name).First(&adminGroup)
	if adminGroup.ID == 0 {
		db.Create(&adminGroup)
	}

	// public group
	publicGroup := model.Group{
		Name: "public",
		Desc: "Limited group which has access to some specific sections",
	}
	// create public group if not exists
	db.Where("name = ?", publicGroup.Name).First(&publicGroup)
	if publicGroup.ID == 0 {
		db.Create(&publicGroup)
	}

	// add all access policy
	policyAdmin := model.Policy{
		Section: "*",
		UGO:     "15", // rwud
		GroupID: adminGroup.ID,
	}
	db.Create(&policyAdmin)

	// add all access policy
	policiesPublic := []model.Policy{
		{
			Section: "app.users.*",
			UGO:     "14", // rwu-
			GroupID: publicGroup.ID,
		},
		{
			Section: "app.orders.*",
			UGO:     "8", // r---
			GroupID: publicGroup.ID,
		},
	}
	for _, p := range policiesPublic {
		db.Create(&p)
	}

	// admin test user
	userAdmin := model.User{
		Username: "mjmohebbiAdmin",
		Password: "hashed-secret",

		FirstName: "M. Javad",
		LasttName: "Mohebbi",
	}

	// normal test user
	userNormal := model.User{
		Username: "mjmohebbiNormal",
		Password: "hashed-secret",

		FirstName: "M. Javad",
		LasttName: "Mohebbi",
	}

	// admin
	// create admin user if not exists
	db.Where("username = ?", userAdmin.Username).First(&userAdmin)
	if userAdmin.ID == 0 {
		db.Create(&userAdmin)
	}

	// add admin user to administrator group
	db.Model(&userAdmin).Association("Group").Append([]*model.Group{
		&adminGroup,
	})

	// public
	// create normal user if not exists
	db.Where("username = ?", userNormal.Username).First(&userNormal)
	if userNormal.ID == 0 {
		db.Create(&userNormal)
	}

	// add normal user to public group
	db.Model(&userNormal).Association("Group").Append([]*model.Group{
		&publicGroup,
	})

	return userAdmin, userNormal

}
