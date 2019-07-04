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

var (
	ErrUninitializedDevice = errors.New("error unable to use uninitialized device")
	ErrUnknownDevice       = errors.New("error incorrect or unknown device")
)

type VideoStreamer interface {
	StartStream() error
	StopStream() error
	Close() error
}

type cameraDevice struct {
	name      string // possibly "/dev/video0"
	file      *os.File
	framerate int // in frames-per-second (fps)
	buffer    []frame
}

type frame [BlockSize]byte

/*
NOTE: Using channels for the stream is a bad idea because the runtime schedules
them to run at its convenience so we can't count on a constant stream.

type MJPEG chan []frame

func (c *cameraDevice) intitMJPEG() (MJPEG, error) {
	if c.file == nil {
		return nil, ErrUninitializedDevice
	}
	var buflen = len(c.buffer) * c.framerate
	stream := make(MJPEG, buflen)
	defer func() {
		var nilframe []frame
		for {
			select {
			case stream <- nilframe:

			}
		}
	}()

}
*/

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

type socketDevice struct {
	name   string
	handle *os.File
}

func OpenSocketDevice(device string) (*socketDevice, error) {
	fd, err := os.OpenFile(device, 0755, os.ModeDevice|syscall.O_DIRECT|syscall.O_NONBLOCK)
	if err != nil {
		return nil, err
	}

}

func creatDeviceSocket(device VideoStreamer) (*socketDevice, error) {
	s, err := os.Create(os.TempDir())
	if err != nil {
		return nil, err
	}

	switch t := device.(type) {
	case *cameraDevice:
		d := &socketDevice{
			name:   t.name,
			fd:     t.file,
			socket: s,
		}
	default:
		return nil, ErrUnknownDevice
	}
}
