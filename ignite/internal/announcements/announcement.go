package announcements

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
)

var (
	SurveyLink      = "https://bit.ly/3WZS2uS"
	AnnouncementURL = "https://api.ignite.com/v1/announcements"
)

type announcement struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	User      string    `json:"user"`
}

// Fetch fetches the latest announcements from the API.
func Fetch() string {
	resp, err := http.Get(AnnouncementURL) //nolint:gosec
	if err != nil || resp.StatusCode != 200 {
		return fallbackData()
	}
	defer resp.Body.Close()

	type response struct {
		Announcements []announcement `json:"announcements"`
	}
	var data response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fallbackData()
	}

	if len(data.Announcements) == 0 {
		return fallbackData()
	}

	var out strings.Builder
	fmt.Fprintf(&out, "%s\n\n", "Announcements:")

	for _, msg := range data.Announcements {
		fmt.Fprintf(&out, "%s %s\n", icons.Bullet, msg.Text)
	}

	return out.String()
}

func fallbackData() string {
	return fmt.Sprintf("\n%s Survey: %s\n", icons.Survey, SurveyLink)
}
