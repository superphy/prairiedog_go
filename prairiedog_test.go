package prairiedog

import (
	"fmt"
	"testing"

	"github.com/superphy/prairiedog/kmers"
)

func TestNew(t *testing.T) {
	NewGraph()
	// Output: nil
}
func TestKmers(t *testing.T) {
	km := kmers.New("testdata/172.fa")
	header, _ := km.Next()
	fmt.Println(header)
	// Output: >gi|1062504329|gb|CP014670.1| Escherichia coli strain CFSAN004177, complete genome
}
