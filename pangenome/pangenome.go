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
	sequence     string     `json:"sequence,omitempty"`
	ForwardNodes []KmerNode `json:"forward,omitempty"`
	ReverseNodes []KmerNode `json:"reverse,omitempty"`
}

// NewGraph is the main setup for backends.
func NewGraph() *Graph {
	log.Println("Starting NewGraph().")
	g := &Graph{
		K: 11,
	}
	// Create a connection to Dgraph.
	g.dg, _ = setupDgraph("localhost", "9080")
	log.Println("Dgraph connected OK.")
	// Create a connection to Badger.
	g.bd = setupBadger()
	log.Println("Badger connected OK.")
	return g
}

// DropAll discards everything in Dgraph.
func (g *Graph) DropAll(contextMain context.Context) (bool, error) {
	ctx, cancel := context.WithCancel(contextMain)
	defer cancel()

	err := g.dg.Alter(ctx, &api.Operation{DropAll: true})
	if err != nil {
		return false, err
	}

	// Ensure schema is still setup after dropping.
	setupSchema(g.dg)

	return true, nil
}

// SetKVInt sets the key: value pair in Badger for ints.
func (g *Graph) SetKVInt(key string, value int) (bool, error) {
	err := g.bd.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(strconv.Itoa(value)))
		return err
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// SetKVStr sets the key: value pair in Badger for strings.
func (g *Graph) SetKVStr(key string, value string) (bool, error) {
	err := g.bd.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(value))
		return err
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// SetKV sets the key: value pair in Badger for slices of uin64.
func (g *Graph) SetKVSliceUint64(key string, value []uint64) (bool, error) {
	buf, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	err = g.bd.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), buf)
		return err
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetKVInt get the key: value pair in Badger.
func (g *Graph) GetKVInt(key string) (int, error) {
	var val []byte
	err := g.bd.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, err = item.Value()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return -1, err
	}
	s := string(val[:])
	evaluated, err := strconv.Atoi(s)
	if err != nil {
		return -1, err
	}
	return evaluated, nil
}

// GetKVStr gets the key: value pair in Badger.
func (g *Graph) GetKVStr(key string) (string, error) {
	var val []byte
	err := g.bd.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, err = item.Value()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	s := string(val[:])
	return s, nil
}

// GetKVSliceUint64 gets the key: value pair in Badger.
func (g *Graph) GetKVSliceUint64(key string) ([]uint64, error) {
	var val []byte
	err := g.bd.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, err = item.Value()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var sl []uint64
	err = json.Unmarshal(val, &sl)
	if err != nil {
		return nil, err
	}
	return sl, nil
}

func (g *Graph) CreateNode(seq string, contextMain context.Context) (uint64, error) {
	ctx, cancel := context.WithCancel(contextMain)
	defer cancel()

	node := KmerNode{
		sequence: seq,
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

func (g *Graph) GetNode(seq string, contextMain context.Context) (uint64, error) {
	ctx, cancel := context.WithCancel(contextMain)
	defer cancel()

	txn := g.dg.NewTxn()
	defer txn.Discard(ctx)

	q := fmt.Sprintf(`
		{
			all(func: anyofterms(Sequence, %s)) {
				uid
			}
		}
	`, seq)
	resp, err := txn.Query(ctx, q)
	if err != nil {
		log.Fatal(err)
	}

	var decode struct {
		All []struct {
			Uid string
		}
	}
	if err := json.Unmarshal(resp.GetJson(), &decode); err != nil {
		log.Fatal(err)
	}
	if len(decode.All) == 0 {
		return 0, nil
	}
	i, err := strconv.ParseUint(decode.All[0].Uid, 16, 64)
	if err != nil {
		return 0, nil
	}
	return i, nil
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

	var seq1, seq2, header1 string
	// Initial Kmer.
	header1, seq1 = km.Next()
	// If there exists any kmers left in the genome.
	for km.HasNext() {
		var sl []uint64
		// If there exists any kmers left in the particular contig.
		for km.ContigHasNext() {
			_, seq2 = km.Next()

			uid1, err := g.CreateNode(seq1, ctx)
			if err != nil {
				log.Fatal(err)
				return false, err
			}

			// Always append the first node.
			sl = append(sl, uid1)

			uid2, err := g.CreateNode(seq2, ctx)
			if err != nil {
				log.Fatal(err)
				return false, err
			}

			// If on last kmer in a contig, append the second node.
			if !km.ContigHasNext() {
				sl = append(sl, uid2)
			}

			_, err = g.CreateEdge(uid1, uid2, ctx)
			if err != nil {
				log.Fatal(err)
				return false, err
			}
			seq1 = seq2

		}
		// Store the completed path for the contig.
		g.SetKVSliceUint64(header1, sl)
		// Grab next sequence.
		header1, seq1 = km.Next()
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
