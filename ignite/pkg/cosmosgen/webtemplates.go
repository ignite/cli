package cosmosgen

import (
	"github.com/ignite/cli/v29/ignite/pkg/localfs"
	webtemplates "github.com/ignite/web"
)

// React scaffolds a React app for a chain.
func React(path string) error {
	return localfs.Save(webtemplates.ReactBoilerplate(), path)
}

// Vue scaffolds a Vue.js app for a chain.
func Vue(path string) error {
	return localfs.Save(webtemplates.VueBoilerplate(), path)
}
