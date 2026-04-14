package modulemigration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateModuleAddsInitialMigration(t *testing.T) {
	opts := &Options{
		ModuleName:  "blog",
		ModulePath:  "github.com/test/blog",
		FromVersion: 1,
		ToVersion:   2,
	}

	got, err := updateModule(moduleWithServiceRegistrar(`
func (am AppModule) RegisterServices(registrar grpc.ServiceRegistrar) error {
	types.RegisterMsgServer(registrar, keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(registrar, keeper.NewQueryServerImpl(am.keeper))

	return nil
}

func (AppModule) ConsensusVersion() uint64 { return 1 }
`), opts)
	require.NoError(t, err)

	normalized := normalize(got)
	require.Contains(t, normalized, `migrationv2"github.com/test/blog/x/blog/migrations/v2"`)
	require.Contains(t, normalized, `cfg,ok:=registrar.(module.Configurator)`)
	require.Contains(t, normalized, `if!ok{returnnil}`)
	require.Contains(t, normalized, `cfg.RegisterMigration(types.ModuleName,1,migrationv2.Migrate)`)
	require.Contains(t, normalized, `func(AppModule)ConsensusVersion()uint64{return2}`)
}

func TestUpdateModuleAppendsMigrationToExistingConfiguratorBlock(t *testing.T) {
	opts := &Options{
		ModuleName:  "blog",
		ModulePath:  "github.com/test/blog",
		FromVersion: 2,
		ToVersion:   3,
	}

	got, err := updateModule(moduleWithServiceRegistrar(`
func (am AppModule) RegisterServices(registrar grpc.ServiceRegistrar) error {
	types.RegisterMsgServer(registrar, keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(registrar, keeper.NewQueryServerImpl(am.keeper))

	cfg, ok := registrar.(module.Configurator)
	if !ok {
		return nil
	}

	if err := cfg.RegisterMigration(types.ModuleName, 1, migrationv2.Migrate); err != nil {
		return err
	}

	return nil
}

func (AppModule) ConsensusVersion() uint64 { return 2 }
`, `migrationv2 "github.com/test/blog/x/blog/migrations/v2"`), opts)
	require.NoError(t, err)

	normalized := normalize(got)
	require.Equal(t, 1, strings.Count(normalized, `registrar.(module.Configurator)`))
	require.Contains(t, normalized, `migrationv3"github.com/test/blog/x/blog/migrations/v3"`)
	require.Contains(t, normalized, `cfg.RegisterMigration(types.ModuleName,1,migrationv2.Migrate)`)
	require.Contains(t, normalized, `cfg.RegisterMigration(types.ModuleName,2,migrationv3.Migrate)`)
	require.Contains(t, normalized, `func(AppModule)ConsensusVersion()uint64{return3}`)
}

func TestUpdateModuleSupportsConfiguratorSignatureAndConstantVersion(t *testing.T) {
	opts := &Options{
		ModuleName:  "blog",
		ModulePath:  "github.com/test/blog",
		FromVersion: 2,
		ToVersion:   3,
	}

	got, err := updateModule(moduleWithConfigurator(`
const ConsensusVersion = 2

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}

func (AppModule) ConsensusVersion() uint64 { return ConsensusVersion }
`), opts)
	require.NoError(t, err)

	normalized := normalize(got)
	require.NotContains(t, normalized, `registrar.(module.Configurator)`)
	require.Contains(t, normalized, `cfg.RegisterMigration(types.ModuleName,2,migrationv3.Migrate)`)
	require.Contains(t, normalized, `iferr:=cfg.RegisterMigration(types.ModuleName,2,migrationv3.Migrate);err!=nil{panic(err)}`)
	require.Contains(t, normalized, `constConsensusVersion=3`)

	version, err := ConsensusVersion(got)
	require.NoError(t, err)
	require.EqualValues(t, 3, version)
}

func moduleWithServiceRegistrar(body string, extraImports ...string) string {
	imports := []string{
		`"github.com/cosmos/cosmos-sdk/types/module"`,
		`"google.golang.org/grpc"`,
		`"github.com/test/blog/x/blog/keeper"`,
		`"github.com/test/blog/x/blog/types"`,
	}
	imports = append(imports, extraImports...)

	return "package blog\n\nimport (\n\t" + strings.Join(imports, "\n\t") + "\n)\n\n" + body
}

func moduleWithConfigurator(body string, extraImports ...string) string {
	imports := []string{
		`"github.com/cosmos/cosmos-sdk/types/module"`,
		`"github.com/test/blog/x/blog/keeper"`,
		`"github.com/test/blog/x/blog/types"`,
	}
	imports = append(imports, extraImports...)

	return "package blog\n\nimport (\n\t" + strings.Join(imports, "\n\t") + "\n)\n\n" + body
}

func normalize(content string) string {
	return strings.Join(strings.Fields(content), "")
}
