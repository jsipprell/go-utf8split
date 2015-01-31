package utf8split

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func helperStringSplitter(w io.Writer, t *testing.T) {
	splitter := WithDelimiters("\t🙌\n ")
	for i, j := range splitter.SplitString("a   b c dd_d  a\tb e🙌emoji🙌e") {
		if t != nil {
			t.Logf("string %d ---> %v", i, j)
		}
		if w != nil {
			fmt.Fprintf(w, "string %d ---> %v\n", i, j)
		}
	}
}

func ExampleStringSplitter() {
	buf := &bytes.Buffer{}
	helperStringSplitter(buf, nil)
	os.Stdout.Write(buf.Bytes())
	// Output:
	// string 0 ---> a
	// string 1 ---> b
	// string 2 ---> c
	// string 3 ---> dd_d
	// string 4 ---> a
	// string 5 ---> b
	// string 6 ---> e
	// string 7 ---> emoji
	// string 8 ---> e
}

func TestStringSplitter(t *testing.T) {
	helperStringSplitter(nil, t)
}
