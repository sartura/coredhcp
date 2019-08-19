package clientport

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/coredhcp/coredhcp/handler"
	"github.com/coredhcp/coredhcp/logger"
	"github.com/coredhcp/coredhcp/plugins"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
)

var log = logger.GetLogger()

func init() {
	plugins.RegisterPlugin("server", setupServerV6, setupServerV4)
}

// StaticRecords holds a MAC -> IP address mapping
var StaticRecords map[string]net.IP

// DHCPv6Records and DHCPv4Records are mappings between MAC addresses in
// form of a string, to network configurations.
var (
	serverIP     net.IP
	netmask      net.IP
	ClientSubnet net.IPMask
)

// Handler6 not implemented only IPv4
func Handler6(req, resp dhcpv6.DHCPv6) (dhcpv6.DHCPv6, bool) {
	return resp, true
}

// Handler4 handles DHCPv4 packets for the file plugin
func Handler4(req, resp *dhcpv4.DHCPv4) (*dhcpv4.DHCPv4, bool) {
	if req.OpCode != dhcpv4.OpcodeBootRequest {
		log.Printf("plugins/server: Not a BootRequest!")
	}
	switch mt := req.MessageType(); mt {
	case dhcpv4.MessageTypeDiscover:
		resp.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
		resp.UpdateOption(dhcpv4.OptSubnetMask(ClientSubnet))
		resp.UpdateOption(dhcpv4.OptServerIdentifier(serverIP))
		resp.UpdateOption(dhcpv4.OptRouter(serverIP))
	case dhcpv4.MessageTypeRequest:
		resp.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
		resp.UpdateOption(dhcpv4.OptSubnetMask(ClientSubnet))
		resp.UpdateOption(dhcpv4.OptServerIdentifier(serverIP))
		resp.UpdateOption(dhcpv4.OptRouter(serverIP))
	default:
		log.Printf("plugins/server: Unhandled message type: %v", mt)
	}
	return resp, true
}

// setupIPV6 not implemented only IPv4
func setupServerV6(args ...string) (handler.Handler6, error) {
	return nil, nil
}

func setupServerV4(args ...string) (handler.Handler4, error) {
	_, h4, err := setupServer(false, args...)
	return h4, err
}

func setupServer(v6 bool, args ...string) (handler.Handler6, handler.Handler4, error) {
	if v6 {

	} else {
		if len(args) < 2 {
			return nil, nil, errors.New("plugins/server: need a file name, server IP, netmask and a DNS server")
		}
		ip, network, err := net.ParseCIDR(args[0])
		if err != nil {
			return Handler6, Handler4, errors.New("plugins/server: expected an IPv4 address, got: " + args[0])
		}
		serverIP = ip

		netmask = net.IPv4(network.Mask[0], network.Mask[0], network.Mask[0], network.Mask[0])

		if !checkValidNetmask(netmask) {
			return Handler6, Handler4, errors.New("plugins/server: netmask is not valid, got: " + args[0])
		}

		subnet := net.ParseIP(args[1])
		if subnet.To16() == nil {
			return Handler6, Handler4, errors.New("plugins/server: expected an IPv4 address, got:" + ClientSubnet.String())
		}
		subnet = subnet.To4()
		ClientSubnet = net.IPv4Mask(subnet[0], subnet[1], subnet[2], subnet[3])
		log.Printf("plugins/DNS: loaded plugin for DHCPv4.")

	}

	return Handler6, Handler4, nil
}
func checkValidNetmask(netmask net.IP) bool {
	if netmask.To16() == nil {
		return false
	}
	if netmask.IsUnspecified() {
		return false
	}
	netmaskInt := binary.BigEndian.Uint32(netmask.To4())
	x := ^netmaskInt
	y := x + 1
	return (y & x) == 0
}