// Pusher is a real-time messaging and communication service that allows developers to add real-time functionality
// to their applications. It provides features such as real-time event broadcasting, presence channels, and private channels.
package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/pusher/pusher-http-go"
)

// PusherAuth is used to authenicate to Pusher server
func (repo *DBRepo) PusherAuth(w http.ResponseWriter, r *http.Request) {
	userID := repo.App.Session.GetInt(r.Context(), "userID")

	u, _ := repo.DB.GetUserById(userID)

	params, _ := ioutil.ReadAll(r.Body)

	presenceData := pusher.MemberData{
		UserID: strconv.Itoa(userID),
		UserInfo: map[string]string{
			"name": u.FirstName,
			"id":   strconv.Itoa(userID),
		},
	}

	resp, err := app.WsClient.AuthenticatePresenceChannel(params, presenceData)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

func (repo *DBRepo) TestPusher(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["message"] = "Hello world"

	err := repo.App.WsClient.Trigger("public-channel", "test-event", data)
	if err != nil {
		log.Println(err)
	}

}