package main

import (
	"log"

	"golang.org/x/sys/windows/svc"

	"github.com/vbossica/golang-winservice/internal/service"
)

func main() {
	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in service: %v", err)
	}
	if inService {
		service.RunService(false)
		return
	}

	log.Fatal("Can only be run as a Windows service")
}
