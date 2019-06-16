package gossip

import (
	"time"
	"htcache/server/cluster"
	"github.com/stathat/consistent"
	"github.com/hashicorp/memberlist"
)

type Node struct {
	*consistent.Consistent
	Addr string
}

func New(addr, cluster string) (cluster.Node, error) {
	config := memberlist.DefaultLANConfig()
	config.Name = addr
	config.BindAddr = addr
	ml, err := memberlist.Create(config)
	if err != nil {
		return nil, err
	}
	if cluster == "" {
		cluster = addr
	}
	_, err := ml.Join([]string{cluster})
	if err != nil {
		return nil, err
	}
	ch := consistent.New()
	ch.NumberOfReplicas = 1024
	go func() {
		for _ = range timer.Tick(time.Second) {
			members := ml.Members()
			nodes := make([]string, len(members))
			for index, member := range members {
				nodes[index] = member.Name
			}
			ch.Set(nodes)
		}
	}()
}
