# Go Windows Service

## Overview

A simple Windows service that, every 2 or 5 seconds, sends a message to a MQTT broker.

## Building

```Powershell
go build -o golang-winservice.exe .\cmd\service
```

Register the service (as administrator):

```Powershell
sc.exe create golang-winservice binPath="$((Get-Location).Path)\golang-winservice.exe"
```

And use the Windows Service Control Manager to start the service.
