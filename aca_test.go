package aca

import (
	"reflect"
	"testing"
)

func TestACA(t *testing.T) {
	aca := New(NewRuneSet("-|}+=)(&"), false)
	aca.Add("say", "she", "he", "her", "shr").Build()

	if !aca.HasContainedIn("yasherhs") {
		t.Fatal("ACA check contained in failed")
	}

	if !reflect.DeepEqual([]string{"she", "he", "her"}, aca.Match("yasherhs")) {
		t.Fatal("ACA match failed")
	}

	var replacement = '*'
	var replaceSkip = false
	if aca.Replace("yasherhs", replacement, replaceSkip) != "ya****hs" {
		t.Fatal("ACA replace failed")
	}

	if aca.Replace("-y|a}s+h=e)r(h&s", replacement, replaceSkip) != "-y|a}*+*=*)*(h&s" {
		t.Fatal("ACA replace with skips failed")
	}

}
