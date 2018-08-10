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
	_, err := setupDgraph("localhost", "9080")
	fmt.Println(err)
	// Output:
	// <nil>
}
