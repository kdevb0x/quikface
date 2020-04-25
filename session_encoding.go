package quikface

import (
	"bytes"
	"encoding/binary"
	"log"
)

type PreviousClientSession struct {
	Client *Client
	// checksum is the crc32sum
	Checksum []byte
}

func (ps *PreviousClientSession) MarshalBinary() ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.LittleEndian, ps.Client); err != nil {
		log.Printf("error marshaling previous client session: %w\n", err)
		return nil, err
	}

	return buff.Bytes, nil

}

func (ps *PreviousClientSession) UnmarshalBinary(data []byte) error {
	if ps.values == nil || len(ps.values) == 0 {

	}
}
