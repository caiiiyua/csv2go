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
	csv.Reader
	fd *os.File
}

func (d *Decoder) Read() (record []string, err error) {
	return d.Reader.Read()
}

func (d *Decoder) Close() {
	d.fd.Close()
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
