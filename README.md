tableCrud
=========

Take mysql tables and auto generate crud structs and methods.

Take a look at the following two code samples on how to use tableCrud. The first example is of its use. The second example is how to generate the table cruds. 

Currently, the package creates getters (for any column) and a transactional post, put, and delete (ie, insert, update, and delete).  Currently, it only recognizes columns that are varchar, int, and tinyint (bool) and has issues with NULL.  

### Usage

```
package main

import (
    "log"
    
    // wherever you put your crud (see the next code section for an example)
    "crud"
)

func main(){
	// set up your table's crud
	// you can access any table in your mysql db with `crud.<Table>{}`
	// here we are setting up two entries for the user table. This could 
	// easily be a user table and a settings table.
	u1 := crud.User{}
	u2 := crud.User{}

	// set up the db so the crud methods work
	// this currently is set to only use a local mysql instance with no username or password
	crud.SetDB()
	
	// ***
	// use the crud
	// ***
	
	// Getters: the crud.<Table> struct will have getters for all columns
	users := u1.GetById(1)
	log.Println(users[0].Email)
	
	// Post: insert a new record as a transaction
	u1.Post(&crud.UserRecord{Name: "Abbot"})
	u2.Post(&crud.UserRecord{Name: "Costello", Email: "ab@example.com"})
	
	// Put: update a record as a transaction
	users[0].Email = "new_email@example.com"
	u1.Put(users[0])
	
	// Delete: deletes a record as a transaction
	u1.Delete(users[0])
	
	// You can link transactions to happen in a single commit!
	u2.Tx = u1.Tx
	
	u2.Commit() // u1.Commit() would have the same effect of committing the transaction
}
```

### Creating crud files

Current state of this proof of concept. To generate your crud, you can do something like the following:

```
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sethgrid/tableCrud"
)

func main() {
	// connect to the database
	db_user := ""
	db_pw := ""
	db_name := "test_db"
	dataSource := ""
	if len(db_user) == 0 {
		dataSource = fmt.Sprintf("/%s", db_name)
	} else if len(db_pw) == 0 {
		dataSource = fmt.Sprintf("%s@/%s", db_user, db_name)
	} else {
		dataSource = fmt.Sprintf("%s:%s@/%s", db_user, db_pw, db_name)
	}

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatal("unable to open db ", err)
	}
	
	// generate crud for your database and place in src/crud
	tableCrud.GenDBConn("src/crud", true)
	allTables := tableCrud.GetTables(db)
	tableCrud.GenCrud(allTables, "src/crud", true)

	err = db.Close()
	if err != nil {
		log.Fatal("unable to close db", err)
	}

	log.Print("Done")
}

```
