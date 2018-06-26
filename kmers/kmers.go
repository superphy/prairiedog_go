package kmers

import (
	"bufio"
	"fmt"
	"os"
)

type Kmers struct {
	src   string
	lines []string
}

func load(km *Kmers) {
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

func New(s string) *Kmers {
	km := &Kmers{
		src: s,
	}
	load(km)
	return km
}

func (km *Kmers) Next() (string, string) {
	return km.lines[0], km.lines[1]
}
