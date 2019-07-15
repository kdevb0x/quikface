// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import "github.com/pion/webrtc/v2"

var _ webrtc.Configuration

var wrtcCertificate = generateWebRTCCertificate()

// func generateWebRTCCertificate() wrtc.
