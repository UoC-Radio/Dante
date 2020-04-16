package auth

import (
	"github.com/jpillora/ipfilter"
	"strings"
)

func Allowed(ipAddr string) bool {
	f := ipfilter.New(ipfilter.Options{
		AllowedIPs:     []string{"127.0.0.1/24"},
		BlockByDefault: true,
	})

	ip := strings.Split(ipAddr, ":")[0]
	return f.Allowed(ip)
}
