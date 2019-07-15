// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"bytes"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
)

type mjpegFrame []image.Image

type mjpegDevice struct {
	*cameraDevice
	streamStopped bool
}

// StartStream streams mjpeg video to w, which may be http.ResponseWriter.
// Implements VideoStreamer interface.
func (m *mjpegDevice) StartStream(w io.WriteCloser) (func() error, error) {
	const fsize = 1024 * 1024
	m.streamStopped = false
	go func() {
		var jpegbuff [fsize]byte
		var buff = bytes.NewBuffer(jpegbuff[:])

		// m.streamStopped is a sentinal val so we can stop stream
		// from diff goroutine.
		for !m.streamStopped {
			n, err := io.CopyN(buff, m.file, fsize)
			if err != nil || n != fsize {
				log.Printf("encountered error reading a frame from device")
				return
			}

			// img, err := jpeg.Decode(buff)
			// if err != nil {
			// 	return nil, err
			// }
			// buff.Reset()

			const boundary = `frame`
			if httprespwrt, ok := w.(http.ResponseWriter); ok {
				httprespwrt.Header().Set("Content-Type", `multipart/x-mixed-replace;boundary=`+boundary)
				multiwriter := multipart.NewWriter(w)
				multiwriter.SetBoundary(boundary)
				for {
					image := buff.Bytes()
					iw, err := multiwriter.CreatePart(textproto.MIMEHeader{
						"Content-type":   []string{"image/jpeg"},
						"Content-length": []string{strconv.Itoa(len(image))},
					})
					if err != nil {
						log.Println(err)
						break
					}
					_, err = iw.Write(image)
					if err != nil {
						log.Println(err)
						break
					}
				}
			}
		}
	}()
	// the point of this fuckery here is to run the stop stream in another
	// goroutine when called after returning it. Thats all.
	f := func() error {
		ec := make(chan error)
		go func() {
			ec <- m.StopStream()
		}()
		return <-ec
	}
	return f, nil
}

func (m *mjpegDevice) StopStream() error {
	m.streamStopped = true
	return nil
}
