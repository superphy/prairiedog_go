package pangenome

import (
	"fmt"
)

func ExampleBadger() {
	bd := setupBadger()
	defer bd.Close()
	s := bd.Tables()
	fmt.Println(s)
	// Output:
	// []
}

func ExampleDgraph() {
	var schema = `
		sequence: string @index(term) .
	`
	_, err := setupDgraph("localhost", "9080", schema)
	fmt.Println(err)
	// Output:
	// <nil>
}
