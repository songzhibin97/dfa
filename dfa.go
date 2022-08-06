package dfa

import (
	"errors"
	"gopkg.in/yaml.v2"
	"sync"
)

type DFA struct {
	lock sync.Mutex

	// 配置文件
	config string

	// 加载所有状态
	effectiveState map[string]MetaStatus

	start MetaStatus
}

func (d *DFA) Get(key string) *MetaStatus {
	d.lock.Lock()
	defer d.lock.Unlock()
	v := d.effectiveState[key]
	return &v
}

func NewDfa(config string) (*DFA, error) {
	ret := &DFA{
		config:         config,
		effectiveState: make(map[string]MetaStatus),
	}
	var _status []MetaStatus
	err := yaml.Unmarshal([]byte(config), &_status)
	if err != nil {
		return nil, err
	}
	var end bool
	for _, v := range _status {
		if v.InitialState {
			ret.start = v
		}
		if v.FinalState {
			end = true
		}
		ret.effectiveState[v.ID] = v
	}
	if ret.start.ID == "" {
		return nil, errors.New("no start state")
	}
	if !end {
		return nil, errors.New("no final state")
	}
	return ret, nil
}
