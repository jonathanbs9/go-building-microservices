package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := Product{
		Name:  "Jonathan",
		Price: 1000,
		SKU:   "aet-asd-ggg",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal("<> |\n \nError | " + err.Error() + " | \n \n")
	}
}
