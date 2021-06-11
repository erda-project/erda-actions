# copy from androidsdk/android-29
FROM registry.erda.cloud/erda/android-29
ENV SDK_HOME /usr/local
ENV GRADLE_VERSION 6.0.1
ENV GRADLE_SDK_URL https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip
RUN curl -sSL "${GRADLE_SDK_URL}" -o gradle-${GRADLE_VERSION}-bin.zip  \
	&& unzip gradle-${GRADLE_VERSION}-bin.zip -d ${SDK_HOME}  \
	&& rm -rf gradle-${GRADLE_VERSION}-bin.zip
ENV GRADLE_HOME ${SDK_HOME}/gradle-${GRADLE_VERSION}
ENV PATH ${GRADLE_HOME}/bin:$PATH
RUN curl -sL https://deb.nodesource.com/setup_12.x | bash
RUN apt-get update && apt-get install -y nodejs git wget g++ gcc vim net-tools make dnsutils && apt-get clean