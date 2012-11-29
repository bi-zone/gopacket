// Copyright (c) 2012 Google, Inc. All rights reserved.
// Copyright (c) 2009-2012 Andreas Krennmair. All rights reserved.

package gopacket

import (
	"encoding/binary"
	"net"
	"strconv"
)

// LayerType is a unique identifier for each type of layer.  This enumeration
// does not match with any externally available numbering scheme... it's solely
// usable/useful within this library as a means for requesting layer types
// (see Packet.Layer) and determining which types of layers have been decoded.
// A LayerType corresponds 1:1 to a struct type.
type LayerType int

const (
	LayerTypePayload       LayerType = iota // Type: Payload
	LayerTypeDecodeFailure                  // Type: DecodeFailure
	LayerTypeEthernet                       // Type: Ethernet
	LayerTypePPP                            // Type: PPP
	LayerTypeIPv4                           // Type: IPv4
	LayerTypeIPv6                           // Type: IPv6
	LayerTypeTCP                            // Type: TCP
	LayerTypeUDP                            // Type: UDP
	LayerTypeICMP                           // Type: ICMP
	LayerTypeDot1Q                          // Type: Dot1Q
	LayerTypeARP                            // Type: ARP
)

func (l LayerType) String() string {
	switch l {
	case LayerTypePayload:
		return "Payload"
	case LayerTypeDecodeFailure:
		return "DecodeFailure"
	case LayerTypeEthernet:
		return "Ethernet"
	case LayerTypePPP:
		return "PPP"
	case LayerTypeIPv4:
		return "IPv4"
	case LayerTypeIPv6:
		return "IPv6"
	case LayerTypeTCP:
		return "TCP"
	case LayerTypeUDP:
		return "UDP"
	case LayerTypeICMP:
		return "ICMP"
	case LayerTypeDot1Q:
		return "Dot1Q"
	case LayerTypeARP:
		return "ARP"
	}
	return "<Unknown>"
}

// Layer represents a single decoded packet layer (using either the
// OSI or TCP/IP definition of a layer).  When decoding, a packet's data is
// broken up into a number of layers.  The caller may call LayerType() to
// figure out which type of layer he's received from the packet.  Optionally,
// he may then use a type assertion to get the actual layer type for deep
// inspection of the data.
type Layer interface {
	LayerType() LayerType
}

// Payload is a Layer containing the payload of a packet.  The definition of
// what constitutes the payload of a packet depends on previous layers; for
// TCP and UDP, we stop decoding above layer 4 and return the remaining
// bytes as a Payload.  Payload is an ApplicationLayer.
type Payload struct {
	Data []byte
}

// Returns LayerTypePayload
func (p *Payload) LayerType() LayerType { return LayerTypePayload }
func (p *Payload) Payload() []byte      { return p.Data }

// Address is the set of bytes used to address packets at various layers.
// See LinkLayer, NetworkLayer, and TransportLayer specifications.
type Address interface {
	String() string
	Raw() []byte
}

// MACAddress is the set of bytes representing a MAC address
type MACAddress net.HardwareAddr

func (a MACAddress) Raw() []byte    { return a }
func (a MACAddress) String() string { return net.HardwareAddr(a).String() }

// IPAddress is the set of bytes representing an IPv4 or IPv6 address
type IPAddress net.IP

func (a IPAddress) Raw() []byte    { return a }
func (a IPAddress) String() string { return net.IP(a).String() }

// PortAddress is the set of bytes representing a port number.  Users may
// get the port number for TCP/UDP layers by requesting it directly from the
// TCP or UDP layer (SrcPort/DstPort), but this address allows us to treat
// the port as a protocol-agnostic address for an application on a system.
type PortAddress []byte

func (a PortAddress) Raw() []byte    { return a }
func (a PortAddress) String() string { return strconv.Itoa(int(binary.BigEndian.Uint16(a))) }

// These layers correspond to Internet Protocol Suite (TCP/IP) layers, and their
// corresponding OSI layers, as best as possible.

// LinkLayer is the packet layer corresponding to TCP/IP layer 1 (OSI layer 2)
type LinkLayer interface {
	Layer
	SrcLinkAddr() Address
	DstLinkAddr() Address
}

// NetworkLayer is the packet layer corresponding to TCP/IP layer 2 (OSI
// layer 3)
type NetworkLayer interface {
	Layer
	SrcNetAddr() Address
	DstNetAddr() Address
}

// TransportLayer is the packet layer corresponding to the TCP/IP layer 3 (OSI
// layer 4)
type TransportLayer interface {
	Layer
	SrcAppAddr() Address
	DstAppAddr() Address
}

// ApplicationLayer is the packet layer corresponding to the TCP/IP layer 4 (OSI
// layer 7), also known as the packet payload.
type ApplicationLayer interface {
	Layer
	Payload() []byte
}

// ErrorLayer is a packet layer created when decoding of the packet has failed.
// Its payload is all the bytes that we were unable to decode, and the returned
// error details why the decoding failed.
type ErrorLayer interface {
	Layer
	Payload() []byte
	Error() error
}