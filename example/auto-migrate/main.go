package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/javadmohebbi/goAuth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	// command line argument for db path
	dbPath := flag.String("path", "/tmp/go-auth.db", "Path to sqlite database")
	// dbPath := flag.String("path", "/home/mj/Projects/goAuth/go-auth.db", "Path to sqlite database")

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
	db.AutoMigrate(&goAuth.User{}, &goAuth.Group{}, &goAuth.Policy{},
		&goAuth.UserInfoField{}, &goAuth.UserInfoFieldDetail{},
	)

	// migrate default fields
	migraDefaultFields(db)

	// create default user, group, policy
	migrateDefaultUserAndGroupAndPolicy(db)

	// fetch admin
	// u := fetchAdmin(db)
	// log.Println(u)

	// // // updte admin
	// updateUser(db, u)

}

func updateUser(db *gorm.DB, u goAuth.User) {
	u.MetaData = map[string]interface{}{
		"First Name":    "Mamadian",
		"Last Name":     "Mamadian toor",
		"Email Address": "mamad@mam2ad.co",
		"Birth Date":    "1368/12/07",
	}
	u.Password = "newPassword1_"
	db.Model(&u).Updates(
		goAuth.User{
			MetaData: u.MetaData,
			Password: u.Password,
		},
	)

	log.Println(u)
}

func fetchAdmin(db *gorm.DB) goAuth.User {
	var u goAuth.User
	db.Where("username = ?", "admin").First(&u)

	return u
}

// create default fields
// might required by most of users
func migraDefaultFields(db *gorm.DB) {
	customFields := []goAuth.UserInfoField{
		{
			FieldName: "Email Address",
			FieldType: 3,
			NotNull:   true,
			Unique:    true,
			Regex:     "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
		},
		{
			FieldName: "First Name",
			FieldType: 3,
			NotNull:   true,
			Regex:     "^.{3,20}$",
		},
		{
			FieldName: "Last Name",
			FieldType: 3,
			NotNull:   true,
			Regex:     "^.{3,20}$",
		},
		{
			FieldName: "Birth Date",
			FieldType: 5,
			NotNull:   true,

			// YYYY-mm-dd OR YYYY/mm/dd
			Regex: `(\d{4}-\d{2}-\d{2})|((\d{4}\/\d{2}\/\d{2}))`,
		},
	}

	for _, cf := range customFields {
		var _cf goAuth.UserInfoField
		db.Where("field_name = ?", cf.FieldName).First(&_cf)
		if _cf.ID == 0 {
			db.Create(&cf)
		}
	}
}

// create default user, group, policy
func migrateDefaultUserAndGroupAndPolicy(db *gorm.DB) {
	// add super admin group
	// which has access to everything
	adminGroup := goAuth.Group{
		Name: "Administator",
		Desc: "Members of this groups has full access to everything",
	}
	// insert group to db if not exists
	// create admin group if not exists
	db.Where("name = ?", adminGroup.Name).First(&adminGroup)

	// if not exist, insert it
	// or it will adminGroup will be fetched group
	if adminGroup.ID == 0 {
		db.Create(&adminGroup)
	}

	adminPolicy := []goAuth.Policy{

		// administrator group policy
		{
			Section: "*", // everything
			Perm:    15,  // rwud
			GroupID: adminGroup.ID,
		},
	}
	// insert policies if not exist
	for _, p := range adminPolicy {
		var _p goAuth.Policy
		db.Where("section = ? AND perm = ? AND group_id = ?", p.Section, p.Perm, p.GroupID).First(&_p)
		if _p.ID == 0 {
			db.Create(&p)
		}
	}

	// admin test user
	userAdmin := goAuth.User{
		Username: "admin",
		Password: "P4$$word", // password will be saved as hashed string
		MetaData: map[string]interface{}{
			"Email Address": "admin@example.com",
			"First Name":    "Dennis",
			"Last Name":     "Ritchie",
			"Birth Date":    "1941/09/09",
		},
	}

	// append group
	userAdmin.Group = append(userAdmin.Group, &adminGroup)

	// create admin user if not exists
	db.Where("username = ?", userAdmin.Username).First(&userAdmin)
	if userAdmin.ID == 0 {
		db.Create(&userAdmin)
	}
}
