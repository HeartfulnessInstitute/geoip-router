package geoip

// ISOCountryCode reprents the set of countries supported by HFN GeoRouter
type ISOCountryCode string

func (c ISOCountryCode) String() string {
	return string(c)
}

const (
	ISOCountryCodeUS ISOCountryCode = "us"
	ISOCountryCodeIN ISOCountryCode = "in"
)

// DefaultISOCountryCode for HFN routing
const DefaultISOCountryCode = ISOCountryCodeUS

// ParseISOCode parses any given ISO code in string format to a supported ISOCountryCode type
func ParseISOCode(code string) ISOCountryCode {
	switch code {
	case ISOCountryCodeIN.String():
		return ISOCountryCodeIN
	default:
		return ISOCountryCodeUS
	}
}