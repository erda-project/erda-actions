#!/bin/bash
limit_in_bytes=$(cat /sys/fs/cgroup/memory/memory.limit_in_bytes)

export USER_JAVA_OPTS="$JAVA_OPTS"

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
        else
            export JAVA_MAX_MEM_RATIO=50
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

    export JAVA_OPTS="-XX:+UnlockExperimentalVMOptions -XX:+UseCGroupMemoryLimitForHeap -XX:NewRatio=1 -XX:+UseConcMarkSweepGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSCompactAtFullCollection -XX:CMSInitiatingOccupancyFraction=70 $JAVA_OPTS"
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

# spot java profiler
if [ -f /asia/dice_files/ArmsAgent/arms-bootstrap-1.7.0-SNAPSHOT.jar ] && [ "$PACK_TYPE" == "asia" ]; then
    export JAVA_OPTS="$JAVA_OPTS -javaagent:/asia/dice_files/ArmsAgent/arms-bootstrap-1.7.0-SNAPSHOT.jar  -Darms.licenseKey=$LICENSE_KEY -Darms.appName=$DICE_CLUSTER_NAME-$DICE_WORKSPACE-$DICE_APPLICATION_NAME-$DICE_SERVICE_NAME"
fi

# spot java profiler
if [ -f /chngc/dice_files/ArmsAgent/arms-bootstrap-1.7.0-SNAPSHOT.jar ] && [ "$PACK_TYPE" == "chngc" ]; then
    export JAVA_OPTS="$JAVA_OPTS -javaagent:/chngc/dice_files/ArmsAgent/arms-bootstrap-1.7.0-SNAPSHOT.jar  -Darms.licenseKey=$LICENSE_KEY -Darms.appName=$DICE_CLUSTER_NAME-$DICE_WORKSPACE-$DICE_APPLICATION_NAME-$DICE_SERVICE_NAME"
fi

exec java $JAVA_OPTS -Djava.security.egd=file:/dev/./urandom -jar \
/app/${APP_DIR}/app.jar --spring.profiles.active=${SPRING_PROFILES_ACTIVE}
