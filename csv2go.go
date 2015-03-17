package csv2go

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strings"
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

type funcBoolean func(string) bool

var header csvHeader
var origHeader map[int]string

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

func (d *Decoder) Decode(i interface{}) error {
	var err error
	record, err := d.Read()
	if err == io.EOF {
		fmt.Println("Read out of file")
		return err
	} else if err != nil {
		return err
	}
	if d.records == 0 {
		initHeader(record)
		initCsvHeader(i)
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
	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Field(i))
		header[t.Field(i).Name] = origHeader[i]
	}
	fmt.Println(header)
	return err
}

func initHeader(r []string) {
	if origHeader == nil {
		origHeader = make(map[int]string)
	}
	for i, field := range r {
		origHeader[i] = field
	}
	fmt.Println("CSV Header:", origHeader)
}

func NewDecoder(in io.ReadCloser) *Decoder {
	decoder := &Decoder{reader: csv.NewReader(in)}
	decoder.csvReader = in
	decoder.records = 0
	return decoder
}

func init() {
	origHeader = make(map[int]string)
	header = make(csvHeader)

}
