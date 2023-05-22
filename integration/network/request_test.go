package network_test

import (
	"testing"
)

func TestNetworkRequestParam(t *testing.T) {
	t.Skip("Skipped until SPN app is migrated to V2")

	// var (
	// 	env     = envtest.New(t)
	// 	spnPath = setupSPN(env)
	// 	spn     = env.App(
	// 		spnPath,
	// 		envtest.AppHomePath(t.TempDir()),
	// 		envtest.AppConfigPath(path.Join(spnPath, spnConfigFile)),
	// 	)
	// )
	//
	// var (
	// 	ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
	// 	isBackendAliveErr error
	// )
	//
	// // Make sure that the SPN config file is at the latest version
	// migrateSPNConfig(t, spnPath)
	//
	// validator := spn.Config().Validators[0]
	// servers, err := validator.GetServers()
	// require.NoError(t, err)
	//
	// go func() {
	// 	defer cancel()
	//
	// 	if isBackendAliveErr = env.IsAppServed(ctx, servers.API.Address); isBackendAliveErr != nil {
	// 		return
	// 	}
	// 	var b bytes.Buffer
	// 	env.Exec("publish planet chain to spn",
	// 		step.NewSteps(step.New(
	// 			step.Exec(
	// 				envtest.IgniteApp,
	// 				"network", "chain", "publish",
	// 				"https://github.com/ignite/example",
	// 				"--local",
	// 				// The hash is used to be sure the test uses the right config
	// 				// version. Hash value must be updated to the latest when the
	// 				// config version in the repository is updated to a new version.
	// 				"--hash", "b8b2cc2876c982dd4a049ed16b9a6099eca000aa",
	// 			),
	// 			step.Stdout(&b),
	// 		),
	// 			step.New(
	// 				step.Exec(
	// 					envtest.IgniteApp,
	// 					"network", "request", "change-param",
	// 					"1", "mint", "mint_denom", "\"bar\"",
	// 					"--local",
	// 				),
	// 				step.Stdout(&b),
	// 			),
	// 			step.New(
	// 				step.Exec(
	// 					envtest.IgniteApp,
	// 					"network", "chain", "show", "genesis",
	// 					"1",
	// 					"--local",
	// 				),
	// 				step.Stdout(&b),
	// 			),
	// 		),
	// 	)
	// 	require.False(t, env.HasFailed(), b.String())
	// 	t.Log(b.String())
	// }()
	//
	// env.Must(spn.Serve("serve spn chain", envtest.ExecCtx(ctx)))
	//
	// require.NoError(t, isBackendAliveErr, "spn cannot get online in time")
}
