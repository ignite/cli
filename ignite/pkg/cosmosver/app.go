package cosmosver

type AppVersion string

const (
	AppV1 AppVersion = "app_v1"
	AppV2 AppVersion = "app_v2"
)

// DefaultVersion returns the default app version for build.
// TODO change the default value after update the version.
func DefaultVersion() AppVersion {
	return AppV1
}

// String returns the string value.
func (v AppVersion) String() string {
	return string(v)
}
