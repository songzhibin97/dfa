package dfa

import (
	"encoding/json"
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
		list = append(list, automaton.ID)
	}
	assert.Equal(t, []string{"start", "1", "1", "1", "1", "1", "2", "1", "3", "end"}, list)

	ns, _ := json.Marshal(s.Circulation())
	t.Log(string(ns))
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

	s := NewStatusWithStatus(dfa, &MetaStatus{
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
		list = append(list, automaton.ID)
	}
	assert.Equal(t, []string{"10", "1", "1", "1", "1", "1", "2", "1", "3", "end"}, list)

	ns, _ := json.Marshal(s.Circulation())
	t.Log(string(ns))
}

func Test_Demo(t *testing.T) {
	Init(t)
	config, err := ioutil.ReadFile("./meta.yaml")
	if err != nil {
		panic("read meta.yaml error")
	}

	dfa, err := NewDfa(string(config))
	if err != nil {
		panic("new dfa error")
	}
	s := NewStatus(dfa)

	// 查看当前状态下可转移的选项
	s.Peek() // [1]

	s.Transfer("bad") // false

	// 转移到状态1
	s.Transfer("1")
	s.Peek() // [1,2,3]

	// 转移到状态3
	s.Transfer("3")
	s.Peek() // [1,2,3,end]

	s.Transfer("end")
	s.Peek() // []

	ns, _ := json.Marshal(s.Circulation())
	t.Log(string(ns))
}
