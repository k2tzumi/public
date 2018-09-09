// Copyright 2018 github.com/ucirello and https://cirello.io. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

// Command mluqm runs a AI-powered Ur-Quan Masters compatible client.
package main

// From ./src-c/uqm-0.7.0-1/src/uqm/supermelee/netplay/packet.h

type PacketType uint16

const (
	PacketTypeInit PacketType = iota
	PacketTypePing
	PacketTypeAck
	PacketTypeReady
	PacketTypeFleet
	PacketTypeTeamname
	PacketTypeHandshake0
	PacketTypeHandshake1
	PacketTypeHandshakecancel
	PacketTypeHandshakecancelack
	PacketTypeSeedrandom
	PacketTypeInputdelay
	PacketTypeSelectship
	PacketTypeBattleinput
	PacketTypeFramecount
	PacketTypeChecksum
	PacketTypeAbort
	PacketTypeReset
	PacketTypeNum //Number of packet types
)

func (p PacketType) String() string {
	switch p {
	case PacketTypeInit:
		return "PacketTypeInit"
	case PacketTypePing:
		return "PacketTypePing"
	case PacketTypeAck:
		return "PacketTypeAck"
	case PacketTypeReady:
		return "PacketTypeReady"
	case PacketTypeFleet:
		return "PacketTypeFleet"
	case PacketTypeTeamname:
		return "PacketTypeTeamname"
	case PacketTypeHandshake0:
		return "PacketTypeHandshake0"
	case PacketTypeHandshake1:
		return "PacketTypeHandshake1"
	case PacketTypeHandshakecancel:
		return "PacketTypeHandshakecancel"
	case PacketTypeHandshakecancelack:
		return "PacketTypeHandshakecancelack"
	case PacketTypeSeedrandom:
		return "PacketTypeSeedrandom"
	case PacketTypeInputdelay:
		return "PacketTypeInputdelay"
	case PacketTypeSelectship:
		return "PacketTypeSelectship"
	case PacketTypeBattleinput:
		return "PacketTypeBattleinput"
	case PacketTypeFramecount:
		return "PacketTypeFramecount"
	case PacketTypeChecksum:
		return "PacketTypeChecksum"
	case PacketTypeAbort:
		return "PacketTypeAbort"
	case PacketTypeReset:
		return "PacketTypeReset"
	case PacketTypeNum:
		return "PacketTypeNum"
	}
	return "Unknown"
}

type netplaySide uint8

// This enum is used to indicate that a packet containing it relates to
// either the local or the remote player, from the perspective of the
// sender of the message
const (
	netplaySideLocal netplaySide = iota
	netplaySideRemote
)

// From ./src-c/uqm-0.7.0-1/src/uqm/supermelee/meleeship.h
type ship uint8

const (
	meleeAndrosynth ship = iota
	meleeArilou
	meleeChenjesu
	meleeChmmr
	meleeDruuge
	meleeEarthling
	meleeIlwrath
	meleeKohrAh
	meleeMelnorme
	meleeMmrnmhrm
	meleeMycon
	meleeOrz
	meleePkunk
	meleeShofixti
	meleeSlylandro
	meleeSpathi
	meleeSupox
	meleeSyreen
	meleeThraddash
	meleeUmgah
	meleeUrquan
	meleeUtwig
	meleeVux
	meleeYehat
	meleeZoqfotpik

	meleeNone ship = 0xff
)

func (s ship) String() string {
	switch s {
	case meleeAndrosynth:
		return "meleeAndrosynth"
	case meleeArilou:
		return "meleeArilou"
	case meleeChenjesu:
		return "meleeChenjesu"
	case meleeChmmr:
		return "meleeChmmr"
	case meleeDruuge:
		return "meleeDruuge"
	case meleeEarthling:
		return "meleeEarthling"
	case meleeIlwrath:
		return "meleeIlwrath"
	case meleeKohrAh:
		return "meleeKohrAh"
	case meleeMelnorme:
		return "meleeMelnorme"
	case meleeMmrnmhrm:
		return "meleeMmrnmhrm"
	case meleeMycon:
		return "meleeMycon"
	case meleeOrz:
		return "meleeOrz"
	case meleePkunk:
		return "meleePkunk"
	case meleeShofixti:
		return "meleeShofixti"
	case meleeSlylandro:
		return "meleeSlylandro"
	case meleeSpathi:
		return "meleeSpathi"
	case meleeSupox:
		return "meleeSupox"
	case meleeSyreen:
		return "meleeSyreen"
	case meleeThraddash:
		return "meleeThraddash"
	case meleeUmgah:
		return "meleeUmgah"
	case meleeUrquan:
		return "meleeUrquan"
	case meleeUtwig:
		return "meleeUtwig"
	case meleeVux:
		return "meleeVux"
	case meleeYehat:
		return "meleeYehat"
	case meleeZoqfotpik:
		return "meleeZoqfotpik"
	case meleeNone:
		return "meleeNone"
	}
	return "Unknown"
}
