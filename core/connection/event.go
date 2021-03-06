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

// Topic represents the different topics a consumer can subscribe to
const (
	// StateEventTopic represents the connection state change topic
	StateEventTopic = "State"
	// StatisticsEventTopic represents the connection stats topic
	StatisticsEventTopic = "Statistics"
	// SessionEventTopic represents the session event
	SessionEventTopic = "Session"
)

// StateEvent is the struct we'll emit on a StateEvent topic event
type StateEvent struct {
	State       State
	SessionInfo SessionInfo
}

const (
	// SessionCreatedStatus represents a session creation event
	SessionCreatedStatus = "Created"
	// SessionEndedStatus represents a session end
	SessionEndedStatus = "Ended"
)

// SessionEvent represents a session related event
type SessionEvent struct {
	Status      string
	SessionInfo SessionInfo
}
