package dfa

import "sync"

type OptionalAutomaton struct {
	lock       sync.Mutex
	selectList map[string]Automaton
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
	Circulation() []Automaton
}

type Status struct {
	lock   sync.Mutex
	status *status
	dfa    *DFA
	Record []Automaton `json:"record"`
}

func NewStatus(dfa *DFA) *Status {
	return &Status{
		dfa: dfa,
	}
}

func NewStatusWithStatus(dfa *DFA, s *status, records ...Automaton) *Status {
	return &Status{
		dfa:    dfa,
		status: s,
		Record: records,
	}
}

func (s *Status) load() {
	if s.status == nil {
		s.status = s.dfa.Get(s.dfa.start.ID)
	}
	if s.status.load {
		return
	}
	s.status.load = true
	s.status.next = make(map[string]*status)
	for _, nid := range s.status.Next {
		s.status.next[nid] = s.dfa.Get(nid)
	}
}

type status struct {
	ID           string   `json:"id" yaml:"id"`
	Next         []string `json:"next" yaml:"next"`
	next         map[string]*status
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

	ret, ok := s.status.next[id]
	if !ok {
		return false
	}
	for i := len(s.status.AfterCall) - 1; i >= 0; i-- {
		after, ok := Registry.Get(s.status.AfterCall[i])
		if !ok {
			continue
		}
		after(s)
	}

	s.Record = append(s.Record, &Status{
		status: s.status,
		Record: s.Record,
	})

	s.status = ret

	for i := len(s.status.BeforeCall) - 1; i >= 0; i-- {
		before, ok := Registry.Get(s.status.BeforeCall[i])
		if !ok {
			continue
		}
		before(s)
	}

	if s.status.FinalState {
		s.Record = append(s.Record, &Status{
			status: s.status,
			Record: s.Record,
		})
	}

	return true
}

func (s *Status) Peek() *OptionalAutomaton {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.load()

	oa := NewOptionalStatus()
	for k, v := range s.status.next {
		oa.Add(k, &Status{
			status: v,
			Record: s.Record,
		})
	}
	return oa
}

func (s *Status) IsFinalState() bool {
	return s.status.FinalState
}

func (s *Status) Circulation() []Automaton {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.Record
}
