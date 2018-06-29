package pangenome

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

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
	UID          uint64     `json:"uid,omitempty"`
	Sequence     string     `json:"sequence,omitempty"`
	ForwardNodes []KmerNode `json:"forward,omitempty"`
	ReverseNodes []KmerNode `json:"reverse,omitempty"`
}

var Schema = `
	sequence: string @index(term) .
`

func NewGraph() *Graph {
	log.Println("Starting NewGraph().")
	g := &Graph{
		K: 11,
	}
	// Create a connection to Dgraph.
	g.dg = setupDgraph("localhost", "9080", Schema)
	log.Println("Dgraph connected OK.")
	// Create a connection to Badger.
	g.bd = setupBadger()
	log.Println("Badger connected OK.")
	return g
}

func (g *Graph) CreateNode(seq string) (uint64, error) {
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
	uid, err := strconv.ParseUint(assigned.Uids["blank-0"][2:], 16, 64)
	if err != nil {
		log.Fatal(err)
	}
	return uid, err
}

func (g *Graph) CreateEdge(src uint64, dst uint64) (*api.Assigned, error) {
	ctx := context.Background()

	srcNode := KmerNode{
		UID: src,
		ForwardNodes: []KmerNode{{
			UID: dst,
		}},
	}

	nb, err := json.Marshal(srcNode)
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
	return assigned, err
}

// CreateAll Nodes+Edges for all kmers in km.
func (g *Graph) CreateAll(km *kmers.Kmers) (bool, error) {
	var seq1, seq2 string
	_, seq1 = km.Next()
	for km.HasNext() {
		for km.ContigHasNext() {
			_, seq2 = km.Next()
			uid1, err := g.CreateNode(seq1)
			if err != nil {
				log.Fatal(err)
				return false, err
			}
			uid2, err := g.CreateNode(seq2)
			if err != nil {
				log.Fatal(err)
				return false, err
			}
			_, err = g.CreateEdge(uid1, uid2)
			if err != nil {
				log.Fatal(err)
				return false, err
			}
		}
	}
	return true, nil
}

func Run() {
	// Databases.
	g := NewGraph()
	defer g.Close()

	// Load a genome file.
	km := kmers.New("testdata/ECI-2523.fsa")

	h, seq := km.Next()
	fmt.Println(h, seq)

	g.CreateNode(seq)

	// TODO: remove
	_ = g
	_ = km
}

// Close handles teardown.
func (g *Graph) Close() {
	defer g.bd.Close()
}
