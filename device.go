// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"errors"
	"io"
	"os"
	"syscall"
)

const BlockSize = 1024 * 1024 // Blocksize 1M

var (
	ErrUninitializedDevice = errors.New("error unable to use uninitialized device")
	ErrUnknownDevice       = errors.New("error incorrect or unknown device")
)

type VideoStreamer interface {
	// StartStream starts a video stream, and returns its a func that stops
	// the stream when called, and an error. If a non-nil error is returned,
	// the returned function will be nil, so DONT CALL IT!!! IT WILL PANIC!
	//
	// Implementations are encouraged to return their own StopStream func,
	// so that the caller can use the returned func to cancell the stream.
	//
	// Example:
	//
	// 	var s Videostreamer
	// 	stopStream, err := s.StartStream()
	// 	/* check err */
	//
	// 	// *** Do stuff here then ***`
	//
	// 	if err := stopStream(); err != nil {
	//		log.Println(err.Error())
	// 	}
	StartStream(io.WriteCloser) (func() error, error)
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
	d := &cameraDevice{
		name:   device,
		file:   cam,
		buffer: make([]frame, 5),
	}
	return d, nil

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
	d := &socketDevice{
		name:   device,
		handle: fd,
	}
	return d, nil
}
