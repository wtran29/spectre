package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
	"github.com/robfig/cron/v3"
	"github.com/wtran29/spectre/internal/channeldata"
	"github.com/wtran29/spectre/internal/driver"
	"github.com/wtran29/spectre/internal/models"
)

// AppConfig holds application configuration
type AppConfig struct {
	DB            *driver.DB
	Session       *scs.SessionManager
	InProduction  bool
	Domain        string
	MonitorMap    map[int]cron.EntryID
	PreferenceMap map[string]string
	Scheduler     *cron.Cron
	// WsClient      pusher.Client
	WsClient      models.WSClient
	PusherSecret  string
	TemplateCache map[string]*template.Template
	MailQueue     chan channeldata.MailJob
	Version       string
	Identifier    string
}
