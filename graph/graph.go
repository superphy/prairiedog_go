package graph

import (
	"log"
	"os"
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
		g.Address = s
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

// Graph is the main struct.
type Graph struct {
	inputs string
}

// New creates a new graph.
func New(options ...OptionFn) (*Graph, error) {
	g := &Graph{}

	for _, optionFn := range options {
		optionFn(g)
	}

	return g, nil
}

// Run is what's carried out by the program.
func (g *Graph) Run() {
	log.Printf("Run started.")

	log.Printf("Run stopped.")
}
