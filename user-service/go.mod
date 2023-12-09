module user-service

go 1.21

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.9.1 // or the latest stable version
	github.com/lib/pq v1.10.9
	github.com/spf13/viper v1.18.1 // or the latest stable version
	github.com/swaggo/swag v1.7.1 // or the latest stable version
	golang.org/x/crypto v0.16.0 // or the latest stable version
	golang.org/x/net v0.19.0 // indirect
	gorm.io/driver/postgres v1.5.4 // or the latest stable version
	gorm.io/gorm v1.25.5 // or the latest stable version
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/bytedance/sonic v1.10.2 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.1 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/spec v0.20.11 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.16.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgx/v5 v5.5.1 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.6.0 // indirect
	golang.org/x/exp v0.0.0-20231206192017-f3f8817b8deb // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.16.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/bytedance/sonic => github.com/bytedance/sonic v1.9.1
//
//replace github.com/go-playground/validator/v10 => github.com/go-playground/validator/v10 v10.16.0
//
//replace github.com/mattn/go-isatty => github.com/mattn/go-isatty v0.0.20 // indirect
//
//replace github.com/pelletier/go-toml/v2 => github.com/pelletier/go-toml/v2 v2.1.0
//
//replace github.com/stretchr/testify => github.com/stretchr/testify v1.8.3
//
//replace github.com/ugorji/go/codec => github.com/ugorji/go/codec v1.2.12
//
//replace golang.org/x/net => golang.org/x/net v0.19.0
//
//replace google.golang.org/protobuf => google.golang.org/protobuf v1.30.0
//
//replace github.com/gabriel-vasile/mimetype => github.com/gabriel-vasile/mimetype v1.4.2
//
//replace github.com/klauspost/cpuid/v2 => github.com/klauspost/cpuid/v2 v2.2.4
//
//replace golang.org/x/arch => golang.org/x/arch v0.3.0
//
//replace golang.org/x/crypto => golang.org/x/crypto v0.9.0
//
//replace github.com/spf13/viper => github.com/spf13/viper v1.17.0
//
//replace github.com/sagikazarmark/crypt => github.com/sagikazarmark/crypt v0.15.0
//
//replace cloud.google.com/go => cloud.google.com/go v0.110.7
