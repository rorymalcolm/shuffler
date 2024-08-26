package lattice_test

import (
	"testing"

	"github.com/rorymalcolm/shuffler/v2/pkg/lattice"
	"github.com/stretchr/testify/assert"
)

func TestSingleCellLattice(t *testing.T) {
	// Create a lattice with a single cell and string endpoints
	lat := lattice.NewLattice[string]([]string{"cell"})

	// Add endpoints
	err := lat.AddEndpointsForSector([]string{"cell"}, []string{"A"})
	assert.NoError(t, err)

	err = lat.AddEndpointsForSector([]string{"cell"}, []string{"B", "C", "D"})
	assert.NoError(t, err)

	// Check that all endpoints are in with correct ordering
	allEndpoints := lat.GetAllEndpoints()
	expectedEndpoints := []string{"A", "B", "C", "D"}

	assert.ElementsMatch(t, expectedEndpoints, allEndpoints)
}

func TestOneDimensionalLattice(t *testing.T) {
	// Create a one-dimensional lattice
	lat := lattice.NewLattice[string]([]string{"AZ"})

	// Add endpoints to different sectors
	err := lat.AddEndpointsForSector([]string{"us-east-1a"}, []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"})
	assert.NoError(t, err)

	err = lat.AddEndpointsForSector([]string{"us-east-1b"}, []string{"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T"})
	assert.NoError(t, err)

	// Validate total number of endpoints
	allEndpoints := lat.GetAllEndpoints()
	assert.Equal(t, 20, len(allEndpoints))

	subLat, err := lat.SimulateFailure("AZ", "us-east-1a")
	assert.NoError(t, err)
	assert.Equal(t, 10, len(subLat.GetAllEndpoints()))

	subLat, err = lat.SimulateFailure("AZ", "us-east-1b")
	assert.NoError(t, err)
	assert.Equal(t, 10, len(subLat.GetAllEndpoints()))
}

func TestTwoDimensionalLattice(t *testing.T) {
	// Create a two-dimensional lattice
	lat := lattice.NewLattice[string]([]string{"AZ", "Version"})

	// Add endpoints to various sectors
	err := lat.AddEndpointsForSector([]string{"us-east-1a", "1"}, []string{"A", "B", "C", "D", "E"})
	assert.NoError(t, err)

	err = lat.AddEndpointsForSector([]string{"us-east-1a", "2"}, []string{"F", "G", "H", "I", "J"})
	assert.NoError(t, err)

	err = lat.AddEndpointsForSector([]string{"us-east-1b", "1"}, []string{"K", "L", "M", "N", "O"})
	assert.NoError(t, err)

	err = lat.AddEndpointsForSector([]string{"us-east-1b", "2"}, []string{"P", "Q", "R", "S", "T"})
	assert.NoError(t, err)

	// Validate total number of endpoints
	allEndpoints := lat.GetAllEndpoints()
	assert.Equal(t, 20, len(allEndpoints))

	subLat, err := lat.SimulateFailure("AZ", "us-east-1a")
	assert.NoError(t, err)
	assert.Equal(t, 10, len(subLat.GetAllEndpoints()))

	subLat, err = lat.SimulateFailure("AZ", "us-east-1b")
	assert.NoError(t, err)
	assert.Equal(t, 10, len(subLat.GetAllEndpoints()))

	subLat, err = lat.SimulateFailure("Version", "1")
	assert.NoError(t, err)
	assert.Equal(t, 10, len(subLat.GetAllEndpoints()))

	subLat, err = lat.SimulateFailure("Version", "2")
	assert.NoError(t, err)
	assert.Equal(t, 10, len(subLat.GetAllEndpoints()))

	subLat, err = lat.SimulateFailure("AZ", "us-east-1a")
	assert.NoError(t, err)
	subLat, err = subLat.SimulateFailure("Version", "1")
	assert.NoError(t, err)
	assert.Equal(t, 5, len(subLat.GetAllEndpoints()))

	subLat, err = lat.SimulateFailure("AZ", "us-east-1a")
	assert.NoError(t, err)
	subLat, err = subLat.SimulateFailure("Version", "2")
	assert.NoError(t, err)
	assert.Equal(t, 5, len(subLat.GetAllEndpoints()))

	subLat, err = lat.SimulateFailure("AZ", "us-east-1b")
	assert.NoError(t, err)
	subLat, err = subLat.SimulateFailure("Version", "1")
	assert.NoError(t, err)
	assert.Equal(t, 5, len(subLat.GetAllEndpoints()))
}
