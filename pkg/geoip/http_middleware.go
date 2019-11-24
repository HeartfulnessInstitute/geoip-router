package geoip

import (
	"bytes"
	"net"
	"net/http"
	"strings"
)

//ipRange - a structure that holds the start and end of a range of ip addresses
type ipRange struct {
	start net.IP
	end   net.IP
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

var privateRanges = []ipRange{
	{
		start: net.ParseIP("10.0.0.0"),
		end:   net.ParseIP("10.255.255.255"),
	},
	{
		start: net.ParseIP("100.64.0.0"),
		end:   net.ParseIP("100.127.255.255"),
	},
	{
		start: net.ParseIP("172.16.0.0"),
		end:   net.ParseIP("172.31.255.255"),
	},
	{
		start: net.ParseIP("192.0.0.0"),
		end:   net.ParseIP("192.0.0.255"),
	},
	{
		start: net.ParseIP("192.168.0.0"),
		end:   net.ParseIP("192.168.255.255"),
	},
	{
		start: net.ParseIP("198.18.0.0"),
		end:   net.ParseIP("198.19.255.255"),
	},
}

// isPrivateSubnet - check to see if this ip is in a private subnet
func isPrivateSubnet(ipAddress net.IP) bool {
	// my use case is only concerned with ipv4 atm
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		// iterate over all our ranges
		for _, r := range privateRanges {
			// check if this ip is in a private range
			if inRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

func getIPAdress(r *http.Request) net.IP {
	for _, h := range []string{"x-forwarded-for", "x-real-ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// The left most IP is the one closest to the client.
		for _, address := range addresses {
			ip := strings.TrimSpace(address)
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
				// bad address, go to next
				continue
			}
			return realIP
		}
	}
	return net.IPv4zero
}

// HTTPResolverHandler from an HTTP Request
func HTTPResolverHandler(resolver Resolver, rules CountryLocationRoutingRules) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := getIPAdress(r)
		resolvedCountry, err := resolver.ResolveCountryCode(r.Context(), clientIP)
		if err != nil {
			resolvedCountry = DefaultISOCountryCode
		}
		loc, ok := rules[resolvedCountry]
		if !ok {
			// Default to a relative location if rules don't match
			loc = "/" + resolvedCountry.String()
		}
		http.Redirect(w, r, loc, http.StatusFound)
	}
}
