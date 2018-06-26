package prairiedog

import (
	"github.com/dgraph-io/dgo"
	"github.com/superphy/prairiedog/kmers"
)

type Graph struct {
	dg *dgo.Dgraph
	k  int
}

func NewGraph() *Graph {
	g := &Graph{}
	g.dg = newClient("localhost", "9080")
	g.k = 11
	// Sets up the database.
	setup(g.dg)
	return g
}

func main() {
	g := NewGraph()
	km := kmers.New("somefile.fasta")

	// TODO: remove
	_ = g
	_ = km
}
