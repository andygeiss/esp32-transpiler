package ino

import (
	"encoding/json"
	"github.com/andygeiss/esp32-transpiler/business/worker"
	"io/ioutil"
)

// Mapping specifies the business logic to apply transformation to a specific Golang identifier by reading simple JSON map.
type Mapping struct {
	Filename string            `json:"filename"`
	Rules    map[string]string `json:"rules"`
}

// NewMapping creates a new mapping and returns its address.
func NewMapping(filename string) worker.Mapping {
	rules := make(map[string]string, 0)
	return &Mapping{filename, rules}
}

// Apply checks the Golang identifier and transforms it to a specific representation.
func (m *Mapping) Apply(ident string) string {
	for wanted := range m.Rules {
		if ident == wanted {
			ident = m.Rules[ident]
		}
	}
	return ident
}

// Read gets the mapping rules from the local filesystem.
func (m *Mapping) Read() error {
	bytes, err := ioutil.ReadFile(m.Filename)
	if err != nil {
		return err
	}
	var rules map[string]string
	json.Unmarshal(bytes, &rules)
	m.Rules = rules
	return nil
}
