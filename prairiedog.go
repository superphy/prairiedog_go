package prairiedog

import (
	"log"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/dgo"
	"github.com/superphy/prairiedog/kmers"
)

type Graph struct {
	dg *dgo.Dgraph
	bd *badger.DB
	k  int
}

func (g *Graph) setupBadger() {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions
	opts.Dir = "/tmp/badger"
	opts.ValueDir = "/tmp/badger"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	g.bd = db
}

func NewGraph() *Graph {
	g := &Graph{}
	g.dg = newClient("localhost", "9080")
	g.k = 11
	// Sets up Dgraph.
	setup(g.dg)
	// Sets up Badger.
	g.setupBadger()
	return g
}

func main() {
	// Databases.
	g := NewGraph()
	defer g.bd.Close()

	// Load a genome file.
	km := kmers.New("somefile.fasta")

	// TODO: remove
	_ = g
	_ = km
}
