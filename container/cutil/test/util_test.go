package test

import (
	"fmt"
	"github.com/karosown/katool-go/container/cutil"
	"testing"
)

func TestUtil(t *testing.T) {
	isNil := cutil.IsBlank("")
	fmt.Println(isNil)
}
