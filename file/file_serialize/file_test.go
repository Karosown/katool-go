package file_serialize

import (
	"fmt"
	"testing"

	"github.com/karosown/katool-go/file/file_serialize/csv"
)

type TestCSV struct {
	UserName string `csv:"user_name"`
	Age      int    `csv:"age"`
	Sex      string `csv:"sex"`
}

func Test_CSV(t *testing.T) {
	path, err := csv.ReadPath[TestCSV]("/Users/karos/GolandProjects/katool/resources/test.csv")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(path)
}
