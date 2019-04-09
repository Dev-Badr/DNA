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

package connection

import (
	"github.comskytells-research/DNA/network/node/identity"
	"github.comskytells-research/DNA/network/node/market"
	"github.comskytells-research/DNA/network/node/session"
)

// ConnectParams holds plugin specific params
type ConnectParams struct {
	// kill switch option restricting communication only through VPN
	DisableKillSwitch bool
}

// ConnectOptions represents the params we need to ensure a successful connection
type ConnectOptions struct {
	ConsumerID    identity.Identity
	ProviderID    identity.Identity
	Proposal      market.ServiceProposal
	SessionID     session.ID
	SessionConfig []byte
}
