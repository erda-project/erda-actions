SHELL=/bin/bash
GitBranch=$(shell git rev-parse --abbrev-ref HEAD)
GitCommit=$(shell git rev-parse --short HEAD)
Date=$(shell date '+%Y%m%d%H%M%S')
BuildTime=$(shell date '+%Y-%m-%d %T%z')
Registry="registry.erda.cloud/erda-actions"
RegistryForPush="registry.erda.cloud/erda-actions"
DevelopRegistry="registry.cn-hangzhou.aliyuncs.com/dice"
ARCH ?= $(shell go env GOARCH)
PLATFORMS ?= linux/arm64,linux/amd64

.ONESHELL:
echo \
custom-script java-agent \
email mysql-assert app-create app-run git-checkout assert jsonparse redis-cli mysql-cli git-push release dice dice-deploy dice-deploy-addon dice-deploy-service \
dice-deploy-domain dice-deploy-release dice-deploy-redeploy dice-deploy-rollback buildpack buildpack-aliyun java java-build js js-build manual-review js-deploy \
dockerfile docker-push php gitbook js-script sonar integration-test unit-test api-test  java-lint testplan java-dependency-check golang java-unit android ios \
mobile-template lib-publish mobile-publish java-deploy extract-repo-version release-fetch dingtalk-robot-msg oss-upload delete-nodes ess-info loop api-register \
api-publish publish-api-asset mysqldump archive-release erda-get-addon-info erda-get-service-addr erda-mysql-migration push-extensions archive-extensions \
testscene-run testplan-run contrast-security erda-create-custom-addon project-artifacts project-package semgrep:

	@set -eo pipefail

	@echo start make action: $@

	version="$${VERSION}"
	if [[ "$${version}" == "" ]]; then
		if [[ "`ls -l actions/$@ | wc -l | tr -d ' '`" != "2" ]]; then
			echo Multi version [$$(echo `ls actions/$@` | sed 's/ /, /g')] detected, which version you want to make? \
				 specify by env: VERSION=1.0
			exit 1
		fi
		version=`ls actions/$@`
		echo Auto select the only version: $${version}
	fi

	@echo use version: $${version}

	@dockerfile="actions/$@/$${version}/Dockerfile"
	@echo expected Dockerfile: $${dockerfile}
	if [[ ! -f $${dockerfile} ]]; then echo "expected Dockerfile not exist, stop." && exit 1; fi

	@builder="$(shell docker buildx ls | grep erda-actions | awk '{print $1}' | head -n 1)"
	if [[ "$${builder}" == "" ]]; then docker buildx create --name erda-actions; fi

	docker buildx use erda-actions
	image="$(Registry)/erda-actions/$@-action:$${version}-$(Date)-${GitCommit}"
	imageForPush="$(RegistryForPush)/erda-actions/$@-action:$${version}-$(Date)-${GitCommit}"

	@echo image=$${image}

	dockerbuild="docker buildx build --push . --platform=$(PLATFORMS) -f $${dockerfile} -t $${imageForPush} \
				 --label 'branch=$(GitBranch)' --label 'commit=$(GitCommit)' --label 'build-time=$(BuildTime)'"
	# --pull
	if [[ $@ == "java-dependency-check" ]]; then dockerbuild="$${dockerbuild} --pull"; fi
	# --no-cache
	if [[ $@ == "buildpack" || $@ == "java" || $@ == "java-build" || $@ == "java-agent" ]]; then
		dockerbuild="$${dockerbuild} --no-cache"
	fi
	if [[ $@ == "golang" ]]; then dockerbuild="$${dockerbuild} --build-arg GO_VERSION=${GO_VERSION}"; fi
	@echo $${dockerbuild}
	eval "$${dockerbuild}"

	echo "action meta: $@($${version})=$${image}"

	# auto replace dice.yml image
	if [[ "${AUTO_GIT_COMMIT}" == "true" ]]; then
		diceyml="actions/$@/$${version}/dice.yml"
		yq eval ".jobs.[].image = \"$${image}\"" -i $${diceyml}
		git add .
		@git commit -m "Auto update image for action: $@, version: $${version}"
	fi

.ONESHELL:
sonarqube:
	@set -eo pipefail

	@echo start make $@

	version="$${VERSION}"
	sonarDir="actions/sonar/1.0"
	if [[ "$${version}" == "" ]]; then
		if [[ "`ls -l "$${sonarDir}/$@" | wc -l | tr -d ' '`" != "2" ]]; then
			echo Multi version [$$(echo `ls $${sonarDir}/$@` | sed 's/ /, /g')] detected, which version you want to make? \
				 specify by env: VERSION=1.0
			exit 1
		fi
		version=`ls $${sonarDir}/$@`
		echo Auto select the only version: $${version}
	fi

	@echo use version: $${version}
	cd actions/sonar/1.0/sonarqube/$${version}
	docker build . -t registry.cn-hangzhou.aliyuncs.com/dice-third-party/sonar:$${version} -f Dockerfile

all-arch:
	echo "build all arch for $(action)"
	version="$${VERSION}"
	if [[ "$${version}" == "" ]]; then
		if [[ "`ls -l actions/$@ | wc -l | tr -d ' '`" != "2" ]]; then
			echo Multi version [$$(echo `ls actions/$(action)` | sed 's/ /, /g')] detected, which version you want to make? \
				 specify by env: VERSION=1.0
			exit 1
		fi
		version=`ls actions/$(action)`
		echo Auto select the only version: $${version}
	fi

	@echo use version: $${version}

	@dockerfile="actions/$(action)/$${version}/Dockerfile"
	@echo expected Dockerfile: $${dockerfile}
	if [[ ! -f $${dockerfile} ]]; then echo "expected Dockerfile not exist, stop." && exit 1; fi

	arches=(amd64 arm64)
	for arch in $${arches[@]}; do
		if [[ "$(DEVELOP_MODE)" == 'true' ]]; then
			echo "DEVELOP_MODE == true"
			image="$(DevelopRegistry)/$$arch/$(action)-action:$${version}-$(Date)-${GitCommit}"
			imageForPush=$${image}
		else
			image="$(Registry)/$$arch/$(action)-action:$${version}-$(Date)-${GitCommit}"
			imageForPush="$(RegistryForPush)/$$arch/$(action)-action:$${version}-$(Date)-${GitCommit}"
		fi

		@echo image=$${image}

		dockerbuild="DOCKER_BUILDKIT=1 docker build . --platform "linux/$$arch" --build-arg ARCH=$$arch -f $${dockerfile} -t $${imageForPush} \
					 --label 'branch=$(GitBranch)' --label 'commit=$(GitCommit)' --label 'build-time=$(BuildTime)'"
		# --pull
		if [[ $(action) == "java-dependency-check" ]]; then dockerbuild="$${dockerbuild} --pull"; fi
		# --no-cache
		if [[ $(action) == "buildpack" || $(action) == "java" || $(action) == "java-build" || $(action) == "java-agent" ]]; then
			dockerbuild="$${dockerbuild} --no-cache"
		fi
		if [[ $@ == "golang" ]]; then dockerbuild="$${dockerbuild} --build-arg GO_VERSION=${GO_VERSION}"; fi
		@echo $${dockerbuild}
		eval "$${dockerbuild}"

		docker push $${imageForPush}

		echo "action meta: $(action)($${version})-$$arch=$${image}"
	done
