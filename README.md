# Go Windows Service

## Overview

A simple Windows service that, every 2 or 5 seconds, sends a message to a MQTT broker.

## Building and Running

1. Build the application for Windows:

    ```Powershell
    go build -o golang-winservice.exe .\cmd\service
    ```

1. Register the service (as administrator):

    ```Powershell
    sc.exe create golang-winservice \
        binPath="$((Get-Location).Path)\golang-winservice.exe"
    ```

1. Use the Windows Service Control Manager to start the service.

1. Deregister the service (as administrator):

    ```Powershell
    sc.exe delete golang-winservice
    ```
