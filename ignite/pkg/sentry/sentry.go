package sentry

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/version"
)

const IgniteDNS = "https://1d862300ead01c5814d8ead3732fd41f@o1152630.ingest.us.sentry.io/4507891348930560"

func InitSentry(ctx context.Context) (deferMe func(), err error) {
	sentrySyncTransport := sentry.NewHTTPSyncTransport()
	sentrySyncTransport.Timeout = time.Second * 3

	igniteInfo, err := version.GetInfo(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to init sentry: %w", err)
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         IgniteDNS,
		Transport:   sentrySyncTransport,
		Environment: getEnvironment(igniteInfo.CLIVersion),
		Release:     fmt.Sprintf("ignite@%s", igniteInfo.CLIVersion),
		SampleRate:  1.0, // get all events
	}); err != nil {
		return nil, errors.Errorf("failed to init sentry: %w", err)
	}

	return func() {
		sentry.Recover()
		sentry.Flush(time.Second * 2)
	}, nil
}

func getEnvironment(igniteVersion string) string {
	if strings.Contains(igniteVersion, "dev") {
		return "development"
	}

	return "production"
}
