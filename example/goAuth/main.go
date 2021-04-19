package main

import (
	"flag"
	"fmt"

	"github.com/javadmohebbi/goAuth"
	"github.com/javadmohebbi/goAuth/model"
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

	/**
	Above steps is an example of preparing datatbase
	model.User{}, model.Group{}, model.Policy{} are mandatory
	fields in your db model, you must first have this models &
	then you can use this package, OR use the above sample
	to make your own!
	**/

	// fetch user's policies from database
	var adminPolicies []model.Policy
	var userPolicies []model.Policy

	// admin
	db.Raw(`
		SELECT p.section, p.perm FROM policies as p
		JOIN groups as g ON g.id = p.group_id
		JOIN group_users as gu ON g.id = gu.group_id
		WHERE gu.user_id = ?
	`, adminUser.ID).Scan(&adminPolicies)

	// public
	db.Debug().Raw(`
		SELECT p.section, p.perm FROM policies as p
		JOIN groups as g ON g.id = p.group_id
		JOIN group_users as gu ON g.id = gu.group_id
		WHERE gu.user_id = ?
	`, normalUser.ID).Scan(&userPolicies)

	// prepare goAuth
	//
	// and feed it with the fetched policies
	var gaAdminPS []goAuth.GoAuthPolicy
	var gaUserPS []goAuth.GoAuthPolicy

	// initialize admin goAuthPolicy
	for _, p := range adminPolicies {
		gaAdminPS = append(gaAdminPS, goAuth.GoAuthPolicy{
			Section: p.Section, Perm: goAuth.Perm(p.Perm),
		})
	}

	// initialize user goAuthPolicy
	for _, p := range userPolicies {
		gaUserPS = append(gaUserPS, goAuth.GoAuthPolicy{
			Section: p.Section, Perm: goAuth.Perm(p.Perm),
		})
	}

	// test permissions
	// testPermissions := []string{"app.users", "admin.dashboard", "app.users.orders", "app.networks", "app.infrastructures.datacenter"}
	testPermissions := []string{"app.users.orders", "app.upload", "admin.dashboard.create"}

	fmt.Printf("\n\n\n\n- - - - - - - - - - - - - - - - - - - - - - - - -\n\n\n\n")
	fmt.Printf("Permissions for 'Admin' are: \n\n")
	// test admin permission
	for _, tp := range testPermissions {
		r, w, u, d := goAuth.Init(gaAdminPS).GetPermissions(tp)
		fmt.Printf("'%v' permission for %v:\n\tRead:%v\tWrite:%v\tUpdate:%v\tDelete:%v\n\n",
			adminUser.Username, tp, r, w, u, d,
		)
	}

	fmt.Printf("\n- - - - - - - - - - - - - - - - - - - - - - - - -\n\n\n\n")
	fmt.Printf("Permissions for 'Users' are: \n\n")
	// test user permission
	for _, tp := range testPermissions {
		r, w, u, d := goAuth.Init(gaUserPS).GetPermissions(tp)
		fmt.Printf("'%v' permission for %v:\n\tRead:%v\tWrite:%v\tUpdate:%v\tDelete:%v\n\n",
			normalUser.Username, tp, r, w, u, d,
		)
	}

}

// insert some test data into database
// an example of preparing datatbase
// model.User{}, model.Group{}, model.Policy{} are mandatory
// fields in your db model, you must first have this models &
// then you can use this package, OR use the above sample
// to make your own!
func insertTestData(db *gorm.DB) (model.User, model.User) {

	// admin group
	adminGroup := model.Group{
		Name: "administrator",
		Desc: "Unlimited group which has access to every thing",
	}
	// create admin group if not exists
	db.Debug().Where("name = ?", adminGroup.Name).First(&adminGroup)
	if adminGroup.ID == 0 {
		db.Create(&adminGroup)
	}

	// public group
	publicGroup := model.Group{
		Name: "public",
		Desc: "Limited group which has access to some specific sections",
	}
	// create public group if not exists
	db.Debug().Where("name = ?", publicGroup.Name).First(&publicGroup)
	if publicGroup.ID == 0 {
		db.Create(&publicGroup)
	}

	plcs := []model.Policy{

		// admin
		{
			Section: "*",
			Perm:    15, // rwud
			GroupID: adminGroup.ID,
		},

		// normal
		{
			Section: "app.users",
			Perm:    14, // rwu-
			GroupID: publicGroup.ID,
		},
		{
			Section: "app.users.orders",
			Perm:    8, // r---
			GroupID: publicGroup.ID,
		},
	}

	for _, p := range plcs {
		var _p model.Policy
		db.Where("section = ? AND perm = ? AND group_id = ?", p.Section, p.Perm, p.GroupID).First(&_p)
		if _p.ID == 0 {
			db.Create(&p)
		}
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
	db.Debug().Where("username = ?", userAdmin.Username).First(&userAdmin)
	if userAdmin.ID == 0 {
		db.Create(&userAdmin)
	}

	// add admin user to administrator group
	db.Debug().Model(&userAdmin).Association("Group").Append([]*model.Group{
		&adminGroup,
	})

	// public
	// create normal user if not exists
	db.Debug().Where("username = ?", userNormal.Username).First(&userNormal)
	if userNormal.ID == 0 {
		db.Create(&userNormal)
	}

	// add normal user to public group
	db.Debug().Model(&userNormal).Association("Group").Append([]*model.Group{
		&publicGroup,
	})

	return userAdmin, userNormal

}
