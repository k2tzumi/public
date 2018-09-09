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

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
)

// From ./src-c/uqm-0.7.0-1/src/uqm/supermelee/netplay/packet.h

func parsePackets(dir string, r io.Reader) (string, error) {
	var ret string
	var ph packetHeader
	err := binary.Read(r, binary.BigEndian, &ph)
	ret += fmt.Sprintf("%s: %#v (%s)\n", dir, ph, ph.PacketType)
	if err != nil {
		return ret, err
	}
	bodyBuf := make([]byte, ph.Length-packetHeaderLength)
	r.Read(bodyBuf)
	ret += fmt.Sprintf("%s: %s\n", dir, hex.Dump(bodyBuf))
	ret += "----\n"
	return ret, nil
}

const packetHeaderLength = 4

type packetHeader struct {
	Length     uint16
	PacketType PacketType
}

type basePacket struct {
	header packetHeader
}

type PacketType uint16

const (
	PacketInit PacketType = iota
	PacketPing
	PacketAck
	PacketReady
	PacketFleet
	PacketTeamname
	PacketHandshake0
	PacketHandshake1
	PacketHandshakecancel
	PacketHandshakecancelack
	PacketSeedrandom
	PacketInputdelay
	PacketSelectship
	PacketBattleinput
	PacketFramecount
	PacketChecksum
	PacketAbort
	PacketReset
	PacketNum //Number of packet types
)

func (p PacketType) String() string {
	switch p {
	case PacketInit:
		return "PacketInit"
	case PacketPing:
		return "PacketPing"
	case PacketAck:
		return "PacketAck"
	case PacketReady:
		return "PacketReady"
	case PacketFleet:
		return "PacketFleet"
	case PacketTeamname:
		return "PacketTeamname"
	case PacketHandshake0:
		return "PacketHandshake0"
	case PacketHandshake1:
		return "PacketHandshake1"
	case PacketHandshakecancel:
		return "PacketHandshakecancel"
	case PacketHandshakecancelack:
		return "PacketHandshakecancelack"
	case PacketSeedrandom:
		return "PacketSeedrandom"
	case PacketInputdelay:
		return "PacketInputdelay"
	case PacketSelectship:
		return "PacketSelectship"
	case PacketBattleinput:
		return "PacketBattleinput"
	case PacketFramecount:
		return "PacketFramecount"
	case PacketChecksum:
		return "PacketChecksum"
	case PacketAbort:
		return "PacketAbort"
	case PacketReset:
		return "PacketReset"
	case PacketNum:
		return "PacketNum"
	}
	return "Unknown"
}
