package main

import (
	"encoding/json"
	"os"
)

type PoolBlock map[string][]Block //从ipam那边获到block分配的情况 第一个存总共的第二个存现在分配的

type Block struct {
	Total        float64
	Use          float64
	NodeSelector string
}

func (p *PoolBlock) Add(name string, total, use float64, nodeSelector string) {
	pr := *p

	if pr[name] == nil {
		pr = make(PoolBlock)
	}

	pr[name] = append(pr[name], Block{
		Total:        total,
		Use:          use,
		NodeSelector: nodeSelector,
	})
}

func New(name string, total, use float64, nodeSelector string) PoolBlock {
	var p PoolBlock
	if p == nil {
		p = make(map[string][]Block)
	}
	p[name] = append(p[name], Block{
		Total:        total,
		Use:          use,
		NodeSelector: nodeSelector,
	})
	return p
}

func main() {
	//x := make(PoolBlock)
	//x := make(PoolBlock)
	var x PoolBlock
	x.Add("zxz", 18, 18, "yyds")
	json.NewEncoder(os.Stdout).Encode(x)
}
