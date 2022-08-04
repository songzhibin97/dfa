package dfa

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func Init(t *testing.T) {
	Registry.Add("after1", func(s *Status) {
		t.Log("after1", s.Record)
	})
	Registry.Add("after2", func(s *Status) {
		t.Log("after2", s.Record)
	})
	Registry.Add("after3", func(s *Status) {
		t.Log("after3", s.Record)
	})

	Registry.Add("before1", func(s *Status) {
		t.Log("before1", s.Record)
	})
	Registry.Add("before2", func(s *Status) {
		t.Log("before2", s.Record)
	})
	Registry.Add("before3", func(s *Status) {
		t.Log("before3", s.Record)
	})
}

func Test_Status(t *testing.T) {
	Init(t)
	config, err := ioutil.ReadFile("./meta.yaml")
	if err != nil {
		t.Error(err)
	}

	dfa, err := NewDfa(string(config))
	if err != nil {
		t.Error(err)
	}
	s := NewStatus(dfa)

	assert.Equal(t, len(s.Peek().selectList), 1)
	assert.Equal(t, s.Transfer("bad"), false)
	assert.Equal(t, s.Transfer("1"), true)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("1"), true)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("1"), true)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("1"), true)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("1"), true)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("2"), true)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("1"), true)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("end"), false)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("3"), true)
	assert.Equal(t, len(s.Peek().selectList), 4)
	assert.Equal(t, s.Transfer("end"), true)
	assert.Equal(t, len(s.Peek().selectList), 0)
	var list []string
	for _, automaton := range s.Circulation() {
		list = append(list, automaton.(*Status).status.ID)
	}
	assert.Equal(t, []string{"start", "1", "1", "1", "1", "1", "2", "1", "3", "end"}, list)
}

func Test_StatusWithStatus(t *testing.T) {
	Init(t)
	config, err := ioutil.ReadFile("./meta.yaml")
	if err != nil {
		t.Error(err)
	}
	dfa, err := NewDfa(string(config))
	if err != nil {
		t.Error(err)
	}

	s := NewStatusWithStatus(dfa, &status{
		ID:   "10",
		Next: []string{"start", "1", "2", "3", "end"},
	})

	assert.Equal(t, len(s.Peek().selectList), 5)
	assert.Equal(t, s.Transfer("bad"), false)
	assert.Equal(t, s.Transfer("1"), true)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("1"), true)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("1"), true)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("1"), true)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("1"), true)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("2"), true)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("1"), true)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("end"), false)
	assert.Equal(t, len(s.Peek().selectList), 3)
	assert.Equal(t, s.Transfer("3"), true)
	assert.Equal(t, len(s.Peek().selectList), 4)
	assert.Equal(t, s.Transfer("end"), true)
	assert.Equal(t, len(s.Peek().selectList), 0)
	var list []string
	for _, automaton := range s.Circulation() {
		list = append(list, automaton.(*Status).status.ID)
	}
	assert.Equal(t, []string{"10", "1", "1", "1", "1", "1", "2", "1", "3", "end"}, list)

}
