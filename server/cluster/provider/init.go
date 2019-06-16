package provider

import (
	"htcache/server/cluster/provider"
	"htcache/server/cluster/provider/gossip"
)

func init() {
	provider.Register("gossip", gossip.New)
}
