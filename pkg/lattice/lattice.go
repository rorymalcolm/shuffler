package lattice

import (
	"errors"
	"fmt"
)

type DimensionMismatchError struct{}

func (e *DimensionMismatchError) Error() string {
	return "dimension mismatch in lattice"
}

type Lattice[T any] struct {
	dimensions        []string
	endpointsByCoord  map[string][]T
	valuesByDimension map[string]map[string]bool
}

type Coordinate[T any] struct {
	coords []string
	value  T
}

func (c *Coordinate[T]) NewCoordinate(coords []string, value T) *Coordinate[T] {
	c.coords = coords
	c.value = value
	return c
}

func (c *Coordinate[T]) GetCoords() []string {
	return c.coords
}

func (c *Coordinate[T]) GetValue() T {
	return c.value
}

func (c *Coordinate[T]) GetCoordString() string {
	return fmt.Sprintf("%v", c.coords)
}

func NewLattice[T any](dimensions []string) *Lattice[T] {
	return &Lattice[T]{
		dimensions:        dimensions,
		endpointsByCoord:  make(map[string][]T),
		valuesByDimension: make(map[string]map[string]bool),
	}
}

func (l *Lattice[T]) AddEndpointsForSector(sectorCoordinates []string, endpoints []T) error {
	if len(sectorCoordinates) != len(l.dimensions) {
		return &DimensionMismatchError{}
	}

	sectorKey := fmt.Sprintf("%v", sectorCoordinates)
	existing := l.endpointsByCoord[sectorKey]
	toBeAdded := append(existing, endpoints...)

	l.endpointsByCoord[sectorKey] = toBeAdded

	for i, dimension := range l.dimensions {
		value := sectorCoordinates[i]

		if _, exists := l.valuesByDimension[dimension]; !exists {
			l.valuesByDimension[dimension] = make(map[string]bool)
		}

		l.valuesByDimension[dimension][value] = true
	}

	return nil
}

func (l *Lattice[T]) GetEndpointsForSector(sectorCoordinates []string) ([]T, error) {
	if len(sectorCoordinates) != len(l.dimensions) {
		return nil, &DimensionMismatchError{}
	}
	sectorKey := fmt.Sprintf("%v", sectorCoordinates)
	endpoints, ok := l.endpointsByCoord[sectorKey]
	if !ok {
		return nil, errors.New("no endpoints found for sector")
	}
	return endpoints, nil
}

func (l *Lattice[T]) GetAllEndpoints() []T {
	allEndpoints := []T{}
	for _, endpoints := range l.endpointsByCoord {
		allEndpoints = append(allEndpoints, endpoints...)
	}
	return allEndpoints
}

func (l *Lattice[T]) GetAllCoordinates() []map[string]T {
	coordinates := []map[string]T{}
	dimensionsPerDimension := map[string]int{}
	for _, dimension := range l.dimensions {
		dimensionsPerDimension[dimension] = len(l.valuesByDimension[dimension])
	}
	return coordinates
}

func (l *Lattice[T]) GetDimensionality() int {
	dimensionsPerDimension := map[string]int{}
	for _, dimension := range l.dimensions {
		dimensionsPerDimension[dimension] = len(l.valuesByDimension[dimension])
	}
	return len(dimensionsPerDimension)
}

func (l *Lattice[T]) GetDimensionNames() []string {
	return l.dimensions
}

func (l *Lattice[T]) GetDimensionName(dimension int) string {
	return l.dimensions[dimension]
}

func (l *Lattice[T]) GetDimensionValues(dimension string) []string {
	values := []string{}
	for value := range l.valuesByDimension[dimension] {
		values = append(values, value)
	}
	return values
}

func (l *Lattice[T]) GetDimensionSize(dimension string) int {
	return len(l.valuesByDimension[dimension])
}
