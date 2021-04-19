# goAuth
Simple &amp; Easy go package for **authorization**.

Implementing Authentication & Authorization in a right way, some time took time from developers. I have developed a simple & easy-understandig `go module` for **Authorization**.


# Dependencies
There is no specific dependencies but I use  **gorm** (The great ORM library for Golang) as ORM library for **CRUD** in test DB. You can read it's docs [here]([here](https://gorm.io/index.html)).



# Usage
To use this `go module` you must run the further command.

```
go get github.com/javadmohebbi/goAuth
```


### Basic Example
```
    //...

    "github.com/javadmohebbi/goAuth"

    //...

    func main() {
        // ...


        /*
        struct to initialize user policies
        type GoAuthPolicy struct {


            // section is a place holder which can accept *
            // for all characters

            //      - * (everything)
            //      - app.dashboard
            //      - app.sales.order
            //      - admin.settings.*
            Section string `json:"section"`


            // Perm is an unsigned integer and must be between 0, 15
            // This is like mode system in extened file system in linux:
            //      - check 'man chmod'
            /***********************************************************
            // Permission Description
            // ----------------------------------------------
            //   #   Permission           rwud*      Binary
            // ----------------------------------------------
            //   0   none                 ----       0000
            //   1                        ---d       0001
            //   2                        --u-       0010
            //   3                        --ud       0011
            //   4                        -w--       0100
            //   5                        -w-d       0101
            //   6                        -wu-       0110
            //   7                        -wud       0111
            //   8                        r---       1000
            //   9                        r--d       1001
            //   10                       r-u-       1010
            //   11                       r-ud       1011
            //   12                       rw--       1100
            //   13                       rw-d       1101
            //   14                       rwu-       1110
            //   15                       rwud       1111
            // ----------------------------------------------
            // *rwdu => Read Write Update Delete
            // ----------------------------------------------
            Perm     Perm    `json:"perm"`

        }

        You yourself must fetch policies from your database &
        initialize the user's policies using this struct from goAuth module
        */

        var userPolicies []goAuth.GoAuthPolicy

        // initialize goAuthPolicy
        for _, p := range userFetchPoliciesFromDatavase {
            userPolicies = append(userPolicies, goAuth.GoAuthPolicy{
                Section: p.Section, Perm: goAuth.Perm(p.Perm),
            })
        }



        /**
            Now it's time to check if user has that policy or not
            using goAuth.GetPermissions(neededSection string) method

            This method will return bool, bool, bool, bool for Read, Write, Update, Delete
            permission & if user has that, will return true for that section.

            Also neededSection is the section that your user needs to know if has access to or not.
        */

        neededSection := "app.admin.dashoard"

        r, w, u, d := goAuth.Init(userPolicies).GetPermissions(neededSection)

        // check if your user has access to
        // the section you they need
        // ...
    }

```
