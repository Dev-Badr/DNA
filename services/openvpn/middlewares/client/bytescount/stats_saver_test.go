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

package bytescount

import (
	"testing"

	"github.comskytells-research/DNA/network/go-openvpn/openvpn/middlewares/client/bytescount"
	"github.comskytells-research/DNA/network/node/consumer"

	"github.com/stretchr/testify/assert"
)

func TestNewSessionStatsSaver(t *testing.T) {
	channel := make(chan consumer.SessionStatistics, 1)
	saver := NewSessionStatsSaver(channel)
	stats := consumer.SessionStatistics{BytesSent: 1, BytesReceived: 2}
	saver(bytescount.Bytecount{BytesOut: 1, BytesIn: 2})
	assert.Equal(t, stats, <-channel)
}
