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

package wireguard

import (
	"encoding/json"
	"testing"

	"github.com/skytells-research/DNA/network/node/money"
	"github.com/stretchr/testify/assert"
)

func Test_PaymentMethod_Serialize(t *testing.T) {
	price := money.NewMoney(0.5, money.CurrencyMyst)

	var tests = []struct {
		model        Payment
		expectedJSON string
	}{
		{
			Payment{
				Price: price,
			},
			`{
				"price": {
					"amount": 50000000,
					"currency": "MYST"
				}
			}`,
		},
		{
			Payment{},
			`{
				"price": {}
			}`,
		},
	}

	for _, test := range tests {
		jsonBytes, err := json.Marshal(test.model)

		assert.Nil(t, err)
		assert.JSONEq(t, test.expectedJSON, string(jsonBytes))
	}
}

func Test_PaymentMethod_Unserialize(t *testing.T) {
	price := money.NewMoney(0.5, money.CurrencyMyst)

	var tests = []struct {
		json          string
		expectedModel Payment
		expectedError error
	}{
		{
			`{
				"price": {
					"amount": 50000000,
					"currency": "MYST"
				}
			}`,
			Payment{
				Price: price,
			},
			nil,
		},
		{
			`{
				"price": {}
			}`,
			Payment{},
			nil,
		},
		{
			`{}`,
			Payment{},
			nil,
		},
	}

	for _, test := range tests {
		var model Payment
		err := json.Unmarshal([]byte(test.json), &model)

		assert.Equal(t, test.expectedModel, model)
		assert.Equal(t, test.expectedError, err)
	}
}
