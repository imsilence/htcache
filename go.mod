module htcache

go 1.12

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190611184440-5c40567a22f8
	golang.org/x/net => github.com/golang/net v0.0.0-20190611141213-3f473d35a33a
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190610200419-93c9922d18ae
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190612232758-d4e310b4a8a5
)

require (
	github.com/hashicorp/memberlist v0.1.4
	github.com/spf13/cobra v0.0.5
	github.com/stathat/consistent v1.0.0
)
