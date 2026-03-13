package store

import (
	"crypto/md5"
	"encoding/hex"
	"sort"
)

type HashRing struct {
	Nodes []string
	Ring  map[string]string
	Keys  []string
}

func NewHashRing(nodes []string) *HashRing {
	hr := &HashRing{
		Nodes: nodes,
		Ring:  make(map[string]string),
	}
	for _, node := range nodes {
		hr.AddNode(node)
	}
	return hr
}

func (hr *HashRing) hash(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

func (hr *HashRing) AddNode(nodePort string) {
	h := hr.hash(nodePort)
	hr.Ring[h] = nodePort
	hr.Keys = append(hr.Keys, h)
	sort.Strings(hr.Keys)
}

func (hr *HashRing) GetNode(key string) string {
	if len(hr.Keys) == 0 { return "" }
	h := hr.hash(key)
	idx := sort.Search(len(hr.Keys), func(i int) bool {
		return hr.Keys[i] >= h
	})
	if idx == len(hr.Keys) { idx = 0 }
	return hr.Ring[hr.Keys[idx]]
}