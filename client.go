// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package vidchat

import (
	"errors"
	"os"
	"syscall"
)

const BlockSize = 1024 * 1024 // Blocksize 1M

var ErrUninitializedDevice = errors.New("error unable to use uninitialized device")

type cameraDevice struct {
	name      string // possibly "/dev/video0"
	file      *os.File
	framerate int // in frames-per-second (fps)
	buffer    []frame
}

type frame [Blocksize]byte

func (c *cameraDevice) mjpegStream() (<-chan []frame, error) {
	if c.file == nil {
		return nil, ErrUninitializedDevice
	}
	var buflen = len(c.buffer) * c.framerate
	stream := make([]frame, buflen)

}

func OpenCamera(device string) (*cameraDevice, error) {
	if device == "" {
		device = "/dev/video0"
	}
	cam, err := os.OpenFile(device, syscall.O_RDWR|syscall.O_DIRECT|syscall.O_NONBLOCK, 0755)
	if err != nil {
		return nil, err
	}
	return &cameraDevice{
		name:   device,
		file:   cam,
		buffer: make([]frame, 5),
	}, nil

}
