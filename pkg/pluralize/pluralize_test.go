package pluralize

import "testing"

func TestPluralize(t *testing.T) {
	slice1 := []string{}
	slice2 := []string{"1"}
	slice3 := []string{"1", "1", "1", "1", "1"}

	word := "car"

	test1 := Pluralize(len(slice1), word)
	test2 := Pluralize(len(slice2), word)
	test3 := Pluralize(len(slice3), word)

	if test1 != "cars" {
		t.Errorf("Pluralize was incorrect, got: %s, want: %s.", test1, "cars")
	}

	if test2 != "car" {
		t.Errorf("Pluralize was incorrect, got: %s, want: %s.", test2, "car")
	}

	if test3 != "cars" {
		t.Errorf("Pluralize was incorrect, got: %s, want: %s.", test3, "cars")
	}
}

func TestPluralizeWithArg(t *testing.T) {
	slice1 := []string{}
	slice2 := []string{"1"}
	slice3 := []string{"1", "1", "1", "1", "1"}

	word := "class"

	test1 := Pluralize(len(slice1), word, "es")
	test2 := Pluralize(len(slice2), word, "es")
	test3 := Pluralize(len(slice3), word, "es")

	if test1 != "classes" {
		t.Errorf("Pluralize was incorrect, got: %s, want: %s.", test1, "classes")
	}

	if test2 != "class" {
		t.Errorf("Pluralize was incorrect, got: %s, want: %s.", test2, "class")
	}

	if test3 != "classes" {
		t.Errorf("Pluralize was incorrect, got: %s, want: %s.", test3, "classes")
	}
}

func TestPluralizeWithArgSingularPrural(t *testing.T) {
	slice1 := []string{}
	slice2 := []string{"1"}
	slice3 := []string{"1", "1", "1", "1", "1"}

	word := "cand"

	test1 := Pluralize(len(slice1), word, "y,ies")
	test2 := Pluralize(len(slice2), word, "y,ies")
	test3 := Pluralize(len(slice3), word, "y,ies")

	if test1 != "candies" {
		t.Errorf("Pluralize was incorrect, got: %s, want: %s.", test1, "candies")
	}

	if test2 != "candy" {
		t.Errorf("Pluralize was incorrect, got: %s, want: %s.", test2, "candy")
	}

	if test3 != "candies" {
		t.Errorf("Pluralize was incorrect, got: %s, want: %s.", test3, "candies")
	}
}
