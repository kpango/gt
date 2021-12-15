package gt

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	Align  int
	Width  uint64
	Height uint64

	Table struct {
		Lines     Lines
		header    Header
		rows      []Row
		emptyText string
		output    string
	}

	Header Row
	Row    struct {
		width   Width
		height  Height
		Columns []Column
	}
	Column struct {
		header     string
		index      int
		width      Width
		height     Height
		value      string
		align      *Align
		innerTable *Table
	}
	Columns []Column
	Lines   struct {
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
	fmt.Println(new(Table).parse(obj).format().toString())
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

func (t *Table) toString() string {
	return ""
}

func (t *Table) format() *Table {
	if t == nil || t.rows == nil {
		return t
	}
	mRHs := make([]Height, len(t.rows)+1)        // max rows height
	mCWs := make([]Width, len(t.header.Columns)) // max columns width
	for x, c := range t.header.Columns {
		if mRHs[0] < c.height {
			mRHs[0] = c.height
		}
		if mCWs[x] < c.width {
			mCWs[x] = c.width
		}
		c.innerTable.format()
	}

	t.header.height = mRHs[0]
	for y, row := range t.rows {
		if len(mCWs) == 0 {
			mCWs = make([]Width, len(row.Columns)) // max columns width
		}
		for x, c := range row.Columns {
			t.rows[y].Columns[x].innerTable.format()
			if mRHs[y+1] < c.height {
				mRHs[y+1] = c.height
			}
			if mCWs[x] < c.width {
				mCWs[x] = c.width
			}
		}
		t.rows[y].height = mRHs[y+1]
	}

	for x := range t.header.Columns {
		t.header.Columns[x].height = mRHs[0]
		t.header.Columns[x].width = mCWs[x]
	}

	for y, row := range t.rows {
		for x := range row.Columns {
			t.rows[y].Columns[x].height = mRHs[y+1]
			t.rows[y].Columns[x].width = mCWs[x]
		}
		t.rows[y].height = mRHs[y+1]
	}
	return t
}

func (t *Table) parseRow(rt reflect.Type, rv reflect.Value) Row {
	row := Row{
		height:  1,
		width:   0,
		Columns: make([]Column, 0, rt.NumField()),
	}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tag := field.Tag.Get(tagKey)
		value := rv.Field(i)
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
			row.Columns = append(row.Columns, Columns(ir.Columns).Serialize())
		default:
			val, err := IFtoa(value.Interface())
			if err != nil {
				continue
			}
			str := strings.ReplaceAll(val, "\r\n", "\n")
			row.Columns = append(row.Columns, Column{
				header:     tag,
				index:      i,
				value:      str,
				width:      Width(len(str)),
				height:     Height(1 + strings.Count(str, "\n")),
				align:      nil,
				innerTable: nil,
			})
		}
	}
	fmt.Printf("%+v", row)
	return row
}

func (c Columns) Serialize() Column {
	for _, pc := range c {
		pc.innerTable.format()
	}
	col := Column{
		// header  string
		// index   int
		// width   Width
		// height  Height
		// align   Align
		// value   string
		// Columns []Column
	}
	return col
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
