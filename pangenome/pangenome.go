package pangenome

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/superphy/prairiedog/kmers"
)

type Graph struct {
	dg *dgo.Dgraph
	bd *badger.DB
	K  int
}

type KmerNode struct {
	Uid          string     `json:"uid,omitempty"`
	Sequence     string     `json:"sequence,omitempty"`
	ForwardNodes []KmerNode `json:"forward,omitempty"`
	ReverseNodes []KmerNode `json:"reverse,omitempty"`
}

var Schema = `
	sequence: string @index(term) .
`

func NewGraph() *Graph {
	g := &Graph{
		K: 11,
	}
	// Create a connection to Dgraph.
	g.dg = setupDgraph("localhost", "9080", Schema)
	// Create a connection to Badger.
	g.bd = setupBadger()
	return g
}

func (g *Graph) CreateNode(seq string) *api.Assigned {
	ctx := context.Background()

	node := KmerNode{
		Sequence: seq,
	}

	nb, err := json.Marshal(node)
	if err != nil {
		log.Fatal(err)
	}

	mu := &api.Mutation{
		CommitNow: true,
	}

	mu.SetJson = nb
	assigned, err := g.dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}

	// Return the UID assigned by Dgraph.
	return assigned
}

func Run() {
	// Databases.
	g := NewGraph()
	defer g.bd.Close()

	// Load a genome file.
	km := kmers.New("testdata/ECI-2523.fsa")

	h, seq := km.Next()
	fmt.Println(h, seq)

	g.CreateNode(seq)

	// TODO: remove
	_ = g
	_ = km
}
