tableCrud
=========

Take mysql tables and auto generate crud structs

Current implemenation

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
    
    	allTables := tableCrud.GetTables(db)
    	tableCrud.GenCrud(allTables, "src/github.com/sethgrid/tableCrud", false)
    
    	log.Print("user's first column is a ", allTables["user"].Cols[0].COLUMN_TYPE)
    
    	err = db.Close()
    	if err != nil {
    		log.Fatal("unable to close db", err)
    	}
    
    	log.Print("Closed")
    }
```
