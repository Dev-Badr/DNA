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
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseJSONOptions_HandlesNil(t *testing.T) {
	options, err := ParseJSONOptions(nil)

	assert.NoError(t, err)
	assert.Equal(t, DefaultOptions, options)
}

func Test_ParseJSONOptions_HandlesEmptyRequest(t *testing.T) {
	request := json.RawMessage(`{}`)
	options, err := ParseJSONOptions(&request)

	assert.NoError(t, err)
	assert.Equal(t, DefaultOptions, options)
}

func Test_ParseJSONOptions_ValidRequest(t *testing.T) {
	request := json.RawMessage(`{"connectDelay": 3000, "portMin":123, "portMax":1234, "subnet":"10.10.0.0/16"}`)
	options, err := ParseJSONOptions(&request)

	assert.NoError(t, err)
	assert.Equal(t, Options{
		ConnectDelay: 3000,
		PortMin:      123,
		PortMax:      1234,
		Subnet: net.IPNet{
			IP:   net.ParseIP("10.10.0.0").To4(),
			Mask: net.IPv4Mask(255, 255, 0, 0),
		},
	}, options)
}
