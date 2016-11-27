# ACA
ACA is a AC automation implementation for [Go](https://golang.org). 

# Documentation
Documentation can be found at [Godoc](https://godoc.org/github.com/cosiner/aca)

# Example
```Go
func TestACA(t *testing.T) {
	var aca ACA

	aca.Add("say", "she", "he", "her", "shr").Build()

	if !aca.HasContainedIn("yasherhs") {
		t.Fatal("ACA check contained in failed")
	}

	if !reflect.DeepEqual([]string{"she", "he", "her"}, aca.Match("yasherhs")) {
		t.Fatal("ACA match failed")
	}

	if aca.Replace("yasherhs", '*', nil) != "ya****hs" {
		t.Fatal("ACA replace failed")
	}

	if aca.Replace("-y|a}s+h=e)r(h&s", '*', NewRuneSet("-|}+=)(&")) != "-y|a}*+*=*)*(h&s" {
		t.Fatal("ACA replace with skips failed")
	}
}
```

# LICENSE
MIT.
