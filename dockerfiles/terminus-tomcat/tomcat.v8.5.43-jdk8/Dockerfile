# copied some from https://github.com/docker-library/tomcat
FROM registry.erda.cloud/erda/terminus-openjdk:v1.8.0.181

ENV CATALINA_HOME /opt/tomcat
ENV PATH $CATALINA_HOME/bin:$PATH
RUN mkdir -p "$CATALINA_HOME"
WORKDIR $CATALINA_HOME

# let "Tomcat Native" live somewhere isolated
ENV TOMCAT_NATIVE_LIBDIR $CATALINA_HOME/native-jni-lib
ENV LD_LIBRARY_PATH ${LD_LIBRARY_PATH:+$LD_LIBRARY_PATH:}$TOMCAT_NATIVE_LIBDIR

ENV TOMCAT_MAJOR 8
ENV TOMCAT_VERSION 8.5.43

RUN curl -s https://www-us.apache.org/dist/tomcat/tomcat-$TOMCAT_MAJOR/v$TOMCAT_VERSION/bin/apache-tomcat-$TOMCAT_VERSION.tar.gz -o /tmp/tomcat.tar.gz && \
        tar xzf /tmp/tomcat.tar.gz --strip-components=1 -C $CATALINA_HOME && \
        rm -f /tmp/tomcat.tar.gz

EXPOSE 8080
CMD ["catalina.sh", "run"]
