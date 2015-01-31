go-utf8split
============

A utility package to split utf8 encoded strings or slices into arbitrary
fields.



Usage
-----

 * From ``split_test.go``:

```go
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
```
