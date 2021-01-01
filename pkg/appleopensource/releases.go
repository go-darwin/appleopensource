// Copyright 2016 The appleopensource Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package appleopensource

// Platform represents a release version platform type.
type Platform int

const (
	// Unknown unknown platform.
	Unknown Platform = iota
	// MacOS is a macOS platform.
	MacOS
	// Xcode is a Developer Tools(Xcode) platform.
	Xcode
	// IOS is a iOS platform.
	IOS
	// Server is a macOS Server platform.
	Server
)

func (p Platform) String() string {
	switch p {
	case MacOS:
		return "macos"
	case Xcode:
		return "developer-tools"
	case IOS:
		return "ios"
	case Server:
		return "server"
	default:
		return ""
	}
}

// KnownRelease known release versions.
var KnownRelease = [...][]string{
	MacOS:  releaseMacOS,
	Xcode:  releaseXcode,
	IOS:    releaseIOS,
	Server: releaseServer,
}

var (
	releaseMacOS = []string{
		"10.12.3",
		"10.12.2",
		"10.12.1",
		"10.12",
		"10.11.6",
		"10.11.5",
		"10.11.4",
		"10.11.3",
		"10.11.3",
		"10.11.1",
		"10.11",
		"10.10.5",
		"10.10.4",
		"10.10.3",
		"10.10.3",
		"10.10.1",
		"10.10",
		"10.9.5",
		"10.9.4",
		"10.9.3",
		"10.9.3",
		"10.9.1",
		"10.9",
		"10.8.5",
		"10.8.4",
		"10.8.3",
		"10.8.3",
		"10.8.1",
		"10.8",
		"10.7.5",
		"10.7.4",
		"10.7.3",
		"10.7.3",
		"10.7.1",
		"10.7",
		"10.6.8",
		"10.6.7",
		"10.6.6",
		"10.6.5",
		"10.6.4",
		"10.6.3",
		"10.6.3",
		"10.6.1",
		"10.6",
		"10.5.8",
		"10.5.7",
		"10.5.6",
		"10.5.5",
		"10.5.4",
		"10.5.3",
		"10.5.3",
		"10.5.1",
		"10.5",
		"10.4.11.x86",
		"10.4.11.ppc",
		"10.4.10.x86",
		"10.4.10.ppc",
		"10.4.9.x86",
		"10.4.9.ppc",
		"10.4.8.x86",
		"10.4.8.ppc",
		"10.4.7.x86",
		"10.4.7.ppc",
		"10.4.6.x86",
		"10.4.6.ppc",
		"10.4.5.x86",
		"10.4.5.ppc",
		"10.4.4.x86",
		"10.4.4.ppc",
		"10.4.3",
		"10.4.3",
		"10.4.1",
		"10.4",
		"10.3.9",
		"10.3.8",
		"10.3.7",
		"10.3.6",
		"10.3.5",
		"10.3.4",
		"10.3.3",
		"10.3.3",
		"10.3.1",
		"10.3",
		"10.2.8",
		"10.2.8.G5",
		"10.2.7",
		"10.2.6",
		"10.2.5",
		"10.2.4",
		"10.2.3",
		"10.2.3",
		"10.2.1",
		"10.2",
		"10.1.5",
		"10.1.4",
		"10.1.3",
		"10.1.3",
		"10.1.1",
		"10.1",
		"10.0.4",
		"10.0.3",
		"10.0.3",
		"10.0.1",
		"10.0",
	}

	releaseXcode = []string{
		"8.2.1",
		"8.1",
		"8.0",
		"7.3.1",
		"7.3",
		"7.2",
		"7.1",
		"7.0",
		"6.4",
		"6.3",
		"6.2",
		"6.1",
		"6.0",
		"5.1",
		"5.0",
		"4.6",
		"4.5",
		"4.4",
		"4.3",
		"4.2",
		"4.1",
		"4.0",
		"3.2.6",
		"3.2.5",
		"3.2.4",
		"3.2.3",
		"3.2.2",
		"3.2.1",
		"3.2",
		"3.1.4",
		"3.1.3",
		"3.1.2",
		"3.1.1",
		"3.1",
		"3.1b",
		"3.0",
		"2.5",
		"2.4.1",
		"2.4",
		"2.3",
		"2.2",
		"2.1",
		"WWDC2004DP",
		"WWDC2003DP",
		"Nov2004",
		"1.5",
		"1.2",
		"Jun2003",
		"Dec2002",
		"May2002",
		"Dec2001",
	}

	releaseIOS = []string{
		"10.2.1",
		"10.2",
		"10.1",
		"10.0",
		"9.2",
		"9.1",
		"9.0",
		"8.4.1",
		"8.4",
		"8.3",
		"8.2",
		"8.1.3",
		"8.1.2",
		"8.1.1",
		"8.1",
		"8.0.1",
		"8.0",
		"7.1.2",
		"7.1.1",
		"7.1",
		"7.0.3",
		"7.0",
		"6.1.3",
		"6.1",
		"6.0.1",
		"6.0",
		"5.1.1",
		"5.1",
		"5.0",
		"4.3",
		"4.3.3",
		"4.3.2",
		"4.3.1",
		"4.2",
		"4.1",
		"4.0",
		"3.2",
		"3.1.3",
		"3.1.2",
		"3.1.1",
		"3.1",
		"3.0",
		"2.2.1",
		"2.2",
		"2.1",
		"2.0",
		"SDKb8",
		"SDKb7",
		"SDKb6",
		"SDKb5",
		"SDKb4",
		"SDKb3",
		"SDKb2",
		"1.1.4",
		"1.1.3",
		"1.1.2",
		"1.1.1",
		"1.0.1",
		"1.0",
	}

	releaseServer = []string{
		"3.0.2",
		"2.2.2",
	}
)
