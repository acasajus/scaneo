package main

const (
	scansText = `{{define "scans"}}// DON'T EDIT *** generated by scaneo *** DON'T EDIT //

package {{.PackageName}}

import "database/sql"

{{range .Tokens}}func {{$.Visibility}}can{{title .Name}}(r *sql.Row) (*{{.Name}}, error) {
	var s {{.Name}}
	if err := r.Scan({{range .Fields}}
		&s.{{.Name}},{{end}}
	); err != nil {
		return nil, err
	}
	return &s, nil
}

func {{$.Visibility}}can{{title .Name}}s(rs *sql.Rows) ([]*{{.Name}}, error) {
	structs := make([]*{{.Name}}, 0, 16)
	var err error
	for rs.Next() {
		var s {{.Name}}
		if err = rs.Scan({{range .Fields}}
			&s.{{.Name}},{{end}}
		); err != nil {
			return nil, err
		}
		structs = append(structs, &s)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return structs, nil
}
{{$fLen := len .Fields}}
var (
	insert{{title .Name}}Fields = "({{range $key, $value := .Fields}}{{if $i}}, {{end}}{{$value.Name}}{{end}})"
	select{{title .Name}}Fields = "{{range $key, $value := .Fields}}{{if $i}}, {{end}}{{$value.Name}}{{end}}"
)

{{end}}{{end}}`
)
