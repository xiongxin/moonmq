package proto

import (
	"encoding/binary"
	"encoding/json"
	"errors"
)

var (
	ErrInvalidBuf = errors.New("invalid decode buf")
	ErrBufShort   = errors.New("decode buf is too short")
)

/*
   Proto binary format is

   |total length(4 bytes)|header length(4 bytes)|header json|body|

   total length = 4 + len(header json) + len(body)
   header length = len(header json)
*/

type Proto struct {
	Version uint32 `json:"version"`

	Method uint32 `json:"method"`

	Fields map[string]string `json:"fields"`

	Body []byte `json:"-"`
}

func NewProto() *Proto {
	p := new(Proto)

	p.Fields = map[string]string{}
	p.Body = []byte{}
	return p
}

func Marshal(p *Proto) ([]byte, error) {
	header, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	length := 4 + len(header) + len(p.Body)
	buf := make([]byte, 8, length)

	binary.BigEndian.PutUint32(buf[0:4], uint32(length))
	binary.BigEndian.PutUint32(buf[4:8], uint32(len(header)))

	buf = append(buf, header...)
	buf = append(buf, p.Body...)

	return buf, nil
}

func Unmarshal(buf []byte, p *Proto) error {
	if len(buf) <= 8 {
		return ErrBufShort
	}

	totalLen := binary.BigEndian.Uint32(buf[0:4])
	headerLen := binary.BigEndian.Uint32(buf[4:8])

	if uint32(len(buf)) != (totalLen + 4) {
		return ErrInvalidBuf
	}

	if headerLen > (totalLen - 4) {
		return ErrInvalidBuf
	}

	if err := json.Unmarshal(buf[8:8+headerLen], p); err != nil {
		return err
	}

	if p.Fields == nil {
		p.Fields = map[string]string{}
	}

	if p.Body == nil {
		p.Body = []byte{}
	}

	p.Body = buf[8+headerLen : 4+totalLen]
	return nil
}