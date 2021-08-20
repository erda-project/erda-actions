GOLANG_VERSION=$1
GO_REL_SHA256=$2

image=registry.erda.cloud/erda/terminus-golang:${GOLANG_VERSION}

docker build --build-arg GOLANG_VERSION=${GOLANG_VERSION} --build-arg GO_REL_SHA256=${GO_REL_SHA256} -t ${image} .
docker push ${image}
