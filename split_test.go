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

func helperByteSplitter(w io.Writer, t *testing.T) {
	splitter := New([]byte{226, 153, 150, 44, 226, 153, 158})
	for i, j := range splitter.Split([]byte("a,b c;d♞rook♖castle;4")) {
		if t != nil {
			t.Logf("slice %d ---> %v (%v)", i, j, string(j))
		}
		if w != nil {
			fmt.Fprintf(w, "slice %d ---> %v (%v)\n", i, j, string(j))
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

func ExampleByteSplitter() {
	buf := &bytes.Buffer{}
	helperByteSplitter(buf, nil)
	os.Stdout.Write(buf.Bytes())
	// Output:
	// slice 0 ---> [97] (a)
	// slice 1 ---> [98 32 99 59 100] (b c;d)
	// slice 2 ---> [114 111 111 107] (rook)
	// slice 3 ---> [99 97 115 116 108 101 59 52] (castle;4)
}

func TestStringSplitter(t *testing.T) {
	helperStringSplitter(nil, t)
}

func TestByteSplitter(t *testing.T) {
	helperByteSplitter(nil, t)
}

func TestStandaloneBytes(t *testing.T) {
	src := []byte("a,b c;d♞rook♖castle;4")
	for i, b := range Bytes(src, []byte("♖,;"), []byte{226, 153, 158}) {
		t.Logf("slice %d ---> %v", i, string(b))
	}
}
