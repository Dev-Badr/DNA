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

package discovery

import (
	"testing"
	"time"

	"github.com/skytells-research/DNA/network/node/market"
	"github.com/skytells-research/DNA/network/node/money"
	"github.com/skytells-research/DNA/network/node/services/openvpn/discovery/dto"
	"github.com/stretchr/testify/assert"
)

var (
	locationLTTelia = market.Location{"LT", "Vilnius", "AS8764"}
	protocol        = "tcp"
)

func Test_NewServiceProposalWithLocation(t *testing.T) {
	proposal := NewServiceProposalWithLocation(locationLTTelia, protocol)

	assert.Exactly(
		t,
		market.ServiceProposal{
			ServiceType: "openvpn",
			ServiceDefinition: dto.ServiceDefinition{
				Location:          locationLTTelia,
				LocationOriginate: locationLTTelia,
				SessionBandwidth:  83886080,
				Protocol:          "tcp",
			},

			PaymentMethodType: "PER_TIME",
			PaymentMethod: dto.PaymentPerTime{
				Price:    money.Money{12500000, money.Currency("MYST")},
				Duration: 60 * time.Minute,
			},
		},
		proposal,
	)
}
