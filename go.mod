module github.com/erda-project/erda-actions

go 1.14

replace (
	github.com/influxdata/influxql => github.com/erda-project/influxql v1.1.0-ex

	k8s.io/api => k8s.io/api v0.18.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.3
	k8s.io/apiserver => k8s.io/apiserver v0.18.3
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.3
	k8s.io/client-go => k8s.io/client-go v0.18.3
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.18.3
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.18.3
	k8s.io/code-generator => k8s.io/code-generator v0.18.3
	k8s.io/component-base => k8s.io/component-base v0.18.3
	k8s.io/component-helpers => k8s.io/component-helpers v0.18.3
	k8s.io/controller-manager => k8s.io/controller-manager v0.18.3
	k8s.io/cri-api => k8s.io/cri-api v0.18.3
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.18.3
	k8s.io/klog => k8s.io/klog v1.0.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.3
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.18.3
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.18.3
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.18.3
	k8s.io/kubectl => k8s.io/kubectl v0.18.3
	k8s.io/kubelet => k8s.io/kubelet v0.18.3
	k8s.io/kubernetes => k8s.io/kubernetes v1.18.3
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.18.3
	k8s.io/metrics => k8s.io/metrics v0.18.3
	k8s.io/mount-utils => k8s.io/mount-utils v0.18.3
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.18.3
)

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.426
	github.com/andrianbdn/iospng v0.0.0-20180730113000-dccef1992541
	github.com/bitly/go-simplejson v0.5.1-0.20181114203107-9db4a59bd4d8
	github.com/caarlos0/env v3.3.1-0.20180521112546-3e0f30cbf50b+incompatible
	github.com/cespare/trie v0.0.0-20150610204604-3fe1a95cbba9 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/erda-project/erda v0.0.0-20210518111439-5e375de156b8
	github.com/hashicorp/go-multierror v1.0.0
	github.com/labstack/gommon v0.3.0
	github.com/machinebox/progress v0.2.0
	github.com/matryer/is v1.4.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/otiai10/copy v1.5.0
	github.com/parnurzeal/gorequest v0.2.16
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.12.1-0.20201118115123-7230c61342c8
	github.com/robertkrimen/terst v0.0.0-20140908162406-4b1c60b7cc23
	github.com/sabhiram/go-gitignore v0.0.0-20180611051255-d3107576ba94
	github.com/shogo82148/androidbinary v1.0.2
	github.com/sirupsen/logrus v1.8.0
	github.com/stretchr/testify v1.7.0
	github.com/thedevsaddam/gojsonq/v2 v2.5.2
	github.com/toqueteos/trie v0.0.0-20150530104557-56fed4a05683 // indirect
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/net v0.0.0-20210421230115-4e50805a0758
	google.golang.org/grpc/examples v0.0.0-20210518222651-23a83dd097ec // indirect
	gopkg.in/src-d/enry.v1 v1.6.4
	gopkg.in/stretchr/testify.v1 v1.2.2
	gopkg.in/toqueteos/substring.v1 v1.0.2 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gotest.tools v2.2.0+incompatible
	howett.net/plist v0.0.0-20201203080718-1454fab16a06
)
