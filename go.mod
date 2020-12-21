module github.com/tendermint/starport

go 1.15

require (
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/Pallinder/go-randomdata v1.2.0
	github.com/briandowns/spinner v1.11.1
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/cosmos/cosmos-sdk v0.40.0-rc3
	github.com/cosmos/go-bip39 v0.0.0-20200817134856-d632e0d11689
	github.com/dariubs/percent v0.0.0-20200128140941-b7801cf1c7e2
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/fatih/color v1.9.0
	github.com/gertd/go-pluralize v0.1.7
	github.com/go-git/go-git/v5 v5.1.0
	github.com/gobuffalo/genny v0.6.0
	github.com/gobuffalo/packr/v2 v2.8.1
	github.com/gobuffalo/plush v3.8.3+incompatible
	github.com/gobuffalo/plushgen v0.1.2
	github.com/goccy/go-yaml v1.8.0
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/google/uuid v1.1.2
	github.com/gookit/color v1.2.7
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/grpc-gateway v1.15.2
	github.com/imdario/mergo v0.3.11
	github.com/improbable-eng/grpc-web v0.13.0
	github.com/jpillora/chisel v1.7.3
	github.com/manifoldco/promptui v0.8.0
	github.com/mwitkow/grpc-proxy v0.0.0-20181017164139-0f1106ef9c76
	github.com/olekukonko/tablewriter v0.0.4
	github.com/otiai10/copy v1.2.0
	github.com/pelletier/go-toml v1.8.0
	github.com/pkg/errors v0.9.1
	github.com/radovskyb/watcher v1.0.7
	github.com/rakyll/statik v0.1.7
	github.com/rdegges/go-ipify v0.0.0-20150526035502-2d94a6a86c40
	github.com/regen-network/cosmos-proto v0.3.0
	github.com/rs/cors v1.7.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/spn v0.0.0-20201215081711-b9ec9286ed83
	github.com/tendermint/tendermint v0.34.0-rc6
	golang.org/x/crypto v0.0.0-20201217014255-9d1352758620 // indirect
	golang.org/x/mod v0.4.0
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/sys v0.0.0-20201221093633-bc327ba9c2f0 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	google.golang.org/genproto v0.0.0-20201019141844-1ed22bb0c154 // indirect
	google.golang.org/grpc v1.33.0
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
