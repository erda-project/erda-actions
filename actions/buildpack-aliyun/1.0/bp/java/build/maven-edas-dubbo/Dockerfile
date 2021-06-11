FROM registry.erda.cloud/erda/terminus-maven:3-jdk-8-alpine

RUN echo "http://mirrors.aliyun.com/alpine/v3.8/main/" > /etc/apk/repositories && \
    echo "http://mirrors.aliyun.com/alpine/v3.8/community/" >> /etc/apk/repositories

RUN apk add --no-cache python

COPY /bp/build/maven-settings.xml /root/.m2/settings.xml

ARG MAVEN_OPTS

ARG PARENT_POM_DIR

COPY ${PARENT_POM_DIR} /tmp/parent_pom/

RUN cd /tmp/parent_pom && MAVEN_OPTS=${MAVEN_OPTS} mvn dependency:go-offline --fail-never --non-recursive -e -B

ARG ALL_POM_DIR

COPY ${ALL_POM_DIR} /tmp/all_pom/

RUN cd /tmp/all_pom && MAVEN_OPTS=${MAVEN_OPTS} mvn dependency:go-offline --fail-never -B

# compile

COPY /code /code

WORKDIR /code

COPY /bp/build/rewrite_pom.py /root/rewrite_pom.py

RUN find . -name 'pom.xml' -exec python /root/rewrite_pom.py {} \;

ARG FORCE_UPDATE_SNAPSHOT

RUN MAVEN_OPTS=${MAVEN_OPTS} mvn clean package -e -B -U -Dmaven.test.skip=true

RUN for module in `find ${PWD} -type d -name "target"`; do \
      find ${module} -type f -name "*.jar" -exec ls -l {} \; | sort -k5 -n -r | head -n1 | awk '{print $9}' | xargs -I {} sh -c 'cd $(dirname "{}") && mv "{}" ../app.jar' ; \
    done

RUN find ${PWD} -not -name "app.jar" -delete >/dev/null 2>&1 || true
