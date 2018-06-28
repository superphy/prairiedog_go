package pangenome

import (
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/dgo"
	"github.com/superphy/prairiedog/kmers"
)

type Graph struct {
	dg *dgo.Dgraph
	bd *badger.DB
	K  int
}

func NewGraph() *Graph {
	g := &Graph{
		K: 11,
	}
	// Create a connection to Dgraph.
	g.dg = setupDgraph("localhost", "9080")
	// Create a connection to Badger.
	g.bd = setupBadger()
	return g
}

func Run() {
	// Databases.
	g := NewGraph()
	defer g.bd.Close()

	// Load a genome file.
	km := kmers.New("testdata/ECI-2523.fsa")

	h, s := km.Next()
	fmt.Println(h, s)

	// TODO: remove
	_ = g
	_ = km
}
