package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wtran29/spectre/internal/models"
)

const (
	HTTP           = 1
	HTTPS          = 2
	SSLCertificate = 3
)

// jsonResp is the JSON response that is sent back to client
type jsonResp struct {
	OK            bool      `json:"ok"`
	Message       string    `json:"message"`
	ServiceID     int       `json:"service_id"`
	HostServiceID int       `json:"host_service_id"`
	HostID        int       `json:"host_id"`
	OldStatus     string    `json:"old_status"`
	NewStatus     string    `json:"new_status"`
	LastCheck     time.Time `json:"last_check"`
}

// ScheduledCheck performs scheduled check on a host service by id
func (repo *DBRepo) ScheduledCheck(hostServiceID int) {
	log.Println("********* Running check for", hostServiceID)

	hs, err := repo.DB.GetHostServiceByID(hostServiceID)
	if err != nil {
		log.Println(err)
		return
	}

	h, err := repo.DB.GetHostByID(hs.HostID)
	if err != nil {
		log.Println(err)
		return
	}
	// tests the service
	newStatus, msg := repo.testServiceForHost(h, hs)

	if newStatus != hs.Status {
		repo.updateHostServiceStatusCount(h, hs, newStatus, msg)
	}

}

func (repo *DBRepo) updateHostServiceStatusCount(h models.Host, hs models.HostService, newStatus, msg string) {

	// update host service record in db with status (if changed) and last check
	hs.Status = newStatus
	hs.LastCheck = time.Now()
	err := repo.DB.UpdateHostService(hs)
	if err != nil {
		log.Println(err)
		return
	}

	pending, healthy, warning, problem, err := repo.DB.GetAllServiceStatusCounts()
	if err != nil {
		log.Println(err)
		return
	}
	data := make(map[string]string)
	data["healthy_count"] = strconv.Itoa(healthy)
	data["pending_count"] = strconv.Itoa(pending)
	data["problem_count"] = strconv.Itoa(problem)
	data["warning_count"] = strconv.Itoa(warning)
	log.Println(data)
	repo.broadcastMessage("public_channel", "host-service-count-changed", data)

	log.Println("New status is", newStatus, "and msg is", msg)
}

func (repo *DBRepo) broadcastMessage(channel, messageType string, data map[string]string) {

	err := app.WsClient.Trigger(channel, messageType, data)
	if err != nil {
		log.Println(err)
	}
}

func (repo *DBRepo) TestCheck(w http.ResponseWriter, r *http.Request) {
	hostServiceID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	oldStatus := chi.URLParam(r, "oldStatus")
	ok := true

	// get host service
	hs, err := repo.DB.GetHostServiceByID(hostServiceID)
	if err != nil {
		log.Println(err)
		ok = false
		return
	}

	// get host
	h, err := repo.DB.GetHostByID(hs.HostID)
	if err != nil {
		log.Println(err)
		ok = false
		return
	}

	// test service
	newStatus, msg := repo.testServiceForHost(h, hs)

	// update the host service in the db with status if there is a change and last check
	hs.Status = newStatus
	hs.LastCheck = time.Now()
	hs.UpdatedAt = time.Now()

	err = repo.DB.UpdateHostService(hs)
	if err != nil {
		log.Println(err)
		ok = false
	}
	// broadcast service status changed event

	var resp jsonResp
	// create json
	if ok {
		resp = jsonResp{
			OK:            true,
			Message:       msg,
			ServiceID:     hs.ServiceID,
			HostServiceID: hs.ID,
			HostID:        hs.HostID,
			OldStatus:     oldStatus,
			NewStatus:     newStatus,
			LastCheck:     time.Now(),
		}
	} else {
		resp.OK = false
		resp.Message = "Something went wrong"
	}

	// send json to client

	out, _ := json.MarshalIndent(resp, "", "	")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (repo *DBRepo) testServiceForHost(h models.Host, hs models.HostService) (string, string) {
	var msg, newStatus string

	switch hs.ServiceID {
	case HTTP:
		msg, newStatus = testHTTPForHost(h.URL)
		break
	}
	// TODO - broadcast to clients if appropriate
	if hs.Status != newStatus {
		data := make(map[string]string)
		data["host_id"] = strconv.Itoa(hs.HostID)
		data["host_service_id"] = strconv.Itoa(hs.ID)
		data["host_name"] = h.HostName
		data["service_name"] = hs.Service.ServiceName
		data["icon"] = hs.Service.Icon
		data["status"] = newStatus
		data["message"] = fmt.Sprintf("%s on %s reports %s", hs.Service.ServiceName, h.HostName, newStatus)
		data["last_check"] = time.Now().Format("01-02-2006, 3:04:06 PM")

		repo.broadcastMessage("public-channel", "host-service-status-changed", data)
	}
	// TODO - send email/sms if appropriate

	return newStatus, msg
}

func testHTTPForHost(url string) (string, string) {
	if strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}
	url = strings.Replace(url, "https://", "http://", -1)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("%s - %s", url, "error connecting"), "problem"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("%s - %s", url, resp.Status), "problem"
	}

	return fmt.Sprintf("%s - %s", url, resp.Status), "healthy"
}
