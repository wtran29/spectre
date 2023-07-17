package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pusher/pusher-http-go"
	"github.com/robfig/cron/v3"
	"github.com/wtran29/spectre/internal/channeldata"
	"github.com/wtran29/spectre/internal/config"
	"github.com/wtran29/spectre/internal/driver"
	"github.com/wtran29/spectre/internal/helpers"
	"github.com/wtran29/spectre/internal/repository/dbrepo"
)

var testSession *scs.SessionManager

func TestMain(m *testing.M) {
	testSession = scs.New()
	testSession.Lifetime = 24 * time.Hour
	testSession.Cookie.Persist = true
	testSession.Cookie.SameSite = http.SameSiteLaxMode
	testSession.Cookie.Secure = false

	mailQueue := make(chan channeldata.MailJob, 5)

	a := config.AppConfig{
		DB:           &driver.DB{},
		Session:      testSession,
		InProduction: false,
		Domain:       "localhost",
		MailQueue:    mailQueue,
	}

	app = &a

	preferenceMap := make(map[string]string)
	app.PreferenceMap = preferenceMap

	// create pusher client
	dws := dummyWS{
		AppID:  "1",
		Secret: "123abc",
		Key:    "abc123",
		Secure: false,
		Host:   "localhost:4001",
	}

	app.WsClient = &dws

	monitorMap := make(map[int]cron.EntryID)
	app.MonitorMap = monitorMap

	localZone, _ := time.LoadLocation("Local")
	scheduler := cron.New(cron.WithLocation(localZone), cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))

	app.Scheduler = scheduler

	repo := NewTestHandlers(app)
	NewHandlers(repo, app)

	helpers.NewHelpers(app)

	helpers.SetViews("./../../views")

	os.Exit(m.Run())
}

// gets the context with session added
func getCtx(r *http.Request) context.Context {
	ctx, err := testSession.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

// NewTestHandlers creates a new repository
func NewTestHandlers(a *config.AppConfig) *DBRepo {
	return &DBRepo{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

type dummyWS struct {
	AppID                        string
	Key                          string
	Secret                       string
	Host                         string // host or host:port pair
	Secure                       bool   // true for HTTPS
	Cluster                      string
	HTTPClient                   *http.Client
	EncryptionMasterKey          string  // deprecated
	EncryptionMasterKeyBase64    string  // for E2E
	validatedEncryptionMasterKey *[]byte // parsed key for use
}

func (c *dummyWS) Trigger(channel string, eventName string, data interface{}) error {
	return nil
}

func (c *dummyWS) TriggerMulti(channels []string, eventName string, data interface{}) error {
	return nil
}

func (c *dummyWS) TriggerExclusive(channel string, eventName string, data interface{}, socketID string) error {
	return nil
}

func (c *dummyWS) TriggerMultiExclusive(channels []string, eventName string, data interface{}, socketID string) error {
	return nil
}

func (c *dummyWS) TriggerBatch(batch []pusher.Event) error {
	return nil
}

func (c *dummyWS) Channels(additionalQueries map[string]string) (*pusher.ChannelsList, error) {
	var cl pusher.ChannelsList
	return &cl, nil
}

func (c *dummyWS) Channel(name string, additionalQueries map[string]string) (*pusher.Channel, error) {
	var cl pusher.Channel
	return &cl, nil
}

func (c *dummyWS) GetChannelUsers(name string) (*pusher.Users, error) {
	var cl pusher.Users
	return &cl, nil
}

func (c *dummyWS) AuthenticatePrivateChannel(params []byte) (response []byte, err error) {

	return []byte("Hello"), nil
}

func (c *dummyWS) AuthenticatePresenceChannel(params []byte, member pusher.MemberData) (response []byte, err error) {
	jsonStr := `{"auth": "abc123:e632c9f9291250d20e8067da6ec263b31cf9a7a59491b26d90cb15267e226520",
    "channel_data": "{\"user_id\":\"1\",\"user_info\":{\"id\":\"1\",\"name\":\"Admin\"}}"}`

	return []byte(jsonStr), nil
}

func (c *dummyWS) Webhook(header http.Header, body []byte) (*pusher.Webhook, error) {
	var wh pusher.Webhook
	return &wh, nil
}
