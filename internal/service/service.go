package service

import (
	"fmt"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"

	"golang-winservice/internal/core"
)

const (
	ServiceName  = "golang-winservice"
	MQTTBroker   = "tcp://localhost:1883" // Default local MQTT broker
	MQTTClientID = "windows-service-client"
	MQTTTopic    = "service/status"
)

var eventLog debug.Log

type WindowsService struct {
	mqttClient *core.MQTTClient
}

func (m *WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}

	// Initialize the MQTT client
	m.mqttClient = core.NewMQTTClient(MQTTClientID, MQTTBroker, MQTTTopic)
	err := m.mqttClient.Connect()
	if err != nil {
		eventLog.Error(1, fmt.Sprintf("Failed to connect to MQTT broker: %v", err))
	} else {
		eventLog.Info(1, "Successfully connected to MQTT broker")
	}

	// Initialize the TickManager
	tickManager := core.NewTickManager()

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case <-tickManager.CurrentTick:
			eventLog.Info(1, "Tick processed successfully")
			if m.mqttClient.IsConnected() {
				err := m.mqttClient.SendStatusUpdate()
				if err != nil {
					eventLog.Error(1, fmt.Sprintf("Failed to send MQTT message: %v", err))
				} else {
					eventLog.Info(1, "MQTT status message sent successfully")
				}
			} else {
				eventLog.Warning(1, "MQTT client not connected, attempting reconnect")
				// Try to reconnect
				err := m.mqttClient.Connect()
				if err != nil {
					eventLog.Error(1, fmt.Sprintf("MQTT reconnect failed: %v", err))
				}
			}
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				tickManager.UseSlowTick()
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				tickManager.UseFastTick()
			default:
				eventLog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func RunService(isDebug bool) {
	var err error
	if isDebug {
		eventLog = debug.New(ServiceName)
	} else {
		eventLog, err = eventlog.Open(ServiceName)
		if err != nil {
			return
		}
	}
	defer eventLog.Close()

	eventLog.Info(1, fmt.Sprintf("starting %s service", ServiceName))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err = run(ServiceName, &WindowsService{})
	if err != nil {
		eventLog.Error(1, fmt.Sprintf("%s service failed: %v", ServiceName, err))
		return
	}
	eventLog.Info(1, fmt.Sprintf("%s service stopped", ServiceName))
}
