package iputils

import (
	"errors"
	"fmt"
	"net"
)

type Cidr struct {
	Ip         net.IP
	NetAddress *net.IPNet
}

func (c *Cidr) PrefixToMask() string {
	if len(m) != 4 {
		return errors.New("Error: IPv4 only.")
	}
	return fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
}

func (c *Cidr) GetIpv4WithMask(s string) string {
	m := c.NetAddress.Mask
	if len(m) != 4 {
		return errors.New("Error: IPv4 only.")
	}
	return fmt.Sprintf("%s %d.%d.%d.%d", c.Ip.String(), m[0], m[1], m[2], m[3])
}

func StringToCidr(s string) (*Cidr, error) {
	ip, netaddr, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	cidr := &Cidr{Ip: ip, NetAddress: netaddr}
	return cidr
}
