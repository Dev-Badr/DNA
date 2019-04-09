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
	"github.comskytells-research/DNA/network/node/core/service"
	"github.comskytells-research/DNA/network/node/services/noop"
	"github.comskytells-research/DNA/network/node/services/openvpn"
	openvpn_service "github.comskytells-research/DNA/network/node/services/openvpn/service"
	"github.comskytells-research/DNA/network/node/services/wireguard"
	wireguard_service "github.comskytells-research/DNA/network/node/services/wireguard/service"
	"github.com/urfave/cli"
)

var (
	serviceTypes = []string{"openvpn", "wireguard", "noop"}

	serviceTypesFlagsParser = map[string]func(ctx *cli.Context) service.Options{
		noop.ServiceType:      noop.ParseFlags,
		openvpn.ServiceType:   openvpn_service.ParseFlags,
		wireguard.ServiceType: wireguard_service.ParseFlags,
	}
)
