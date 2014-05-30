package tableCrud

import (
	"log"
	"strings"
)

const DB_CONNECT = `package crud

import (
	"log"
	"fmt"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func SetDB(){
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
	var err error
	DB, err = sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatal("unable to open db ", err)
	}
}
`

const OUTPUT_FILE = `{{$name := .Name}}{{$cols := .Cols}}package crud

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type {{ title  .Name }}Record struct{ {{range $cols}}
	{{title .COLUMN_NAME.String }} {{ castType .DATA_TYPE.String }}{{end}}
}

type {{ title .Name }} struct{}

{{ range $cols }}
func (t {{ title $name}}) GetBy{{ title .COLUMN_NAME.String}}({{ .COLUMN_NAME.String }} {{ castType .DATA_TYPE.String }}) []*{{ title $name}}Record{
	r, err := DB.Query("select * from {{ $name }} where {{ .COLUMN_NAME.String }}=?", {{ .COLUMN_NAME.String}})
	if err != nil {
		// TODO: use our ln package
		log.Print("error in Query: ", err)
	}
	res := make([]*{{ title $name }}Record, 0)
	for r.Next() { {{ range $cols }}
		var {{ .COLUMN_NAME.String }} {{ castType .DATA_TYPE.String }}{{end}}
		err = r.Scan({{ join $cols }})
		if err != nil {
			log.Print("error scanning schema ", err)
		}
		s := &{{ title $name }}Record{ {{ structInit $cols }} }
		res = append(res, s)
	}
	return res
}
{{end}}
`

// Helper function to populate Scan arguments in template
func joinComma(cols []*TableSchema) string {
	conversion := make([]string, 0)
	for _, col := range cols {
		conversion = append(conversion, "&"+col.COLUMN_NAME.String)
	}
	return strings.Join(conversion, ",")
}

// Helper function to init a struct
func structInit(cols []*TableSchema) string {
	conversion := make([]string, 0)
	for _, col := range cols {
		conversion = append(conversion, strings.Title(col.COLUMN_NAME.String)+":"+col.COLUMN_NAME.String)
	}
	return strings.Join(conversion, ",")
}

// Helper function to change mysql types to go types
func castType(m string) string {
	switch m {
	// TODO: add more cases
	case "varchar":
		return "string"
	case "int":
		return "int64"
	case "tinyint":
		return "bool"
	default:
		// TODO: change this to panic after all cases are set
		log.Print("Unknown type: ", m)
		return "?"
	}
}
