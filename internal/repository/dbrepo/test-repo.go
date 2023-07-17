package dbrepo

import (
	"github.com/wtran29/spectre/internal/models"
)

// AllUsers returns all users
func (m *testDBRepo) AllUsers() ([]*models.User, error) {
	var users []*models.User

	return users, nil
}

// GetUserById returns a user by id
func (m *testDBRepo) GetUserById(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) InsertUser(u models.User) (int, error) {
	return 1, nil
}
func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}
func (m *testDBRepo) DeleteUser(id int) error {
	return nil
}
func (m *testDBRepo) UpdatePassword(id int, newPassword string) error {
	return nil
}

// Authenticate authenticates user
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {

	return 1, "", nil
}

// InsertRememberMeToken inserts a remember me token into remember_tokens for a user
func (m *testDBRepo) InsertRememberMeToken(id int, token string) error {
	return nil
}

func (m *testDBRepo) DeleteToken(token string) error {
	return nil
}

func (m *testDBRepo) CheckForToken(id int, token string) bool {
	return true
}

func (m *testDBRepo) AllPreferences() ([]models.Preference, error) {
	var pref []models.Preference
	return pref, nil
}
func (m *testDBRepo) SetSystemPref(name, value string) error {
	return nil
}
func (m *testDBRepo) InsertOrUpdateSitePreferences(pm map[string]string) error {
	return nil
}
func (m *testDBRepo) UpdateSystemPref(name, value string) error {
	return nil
}

func (m *testDBRepo) InsertHost(h models.Host) (int, error) {
	return 1, nil
}
func (m *testDBRepo) GetHostByID(id int) (models.Host, error) {
	var host models.Host
	return host, nil
}
func (m *testDBRepo) UpdateHost(h models.Host) error {
	return nil
}
func (m *testDBRepo) AllHosts() ([]models.Host, error) {
	var hosts []models.Host
	return hosts, nil
}

func (m *testDBRepo) UpdateHostServiceStatus(hostID, serviceID, active int) error {
	return nil
}
func (m *testDBRepo) GetAllServiceStatusCounts() (int, int, int, int, error) {
	return 1, 0, 0, 0, nil
}
func (m *testDBRepo) GetServicesByStatus(status string) ([]models.HostService, error) {
	var hs []models.HostService
	return hs, nil
}
func (m *testDBRepo) GetHostServiceByID(id int) (models.HostService, error) {
	var hs models.HostService
	return hs, nil
}
func (m *testDBRepo) UpdateHostService(hs models.HostService) error {
	return nil
}
func (m *testDBRepo) GetServicesToMonitor() ([]models.HostService, error) {
	var hs []models.HostService
	return hs, nil
}
func (m *testDBRepo) GetHostServiceByHostIdServiceId(hostID, serviceID int) (models.HostService, error) {
	var hs models.HostService
	return hs, nil
}
func (m *testDBRepo) GetAllEvents() ([]models.Event, error) {
	var events []models.Event
	return events, nil
}
func (m *testDBRepo) InsertEvent(e models.Event) error {
	return nil
}
