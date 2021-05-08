package main

import (
	"fmt"

	"github.com/javadmohebbi/goAuth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	USER_TABLE_NAME   = "users"
	GROUP_TABLE_NAME  = "groups"
	POLICY_TABLE_NAME = "policies"
)

// user model definition
type user struct {

	// default schema for users
	goAuth.User

	// Field1 string
	// Field2 string
	// FieldN datatype
}

// force table name
func (t user) TableName() string {
	return USER_TABLE_NAME
}

// group model definition
type group struct {

	// default schema for groups
	goAuth.Group

	// Field1 string
	// Field2 string
	// FieldN datatype
}

// force table name
func (t group) TableName() string {
	return GROUP_TABLE_NAME
}

// policy model definition
type policy struct {

	// default schema for policy
	goAuth.Policy

	// Field1 string
	// Field2 string
	// FieldN datatype
}

// force table name
func (t policy) TableName() string {
	return POLICY_TABLE_NAME
}

func main() {

	// openning/creating-new sqlite database
	db, err := gorm.Open(sqlite.Open("/tmp/goAuthExample.sqlite"), &gorm.Config{})

	// if can not open database, stop the execution
	// otherwise continue with openned/created db
	if err != nil {
		panic(err)
	}

	// auto migrate database
	// create/update database tables with the provided models
	// migrate db
	db.AutoMigrate(&user{}, &group{}, &policy{})

	// add super admin group
	// which has access to everything
	adminGroup := group{
		Group: goAuth.Group{
			Name: "Administrator",
			Desc: "This user has access to everything",
		},
	}
	// insert group to db if not exists
	// create admin group if not exists
	db.Where("name = ?", adminGroup.Name).First(&adminGroup)

	// if not exist, insert it
	// or it will adminGroup will be fetched group
	if adminGroup.ID == 0 {
		db.Create(&adminGroup)
	}

	// add Sales group
	// which has full access to sales.orders.*
	// but has read-only access to sales.transfer.ready
	salesGrp := group{
		Group: goAuth.Group{
			Name: "Sales",
			Desc: "All of our sales group could be member of this group",
		},
	}
	// insert group to db if not exists
	// create admin group if not exists
	db.Where("name = ?", salesGrp.Name).First(&salesGrp)

	// if not exist, insert it
	// or it will salesGrp will be fetched group
	if salesGrp.ID == 0 {
		db.Create(&salesGrp)
	}

	// add all the sampled policies to
	// database table
	policies := []policy{
		{
			// administrator group policy
			goAuth.Policy{
				Section: "*", // everything
				Perm:    15,  // rwud
				GroupID: adminGroup.ID,
			},
		},
		{
			// Sales group policies
			goAuth.Policy{
				Section: "sales.order.*",
				Perm:    15, // rwud
				GroupID: salesGrp.ID,
			},
		},
		{
			// Sales group policies
			goAuth.Policy{
				Section: "sales.transfer.ready",
				Perm:    8, // r---
				GroupID: salesGrp.ID,
			},
		},
	}

	// insert policies if not exist
	for _, p := range policies {
		var _p policy
		db.Where("section = ? AND perm = ? AND group_id = ?", p.Section, p.Perm, p.GroupID).First(&_p)
		if _p.ID == 0 {
			db.Create(&p)
		}
	}

	// insert two user with admin accesss and sales accesss

	// admin test user
	userAdmin := goAuth.User{
		Username: "adminUser",
		Password: "hashed-secret", // store password not in clear-text format

		FirstName: "M. Javad",
		LastName:  "Mohebbi",
	}
	// add admin to the adminUser group
	userAdmin.Group = append(userAdmin.Group, &adminGroup.Group)

	// sales test user
	salesUser := goAuth.User{
		Username: "salesUser",
		Password: "hashed-secret", // store password not in clear-text format

		FirstName: "M. Javad",
		LastName:  "Mohebbi",
	}
	// add salesuser to the sales group
	salesUser.Group = append(salesUser.Group, &salesGrp.Group)

	// admin
	// create admin user if not exists
	db.Where("username = ?", userAdmin.Username).First(&userAdmin)
	if userAdmin.ID == 0 {
		db.Create(&userAdmin)
	}

	// public
	// create normal user if not exists
	db.Where("username = ?", salesUser.Username).First(&salesUser)
	if salesUser.ID == 0 {
		db.Create(&salesUser)
	}

	// fetch user's policies from database
	var adminPolicies []goAuth.Policy
	var userPolicies []goAuth.Policy

	// admin
	db.Debug().Raw(`
	SELECT p.section, p.perm FROM `+POLICY_TABLE_NAME+` as p
	JOIN `+GROUP_TABLE_NAME+` as g ON g.id = p.group_id
	JOIN `+GROUP_TABLE_NAME[:len(GROUP_TABLE_NAME)-1]+`_`+USER_TABLE_NAME+` as gu ON g.id = gu.group_id
	WHERE gu.user_id = ?
`, userAdmin.ID).Scan(&adminPolicies)

	// public
	db.Debug().Raw(`
	SELECT p.section, p.perm FROM `+POLICY_TABLE_NAME+` as p
	JOIN `+GROUP_TABLE_NAME+` as g ON g.id = p.group_id
	JOIN `+GROUP_TABLE_NAME[:len(GROUP_TABLE_NAME)-1]+`_`+USER_TABLE_NAME+` as gu ON g.id = gu.group_id
	WHERE gu.user_id = ?
`, salesUser.ID).Scan(&userPolicies)

	// prepare goAuth
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
	testPermissions := []string{"sales.orders.files", "sales.transfer.ready", "app.upload", "admin.dashboard.create"}

	fmt.Printf("\n\n\n\n- - - - - - - - - - - - - - - - - - - - - - - - -\n\n\n\n")
	fmt.Printf("Permissions for 'Admin' are: \n\n")
	// test admin permission
	for _, tp := range testPermissions {
		r, w, u, d := goAuth.Init(gaAdminPS).GetPermissions(tp)
		fmt.Printf("'%v' permission for %v:\n\tRead:%v\tWrite:%v\tUpdate:%v\tDelete:%v\n\n",
			userAdmin.Username, tp, r, w, u, d,
		)
	}

	fmt.Printf("\n- - - - - - - - - - - - - - - - - - - - - - - - -\n\n\n\n")
	fmt.Printf("Permissions for 'Users' are: \n\n")
	// test user permission
	for _, tp := range testPermissions {
		r, w, u, d := goAuth.Init(gaUserPS).GetPermissions(tp)
		fmt.Printf("'%v' permission for %v:\n\tRead:%v\tWrite:%v\tUpdate:%v\tDelete:%v\n\n",
			salesUser.Username, tp, r, w, u, d,
		)
	}

}
