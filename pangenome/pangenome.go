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

var schema = `
	sequence: string @index(term) .
`

// NewGraph is the main setup for backends.
func NewGraph() *Graph {
	log.Println("Starting NewGraph().")
	g := &Graph{
		K: 11,
	}
	// Create a connection to Dgraph.
	g.dg, _ = setupDgraph("localhost", "9080", schema)
	log.Println("Dgraph connected OK.")
	// Create a connection to Badger.
	g.bd = setupBadger()
	log.Println("Badger connected OK.")
	return g
}

// SetKV sets the key: value pair in Badger.
func (g *Graph) SetKV(key string, value int) (bool, error) {
	err := g.bd.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("answer"), []byte("42"))
		return err
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetKV get the key: value pair in Badger.
func (g *Graph) GetKV(key string) (interface{}, error) {
	var val interface{}
	err := g.bd.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("answer"))
		if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}
		fmt.Printf("The answer is: %s\n", val)
		return nil
	})
	if err != nil {
		return false, err
	}
	return val, nil
}

func (g *Graph) CreateNode(seq string, contextMain context.Context) (uint64, error) {
	ctx, cancel := context.WithCancel(contextMain)
	defer cancel()

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

	txn := g.dg.NewTxn()
	defer txn.Discard(ctx)

	assigned, err := txn.Mutate(ctx, mu)
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

func (g *Graph) CreateEdge(src uint64, dst uint64, contextMain context.Context) (*api.Assigned, error) {
	ctx, cancel := context.WithCancel(contextMain)
	defer cancel()

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

	txn := g.dg.NewTxn()
	defer txn.Discard(ctx)

	assigned, err := txn.Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}
	return assigned, err
}

// CreateAll Nodes+Edges for all kmers in km.
func (g *Graph) CreateAll(km *kmers.Kmers, contextMain context.Context) (bool, error) {
	ctx, cancel := context.WithCancel(contextMain)
	defer cancel()

	var seq1, seq2 string
	_, seq1 = km.Next()
	for km.HasNext() {
		log.Println("outerloop")
		for km.ContigHasNext() {
			log.Println("innerloop")
			_, seq2 = km.Next()
			uid1, err := g.CreateNode(seq1, ctx)
			if err != nil {
				log.Fatal(err)
				return false, err
			}
			uid2, err := g.CreateNode(seq2, ctx)
			if err != nil {
				log.Fatal(err)
				return false, err
			}
			_, err = g.CreateEdge(uid1, uid2, ctx)
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
	contextMain, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load a genome file.
	km := kmers.New("testdata/ECI-2523.fsa")

	h, seq := km.Next()
	fmt.Println(h, seq)

	g.CreateNode(seq, contextMain)

	// TODO: remove
	_ = g
	_ = km
}

// Close handles teardown.
func (g *Graph) Close() {
	defer g.bd.Close()
}
