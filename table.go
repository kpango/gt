package gt

import (
	"reflect"
	"strings"
)

type (
	Align  int
	Width  uint64
	Height uint64

	Table struct {
		header    Header
		Lines     Lines
		rows      []Row
		EmptyText string
		output    string
	}

	Header Row
	Row    struct {
		width   Width
		height  Height
		Columns []Colmun
	}
	Colmun struct {
		header  string
		index   int
		width   Width
		height  Height
		align   Align
		value   string
		Columns []Colmun
	}
	Lines struct {
		enabled bool
		Top     string // Top border
		Left    string // Left border
		Right   string // Right border
		Bottom  string // Bottom border
		MVB     string // Middle vertical border
		MHB     string // Middle horizontal border
	}
)

const (
	_ Align = iota
	Center
	Right
	Left

	tagKey = "table"
)

func Print(obj interface{}) {
	(&Table{}).parse(obj)
}

func (t *Table) parse(data interface{}) *Table {
	rt := reflect.TypeOf(data)
	rv := reflect.ValueOf(data)
	switch rt.Kind() {
	case reflect.Ptr:
		return t.parse(rv.Elem().Interface())
	case reflect.Array, reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			v := rv.Index(i)
			if v.Kind() == reflect.Ptr {
				if v.IsNil() {
					continue
				}
				v = v.Elem()
			}
			t.rows = append(t.rows, t.parseRow(v.Type(), v))
		}
	case reflect.Struct:
		t.rows = append(t.rows, t.parseRow(rt, rv))
	}
	return t
}

func (t *Table) print() string {
	return ""
}

func (t *Table) format() *Table {
	var mHH Height = 0                           // minimum header height
	mCWs := make([]Width, len(t.header.Columns)) // minimum column widths
	for x, c := range t.header.Columns {
		if mHH < c.height {
			mHH = c.height
		}
		if mCWs[x] < c.width {
			mCWs[x] = c.width
		}
	}
	for y, row := range t.rows {
		for x, col := range row.Columns {

		}
	}
}

func (t *Table) parseRow(rt reflect.Type, rv reflect.Value) Row {
	row := Row{
		height:  1,
		width:   0,
		Columns: make([]Colmun, 0, rt.NumField()),
	}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag := field.Tag.Get(tagKey)
		name := field.Name
		value := rv.Field(i)
		typ := field.Type.Name()
		if field.Type.Kind() == reflect.Ptr {
			if value.IsNil() {
				continue
			}
			value = value.Elem()
		}
		switch field.Type.Kind() {
		case reflect.Struct:
			ir := t.parseRow(field.Type, value)
			if ir.height > row.height {
				row.height = ir.height
			}
			row.Columns = append(row.Columns, ir.ToColumn())
		default:
			str := strings.ReplaceAll(IFtoa(value.Interface()), "\r\n", "\n")
			row.Columns = append(row.Columns, Colmun{
				value:  str,
				height: 1 + strings.Count(str, "\n"),
			})
		}
		// if tag != "" {
		// 	if strings.Contains(tag, ",") {
		// 		parts := strings.Split(tag, ",")
		// 		header = parts[0]
		// 		if parts[1] == "" {
		// 			width = 0
		// 		} else {
		// 			width, err = strconv.Atoi(parts[1])
		// 			if err != nil {
		// 				panic(fmt.Errorf("width is not an integer: %s", parts[1], err.Error()))
		// 			}
		// 		}
		//
		// 		// 2: format
		// 		if len(parts) > 2 {
		// 			format = parts[2]
		// 		} else {
		// 			format = ""
		// 		}
		// 	} else {
		// 		header = tag
		// 		width = 0
		// 		format = ""
		// 	}
		//
		// 	if header == "" || header == "-" {
		// 		header = field.Name
		// 	}
		//
		// 	t.headers = append(t.headers, header)
		// 	t.minwidths = append(t.minwidths, width)
		// 	t.formats = append(t.formats, format)
		// } else {
		//
		// }
		// fmt.Println(tag, name, value, typ)
	}

	return Row{}
}

func (r Row) ToColumn() Colmun {
	return r.Columns[0]
}

