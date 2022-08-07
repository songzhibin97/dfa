package dfa

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestNewDfa(t *testing.T) {
	config, err := ioutil.ReadFile("./meta.yaml")
	if err != nil {
		t.Error(err)
	}

	dfa, err := NewDfa(string(config))
	if err != nil {
		t.Error(err)
	}
	t.Log(dfa.effectiveState)
}

func TestDFA_checkConnectivity(t *testing.T) {
	tests := []struct {
		config    string
		err       bool
		diagnosis int
	}{
		{
			config: `
- id: "start"
  payload: "start"
  next:
    - "1"
  initial_state: true

- id: "1"
  payload: "1"
  next:
    - "2"

- id: "2"
  payload: "2"
  next:
    - "3"

- id: "3"
  payload: "3"
  next:
    - "end"

- id: "end"
  payload: "end"
  final_state: true

- id: "end1"
  payload: "end1"
  final_state: true		

- id: "end2"
  payload: "end2"
  final_state: true			
`,
			err:       false,
			diagnosis: 1,
		},
		{
			config: `
- id: "start"
  payload: "start"
  next:
    - "1"
  initial_state: true

- id: "1"
  payload: "1"
  next:
    - "2"

- id: "2"
  payload: "2"
  next:
    - "3"

- id: "3"
  payload: "3"
  next:
    - "end"

- id: "4"
  payload: "4"
  next:
    - "end"			

- id: "end"
  payload: "end"
  final_state: true
`,
			err:       false,
			diagnosis: 1,
		},
		{
			config: `
- id: "start"
  payload: "start"
  next:
    - "1"
  initial_state: true

- id: "1"
  payload: "1"
  next:
    - "2"

- id: "2"
  payload: "2"
  next:
    - "3"

- id: "3"
  payload: "3"
  next:
    - "4"
    - "end"

- id: "end"
  payload: "end"
  final_state: true
`,
			err:       true,
			diagnosis: 0,
		},
	}
	for _, tt := range tests {
		d, err := NewDfa(tt.config)
		assert.Equal(t, tt.err, err != nil)
		assert.Equal(t, len(d.diagnosis), tt.diagnosis)
	}
}
