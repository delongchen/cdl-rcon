package rcon

import (
	"bytes"
	"encoding/binary"
)

const (
	SERVERDATA_AUTH           int32 = 3
	SERVERDATA_AUTH_RESPONSE  int32 = 2
	SERVERDATA_EXECCOMMAND    int32 = 2
	SERVERDATA_RESPONSE_VALUE int32 = 0
)

type BasicPacket struct {
	ID   int32
	Type int32
	Body string
}

func EmptyPacket(id int32) *BasicPacket {
	return &BasicPacket{
		ID:   id,
		Type: SERVERDATA_RESPONSE_VALUE,
		Body: "",
	}
}

func CMDPacket(cmd string, id int32) *BasicPacket {
	return &BasicPacket{
		ID:   id,
		Body: cmd,
		Type: SERVERDATA_EXECCOMMAND,
	}
}

type CMDTask struct {
	Vec []*BasicPacket
	Ch  chan *BasicPacket
	CMD string
}

type CMDTaskMap map[int32]*CMDTask

func ReadUInt32(buffer []byte, offset *int) uint32 {
	ret := binary.LittleEndian.Uint32(buffer[*offset:])
	*offset += 4
	return ret
}

func ReadString(bytes []byte, offset *int) string {
	l := 0
	tmp := bytes[*offset:]
	for _, v := range tmp {
		*offset++
		if v == byte(0) {
			break
		}
		l++
	}
	return string(tmp[:l])
}

func (p *BasicPacket) ToBytes() []byte {
	bodyBytes := []byte(p.Body)
	buf := bytes.Buffer{}
	size := int32(len(bodyBytes) + 10)
	_ = binary.Write(&buf, binary.LittleEndian, size)
	_ = binary.Write(&buf, binary.LittleEndian, p.ID)
	_ = binary.Write(&buf, binary.LittleEndian, p.Type)
	buf.Write(bodyBytes)
	buf.WriteByte(0)
	buf.WriteByte(0)
	return buf.Bytes()
}

func FromBytes(b []byte) *BasicPacket {
	offset := 0

	_ = int32(ReadUInt32(b, &offset))
	id := int32(ReadUInt32(b, &offset))
	t := int32(ReadUInt32(b, &offset))
	body := ReadString(b, &offset)

	return &BasicPacket{
		ID:   id,
		Body: body,
		Type: t,
	}
}
