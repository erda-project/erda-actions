ARG DEP_CACHE_IMAGE
FROM ${DEP_CACHE_IMAGE} as cache-image

FROM {{BP_DOCKER_BASE_REGISTRY}}/erda/terminus-maven:3-jdk-8-alpine

ARG MAVEN_OPTS
ARG MAVEN_EXTRA_ARGS

COPY --from=cache-image /root/.m2 /root/.m2
RUN true
COPY /bp/build/maven-settings.xml /root/.m2/settings.xml

# compile
RUN true
COPY /code /code

WORKDIR /code

ARG FORCE_UPDATE_SNAPSHOT
ARG MODULES

RUN MAVEN_OPTS=${MAVEN_OPTS} mvn clean package -e -B -U -Dmaven.test.skip=true ${MODULES} ${MAVEN_EXTRA_ARGS}

RUN for module in `find ${PWD} -type d -name "target"`; do \
      find ${module} -type f -name "*.jar" -exec ls -l {} \; | sort -k5 -n -r | head -n1 | awk '{print $9}' | xargs -I {} sh -c 'cd $(dirname "{}") && mv "{}" ../app.jar' ; \
    done

RUN find ${PWD} -not -name "app.jar" -delete >/dev/null 2>&1 || true

RUN date || :
