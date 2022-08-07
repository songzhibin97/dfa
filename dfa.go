package dfa

import (
	"errors"
	"fmt"
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

	ends []MetaStatus

	diagnosis []string
}

// checkConnectivity 检查连通性
func (d *DFA) checkConnectivity() error {
	// 检查是否起点能到所有终点
	var endMp = make(map[string]bool)
	var stepMp = make(map[string]bool)
	var unknownNode []string

	var dfs func(node MetaStatus)
	dfs = func(node MetaStatus) {
		if node.ID == "" {
			unknownNode = append(unknownNode, node.ID)
			return
		}
		if stepMp[node.ID] {
			return
		}
		stepMp[node.ID] = true

		if node.FinalState {
			endMp[node.ID] = true
			return
		}
		for _, nid := range node.Next {
			if nid == node.ID {
				continue
			}
			dfs(d.effectiveState[nid])
		}
	}

	dfs(d.start)

	if len(unknownNode) != 0 {
		return errors.New(fmt.Sprintf("unknown node: %v", unknownNode))
	}
	// 检查是否所有点都有用到(如果没用到 直接删了)
	if len(endMp) != len(d.ends) {
		unionSet := map[string]bool{}
		for _, end := range d.ends {
			unionSet[end.ID] = true
		}
		d.ends = d.ends[:0]
		for key := range endMp {
			d.ends = append(d.ends, d.effectiveState[key])
			delete(unionSet, key)
		}
		unionList := make([]string, 0, len(unionSet))
		for s := range unionSet {
			unionList = append(unionList, s)
			delete(d.effectiveState, s)
		}
		d.diagnosis = append(d.diagnosis, fmt.Sprintf("not all end nodes are reachable: %v", unionList))
	}

	if len(stepMp) != len(d.effectiveState) {
		unionList := make([]string, 0, len(d.effectiveState)-len(stepMp))
		for key := range d.effectiveState {
			if !stepMp[key] {
				unionList = append(unionList, key)
				delete(d.effectiveState, key)
			}
		}
		d.diagnosis = append(d.diagnosis, fmt.Sprintf("not all nodes are reachable: %v", unionList))
	}

	return nil
}

func (d *DFA) checkNodeEffective(node MetaStatus) error {
	if node.ID == "" {
		return errors.New("node id is empty")
	}

	if node.InitialState && len(node.Next) == 0 {
		return errors.New("initial state must have next")
	}

	if node.FinalState && len(node.Next) != 0 {
		return errors.New("final state must not have next")
	}

	if !node.FinalState && len(node.Next) == 0 {
		return errors.New("non-final state must have next")
	}
	return nil
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
	for _, v := range _status {

		err = ret.checkNodeEffective(v)
		if err != nil {
			return nil, err
		}

		// check 检查 是否是一个有效的状态节点
		if v.InitialState && len(v.Next) != 0 {
			if ret.start.ID != "" {
				return nil, errors.New("duplicate start state")
			}
			ret.start = v
		}
		if v.FinalState {
			// 终态的可执行转态必须为空
			ret.ends = append(ret.ends, v)
		}
		ret.effectiveState[v.ID] = v
	}

	if ret.start.ID == "" {
		return nil, errors.New("no start state")
	}

	if len(ret.ends) == 0 {
		return nil, errors.New("no final state")
	}
	// 检查从 start 到 end 的路径是否存在

	return ret, ret.checkConnectivity()
}
