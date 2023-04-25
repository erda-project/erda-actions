module github.com/erda-project/erda-actions

go 1.17

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5
	github.com/docker/spdystream => github.com/docker/spdystream v0.0.0-20160310174837-449fdfce4d96
	github.com/erda-project/erda-proto-go => ./.erda/api/proto-go
	github.com/google/gnostic => github.com/googleapis/gnostic v0.4.0
	github.com/googlecloudplatform/flink-operator => github.com/erda-project/flink-on-k8s-operator v0.0.0-20220211023605-1a64c463171a
	github.com/influxdata/influxql => github.com/erda-project/influxql v1.1.0-ex
	github.com/olivere/elastic v6.2.35+incompatible => github.com/erda-project/elastic v0.0.1-ex
	github.com/pingcap/pd/v4 => github.com/tikv/pd v1.0.8
	github.com/rancher/remotedialer => github.com/erda-project/remotedialer v0.2.6-0.20210618084817-52c879aadbcb
	go.etcd.io/bbolt => github.com/coreos/bbolt v1.3.5
	go.etcd.io/bbolt v1.3.5 => github.com/coreos/bbolt v1.3.5
	go.opentelemetry.io/otel => go.opentelemetry.io/otel v1.2.0
	go.opentelemetry.io/otel/exporters/stdout => go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.2.0
	go.opentelemetry.io/proto/otlp v0.11.0 => github.com/recallsong/opentelemetry-proto-go/otlp v0.11.1-0.20211202093058-995eca7123f5
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	k8s.io/api => k8s.io/api v0.21.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.2
	k8s.io/apiserver => k8s.io/apiserver v0.21.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.21.2
	k8s.io/client-go => k8s.io/client-go v0.21.2
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.21.2
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.21.2
	k8s.io/code-generator => k8s.io/code-generator v0.21.2
	k8s.io/component-base => k8s.io/component-base v0.21.2
	k8s.io/component-helpers => k8s.io/component-helpers v0.21.2
	k8s.io/controller-manager => k8s.io/controller-manager v0.21.2
	k8s.io/cri-api => k8s.io/cri-api v0.21.2
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.21.2
	k8s.io/klog => k8s.io/klog v1.0.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.21.2
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.21.2
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.21.2
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.21.2
	k8s.io/kubectl => k8s.io/kubectl v0.21.2
	k8s.io/kubelet => k8s.io/kubelet v0.21.2
	k8s.io/kubernetes => k8s.io/kubernetes v1.21.2
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.21.2
	k8s.io/metrics => k8s.io/metrics v0.21.2
	k8s.io/mount-utils => k8s.io/mount-utils v0.21.2
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.21.2
)

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.1600
	github.com/aliyun/aliyun-oss-go-sdk v2.1.4+incompatible
	github.com/andrianbdn/iospng v0.0.0-20180730113000-dccef1992541
	github.com/bitly/go-simplejson v0.5.1-0.20181114203107-9db4a59bd4d8
	github.com/caarlos0/env v3.3.1-0.20180521112546-3e0f30cbf50b+incompatible
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/docker/cli v20.10.20+incompatible
	github.com/docker/distribution v2.8.1+incompatible
	github.com/erda-project/erda v0.0.0-20221214015351-c24776421f50
	github.com/erda-project/erda-proto-go v0.0.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/labstack/gommon v0.3.0
	github.com/machinebox/progress v0.2.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/otiai10/copy v1.5.0
	github.com/parnurzeal/gorequest v0.2.16
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.13.1
	github.com/robertkrimen/terst v0.0.0-20140908162406-4b1c60b7cc23
	github.com/sabhiram/go-gitignore v0.0.0-20201211210132-54b8a0bf510f
	github.com/shogo82148/androidbinary v1.0.2
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.8.1
	github.com/thedevsaddam/gojsonq/v2 v2.5.2
	golang.org/x/crypto v0.1.0
	golang.org/x/net v0.7.0
	gopkg.in/src-d/enry.v1 v1.6.4
	gopkg.in/stretchr/testify.v1 v1.2.2
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
	gotest.tools v2.2.0+incompatible
	howett.net/plist v1.0.0
	k8s.io/api v0.25.3
	sigs.k8s.io/yaml v1.3.0
)

