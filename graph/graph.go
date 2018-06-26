package graph

import (
	"context"
	"log"
	"os"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"
)

// OptionFn are for handling command-line arguments.
type OptionFn func(*Graph)

// InputFiles saves the path to the Graph.
func InputFiles(s string) OptionFn {
	return func(g *Graph) {
		g.inputs = s
	}
}

func Address(s string) OptionFn {
	return func(g *Graph) {
		g.address = s
	}
}

func Port(s string) OptionFn {
	return func(g *Graph) {
		g.port = s
	}
}

func LogFile(s string) OptionFn {
	return func(srvr *Graph) {
		f, err := os.OpenFile(s, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}

		log.SetOutput(f)
	}
}

func newClient() *dgo.Dgraph {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

func setup(c *dgo.Dgraph) {
	// Install a schema into dgraph. Accounts have a `name` and a `balance`.
	err := c.Alter(context.Background(), &api.Operation{
		Schema: `
			sequence: string @index(term) .
			count: int .
		`,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Kmer struct {
	mer   [21]byte
	count int
}

// Graph is the main struct.
type Graph struct {
	inputs  string
	port    string
	address string
	client  *dgo.Dgraph
}

// New creates a new graph.
func New(options ...OptionFn) (*Graph, error) {
	g := &Graph{}

	for _, optionFn := range options {
		optionFn(g)
	}
	g.client = newClient()
	setup(g.client)

	return g, nil
}

// Run is what's carried out by the program.
func (g *Graph) Run() {
	log.Printf("Run started.")

	log.Printf("Run stopped.")
}
