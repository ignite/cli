package cosmosgen

import (
	webtemplates "github.com/ignite/web"

	"github.com/ignite/cli/v29/ignite/pkg/localfs"
)

// Vue scaffolds a Vue.js app for a chain.
func Vue(path string) error {
	return localfs.Save(webtemplates.VueBoilerplate(), path)
}
