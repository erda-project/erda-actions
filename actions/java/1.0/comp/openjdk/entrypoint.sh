#!/bin/bash

if [[ "$CONTAINER_VERSION" == v11* ]]; then
    alternatives --set java $(alternatives --list | grep java_sdk_11  | awk '{print $3}' | head -n 1)/bin/java
    alternatives --set javac $(alternatives --list | grep java_sdk_11  | awk '{print $3}' | head -n 1)/bin/javac
    export JAVA_HOME=/usr/lib/jvm/java-11
    echo export JAVA_HOME=/usr/lib/jvm/java-11 >> /root/.bashrc
    echo export JAVA_HOME=/usr/lib/jvm/java-11 >> /home/dice/.bashrc
fi

limit_in_bytes=$(cat /sys/fs/cgroup/memory/memory.limit_in_bytes)

version=$(java -version 2>&1 | awk -F '"' '/version/ {print $2}')
echo version "$version"
IFS=. read major minor extra <<<"$version";
echo major "$major"

export USER_JAVA_OPTS="$JAVA_OPTS"

# If not default limit_in_bytes in cgroup
if [ "$limit_in_bytes" -ne "9223372036854771712" ] && [ -z "${JAVA_OPTS_DISABLE_PRESET:-}" ]
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

    if (( major < 11 )); then
      export JAVA_OPTS="-Djava.security.egd=file:/dev/./urandom $JAVA_OPTS"
      export JAVA_OPTS="-XX:+UseContainerSupport -XX:+UnlockExperimentalVMOptions -XX:+UseCGroupMemoryLimitForHeap -XX:NewRatio=1 -XX:+UseConcMarkSweepGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSCompactAtFullCollection -XX:CMSInitiatingOccupancyFraction=70 $JAVA_OPTS"
    fi
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


export PROFILING_ENABLED=${PROFILING_ENABLED:-true}
## spot java profiler
if [[ $PROFILING_ENABLED == true ]] && [ -e "/opt/pyroscope/pyroscope.jar" ]; then
  echo "profiling enabled"

  export JAVA_TOOL_OPTIONS="$JAVA_TOOL_OPTIONS -javaagent:/opt/pyroscope/pyroscope.jar"

  export PYROSCOPE_APPLICATION_NAME=${PYROSCOPE_APPLICATION_NAME:-$DICE_SERVICE}
  export PYROSCOPE_SERVER_ADDRESS=${PYROSCOPE_SERVER_ADDRESS:-$DICE_COLLECTOR_URL}
  export PYROSCOPE_FORMAT=${PYROSCOPE_FORMAT:-jfr}
  export PYROSCOPE_PROFILER_LOCK=${PYROSCOPE_PROFILER_LOCK:-10ms}
  export PYROSCOPE_PROFILER_ALLOC=${PYROSCOPE_PROFILER_ALLOC:-2m}
  export PYROSCOPE_UPLOAD_INTERVAL=${PYROSCOPE_UPLOAD_INTERVAL:-100s}
  export PYROSCOPE_SAMPLING_DURATION=${PYROSCOPE_SAMPLING_DURATION:-1s}

  export PYROSCOPE_LABELS="$PYROSCOPE_LABELS,DICE_SERVICE=${DICE_SERVICE},DICE_WORKSPACE=${DICE_WORKSPACE},DICE_CLUSTER_NAME=${DICE_CLUSTER_NAME},POD_IP=${POD_IP}"
  export PYROSCOPE_LABELS="$PYROSCOPE_LABELS,DICE_ORG_ID=${DICE_ORG_ID},DICE_PROJECT_ID=${DICE_PROJECT_ID},DICE_APPLICATION_ID=${DICE_APPLICATION_ID},DICE_APPLICATION_NAME=${DICE_APPLICATION_NAME}"
else
  echo "profiling disabled"
fi

if [ "${OPEN_JACOCO_AGENT}" = "true" ]
then
  echo "OPEN_JACOCO_AGENT"
  export JACOCO_PORT=${JACOCO_PORT-"6300"}
  export JACOCO_INCLUDES=${JACOCO_INCLUDES-"*"}
  export JACOCO_EXCLUDES=${JACOCO_EXCLUDES}
  export JACOCO_INCLBOOTSTRAPCLASSES=${JACOCO_INCLBOOTSTRAPCLASSES-"false"}
  export JACOCO_INCLNOLOCATIONCLASSES=${JACOCO_INCLNOLOCATIONCLASSES-"false"}
  export JAVA_OPTS="$JAVA_OPTS -javaagent:/opt/jacoco/jacocoagent.jar=address=*,port=$JACOCO_PORT,dumponexit=false,output=tcpserver,includes=$JACOCO_INCLUDES,excludes=$JACOCO_EXCLUDES,inclbootstrapclasses=$JACOCO_INCLBOOTSTRAPCLASSES,inclnolocationclasses=$JACOCO_INCLNOLOCATIONCLASSES"
fi

# print JAVA_OPTS、JAVA_TOOL_OPTIONS
echo JAVA_OPTS=${JAVA_OPTS}
echo JAVA_TOOL_OPTIONS="${JAVA_TOOL_OPTIONS}"

export FILEBEAT_CONFIG=${FILEBEAT_CONFIG:-/assets/filebeat-$DICE_WORKSPACE.yml}
# TODO: for ibm filebeat baseimage
if [[ -f /opt/filebeat/filebeat && -f ${FILEBEAT_CONFIG} && -z "${FILEBEAT_DISABLE:-}" ]]; then
    echo "Run filebeat (${FILEBEAT_CONFIG})"
    /opt/filebeat/filebeat -c ${FILEBEAT_CONFIG} &
fi

bash /pre_start.sh ${SCRIPT_ARGS}

if [ $# -eq 0 ]; then
  if [ -d /app/app ]; then
    run_cmd=$(find -L /app/app -type f -perm -u=x -not -name '*.bat' | head -n 1)
    if [ -z "$run_cmd" ]; then
      echo "not found run bin!!!"
      exit 1
    else
      eval $run_cmd
    fi
  else
    # TODO: not only jar
    exec java $JAVA_OPTS -jar /app/app.jar
  fi
else
  exec "$@"
fi
