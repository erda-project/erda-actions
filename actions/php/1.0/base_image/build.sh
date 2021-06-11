registry="registry.erda.cloud/erda"
version="7.2"
image="${registry}/terminus-php-apache:${version}"
dockerfile="${version}/Dockerfile"
echo image=${image}

docker build . -f ${dockerfile} -t ${image}
docker push ${image}