package csv2go

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"
)

type Account struct {
	CardId              string
	Name, Mobile        string
	CreateDate          time.Time
	SaveAmt, ConsumeAmt float64
	LeftAmt             float64
	Items               []string
}

type Decoder struct {
	reader     *csv.Reader
	header     csvHeader
	lastRecord []string
	trim       string
	csvReader  io.ReadCloser
	bHandler   funcBoolean
}

type funcBoolean func(string) bool

// setBooleanHandler custom boolean handler for fields
func (d *Decoder) setBooleanHandler(f funcBoolean) {
	d.bHandler = f
}

// DoBoolean return boolean of this field
func (d *Decoder) DoBoolean(field string) bool {
	if d.bHandler != nil {
		return d.bHandler()
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

func (d *Decoder) decode(i interface{}) {
	for {
		record, err := d.Read()
		if err == io.EOF {
			fmt.Println("Read out of file")
			break
		} else if err != nil {
			return
		}
		for i, field := range record {
			fmt.Println(i, field)
		}
	}
}

func Test() {
	Test2(Account{})
}

func Test2(i interface{}) {
	fmt.Println("Test")
	typ := reflect.TypeOf(i)
	fmt.Println("type of account:", typ, "contains", typ.NumField(), "fields", "with kind", typ.Kind())
	for i := 0; i < typ.NumField(); i++ {
		st := typ.Field(i)
		fmt.Println("the", i, "field is", st)
	}

}

func NewDecoder(file string) *Decoder {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Open", file, "failed")
		return nil
	}
	decoder := new(Decoder)
	decoder.fd = f
	decoder.Reader = csv.NewReader(f)
	return decoder
}
