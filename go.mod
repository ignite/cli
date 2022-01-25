module github.com/tendermint/starport

go 1.16

require (
	github.com/99designs/keyring v1.1.6
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/blang/semver v3.5.1+incompatible
	github.com/briandowns/spinner v1.11.1
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/charmbracelet/glow v1.4.0
	github.com/containerd/containerd v1.5.8 // indirect
	github.com/cosmos/cosmos-sdk v0.44.5
	github.com/cosmos/go-bip39 v1.0.0
	github.com/cosmos/ibc-go/v2 v2.0.2
	github.com/docker/docker v20.10.7+incompatible
	github.com/emicklei/proto v1.9.0
	github.com/fatih/color v1.12.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-git/go-git/v5 v5.1.0
	github.com/gobuffalo/envy v1.9.0 // indirect
	github.com/gobuffalo/genny v0.6.0
	github.com/gobuffalo/logger v1.0.3
	github.com/gobuffalo/packd v0.3.0
	github.com/gobuffalo/plush v3.8.3+incompatible
	github.com/gobuffalo/plushgen v0.1.2
	github.com/goccy/go-yaml v1.9.4
	github.com/gogo/protobuf v1.3.3
	github.com/google/go-github/v37 v37.0.0
	github.com/gookit/color v1.5.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/rpc v1.2.0
	github.com/iancoleman/strcase v0.2.0
	github.com/imdario/mergo v0.3.12
	github.com/jpillora/chisel v1.7.6
	github.com/manifoldco/promptui v0.9.0
	github.com/mattn/go-zglob v0.0.3
	github.com/moby/sys/mount v0.3.0 // indirect
	github.com/otiai10/copy v1.6.0
	github.com/pelletier/go-toml v1.9.4
	github.com/pkg/errors v0.9.1
	github.com/radovskyb/watcher v1.0.7
	github.com/rdegges/go-ipify v0.0.0-20150526035502-2d94a6a86c40
	github.com/rs/cors v1.7.0
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/flutter/v2 v2.0.2
	github.com/tendermint/spn v0.1.1-0.20220125123630-1268765f1209
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.4
	github.com/tendermint/vue v0.1.58
	golang.org/x/mod v0.4.2
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf
	google.golang.org/grpc v1.43.0
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