func rawObject(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// func (t *Table) Print() {
// 	fullwidth := 0
// 	fullformat := ""
//
// 	realwidths := make([]int, len(t.headers))
//
// 	// calc realwidths for each column
// 	for _, row := range t.rows {
// 		for i, v := range row {
// 			if realwidths[i] < len(v) {
// 				realwidths[i] = len(v)
// 			}
// 			if realwidths[i] < len(t.headers[i]) {
// 				realwidths[i] = len(t.headers[i])
// 			}
// 			if realwidths[i] < t.minwidths[i] {
// 				realwidths[i] = t.minwidths[i]
// 			}
// 		}
// 	}
//
// 	for i, w := range realwidths {
// 		margin := t.Margin
// 		if i == len(realwidths)-1 {
// 			margin = 0 // no margin for last column
// 		}
// 		fullformat += "%-" + strconv.FormatInt(int64(w), 10) + "s" + strings.Repeat(" ", margin)
// 		fullwidth += w + margin
// 	}
// 	fullformat += "\n"
//
// 	// print headers
// 	fmt.Printf(fullformat, asInterfaces(t.headers)...)
//
// 	// print split line
// 	if t.SplitLine {
// 		fmt.Println(strings.Repeat("-", int(fullwidth)))
// 	}
//
// 	// print rows
// 	for _, row := range t.rows {
// 		fmt.Printf(fullformat, asInterfaces(row)...)
// 	}
//
// 	// print empty text if no rows
// 	if len(t.rows) == 0 && len(t.EmptyText) > 0 {
// 		fmt.Println(t.EmptyText)
// 	}
// }
//
// func asInterfaces(list []string) []interface{} {
// 	vals := make([]interface{}, len(list))
// 	for i, v := range list {
// 		vals[i] = v
// 	}
// 	return vals
// }
//

//
// // Output formats slice of structs data and writes to standard output.(Using box drawing characters)
// func Output(slice interface{}) {
// 	fmt.Println(Table(slice))
// }
//
// // OutputA formats slice of structs data and writes to standard output.(Using standard ascii characters)
// func OutputA(slice interface{}) {
// 	fmt.Println(AsciiTable(slice))
// }
//
// // Table formats slice of structs data and returns the resulting string.(Using box drawing characters)
// func Table(slice interface{}) string {
// 	coln, colw, rows, err := parse(slice)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	table := table(coln, colw, rows, m["box-drawing"])
// 	return table
// }
//
// // AsciiTable formats slice of structs data and returns the resulting string.(Using standard ascii characters)
// func AsciiTable(slice interface{}) string {
// 	coln, colw, rows, err := parse(slice)
// 	if err != nil {
// 		return err.Error()
// 	}
// 	table := table(coln, colw, rows, m["ascii"])
// 	return table
// }
//
// func parse(slice interface{}) (
// 	coln []string, // name of columns
// 	colw []int, // width of columns
// 	rows [][]string, // rows of content
// 	err error,
// ) {
//
// 	s, err := sliceconv(slice)
// 	if err != nil {
// 		return
// 	}
// 	for i, u := range s {
// 		v := reflect.ValueOf(u)
// 		t := reflect.TypeOf(u)
// 		if v.Kind() == reflect.Ptr {
// 			v = v.Elem()
// 			t = t.Elem()
// 		}
// 		if v.Kind() != reflect.Struct {
// 			err = errors.New("warning: table: items of slice should be on struct value")
// 			return
// 		}
// 		var row []string
//
// 		m := 0 // count of unexported field
// 		for n := 0; n < v.NumField(); n++ {
// 			if t.Field(n).PkgPath != "" {
// 				m++
// 				continue
// 			}
// 			cn := t.Field(n).Name
// 			ct := t.Field(n).Tag.Get("table")
// 			if ct == "" {
// 				ct = cn
// 			}
// 			cv := fmt.Sprintf("%+v", v.FieldByName(cn).Interface())
//
// 			if i == 0 {
// 				coln = append(coln, ct)
// 				colw = append(colw, len(ct))
// 			}
// 			if colw[n-m] < len(cv) {
// 				colw[n-m] = len(cv)
// 			}
//
// 			row = append(row, cv)
// 		}
// 		rows = append(rows, row)
// 	}
// 	return coln, colw, rows, nil
// }
//
// func table(coln []string, colw []int, rows [][]string, b bd) (table string) {
// 	if len(rows) == 0 {
// 		return ""
// 	}
// 	head := [][]rune{{b.DR}, {b.V}, {b.VR}}
// 	bttm := []rune{b.UR}
// 	for i, v := range colw {
// 		head[0] = append(head[0], []rune(repeat(v+2, b.H)+string(b.HD))...)
// 		head[1] = append(head[1], []rune(" "+coln[i]+repeat(v-StringLength([]rune(coln[i]))+1, ' ')+string(b.V))...)
// 		head[2] = append(head[2], []rune(repeat(v+2, b.H)+string(b.VH))...)
// 		bttm = append(bttm, []rune(repeat(v+2, b.H)+string(b.HU))...)
// 	}
// 	head[0][len(head[0])-1] = b.DL
// 	head[2][len(head[2])-1] = b.VL
// 	bttm[len(bttm)-1] = b.UL
//
// 	var body [][]rune
// 	for _, r := range rows {
// 		row := []rune{b.V}
// 		for i, v := range colw {
// 			// handle non-ascii character
// 			l := StringLength([]rune(r[i]))
//
// 			row = append(row, []rune(" "+r[i]+repeat(v-l+1, ' ')+string(b.V))...)
// 		}
// 		body = append(body, row)
// 	}
//
// 	for _, v := range head {
// 		table += string(v) + "\n"
// 	}
// 	for _, v := range body {
// 		table += string(v) + "\n"
// 	}
// 	table += string(bttm)
// 	return table
// }
//
// func sliceconv(slice interface{}) ([]interface{}, error) {
// 	v := reflect.ValueOf(slice)
// 	if v.Kind() != reflect.Slice {
// 		return nil, errors.New("warning: sliceconv: param \"slice\" should be on slice value")
// 	}
//
// 	l := v.Len()
// 	r := make([]interface{}, l)
// 	for i := 0; i < l; i++ {
// 		r[i] = v.Index(i).Interface()
// 	}
// 	return r, nil
// }
//
// func repeat(time int, char rune) string {
// 	var s = make([]rune, time)
// 	for i := range s {
// 		s[i] = char
// 	}
// 	return string(s)
// }
//
// func StringTrueLength(r ...rune) (l int) {
// 	l = len(r)
// 	for _, v := range r {
// 		switch {
// 		// http://www.tamasoft.co.jp/en/general-info/unicode.html
// 		case v >= 0x2E80 && v <= 0x9FD0,
// 			v >= 0xAC00 && v <= 0xD7A3,
// 			v >= 0xF900 && v <= 0xFACE,
// 			v >= 0xFE00 && v <= 0xFE6C,
// 			v >= 0xFF00 && v <= 0xFF60,
// 			v >= 0x20000 && v <= 0x2FA1D:
// 			l++
// 		}
// 	}
// 	return l
// }
