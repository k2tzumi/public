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
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"cirello.io/errors"
	"github.com/davecgh/go-spew/spew"
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

	switch ph.PacketType {
	case PacketTypeInit:
		err := parsePacketInit(bodyBuf)
		spew.Dump("err PacketTypeInit:", err)
	case PacketTypePing:
		err := parsePacketPing(bodyBuf)
		spew.Dump("err PacketTypePing:", err)
	case PacketTypeAck:
		err := parsePacketAck(bodyBuf)
		spew.Dump("err PacketTypeAck:", err)
	case PacketTypeReady:
		err := parsePacketReady(bodyBuf)
		spew.Dump("err PacketTypeReady:", err)
	case PacketTypeFleet:
		out, err := parsePacketFleet(bodyBuf)
		spew.Dump("err PacketTypeFleet:", err)
		ret += out + "\n----\n"
	case PacketTypeTeamName:
		err := parsePacketTeamName(bodyBuf)
		spew.Dump("err PacketTypeTeamName:", err)
	case PacketTypeHandshake0:
		err := parsePacketHandshake0(bodyBuf)
		spew.Dump("err PacketTypeHandshake0:", err)
	case PacketTypeHandshake1:
		err := parsePacketHandshake1(bodyBuf)
		spew.Dump("err PacketTypeHandshake1:", err)
	case PacketTypeHandshakeCancel:
		err := parsePacketHandshakeCancel(bodyBuf)
		spew.Dump("err PacketTypeHandshakeCancel:", err)
	case PacketTypeHandshakeCancelAck:
		err := parsePacketHandshakeCancelAck(bodyBuf)
		spew.Dump("err PacketTypeHandshakeCancelAck:", err)
	case PacketTypeSeedRandom:
		err := parsePacketSeedRandom(bodyBuf)
		spew.Dump("err PacketTypeSeedRandom:", err)
	case PacketTypeInputDelay:
		err := parsePacketInputDelay(bodyBuf)
		spew.Dump("err PacketTypeInputDelay:", err)
	case PacketTypeSelectShip:
		err := parsePacketSelectShip(bodyBuf)
		spew.Dump("err PacketTypeSelectShip:", err)
	default:
		spew.Dump("unknown packet:", ph.PacketType)
	}
	return ret, nil
}

const packetHeaderLength = 4

type packetHeader struct {
	Length     uint16
	PacketType PacketType
}

type packetInit struct {
	// header packetHeader
	ProtoVersionMajor uint8
	ProtoVersionMinor uint8
	Padding0          uint16
	UqmVersionMajor   uint8
	UqmVersionMinor   uint8
	UqmVersionPatch   uint8
	Padding1          uint8
}

func parsePacketInit(buf []byte) error {
	r := bytes.NewBuffer(buf)
	var p packetInit
	err := binary.Read(r, binary.BigEndian, &p)
	if err != nil {
		return errors.E(err, "cannot parse init packet")
	}
	spew.Dump(p)
	return nil
}

type packetPing struct {
	// header packetHeader
	ID uint32
}

func parsePacketPing(buf []byte) error {
	r := bytes.NewBuffer(buf)
	var p packetPing
	err := binary.Read(r, binary.BigEndian, &p)
	if err != nil {
		return errors.E(err, "cannot parse ping packet")
	}
	spew.Dump(p)
	return nil
}

type packetAck struct {
	// header packetHeader
	ID uint32
}

func parsePacketAck(buf []byte) error {
	r := bytes.NewBuffer(buf)
	var p packetAck
	err := binary.Read(r, binary.BigEndian, &p)
	if err != nil {
		return errors.E(err, "cannot parse ack packet")
	}
	spew.Dump(p)
	return nil
}

type packetReady struct {
	// header packetHeader
	// no contents.
}

func parsePacketReady(buf []byte) error {
	spew.Dump("READY")
	return nil
}

type packetSeedRandom struct {
	// header packetHeader
	Seed uint32
}

func parsePacketSeedRandom(buf []byte) error {
	r := bytes.NewBuffer(buf)
	var p packetSeedRandom
	err := binary.Read(r, binary.BigEndian, &p)
	if err != nil {
		return errors.E(err, "cannot parse seed random packet")
	}
	spew.Dump(p)
	return nil
}

