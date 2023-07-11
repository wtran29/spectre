package handlers

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type job struct {
	HostServiceID int
}

// Run will perform scheduled job
func (j job) Run() {
	Repo.ScheduledCheck(j.HostServiceID)
}

// StartMonitoring starts the monitoring process
func (repo *DBRepo) StartMonitoring() {
	if app.PreferenceMap["monitoring_live"] == "1" {

		data := make(map[string]string)
		data["message"] = "Monitoring is starting..."
		err := app.WsClient.Trigger("public-channel", "app-starting", data)
		if err != nil {
			log.Println(err)
		}

		// TODO trigger a message to broadcast all clients that app is starting to monitor
		// get all services that we want to monitor
		servicesToMonitor, err := repo.DB.GetServicesToMonitor()
		if err != nil {
			log.Println(err)
		}

		// range through the services
		for _, x := range servicesToMonitor {
			log.Println("*** Service to monitor on", x.HostName, "is", x.Service.ServiceName)
			// get the schedule unit and number
			var sch string
			if x.ScheduleUnit == "d" {
				sch = fmt.Sprintf("@every %d%s", x.ScheduleNumber*24, "h")
			} else {
				sch = fmt.Sprintf("@every %d%s", x.ScheduleNumber, x.ScheduleUnit)
			}

			// create job
			var j job
			j.HostServiceID = x.ID
			scheduleID, err := app.Scheduler.AddJob(sch, j)
			if err != nil {
				log.Println(err)
			}

			// save id of job to start/stop schedule
			app.MonitorMap[x.ID] = scheduleID

			// broadcast over web sockets that service is scheduled
			payload := make(map[string]string)
			payload["message"] = "scheduling"
			payload["host_service_id"] = strconv.Itoa(x.ID)
			yearOne := time.Date(0001, 11, 17, 20, 34, 58, 65138737, time.UTC)
			if app.Scheduler.Entry(app.MonitorMap[x.ID]).Next.After(yearOne) {
				data["next_run"] = app.Scheduler.Entry(app.MonitorMap[x.ID]).Next.Format("01-02-2006 3:04:05 PM")
			} else {
				data["next_run"] = "Pending..."
			}

			payload["host"] = x.HostName
			payload["service"] = x.Service.ServiceName
			if x.LastCheck.After(yearOne) {
				payload["last_run"] = x.LastCheck.Format("01-02-2006 3:04:05 PM")
			} else {
				payload["last_run"] = "Pending..."
			}
			payload["schedule"] = fmt.Sprintf("@every %d%s", x.ScheduleNumber, x.ScheduleUnit)

			err = app.WsClient.Trigger("public-channel", "next-run-event", payload)
			if err != nil {
				log.Println(err)
			}

			err = app.WsClient.Trigger("public-channel", "schedule-changed-event", payload)
			if err != nil {
				log.Println(err)
			}
			// end range
		}

	}
}
