# ACA
ACA is a AC automation implementation for [Go](https://golang.org). 

# Documentation
Documentation can be found at [Godoc](https://godoc.org/github.com/cosiner/aca)

# Example
```Go

func TestACA(t *testing.T) {
	aca := New(NewSkipsCleaner([]rune("-|}+=)(&")), NewIgnoreCaseCleaner())
	aca.Add("say", "she", "he", "her", "shr").Build()

	if !aca.HasContainedIn("yaSherhs") {
		t.Fatal("ACA check contained in failed")
	}

	if !reflect.DeepEqual([]string{"ShE", "hE", "hEr"}, aca.Match("yaShErhs")) {
		t.Fatal("ACA match failed")
	}

	var replacement = '*'
	if aca.Replace("yasherhs", replacement) != "ya****hs" {
		t.Fatal("ACA replace failed")
	}

	if aca.Replace("-y|a}s+h=e)r(h&s", replacement) != "-y|a}*+*=*)*(h&s" {
		t.Fatal("ACA replace with skipsCleaner failed")
	}
}
```

# LICENSE
MIT.
