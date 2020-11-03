module github.com/tendermint/starport

go 1.15

require (
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/Pallinder/go-randomdata v1.2.0
	github.com/briandowns/spinner v1.11.1
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/cosmos/cosmos-sdk v0.40.0-rc0
	github.com/cosmos/go-bip39 v0.0.0-20200817134856-d632e0d11689
	github.com/fatih/color v1.9.0
	github.com/go-git/go-git/v5 v5.1.0
	github.com/gobuffalo/envy v1.9.0 // indirect
	github.com/gobuffalo/genny v0.6.0
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/gobuffalo/plush v3.8.3+incompatible
	github.com/gobuffalo/plushgen v0.1.2
	github.com/goccy/go-yaml v1.8.0
	github.com/google/uuid v1.1.2
	github.com/gookit/color v1.2.7
	github.com/gorilla/mux v1.8.0
	github.com/imdario/mergo v0.3.11
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/manifoldco/promptui v0.8.0
	github.com/pelletier/go-toml v1.8.0
	github.com/pkg/errors v0.9.1
	github.com/radovskyb/watcher v1.0.7
	github.com/rakyll/statik v0.1.7
	github.com/regen-network/cosmos-proto v0.3.0
	github.com/rogpeppe/go-internal v1.6.2 // indirect
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	golang.org/x/mod v0.3.0
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	golang.org/x/sys v0.0.0-20201101102859-da207088b7d1 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
