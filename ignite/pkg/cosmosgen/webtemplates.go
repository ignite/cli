package cosmosgen

import (
	webtemplates "github.com/ignite/web"

	"github.com/ignite/cli/ignite/pkg/localfs"
)

// React scaffolds a React app for a chain.
func React(path string) error {
	return localfs.Save(webtemplates.ReactBoilerplate(), path)
}

// Vue scaffolds a Vue.js app for a chain.
func Vue(path string) error {
	return localfs.Save(webtemplates.VueBoilerplate(), path)
}
