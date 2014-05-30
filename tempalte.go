package tableCrud

import (
	"log"
	"strings"
)

const OUTPUT_FILE = `{{$name := .Name}}{{$cols := .Cols}}package apid

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type {{ title  .Name}} struct{ {{range $cols}}
	{{title .COLUMN_NAME.String }} {{ castType .DATA_TYPE.String }}{{end}}
}
{{ range $cols }}
func (t *{{ title $name}}) GetBy{{ title .COLUMN_NAME.String}}({{ .COLUMN_NAME.String }} {{ castType .DATA_TYPE.String }}) ([]*{{ title $name}}]){
	r, err := db.Query("select * from {{ $name }} where {{ .COLUMN_NAME.String }}=?", {{ .COLUMN_NAME.String}})
	if err != nil {
		// TODO: use our ln package
		log.Print("error in Query: ", err)
	}
	res := make([]*{{ title $name }})
	for r.Next() { {{ range $cols }}
		var {{ .COLUMN_NAME.String }} {{ castType .DATA_TYPE.String }}{{end}}
		err = r.Scan({{ join $cols }})
		if err != nil {
			log.Print("error scanning schema ", err)
		}
		s := &{{ title $name }}{ {{ structInit $cols }} }
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
		conversion = append(conversion, col.COLUMN_NAME.String+":"+col.COLUMN_NAME.String)
	}
	return strings.Join(conversion, ",")
}

// Helper function to change mysql types to go types
func castType(m string) string {
	switch m {
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
