# Go Windows Service

## Building

```Powershell
go build -o golang-winservice.exe .\cmd\service
```

Register the service (as administrator):

```Powershell
sc.exe create golang-winservice binPath="$((Get-Location).Path)\golang-winservice.exe"
```

And use the Windows Service Control Manager to start the service.
