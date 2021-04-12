FROM openjdk:8u151-jdk-alpine

RUN apk add --no-cache curl tar bash jq libxml2-utils

# https://github.com/concourse/concourse/issues/2042
RUN unlink  $JAVA_HOME/jre/lib/security/cacerts && \
    cp /etc/ssl/certs/java/cacerts $JAVA_HOME/jre/lib/security/cacerts

ADD assets /opt/action/
ADD test /opt/action/test/
ADD itest /opt/action/itest/

# Run tests (also pre-seeds .m2/repository)
#RUN /opt/action/test/all.sh
