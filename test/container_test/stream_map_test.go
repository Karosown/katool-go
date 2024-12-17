package container

import (
	"fmt"
	"testing"

	"github.com/karosown/katool/container/stream"
)

func Test_Map(t *testing.T) {
	m := map[string]string{
		"a": "b",
		"c": "d",
		"e": "f",
	}
	stream.EntrySet(m).ToStream().ForEach(func(e stream.Entry[string, string]) {
		fmt.Println(e)
	})
}