require (
	cloud.google.com/go/compute/metadata v0.2.2 // indirect
	cloud.google.com/go/iam v0.8.0 // indirect
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/XSAM/otelsql v0.9.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/brahma-adshonor/gohook v1.1.9 // indirect
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/cespare/trie v0.0.0-20150610204604-3fe1a95cbba9 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548 // indirect
	github.com/danjacques/gofslock v0.0.0-20191023191349-0a45f885bc37 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/docker/docker v20.10.24+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.7.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.1.0 // indirect
	github.com/erda-project/erda-infra v1.0.8 // indirect
	github.com/erda-project/erda-infra/tools v0.0.0-20220922112628-1965c0260662 // indirect
	github.com/extrame/ole2 v0.0.0-20160812065207-d69429661ad7 // indirect
	github.com/extrame/xls v0.0.1 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/frankban/quicktest v1.14.3 // indirect
	github.com/getkin/kin-openapi v0.76.0 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/go-eden/routine v1.0.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gogap/errors v0.0.0-20200228125012-531a6449b28c // indirect
	github.com/gogap/stack v0.0.0-20150131034635-fef68dddd4f8 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.1 // indirect
	github.com/googleapis/go-type-adapters v1.0.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/labstack/echo v3.3.10+incompatible // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/matryer/is v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/minio/minio-go v0.0.0-20190308013636-b32976861da0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pingcap/errors v0.11.5-0.20200917111840-a15ef68f753d // indirect
	github.com/pingcap/kvproto v0.0.0-20200818080353-7aaed8998596 // indirect
	github.com/pingcap/log v0.0.0-20200828042413-fce0951f1463 // indirect
	github.com/pingcap/parser v0.0.0-20201022083903-fbe80b0c40bb // indirect
	github.com/pingcap/tidb v1.1.0-beta.0.20200921100526-29e8c0913100 // indirect
	github.com/pingcap/tipb v0.0.0-20200618092958-4fad48b4c8c3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.13.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/rancher/remotedialer v0.2.6-0.20220104192242-f3837f8d649a // indirect
	github.com/rancher/wrangler v0.8.11 // indirect
	github.com/recallsong/go-utils v1.1.2-0.20210826100715-fce05eefa294 // indirect
	github.com/recallsong/unmarshal v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20190728182440-6a916e37a237 // indirect
	github.com/rogpeppe/fastuuid v1.2.0 // indirect
	github.com/rogpeppe/go-internal v1.6.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/shabbyrobe/xmlwriter v0.0.0-20200208144257-9fca06d00ffa // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/t-tiger/gorm-bulk-insert v1.3.0 // indirect
	github.com/tealeg/xlsx v1.0.5 // indirect
	github.com/tealeg/xlsx/v3 v3.2.4-0.20210615062226-d5ce25722f69 // indirect
	github.com/tikv/pd v1.1.0-beta.0.20200907080620-6830f5bb92a2 // indirect
	github.com/toqueteos/trie v0.0.0-20150530104557-56fed4a05683 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200910180754-dd1b699fc489 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.27.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.27.0 // indirect
	go.opentelemetry.io/otel v1.10.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.2.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.2.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.2.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.25.0 // indirect
	go.opentelemetry.io/otel/metric v0.25.0 // indirect
	go.opentelemetry.io/otel/sdk v1.2.0 // indirect
	go.opentelemetry.io/otel/trace v1.10.0 // indirect
	go.opentelemetry.io/proto/otlp v0.11.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.23.0 // indirect
	golang.org/x/arch v0.0.0-20190312162104-788fe5ffcd8c // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/time v0.0.0-20220609170525-579cf78fd858 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.103.0 // indirect
	google.golang.org/genproto v0.0.0-20221201164419-0e50fba7f41c // indirect
	google.golang.org/grpc v1.50.1 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/oauth2.v3 v3.12.0 // indirect
	gopkg.in/toqueteos/substring.v1 v1.0.2 // indirect
	gorm.io/driver/mysql v1.3.2 // indirect
	gorm.io/gorm v1.23.5 // indirect
	k8s.io/apimachinery v0.25.3 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	moul.io/http2curl v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)
