package gossip

import (
	"htcache/server/cluster"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/stathat/consistent"
)

type Node struct {
	*consistent.Consistent
	addr string
}

func (n *Node) IsProcess(key string) (string, bool) {
	addr, _ := n.Get(key)
	return addr, addr == n.Addr()
}

func (n *Node) Addr() string {
	return n.addr
}

func New(addr, cluster string) (cluster.Node, error) {
	host, _, _ := net.SplitHostPort(addr)

	config := memberlist.DefaultLANConfig()
	config.Name = addr
	config.BindAddr = host
	config.AdvertiseAddr = host
	ml, err := memberlist.Create(config)
	if err != nil {
		return nil, err
	}
	nodes := []string{host}
	if cluster != "" {
		nodes = strings.Split(cluster, ",")
	}
	_, err = ml.Join(nodes)
	if err != nil {
		return nil, err
	}
	ch := consistent.New()
	ch.NumberOfReplicas = 1024
	go func() {
		for _ = range time.Tick(time.Second) {
			members := ml.Members()
			nodes := make([]string, len(members))
			for index, member := range members {
				nodes[index] = member.Name
			}
			ch.Set(nodes)
		}
	}()
	return &Node{ch, addr}, nil
}
