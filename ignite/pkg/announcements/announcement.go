package announcements

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
)

const (
	SurveyLink      = "https://bit.ly/3WZS2uS"
	AnnouncementAPI = "http://api.ignite.com/announcements"
)

type announcement struct {
	Announcements []string `json:"announcements"`
}

func GetAnnouncements() string {
	resp, err := http.Get(AnnouncementAPI)
	if err != nil || resp.StatusCode != 200 {
		return fallbackData()
	}
	defer resp.Body.Close()

	var data announcement
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fallbackData()
	}

	// is this needed? or if its empty we don't want to show anything?
	if len(data.Announcements) == 0 {
		return fallbackData()
	}

	var out strings.Builder
	fmt.Fprintf(&out, "\n%s %s\n", icons.Announcement, "Announcements")

	for _, announcement := range data.Announcements {
		fmt.Fprintf(&out, "%s %s\n", icons.Bullet, announcement)
	}

	return out.String()
}

func fallbackData() string {
	return fmt.Sprintf("\n%s Survey: %s\n", icons.Survey, SurveyLink)
}
