package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type Chain []uint64

type RainbowTable struct {
	height     int
	width      int
	table      []Chain
	alphabet   Alphabet
	hashMethod HashType
}

func CreateRaindowTable(height, width int, a Alphabet, hash HashType) RainbowTable {
	table := make([]Chain, height)
	for i := 0; i < height; i++ {
		table[i] = make(Chain, width)
	}

	for i := 0; i < height; i++ {
		index := a.RandomIndex()
		table[i][0] = index
		for j := 1; j < width; j++ {
			index = a.I2i(index, uint64(j))
			table[i][j] = index
		}
	}

	sort.Slice(table[:], func(i, j int) bool {
		return table[i][width-1] < table[j][width-1]
	})
	return RainbowTable{height, width, table, a, hash}
}

func (r *RainbowTable) Export(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(
		fmt.Sprintf("# %s\n# %s %d %d\n# %d %d\n", r.hashMethod, r.alphabet.alphabet, r.alphabet.min, r.alphabet.max, r.height, r.width),
	)
	if err != nil {
		return err
	}

	for i := 0; i < r.height; i++ {
		f.WriteString("\n")
		for j := 0; j < r.width; j++ {
			f.WriteString(fmt.Sprintf("%d ", r.table[i][j]))
		}
	}

	return nil
}

func (r *RainbowTable) Import(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	// Process header
	for i := 0; i < 4; i++ {
		if sc.Scan() {
			line := strings.ReplaceAll(sc.Text(), "# ", "")
			switch i {
			case 0:
				r.hashMethod = HashType(line)
				break
			case 1:
				tokens := strings.Split(line, " ")
				min, _ := strconv.Atoi(tokens[1])
				max, _ := strconv.Atoi(tokens[2])
				r.alphabet = GenerateAlphabet(
					tokens[0],
					min,
					max,
				)
				break
			case 2:
				tokens := strings.Split(line, " ")
				r.height, _ = strconv.Atoi(tokens[0])
				r.width, _ = strconv.Atoi(tokens[1])
				break
			default:
				break
			}
		} else {
			return errors.New("Failed to process the header during rainbow table import")
		}
	}

	if r.table == nil || len(r.table) != r.height || len(r.table[0]) != r.width {
		newTable := CreateRaindowTable(r.height, r.width, r.alphabet, r.hashMethod)
		*r = newTable
	}

	// Read through tokens until EOF
	var i int
	for sc.Scan() {
		tokens := strings.Split(sc.Text(), " ")
		for j := 0; j < r.width; j++ {
			entry, _ := strconv.Atoi(tokens[j])
			r.table[i][j] = uint64(entry)
		}
		i++
	}

	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}

func (r *RainbowTable) Print() {
	fmt.Printf("Hash method: %s\n", r.hashMethod)
	fmt.Printf("Alphabet: %s\n", r.alphabet.alphabet)
	fmt.Printf("Alphabet lenght: %d\n", r.alphabet.length)
	fmt.Printf("Min size: %d\n", r.alphabet.min)
	fmt.Printf("Max size: %d\n", r.alphabet.max)
	fmt.Printf("Possibilities: %d\n", r.alphabet.possibilities)
	fmt.Printf("Height: %d\n", r.height)
	fmt.Printf("Width: %d\n\n", r.width)
	fmt.Println("Content:")
	for i := 0; i < r.height; i++ {
		fmt.Printf("Chain %d: %d --> %d\n", i, r.table[i][0], r.table[i][r.width-1])
	}
}

func (r *RainbowTable) Invert(hash []byte) (out string, err error) {
	var nbCandidates int
	for t := r.width - 1; t < 0; t-- {
		idx := r.alphabet.H2i(hash, uint64(t))
		for i := t + 1; i < r.width; i++ {
			idx = r.alphabet.I2i(idx, uint64(i))
		}
		if a, b, err := r.Search(idx); err == nil {
			for i := a; i <= b; i++ {
				if out, valid := r.CheckCandidate(hash, t, r.table[i][0]); valid {
					return out, nil
				} else {
					nbCandidates++
				}
			}
		}
	}

	return "", errors.New("No candidate found")
}

func (r *RainbowTable) Search(idx uint64) (A int, B int, Err error) {
	A = sort.Search(r.height, func(i int) bool {
		return r.table[i][r.width-1] == idx
	})
	if A < r.height {
		for j := A - 1; j > 0; j-- {
			if r.table[j][r.width-1] != idx {
				break
			}
			A = j
		}
		for j := A + 1; j < r.height; j++ {
			if r.table[j][r.width-1] != idx {
				break
			}
			B = j
		}
		return
	}

	return 0, 0, errors.New("Not found")
}

func (r *RainbowTable) CheckCandidate(hash []byte, t int, idx uint64) (string, bool) {
	for i := 1; i < t; i++ {
		idx = r.alphabet.I2i(idx, uint64(i))
	}
	clair := r.alphabet.I2c(idx)
	h2, err := Hash(clair, r.hashMethod)
	if err != nil {
		logrus.Fatal(err)
	}

	return string(clair), bytes.Equal(h2, hash)
}
