package sharder

import (
	"crypto/md5"
	"encoding/binary"
	"math"
	"math/rand"

	"github.com/pkg/errors"

	"github.com/rorymalcolm/shuffler/v2/pkg/lattice"
)

type Shuffler[T any] struct {
	seed string
}

func NewShuffler[T any](seed string) *Shuffler[T] {
	return &Shuffler[T]{seed: seed}
}

type InvalidDimensionError struct{}

func (e *InvalidDimensionError) Error() string {
	return "invalid dimension"
}

func (s *Shuffler[T]) ShuffleShard(inputLattice *lattice.Lattice[T], identifier []byte, endpointsPerCel int) (*lattice.Lattice[T], error) {
	chosen := lattice.NewLattice[T](inputLattice.GetDimensionNames())

	// calculate a hash of the seed, then the identifier
	hash := md5.New()
	hash.Write([]byte(s.seed))
	hash.Write(identifier)
	// get the hash sum
	hashSum := hash.Sum(nil)

	// then, use the first 64 bytes of the hash sum to shuffle the shards, as a signed int64
	shardSeed := int64(binary.LittleEndian.Uint32(hashSum[:8]))

	// then, create a random seeded value from the shardSeed
	random := rand.New(rand.NewSource(shardSeed))
	var shuffledShards [][]string

	for _, dimName := range inputLattice.GetDimensionNames() {
		vals := inputLattice.GetDimensionValues(dimName)
		// shuffle the coordinates
		random.Shuffle(len(vals), func(i, j int) {
			vals[i], vals[j] = vals[j], vals[i]
		})

		shuffledShards = append(shuffledShards, vals)
	}

	dimensionality := inputLattice.GetDimensionality()

	if dimensionality == 0 {
		return nil, &InvalidDimensionError{}
	}

	if dimensionality == 1 {
		for _, shard := range shuffledShards[0] {
			endpoints, err := inputLattice.GetEndpointsForSector([]string{shard})
			if err != nil {
				return nil, errors.Wrap(err, "failed to get endpoints for sector")
			}
			rand.Shuffle(len(endpoints), func(i, j int) {
				endpoints[i], endpoints[j] = endpoints[j], endpoints[i]
			})
			chosen.AddEndpointsForSector([]string{shard}, endpoints)
		}

		return chosen, nil
	}

	minimumDimensionsSize := math.MaxInt32

	for _, dim := range inputLattice.GetDimensionNames() {
		if len(inputLattice.GetDimensionValues(dim)) < minimumDimensionsSize {
			minimumDimensionsSize = len(inputLattice.GetDimensionValues(dim))
		}
	}

	for i := 0; i < minimumDimensionsSize; i++ {
		var coordinates []string

		for j := 0; j < dimensionality; j++ {
			coordinates = append(coordinates, shuffledShards[j][i])
		}

		endpoints, err := inputLattice.GetEndpointsForSector(coordinates)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get endpoints for sector")
		}

		random.Shuffle(len(endpoints), func(i, j int) {
			endpoints[i], endpoints[j] = endpoints[j], endpoints[i]
		})
		chosen.AddEndpointsForSector(coordinates, endpoints[:endpointsPerCel])
	}

	return chosen, nil
}
