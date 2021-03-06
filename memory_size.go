package main

import (
	"strconv"
)

// MemorySize represents the size in byte
type MemorySize struct {
	value int
}

// UnmarshalJSON unmarshals a JSON string
func (m *MemorySize) UnmarshalJSON(data []byte) error {
	s, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}

	m.value = s * 1048576 // MB to bytes
	return nil
}

// String returns the string representation of MemorySize
func (m *MemorySize) String() string {
	return strconv.Itoa(m.value)
}
