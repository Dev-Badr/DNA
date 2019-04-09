//+build !windows

/*
 * Copyright (C) 2019 2019 Skytells, Inc.
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

package service

import (
	"encoding/json"
	"sync"

	log "github.com/cihub/seelog"
	"github.com/skytells-research/DNA/network/node/core/location"
	"github.com/skytells-research/DNA/network/node/identity"
	"github.com/skytells-research/DNA/network/node/nat"
	wg "github.com/skytells-research/DNA/network/node/services/wireguard"
	"github.com/skytells-research/DNA/network/node/services/wireguard/endpoint"
	"github.com/skytells-research/DNA/network/node/services/wireguard/resources"
	"github.com/skytells-research/DNA/network/node/session"
	"github.com/pkg/errors"
)

// NewManager creates new instance of Wireguard service
func NewManager(
	location location.ServiceLocationInfo,
	natService nat.NATService,
	portMap func(port int) (releasePortMapping func()),
	options Options) *Manager {

	resourceAllocator := resources.NewAllocator(options.PortMin, options.PortMax, options.Subnet)
	return &Manager{
		natService: natService,

		publicIP:        location.PubIP,
		outboundIP:      location.OutIP,
		currentLocation: location.Country,

		connectionEndpointFactory: func() (wg.ConnectionEndpoint, error) {
			return endpoint.NewConnectionEndpoint(location, resourceAllocator, portMap, options.ConnectDelay)
		},
	}
}

// Manager represents an instance of Wireguard service
type Manager struct {
	wg         sync.WaitGroup
	natService nat.NATService

	connectionEndpointFactory func() (wg.ConnectionEndpoint, error)

	publicIP        string
	outboundIP      string
	currentLocation string
}

// ProvideConfig provides the config for consumer
func (manager *Manager) ProvideConfig(publicKey json.RawMessage) (session.ServiceConfiguration, session.DestroyCallback, error) {
	key := &wg.ConsumerConfig{}
	err := json.Unmarshal(publicKey, key)
	if err != nil {
		return nil, nil, err
	}

	connectionEndpoint, err := manager.connectionEndpointFactory()
	if err != nil {
		return nil, nil, err
	}

	if err := connectionEndpoint.Start(nil); err != nil {
		return nil, nil, err
	}

	if err := connectionEndpoint.AddPeer(key.PublicKey, nil); err != nil {
		return nil, nil, err
	}

	config, err := connectionEndpoint.Config()
	if err != nil {
		return nil, nil, err
	}

	natRule := nat.RuleForwarding{SourceAddress: config.Consumer.IPAddress.String(), TargetIP: manager.outboundIP}
	if err := manager.natService.Add(natRule); err != nil {
		return nil, nil, errors.Wrap(err, "failed to add NAT forwarding rule")
	}

	destroy := func() {
		if err := manager.natService.Del(natRule); err != nil {
			log.Error(logPrefix, "failed to delete NAT forwarding rule: ", err)
		}
		if err := connectionEndpoint.Stop(); err != nil {
			log.Error(logPrefix, "failed to stop connection endpoint: ", err)
		}
	}

	return config, destroy, nil
}

// Serve starts service - does block
func (manager *Manager) Serve(providerID identity.Identity) error {
	manager.wg.Add(1)
	log.Info(logPrefix, "Wireguard service started successfully")

	manager.wg.Wait()
	return nil
}

// Stop stops service.
func (manager *Manager) Stop() error {
	manager.wg.Done()

	log.Info(logPrefix, "Wireguard service stopped")
	return nil
}
