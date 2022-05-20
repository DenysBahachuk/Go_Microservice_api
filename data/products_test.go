package data

import "testing"

func TestChackValidation(t *testing.T) {
	p := &Product{Name: "coffee", Price: 1.54, SKU: "dfgd-sdfs"}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