type packetInputDelay struct {
	// header packetHeader
	Delay uint32
}

func parsePacketInputDelay(buf []byte) error {
	r := bytes.NewBuffer(buf)
	var p packetInputDelay
	err := binary.Read(r, binary.BigEndian, &p)
	if err != nil {
		return errors.E(err, "cannot parse input delay packet")
	}
	spew.Dump(p)
	return nil
}

type packetTeamName struct {
	// header packetHeader
	fixed struct {
		Side    netplaySide
		Padding uint8
	}
	dynamic struct {
		// '\0' terminated.
		// Be sure to add padding to this structure to make it a
		// multiple of 4 bytes in length.
		Name []byte
	}
}

func parsePacketTeamName(buf []byte) error {
	r := bytes.NewBuffer(buf)
	var ptn packetTeamName
	err := binary.Read(r, binary.BigEndian, &ptn.fixed)
	if err != nil {
		return errors.E(err, "cannot parse fixed part of team name packet")
	}
	spew.Dump(ptn)
	return nil
}

/*
// Structure describing an update to a player's fleet.
// TODO: use strings as ship identifiers, instead of numbers,
// so that adding of new ships doesn't break this.
typedef struct {
	PacketHeader header;
	uint8 side;
	uint8 padding;
	uint16 numShips;
	FleetEntry ships[];
	// Be sure to add padding to this structure to make it a multiple of
	// 4 bytes in length.
} Packet_Fleet;
typedef struct {
	uint8 index;  // Position in the fleet
	uint8 ship;   // Ship type index; actually MeleeShip
} FleetEntry;
*/
type packetFleetEntry struct {
	Index uint8
	Ship  ship
}
type packetFleet struct {
	// header packetHeader
	fixed struct {
		Side     netplaySide
		Padding  uint8
		NumShips uint16
	}
	dynamic struct {
		FleetEntry []packetFleetEntry
	}
	// Be sure to add padding to this structure to make it a
	// multiple of 4 bytes in length.
}

func parsePacketFleet(buf []byte) (string, error) {
	r := bytes.NewBuffer(buf)
	var pf packetFleet
	err := binary.Read(r, binary.BigEndian, &pf.fixed)
	if err != nil {
		return "", errors.E(err, "cannot parse fixed part of fleet packet")
	}
	for i := uint16(0); i < pf.fixed.NumShips; i++ {
		var pfe packetFleetEntry
		err := binary.Read(r, binary.BigEndian, &pfe)
		if err != nil {
			return "", errors.E(err, "cannot parse dynamic part of fleet packet")
		}
		pf.dynamic.FleetEntry = append(pf.dynamic.FleetEntry, pfe)
	}
	ret := fmt.Sprintf("fleet: %v %v %v %v %v",
		pf.fixed.Side, pf.fixed.Padding, pf.fixed.NumShips,
		pf.dynamic.FleetEntry[0].Index, pf.dynamic.FleetEntry[0].Ship,
	)
	return ret, nil
}

type packetHandshake0 struct{}

func parsePacketHandshake0(buf []byte) error {
	spew.Dump("packetHandshake0")
	return nil
}

type packetHandshake1 struct{}

func parsePacketHandshake1(buf []byte) error {
	spew.Dump("acketHandshake1")
	return nil
}

type packetHandshakeCancel struct{}

func parsePacketHandshakeCancel(buf []byte) error {
	spew.Dump("packetHandshakeCancel")
	return nil
}

type packetHandshakeCancelAck struct{}

func parsePacketHandshakeCancelAck(buf []byte) error {
	spew.Dump("packetHandshakeCancelAck")
	return nil
}

type packetSelectShip struct {
	Ship    uint16
	Padding uint16
}

func parsePacketSelectShip(buf []byte) error {
	r := bytes.NewBuffer(buf)
	var p packetSelectShip
	err := binary.Read(r, binary.BigEndian, &p)
	if err != nil {
		return errors.E(err, "cannot parse select ship packet")
	}
	spew.Dump(p)
	return nil
}
