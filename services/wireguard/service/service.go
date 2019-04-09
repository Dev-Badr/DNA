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
	"github.comskytells-research/DNA/network/node/market"
	"github.comskytells-research/DNA/network/node/money"
	wg "github.comskytells-research/DNA/network/node/services/wireguard"
)

const logPrefix = "[service-wireguard] "

// GetProposal returns the proposal for wireguard service
func GetProposal(country string) market.ServiceProposal {
	return market.ServiceProposal{
		ServiceType: wg.ServiceType,
		ServiceDefinition: wg.ServiceDefinition{
			Location:          market.Location{Country: country},
			LocationOriginate: market.Location{Country: country},
		},
		PaymentMethodType: wg.PaymentMethod,
		PaymentMethod: wg.Payment{
			Price: money.NewMoney(0, money.CurrencyMyst),
		},
	}
}
