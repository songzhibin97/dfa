package dfa

import (
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
