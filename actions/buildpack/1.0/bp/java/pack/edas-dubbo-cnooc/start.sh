#!/bin/bash

JAVA_OPTS="$JAVA_OPTS {{JAVA_OPTS}}"

export USER_JAVA_OPTS="$JAVA_OPTS"

limit_in_bytes=$(cat /sys/fs/cgroup/memory/memory.limit_in_bytes)

# If not default limit_in_bytes in cgroup
if [ "$limit_in_bytes" -ne "9223372036854771712" ]
then
    limit_in_megabytes=$(expr $limit_in_bytes \/ 1048576)

    # JAVA_MAX_MEM_RATIO not given, we calc it based on total mem
    if [ -z "${JAVA_MAX_MEM_RATIO:-}" ]
    then
        if [ "$limit_in_megabytes" -ge 4096 ]
        then
            export JAVA_MAX_MEM_RATIO=75
        elif [ "$limit_in_megabytes" -ge 1024 ]
        then
            export JAVA_MAX_MEM_RATIO=70
        elif [ "$limit_in_megabytes" -ge 300 ]
        then
            export JAVA_MAX_MEM_RATIO=50
        else
            export JAVA_MAX_MEM_RATIO=25
        fi
    fi

    if ! [ "${JAVA_MAX_MEM_RATIO}" -eq 0 ]
    then
        max_size=$(expr $limit_in_megabytes \* ${JAVA_MAX_MEM_RATIO} \/ 100)
        # no -Xmx exist
        if ! echo "${JAVA_OPTS:-}" | grep -q -- "-Xmx"
        then
            export JAVA_OPTS="-Xmx${max_size}m $JAVA_OPTS"
        fi
        # no -Xms exist
        if ! echo "${JAVA_OPTS:-}" | grep -q -- "-Xms"
        then
            export JAVA_OPTS="-Xms${max_size}m $JAVA_OPTS"
        fi
    fi

    export JAVA_OPTS="-XX:+PrintGCDetails -XX:+PrintGCTimeStamps $JAVA_OPTS"
    echo JAVA_OPTS=$JAVA_OPTS
fi

# if user add DISABLE_PRESET_JAVA_OPTS env clear erda JAVA_OPTS
if [ "${DISABLE_PRESET_JAVA_OPTS}" = "true" ]
then
  export JAVA_OPTS="$USER_JAVA_OPTS"
fi

# spot java agent
if [ -f /opt/spot/spot-agent/spot-agent.jar ]; then
    export JAVA_OPTS="$JAVA_OPTS -javaagent:/opt/spot/spot-agent/spot-agent.jar"
fi

# spot java profiler
if [ -f /opt/spot/spot-agent/spot-profiler.jar ]; then
    export JAVA_OPTS="$JAVA_OPTS -javaagent:/opt/spot/spot-agent/spot-profiler.jar"
fi

$JAVA_HOME/bin/java -Dedas.dubbo.subscribe.time=1000 $JAVA_OPTS -jar $CATALINA_OPTS -Djava.security.egd=file:/dev/./urandom -Dalicloud.deployment.mode=EDAS_MANAGED -Dcatalina.logs=$CATALINA_HOME/logs -Dpandora.location=$CATALINA_HOME/deploy/taobao-hsf.sar   "/home/admin/app/${APP_DIR}/app.jar"  --spring.profiles.active=$SPRING_PROFILES_ACTIVE  --server.context-path=/  --server.tomcat.uri-encoding=UTF-8 --server.tomcat.max-threads=400
