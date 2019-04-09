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

package connection

import (
	"github.com/skytells-research/DNA/network/node/services/wireguard/resources"
	"github.com/skytells-research/DNA/network/node/services/wireguard/service"
)

func connectionResourceAllocator() *resources.Allocator {
	// Resource allocator uses config received from the provider. No configuration options required, passing default ones.
	return resources.NewAllocator(service.DefaultOptions.PortMin, service.DefaultOptions.PortMax, service.DefaultOptions.Subnet)
}
