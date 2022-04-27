package protocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"math"
)

//todo 改成自定义包
var (
	ErrInvalidBuffer = errors.New("invalid protocol")
)

const (
	OpHeartbeat      = uint16(0)
	OpHeartbeatReply = uint16(1)

	OpAuth      = uint16(2)
	OpAuthReply = uint16(3)

	OpSendMsg      = uint16(4)
	OpSendMsgReply = uint16(5)

	OpSendRoomMsg      = uint16(6)
	OpSendMsgRoomReply = uint16(7)

	OpSendAreaMsg      = uint16(8)
	OpSendAreaMsgReply = uint16(9)

	OpDisconnect      = uint16(10)
	OpDisconnectReply = uint16(11)
)

const maxPackSize = math.MaxUint32

const (
	packSize     = 4
	headerSize   = 2
	versionSize  = 2
	checksumSize = 2
	opSize       = 4
	seqSize      = 4
	HeaderLen    = packSize + headerSize + versionSize + checksumSize + opSize + seqSize

	packSizeOffset     = packSize
	headerSizeOffset   = packSizeOffset + headerSize
	versionSizeOffset  = headerSizeOffset + versionSize
	checksumSizeOffset = versionSizeOffset + checksumSize
	opSizeOffset       = checksumSizeOffset + opSize
	seqSizeOffset      = opSizeOffset + seqSize
)

type Proto struct {
	Version  uint16 //2 byte 16 bit
	Op       uint16
	Checksum uint16
	Seq      uint16
	Data     []byte
}

// SerializeTo todo checksum
//传进来避免内存逃逸
func (p *Proto) SerializeTo(bytes []byte) (err error) {
	if len(bytes) < HeaderLen {
		return ErrInvalidBuffer
	}
	binary.BigEndian.PutUint32(bytes[0:packSizeOffset], uint32(len(p.Data)+HeaderLen))
	binary.BigEndian.PutUint16(bytes[packSizeOffset:headerSizeOffset], HeaderLen)
	binary.BigEndian.PutUint16(bytes[headerSizeOffset:versionSizeOffset], p.Version)
	binary.BigEndian.PutUint16(bytes[versionSizeOffset:checksumSizeOffset], p.Checksum)
	binary.BigEndian.PutUint16(bytes[checksumSizeOffset:opSizeOffset], p.Op)
	binary.BigEndian.PutUint16(bytes[opSizeOffset:seqSizeOffset], p.Seq)
	copy(bytes[seqSizeOffset:], p.Data)
	return
}

func (p *Proto) DecodeFromBytes(b *bufio.Reader) (err error) {
	bufio.NewReader(bytes.NewReader([]byte{}))
	headBuf := make([]byte, HeaderLen)
	if err = binary.Read(b, binary.BigEndian, &headBuf); err != nil {
		return
	}
	packLen := binary.BigEndian.Uint32(headBuf[0:packSizeOffset])
	//暂时没用到
	headerLen := binary.BigEndian.Uint16(headBuf[packSizeOffset:headerSizeOffset])
	if headerLen != HeaderLen {
		return ErrInvalidBuffer
	}
	p.Version = binary.BigEndian.Uint16(headBuf[headerSizeOffset:versionSizeOffset])
	p.Checksum = binary.BigEndian.Uint16(headBuf[versionSizeOffset:checksumSizeOffset])
	p.Op = binary.BigEndian.Uint16(headBuf[checksumSizeOffset:opSizeOffset])
	p.Seq = binary.BigEndian.Uint16(headBuf[opSizeOffset:seqSizeOffset])
	p.Data = make([]byte, packLen-HeaderLen)
	if err = binary.Read(b, binary.BigEndian, &p.Data); err != nil {
		return
	}
	return
}

// Checksum 默认不开启校验
func Checksum(data []byte, csum uint32) uint16 {
	length := len(data) - 1
	for i := 0; i < length; i += 2 {
		// For our test packet, doing this manually is about 25% faster
		// (740 ns vs. 1000ns) than doing it by calling binary.BigEndian.Uint16.
		csum += uint32(data[i]) << 8
		csum += uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		csum += uint32(data[length]) << 8
	}
	for csum > 0xffff {
		csum = (csum >> 16) + (csum & 0xffff)
	}
	return ^uint16(csum)
}

func ReadTcp() {

}
