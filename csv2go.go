package csv2go

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Decoder struct {
	reader     *csv.Reader
	header     csvHeader
	lastRecord []string
	trim       string
	csvReader  io.ReadCloser
	bHandler   funcBoolean
	records    int
}

type csvHeader map[string]string
type Header map[string]int

type funcBoolean func(string) bool

var header csvHeader
var origHeader Header
var cols map[int]string

var headerCache map[string]structField

type structField reflect.StructField

// Comma set record delimiter, default is ','
func (d *Decoder) Comma(s rune) *Decoder {
	d.reader.Comma = s
	return d
}

// Comment set comment character for start of line
func (d *Decoder) Comment(s rune) *Decoder {
	d.reader.Comment = s
	return d
}

// FieldsPerRecord set number of expect fields per record
func (d *Decoder) FieldsPerRecord(n int) *Decoder {
	d.reader.FieldsPerRecord = n
	return d
}

// LazyQuotes set allow lazy quotes
func (d *Decoder) LazyQuotes(b bool) *Decoder {
	d.reader.LazyQuotes = b
	return d
}

// TrailingComma ignored
func (d *Decoder) TrailingComma(b bool) *Decoder {
	d.reader.TrailingComma = b
	return d
}

// TrimLeadingSpace set trim leading space
func (d *Decoder) TrimLeadingSpace(b bool) *Decoder {
	d.reader.TrimLeadingSpace = b
	return d
}

// setBooleanHandler custom boolean handler for fields
func (d *Decoder) setBooleanHandler(f funcBoolean) *Decoder {
	d.bHandler = f
	return d
}

// DoBoolean return boolean of this field
func (d *Decoder) DoBoolean(field string) bool {
	if d.bHandler != nil {
		return d.bHandler(field)
	}
	if len(field) > 0 {
		return true
	}
	return false
}

// Read() read next record of csv file and take last record
func (d *Decoder) Read() (record []string, err error) {
	d.lastRecord, err = d.reader.Read()
	return d.lastRecord, err
}

func (d *Decoder) Close() error {
	if d.csvReader != nil {
		return d.csvReader.Close()
	}
	return nil
}

// return column value according col name
func (d *Decoder) colByName(col string) string {
	idx := origHeader[col]
	return d.lastRecord[idx]

}

func (d *Decoder) structParser(name string, v string) interface{} {
	var value interface{}
	switch name {
	case "time.Time":
		var date time.Time
		date, _ = time.Parse("2006/1/2", v)
		value = date
		fmt.Println("struct parsed date", date)
	}
	fmt.Println("struct parsed:", value, "on", v)
	return value
}

func (d *Decoder) setValue(v reflect.Value, field, vstr string) error {
	var err error
	f := v.FieldByName(field)
	typ := headerCache[field]
	// fmt.Println(field, typ.Type, typ.Type.Kind(), typ.Name)
	if f.CanSet() {
		// fmt.Println(field, typ.Type)
		switch typ.Type.Kind() {
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			var value float64
			value, err = strconv.ParseFloat(vstr, 64)
			if err != nil {
				value = 0
			}
			f.SetFloat(value)
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			var value int64
			value, err = strconv.ParseInt(vstr, 10, 64)
			if err != nil {
				value = 0
			}
			f.SetInt(value)
		case reflect.Struct:
			f.Set(reflect.ValueOf(d.structParser(fmt.Sprintf("%s", typ.Type), vstr)))
		case reflect.String:
			fmt.Println(field, typ.Type, vstr)
			f.SetString(vstr)
		default:
			fmt.Println(field, typ.Type)
			f.SetString(vstr)
		}
	}

	return err
}

// Decode next record of csv lines
func (d *Decoder) Decode(i interface{}) error {
	var err error
	record, err := d.Read()
	d.lastRecord = record
	if err == io.EOF {
		fmt.Println("Read out of file")
		return err
	} else if err != nil {
		return err
	}
	if d.records == 0 {
		initHeader(record)
		initCsvHeader(i)
		d.header = header
		fmt.Println(d.header)
	} else {
		value := reflect.ValueOf(i).Elem()
		for field, col := range d.header {
			colValue := d.colByName(col)
			d.setValue(value, field, colValue)
			// fmt.Println(field, colValue)
		}
	}
	// for i, field := range record {
	// 	fmt.Println(i, field)
	// }
	d.records += 1
	return err
}

// skip filed of struct contains '-'
func skip(tag reflect.StructTag) bool {
	return strings.HasPrefix(field(tag), "-")
}

func field(tag reflect.StructTag) string {
	return tag.Get("csv")
}

func initCsvHeader(i interface{}) error {
	var err error
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		fmt.Println("kind of i:", t.Kind())
		return err
	}
	j := 0
	custome := false
	for i := 0; i < t.NumField(); i++ {
		headerCache[t.Field(i).Name] = structField(t.Field(i))
		// fmt.Println(t.Field(i))
		// skip any field unexpected
		if skip(t.Field(i).Tag) {
			continue
		} else if len(field(t.Field(i).Tag)) > 0 {
			if !custome {
				for k := range header {
					delete(header, k)
				}
				custome = true
			}
			header[t.Field(i).Name] = field((t.Field(i).Tag))

		} else {
			header[t.Field(i).Name] = t.Field(i).Name
		}
		j += 1
	}
	return err
}

func initHeader(r []string) {
	if origHeader == nil {
		origHeader = make(Header)
	}
	for i, field := range r {
		origHeader[field] = i
		cols[i] = field
	}
	fmt.Println("CSV Header:", origHeader)
}

// NewDecoder get a Decoder for use
func NewDecoder(in io.ReadCloser) *Decoder {
	decoder := &Decoder{reader: csv.NewReader(in)}
	decoder.csvReader = in
	decoder.records = 0
	return decoder
}

func init() {
	origHeader = make(Header)
	header = make(csvHeader)
	cols = make(map[int]string)
	headerCache = make(map[string]structField)
}
