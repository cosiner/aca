# ACA
ACA is a AC automation implementation for [Go](https://golang.org). 

# Documentation
Documentation can be found at [Godoc](https://godoc.org/github.com/cosiner/aca)

# Example
```Go

func TestACA(t *testing.T) {
	strs := []string{"say", "she", "he", "her", "shr"}
	aca := New(
		strs,
		GroupCleaners(
			NewSkipsCleaner([]rune("-|}+=)(&")),
			NewIgnoreCaseCleaner(),
		),
	)

	{
		var contains QueryContainsProcessor
		aca.Process("yaSherhs", &contains)
		if !contains.Result() {
			t.Fatal("ACA check contained in failed")
		}
	}

	{
		var matched QueryMatchedProcessor
		aca.Process("yaShErhs", &matched)
		if !reflect.DeepEqual([]string{"ShE", "hE", "hEr"}, matched.Result()) {
			t.Fatal("ACA match failed")
		}
	}

	var replacement = '*'
	{
		var replace = NewReplaceMatchedHandler(replacement)
		aca.Process("yasherhs", replace)
		if replace.Result() != "ya****hs" {
			t.Fatal("ACA replace failed")
		}
	}

	{
		var replace = NewReplaceMatchedHandler(replacement)
		aca.Process("-y|a}s+h=e)r(h&s", replace)
		if replace.Result() != "-y|a}*+*=*)*(h&s" {
			t.Fatal("ACA replace with skipsCleaner failed")
		}
	}
}
```

# LICENSE
MIT.
