package chunk

import (
	"fmt"
	"io"
	"os"

	"github.com/klauspost/reedsolomon"
)

/*
	var dataShards = flag.Int("data", 4, "Number of shards to split the data into, must be below 257.")
	var parShards = flag.Int("par", 2, "Number of parity shards")
*/
type Fec struct {
	reader    io.Reader
	shards    [][]byte
	blocksize int
	shardsnum int
	index     int
}

// NewRabin creates a new Rabin splitter with the given
// average block size.
func NewFec(r io.Reader, dataShards int, parShards int) *Fec {
	myEnc, err := reedsolomon.New(dataShards, parShards)
	myCheckError(err)

	b, err := io.ReadAll(r)
	shds, err := myEnc.Split(b)
	myCheckError(err)

	err = myEnc.Encode(shds)
	myCheckError(err)

	block_size := len(shds[0])
	shards_num := len(shds)

	return &Fec{
		shards:    shds,
		blocksize: block_size,
		shardsnum: shards_num,
		index:     0,
	}
}

// NextBytes produces a new chunk.
func (splitter *Fec) NextBytes() ([]byte, error) {
	if splitter.index >= splitter.shardsnum {
		return nil, nil
	}
	b := splitter.shards[splitter.index]
	splitter.index += 1
	return b, nil
}

// Reader returns the io.Reader associated to this Splitter.
func (r *Fec) Reader() io.Reader {
	return r.reader
}

func myCheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}
