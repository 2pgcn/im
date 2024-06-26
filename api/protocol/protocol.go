package protocol

import (
	"bufio"
	"encoding/binary"
	"errors"
	"github.com/2pgcn/gameim/pkg/safe"
	"github.com/golang/protobuf/proto"
	"math"
	"sync"
)

// todo 改成自定义包
var (
	ErrInvalidBuffer = errors.New("invalid protocol")
	ErrEOFData       = errors.New("IO EOF")
)
var Version = uint16(1)

// todo 添加到配置里
var ProtoPool *safe.Pool[*Proto]
var ones sync.Once

// todo 将bufio重写,共
// var bytesPool *safe.BytePool
var headerPool *sync.Pool
var dataPool *sync.Pool
var defaultSize = 1024
var defaultNum = 1024

func init() {
	ones.Do(func() {
		ProtoPool = safe.NewPool(func() *Proto {
			return &Proto{
				Version:  1,
				Op:       0,
				Checksum: 0,
				Seq:      0,
				Data:     nil,
			}
		})
		ProtoPool.Grow(defaultNum)
		headerPool = &sync.Pool{
			New: func() any {
				return make([]byte, HeaderLen)
			}}
		dataPool = &sync.Pool{New: func() any {
			return make([]byte, defaultSize)
		}}
	})
}

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

	OpErrReply = uint16(12)
	OpAck      = uint16(13)
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

var (
	AllRecvCount int64 = 0
	AllSendCount int64 = 0
)

func (p *Proto) Reset() {
	p.Version = 1
	p.Op = 0
	p.Checksum = 0
	p.Seq = 0
	p.Data = []byte{}
}

// SerializeTo todo checksum
// 传进来避免内存逃逸
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

// WriteTcp default Retry 3
func (p *Proto) WriteTcp(writer *bufio.Writer) (err error) {
	//operation := func() error {
	//	return p.writeTcp(writer)
	//}
	//err = backoff.Retry(operation, backoff.NewExponentialBackOff())
	return p.writeTcp(writer)
}

func (p *Proto) writeTcp(writer *bufio.Writer) (err error) {
	b := make([]byte, HeaderLen+len(p.Data))
	if len(b) < HeaderLen {
		return ErrInvalidBuffer
	}
	binary.BigEndian.PutUint32(b[0:packSizeOffset], uint32(len(p.Data)+HeaderLen))
	binary.BigEndian.PutUint16(b[packSizeOffset:headerSizeOffset], HeaderLen)
	binary.BigEndian.PutUint16(b[headerSizeOffset:versionSizeOffset], p.Version)
	binary.BigEndian.PutUint16(b[versionSizeOffset:checksumSizeOffset], p.Checksum)
	binary.BigEndian.PutUint16(b[checksumSizeOffset:opSizeOffset], p.Op)
	binary.BigEndian.PutUint16(b[opSizeOffset:seqSizeOffset], p.Seq)
	copy(b[seqSizeOffset:], p.Data)
	if err = binary.Write(writer, binary.BigEndian, b); err != nil {
		return errors.Join(errors.New("writer error"), err)
	}
	return writer.Flush()
}

func (p *Proto) WriteTcpNotFlush(writer *bufio.Writer) (err error) {
	b := make([]byte, HeaderLen+len(p.Data))
	if len(b) < HeaderLen {
		return ErrInvalidBuffer
	}
	binary.BigEndian.PutUint32(b[0:packSizeOffset], uint32(len(p.Data)+HeaderLen))
	binary.BigEndian.PutUint16(b[packSizeOffset:headerSizeOffset], HeaderLen)
	binary.BigEndian.PutUint16(b[headerSizeOffset:versionSizeOffset], p.Version)
	binary.BigEndian.PutUint16(b[versionSizeOffset:checksumSizeOffset], p.Checksum)
	binary.BigEndian.PutUint16(b[checksumSizeOffset:opSizeOffset], p.Op)
	binary.BigEndian.PutUint16(b[opSizeOffset:seqSizeOffset], p.Seq)
	copy(b[seqSizeOffset:], p.Data)
	if err = binary.Write(writer, binary.BigEndian, b); err != nil {
		return errors.Join(errors.New("writer error"), err)
	}
	return
}

func (p *Proto) DecodeFromBytes(b *bufio.Reader) (err error) {
	//headBuf := headerPool.Get().([]byte)
	//defer headerPool.Put(headBuf)
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
	//超出bytesPool长度
	p.Data = make([]byte, packLen-uint32(headerLen))
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

//// SetError todo 循环依赖
//func (p *Proto) SetError(err *error2.Error, op uint16) {
//	p.Op = op
//	p.Data = nil
//}

func NewProtoMsg(op uint16, msg proto.Message) (*Proto, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	m := ProtoPool.Get()
	m.Op = op
	m.Data = data
	return m, nil
}
