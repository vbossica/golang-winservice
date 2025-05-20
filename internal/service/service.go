package service

import (
	"fmt"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

const (
	ServiceName = "golang-winservice"
)

var eventLog debug.Log

type WindowsService struct{}

func (m *WindowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}

	// Initialize the TickManager
	tickManager := core.NewTickManager()

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case <-tickManager.CurrentTick:
			eventLog.Info(1, "Tick processed successfully")
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
