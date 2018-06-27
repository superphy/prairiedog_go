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

	seq := make([]byte, 0)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(s, ">") {
			if len(seq) != 0 {
				km.Sequences = append(km.Sequences, string(seq))
				seq = nil
			}
			km.Headers = append(km.Headers, s)
		} else {
			seq = append(seq, []byte(s)...)
		}
	}
	km.Sequences = append(km.Sequences, string(seq))
}

func New(s string) *Kmers {
	km := &Kmers{
		src: s,
		li:  0,
		pi:  0,
		K:   11,
	}
	km.load()
	return km
}

func (km *Kmers) Next() (string, string) {
	lastOfSequences := km.li == len(km.Sequences)-1
	endOfSeq := km.pi+km.K > len(km.Sequences[km.li])

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
	// log.Printf("%v, %v, %v, %v", km.li, km.pi, lastOfSequences, endOfSeq)
	sl := km.Sequences[km.li][km.pi : km.pi+km.K]

	// Increment.
	km.pi++

	return header, sl
}
