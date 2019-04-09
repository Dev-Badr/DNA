/*
 * Copyright (C) 2019 Skytells, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package endpoint

import (
	"fmt"
	"net"

	log "github.com/cihub/seelog"
	"github.com/skytells-research/DNA/network/node/core/location"
	wg "github.com/skytells-research/DNA/network/node/services/wireguard"
	"github.com/skytells-research/DNA/network/node/services/wireguard/key"
	"github.com/skytells-research/DNA/network/node/services/wireguard/resources"
)

const logPrefix = "[wireguard-connection-endpoint] "

type wgClient interface {
	ConfigureDevice(name string, config wg.DeviceConfig, subnet net.IPNet) error
	ConfigureRoutes(iface string, ip net.IP) error
	DestroyDevice(name string) error
	AddPeer(name string, peer wg.PeerInfo, allowedIP ...string) error
	RemovePeer(name string, publicKey string) error
	PeerStats() (wg.Stats, error)
	Close() error
}

type connectionEndpoint struct {
	iface              string
	privateKey         string
	location           location.ServiceLocationInfo
	ipAddr             net.IPNet
	endpoint           net.UDPAddr
	resourceAllocator  *resources.Allocator
	wgClient           wgClient
	releasePortMapping func()
	mapPort            func(port int) (releasePortMapping func())
	connectDelay       int // connect delay in milliseconds
}

// Start starts and configure wireguard network interface for providing service.
// If config is nil, required options will be generated automatically.
func (ce *connectionEndpoint) Start(config *wg.ServiceConfig) error {
	if err := ce.cleanAbandonedInterfaces(); err != nil {
		return err
	}

	iface, err := ce.resourceAllocator.AllocateInterface()
	if err != nil {
		return err
	}

	port, err := ce.resourceAllocator.AllocatePort()
	if err != nil {
		return err
	}

	ce.iface = iface
	ce.endpoint.Port = port
	ce.endpoint.IP = net.ParseIP(ce.location.PubIP)

	var deviceConfig deviceConfig
	if config == nil {
		// nil config mean its a provider Start
		ce.releasePortMapping = ce.mapPort(port)
		privateKey, err := key.GeneratePrivateKey()
		if err != nil {
			return err
		}
		ipAddr, err := ce.resourceAllocator.AllocateIPNet()
		if err != nil {
			return err
		}
		ce.ipAddr = ipAddr
		ce.ipAddr.IP = providerIP(ce.ipAddr)
		ce.privateKey = privateKey
		deviceConfig.listenPort = ce.endpoint.Port
	} else {
		ce.ipAddr = config.Consumer.IPAddress
		ce.privateKey = config.Consumer.PrivateKey
	}

	deviceConfig.privateKey = ce.privateKey
	return ce.wgClient.ConfigureDevice(ce.iface, deviceConfig, ce.ipAddr)
}

// AddPeer adds new wireguard peer to the wireguard network interface.
func (ce *connectionEndpoint) AddPeer(publicKey string, endpoint *net.UDPAddr, allowedIP ...string) error {
	return ce.wgClient.AddPeer(ce.iface, peerInfo{endpoint, publicKey}, allowedIP...)
}

// RemovePeer removes a wireguard peer from the wireguard network interface.
func (ce *connectionEndpoint) RemovePeer(publicKey string) error {
	return ce.wgClient.RemovePeer(ce.iface, publicKey)
}

// PeerStats returns stats information about connected peer.
func (ce *connectionEndpoint) PeerStats() (wg.Stats, error) {
	return ce.wgClient.PeerStats()
}

// Config provides wireguard service configuration for the current connection endpoint.
func (ce *connectionEndpoint) Config() (wg.ServiceConfig, error) {
	publicKey, err := key.PrivateKeyToPublicKey(ce.privateKey)
	if err != nil {
		return wg.ServiceConfig{}, err
	}

	var config wg.ServiceConfig
	config.Provider.PublicKey = publicKey
	config.Provider.Endpoint = ce.endpoint
	config.Consumer.IPAddress = ce.ipAddr
	config.Consumer.IPAddress.IP = ce.consumerIP(ce.ipAddr)
	if ce.location.OutIP != ce.location.PubIP {
		config.Consumer.ConnectDelay = ce.connectDelay
	}
	return config, nil
}

func (ce *connectionEndpoint) ConfigureRoutes(ip net.IP) error {
	return ce.wgClient.ConfigureRoutes(ce.iface, ip)
}

// Stop closes wireguard client and destroys wireguard network interface.
func (ce *connectionEndpoint) Stop() error {
	ce.releasePortMapping()

	if err := ce.wgClient.Close(); err != nil {
		return err
	}

	if err := ce.resourceAllocator.ReleasePort(ce.endpoint.Port); err != nil {
		return err
	}

	if err := ce.resourceAllocator.ReleaseIPNet(ce.ipAddr); err != nil {
		return err
	}

	return ce.resourceAllocator.ReleaseInterface(ce.iface)
}

func (ce *connectionEndpoint) cleanAbandonedInterfaces() error {
	ifaces, err := ce.resourceAllocator.AbandonedInterfaces()
	if err != nil {
		return err
	}

	for _, iface := range ifaces {
		if err := ce.wgClient.DestroyDevice(iface.Name); err != nil {
			log.Warn(logPrefix, fmt.Sprintf("failed to destroy abandoned interface: %s, error: %v", iface.Name, err))
		}
		log.Info(logPrefix, "abandoned interface destroyed: ", iface.Name)
	}

	return nil
}

type deviceConfig struct {
	privateKey string
	listenPort int
}

func (d deviceConfig) PrivateKey() string {
	return d.privateKey
}

func (d deviceConfig) ListenPort() int {
	return d.listenPort
}

type peerInfo struct {
	endpoint  *net.UDPAddr
	publicKey string
}

func (p peerInfo) Endpoint() *net.UDPAddr {
	return p.endpoint
}
func (p peerInfo) PublicKey() string {
	return p.publicKey
}

func providerIP(subnet net.IPNet) net.IP {
	subnet.IP[len(subnet.IP)-1] = byte(1)
	return subnet.IP
}
