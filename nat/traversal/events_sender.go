/*
 * Copyright (C) 2019 The "MysteriumNetwork/node" Authors.
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

package traversal

import (
	"fmt"

	log "github.com/cihub/seelog"
)

const eventsSenderLogPrefix = "[traversal-events-sender] "

// EventsSender allows subscribing to NAT events and sends them to server
type EventsSender struct {
	metricsSender metricsSender
	ipResolver    ipResolver
	lastIp        string
	lastEventType string
}

type metricsSender interface {
	SendNATMappingSuccessEvent() error
	SendNATMappingFailEvent(err error) error
}

type ipResolver func() (string, error)

// NewEventsSender returns a new instance of events sender
func NewEventsSender(metricsSender metricsSender, ipResolver ipResolver) *EventsSender {
	return &EventsSender{metricsSender: metricsSender, ipResolver: ipResolver, lastIp: ""}
}

// ConsumeNATEvent sends received event to server
func (es *EventsSender) ConsumeNATEvent(event Event) {
	publicIP, err := es.ipResolver()
	if err != nil {
		log.Warnf(eventsSenderLogPrefix, "resolving public ip failed: ", err)
		return
	}
	if publicIP == es.lastIp && event.Type == es.lastEventType {
		return
	}

	err = es.sendNATEvent(event)
	if err != nil {
		log.Warnf(eventsSenderLogPrefix, "sending event failed: ", err)
	}

	es.lastIp = publicIP
	es.lastEventType = event.Type
}

func (es *EventsSender) sendNATEvent(event Event) error {
	switch event.Type {
	case SuccessEventType:
		return es.metricsSender.SendNATMappingSuccessEvent()
	case FailureEventType:
		return es.metricsSender.SendNATMappingFailEvent(event.Error)
	default:
		return fmt.Errorf("unknown event type: %v", event.Type)
	}
}
