// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

// Example service program that beeps.
//
// The program demonstrates how to create Windows service and
// install / remove it on a computer. It also shows how to
// stop / start / pause / continue any service, and how to
// write to event log. It also shows how to use debug
// facilities available in debug package.
package main

import (
	"flag"
	"log"

	"golang.org/x/sys/windows/svc"
)

var svcName = "mysampleservice"

func main() {
	flag.StringVar(&svcName, "name", svcName, "name of the service")
	flag.Parse()

	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in service: %v", err)
	}
	if inService {
		RunService(svcName, false)
		return
	}

	log.Fatal("Can only be run as a Windows service")
}
