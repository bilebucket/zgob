package zgob

import (
	"fmt"
	"os"
	"testing"
)

type person struct {
	Name string
	Age  int
}

func TestZgob(t *testing.T) {
	filename := "testsave"
	RegisterTypes(person{})

	persons := []*person{
		&person{
			Name: "Bob",
			Age:  52,
		},
		&person{
			Name: "Alice",
			Age:  42,
		},
	}
	err := Save(persons, filename)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Successfully saved to %s\n", filename)

	loadedPersons := []*person{}
	bytesRead, err := Load(&loadedPersons, filename)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Successfully read %d bytes from %s\n", bytesRead, filename)

	if len(persons) != len(loadedPersons) {
		t.Errorf("Expected loadedPersons length to be %d, got %d instead", len(persons), len(loadedPersons))
		t.FailNow()
	}

	for index, p := range persons {
		lp := loadedPersons[index]
		if p.Name != lp.Name || p.Age != lp.Age {
			t.Errorf("Expected value %v at index %d, got %v instead", p, index, lp)
		}
	}

	if err := os.Remove(filename); err != nil {
		t.Fatal(err)
	}
}
