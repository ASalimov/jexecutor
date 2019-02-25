package logic

import (
	"sort"
)

// Order is the item for executing
type Order struct {
	Job   string
	Query map[string]interface{}
}

func (o *Order) GetVals() (rsp string) {
	// To store the keys in slice in sorted order
	var keys []string
	for k := range o.Query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// To perform the opertion you want
	for _, k := range keys {
		rsp += o.Query[k].(string) + ", "
	}
	rsp = rsp[0 : len(rsp)-2]
	return
}
