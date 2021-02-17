module github.com/tendermint/starport

go 1.15

require (
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/Pallinder/go-randomdata v1.2.0
	github.com/briandowns/spinner v1.11.1
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/cosmos/cosmos-sdk v0.41.0
	github.com/cosmos/go-bip39 v1.0.0
	github.com/cosmos/relayer v1.0.0-rc1.0.20210205103857-f4b56856caeb
	github.com/dariubs/percent v0.0.0-20200128140941-b7801cf1c7e2
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/fatih/color v1.10.0
	github.com/gertd/go-pluralize v0.1.7
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/go-git/go-git/v5 v5.1.0
	github.com/gobuffalo/genny v0.6.0
	github.com/gobuffalo/packr/v2 v2.8.1
	github.com/gobuffalo/plush v3.8.3+incompatible
	github.com/gobuffalo/plushgen v0.1.2
	github.com/goccy/go-yaml v1.8.0
	github.com/gookit/color v1.3.6
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/imdario/mergo v0.3.11
	github.com/improbable-eng/grpc-web v0.13.0
	github.com/jpillora/chisel v1.7.3
	github.com/magefile/mage v1.11.0 // indirect
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-zglob v0.0.3
	github.com/mwitkow/grpc-proxy v0.0.0-20181017164139-0f1106ef9c76
	github.com/olekukonko/tablewriter v0.0.4
	github.com/otiai10/copy v1.4.2
	github.com/pelletier/go-toml v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/radovskyb/watcher v1.0.7
	github.com/rakyll/statik v0.1.7
	github.com/rdegges/go-ipify v0.0.0-20150526035502-2d94a6a86c40
	github.com/rogpeppe/go-internal v1.7.0 // indirect
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.7.1 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/spn v0.0.0-20201215081711-b9ec9286ed83
	github.com/tendermint/tendermint v0.34.3
	golang.org/x/mod v0.4.1
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	golang.org/x/sys v0.0.0-20210217090653-ed5674b6da4a // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	google.golang.org/grpc v1.35.0
)

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
