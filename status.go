package dfa

import (
	"strings"
	"sync"
)

type OptionalAutomaton struct {
	lock       sync.Mutex
	selectList map[string]Automaton
}

func (o *OptionalAutomaton) String() string {
	o.lock.Lock()
	defer o.lock.Unlock()

	b := strings.Builder{}
	for key := range o.selectList {
		b.WriteString(key + " ")
	}
	return b.String()
}

// Get 获取状态
func (o *OptionalAutomaton) Get(key string) (Automaton, bool) {
	o.lock.Lock()
	defer o.lock.Unlock()
	ret, ok := o.selectList[key]
	return ret, ok
}

// Add 添加状态 如果存在返回false
func (o *OptionalAutomaton) Add(key string, status Automaton) bool {
	o.lock.Lock()
	defer o.lock.Unlock()
	if _, ok := o.selectList[key]; ok {
		return false
	}
	o.selectList[key] = status
	return true
}

func NewOptionalStatus() *OptionalAutomaton {
	return &OptionalAutomaton{
		selectList: make(map[string]Automaton),
	}
}

type Automaton interface {
	// Transfer 转移到下一个状态
	Transfer(id string) bool

	// Peek 查看下一个可能得状态列表
	Peek() *OptionalAutomaton

	// IsFinalState 是否是终态
	IsFinalState() bool

	// Circulation 所有的流转状态
	Circulation() []MetaStatus
}

type Status struct {
	lock sync.Mutex
	*MetaStatus
	dfa    *DFA
	Record []MetaStatus `json:"record"`
}

func NewStatus(dfa *DFA) *Status {
	return &Status{
		dfa: dfa,
	}
}

func NewStatusWithStatus(dfa *DFA, s *MetaStatus, records ...MetaStatus) *Status {
	return &Status{
		dfa:        dfa,
		MetaStatus: s,
		Record:     records,
	}
}

func (s *Status) load() {
	if s.MetaStatus == nil {
		s.MetaStatus = s.dfa.Get(s.dfa.start.ID)
	}
	if s.MetaStatus.load {
		return
	}
	s.MetaStatus.load = true
	s.MetaStatus.next = make(map[string]*MetaStatus)
	for _, nid := range s.MetaStatus.Next {
		s.MetaStatus.next[nid] = s.dfa.Get(nid)
	}
}

type MetaStatus struct {
	ID           string   `json:"id" yaml:"id"`
	Next         []string `json:"next" yaml:"next"`
	next         map[string]*MetaStatus
	Payload      interface{} `json:"payload" yaml:"payload"`
	InitialState bool        `json:"initial_state" yaml:"initial_state"`
	FinalState   bool        `json:"final_state" yaml:"final_state"`
	AfterCall    []string    `json:"after_call" yaml:"after_call"`
	BeforeCall   []string    `json:"before_call" yaml:"before_call"`
	load         bool
}

func (s *Status) Transfer(id string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.load()

	ret, ok := s.MetaStatus.next[id]
	if !ok {
		return false
	}
	for i := len(s.MetaStatus.AfterCall) - 1; i >= 0; i-- {
		after, ok := Registry.Get(s.MetaStatus.AfterCall[i])
		if !ok {
			continue
		}
		after(s)
	}

	s.Record = append(s.Record, *s.MetaStatus)

	s.MetaStatus = ret

	for i := len(s.MetaStatus.BeforeCall) - 1; i >= 0; i-- {
		before, ok := Registry.Get(s.MetaStatus.BeforeCall[i])
		if !ok {
			continue
		}
		before(s)
	}

	if s.MetaStatus.FinalState {
		s.Record = append(s.Record, *s.MetaStatus)
	}

	return true
}

func (s *Status) Peek() *OptionalAutomaton {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.load()

	oa := NewOptionalStatus()
	for k, v := range s.MetaStatus.next {
		oa.Add(k, &Status{
			MetaStatus: v,
			Record:     s.Record,
		})
	}
	return oa
}

func (s *Status) IsFinalState() bool {
	return s.MetaStatus.FinalState
}

func (s *Status) Circulation() []MetaStatus {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.Record
}
