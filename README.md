# dfa

确定有限状态自动机或确定有限自动机（deterministic finite automaton, DFA）

## 解决了什么问题
根据事先定义的配置信息,在每个节点状态下进行对应的状态转移


## demo

![dfa](dfa.png)

```
start -> 1

1 -> 1 
1 -> 2
1 -> 3

2 -> 1
2 -> 2
2 -> 3

3 -> 1 
3 -> 2 
3 -> 3
3 -> end

start 没有入度,因此为初态,所有状态的起始位置
end 没有出度,因此为终态,所有状态的终止位置
```
```go
// config meta.yaml 

func main() {
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
	s.Peek()  // [1]
	
	s.Transfer("bad") // false 
	
	 // 转移到状态1
	s.Transfer("1")
	s.Peek()  // [1,2,3]
	
	// 转移到状态3
	s.Transfer("3")  
	s.Peek()  // [1,2,3,end]
	
	s.Transfer("end")
	s.Peek()  // []
	
	ns, _ := json.Marshal(s.Circulation())
	t.Log(string(ns))
	// 展示了所有的状态转移
	/*
	[
    {
        "id": "start",
        "next": [
            "1"
        ],
        "payload": "start",
        "initial_state": true,
        "final_state": false,
        "after_call": [
            "after1"
        ],
        "before_call": [
            "before1"
        ]
    },
    {
        "id": "1",
        "next": [
            "1",
            "2",
            "3"
        ],
        "payload": "1",
        "initial_state": false,
        "final_state": false,
        "after_call": [
            "after2"
        ],
        "before_call": [
            "before2"
        ]
    },
    {
        "id": "3",
        "next": [
            "1",
            "2",
            "3",
            "end"
        ],
        "payload": "3",
        "initial_state": false,
        "final_state": false,
        "after_call": [
            "after1",
            "after2",
            "after3"
        ],
        "before_call": [
            "before1",
            "before2",
            "before3"
        ]
    },
    {
        "id": "end",
        "next": null,
        "payload": "end",
        "initial_state": false,
        "final_state": true,
        "after_call": null,
        "before_call": null
    }
]
	
	 */
	
}

```
