package csv2go

import (
	"encoding/csv"
	"fmt"
	"io"
)

type Decoder struct {
	reader     *csv.Reader
	header     csvHeader
	lastRecord []string
	trim       string
	csvReader  io.ReadCloser
	bHandler   funcBoolean
}

type csvHeader map[int]string
type funcBoolean func(string) bool

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
	for {
		record, err := d.Read()
		if err == io.EOF {
			fmt.Println("Read out of file")
			break
		} else if err != nil {
			return err
		}
		for i, field := range record {
			fmt.Println(i, field)
		}
	}
	return err
}

func NewDecoder(in io.ReadCloser) *Decoder {
	decoder := &Decoder{reader: csv.NewReader(in)}
	decoder.csvReader = in
	return decoder
}
