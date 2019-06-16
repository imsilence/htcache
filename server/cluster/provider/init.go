package provider

import (
	"htcache/server/cluster"
	"htcache/server/cluster/provider/gossip"
)

func init() {
	cluster.Register("gossip", gossip.New)
}
