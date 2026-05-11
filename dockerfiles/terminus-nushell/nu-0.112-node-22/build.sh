
image=registry.erda.cloud/erda/terminus-debian-nu:0.112-node.22.22.lts

docker build -t ${image} .
docker push ${image}
