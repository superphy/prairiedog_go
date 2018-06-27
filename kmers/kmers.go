package kmers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Kmers struct {
	src       string
	lines     []string
	Headers   []string // location of all the headers in lines.
	Sequences []string // location of all the sequences in lines.
	li        int      // line index in Headers and Sequences.
	pi        int      // position index in a slice of Sequences.
	K         int
}

func (km *Kmers) load() {
	file, err := os.Open(km.src)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		km.lines = append(km.lines, scanner.Text())
	}
}

func (km *Kmers) index() {
	seq := ""
	for _, line := range km.lines {
		if strings.HasPrefix(line, ">") {
			if seq != "" {
				km.Sequences = append(km.Sequences, seq)
				seq = ""
			}
			km.Headers = append(km.Headers, line)
		} else {
			seq = seq + line
		}
	}
	km.Sequences = append(km.Sequences, seq)
	km.lines = nil
}

func New(s string) *Kmers {
	km := &Kmers{
		src: s,
		li:  0,
		pi:  0,
		K:   11,
	}
	km.load()
	km.index()
	return km
}

func (km *Kmers) Next() (string, string) {
	lastOfSequences := km.li == len(km.Sequences)-1
	endOfSeq := km.pi > len(km.Sequences[km.li])-1+km.K

	// Done.
	if lastOfSequences && endOfSeq {
		return "", ""
	}

	// Move to next sequence.
	if endOfSeq {
		km.li++
		km.pi = 0
	}

	// K is greater than the size of the contig.
	if km.K > len(km.Sequences[km.li])-1 {
		log.Printf("WARNING: contig %s is shorter than the chosen k-value of %v. Skipping contig.", km.Headers[km.li], km.K)
		km.li++
		km.pi = 0
	}

	// Fasta header.
	header := km.Headers[km.li]
	// Slice of the sequence.
	sl := km.Sequences[km.li][km.pi : km.pi+km.K]

	// Increment.
	km.pi += km.K

	return header, sl
}
