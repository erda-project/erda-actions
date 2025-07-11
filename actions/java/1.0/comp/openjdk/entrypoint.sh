#!/bin/bash

# Parse container version (default: 8)
VERSION_NUM=${CONTAINER_VERSION#v}  # Remove 'v' prefix
VERSION_NUM=${VERSION_NUM%%.*}      # Extract major version: 11.0.6 -> 11
VERSION_NUM=${VERSION_NUM%%-*}      # Remove suffix: 17-ea -> 17
VERSION_NUM=${VERSION_NUM%%_*}      # Remove underscore: 8u292 -> 8
VERSION_NUM=${VERSION_NUM:-8}       # Default to 8 if empty

# Check and switch Java version if needed
CURRENT_JAVA=$(java -version 2>&1 | awk -F'"' '/version/{print $2}')
# Handle Java 8 version format: 1.8.x -> 8, others: 11.x.x -> 11
if [[ "$CURRENT_JAVA" =~ ^1\.([0-9]+) ]]; then
  CURRENT_MAJOR=${BASH_REMATCH[1]}
else
  CURRENT_MAJOR=${CURRENT_JAVA%%.*}
fi

if [ "$CURRENT_MAJOR" != "$VERSION_NUM" ]; then
  echo "Switching Java: $CURRENT_MAJOR -> $VERSION_NUM"
  
  # Set alternatives for all Java tools
  for tool in java javac javadoc javap jar jarsigner jps jstack jstat jmap; do
    TOOL_PATH=$(update-alternatives --list "$tool" 2>/dev/null | grep "java-$VERSION_NUM" | head -1)
    [ -n "$TOOL_PATH" ] && update-alternatives --set "$tool" "$TOOL_PATH" >/dev/null 2>&1 || echo "✗ $tool"
  done
fi

# Set Java environment
source /usr/local/bin/load-java-env.sh
echo "JAVA_HOME: $JAVA_HOME"

# Get final Java version for configuration
JAVA_VERSION=$(java -version 2>&1 | awk -F'"' '/version/{print $2}')
# Handle Java 8 version format consistently
if [[ "$JAVA_VERSION" =~ ^1\.([0-9]+) ]]; then
  JAVA_MAJOR=${BASH_REMATCH[1]}
else
  JAVA_MAJOR=${JAVA_VERSION%%.*}
fi
echo "Active Java: $JAVA_VERSION (major: $JAVA_MAJOR)"

# Initialize memory_unlimited flag
memory_unlimited=0

if [[ -f /sys/fs/cgroup/cgroup.controllers ]]; then
  # cgroup v2
  echo "Using cgroup v2"
  memory=$(cat /sys/fs/cgroup/memory.max 2>/dev/null)
  if [ "$memory" != "max" ]; then
    limit_in_bytes=$memory
  else
    memory_unlimited=1
  fi
else
  echo "Using cgroup v1"
  # default cgroup v1
  memory=$(cat /sys/fs/cgroup/memory/memory.limit_in_bytes 2>/dev/null)
  if [ "$memory" != "9223372036854771712" ]; then
    limit_in_bytes=$memory
  else
    memory_unlimited=1
  fi
fi

if [ "$memory_unlimited" -eq 1 ]; then
  echo "Memory is unlimited."
else
  echo "Memory limited in bytes: $limit_in_bytes"
fi

export USER_JAVA_OPTS="$JAVA_OPTS"

# If not default limit_in_bytes in cgroup
if [ "${memory_unlimited}" -ne 1 ] && [ -z "${JAVA_OPTS_DISABLE_PRESET:-}" ]
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

  if (( JAVA_MAJOR < 11 )); then
    export JAVA_OPTS="-Djava.security.egd=file:/dev/./urandom $JAVA_OPTS"
    # UseContainerSupport: default is true after version JDK8u191
    export JAVA_OPTS="-XX:NewRatio=1 -XX:+UseConcMarkSweepGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSCompactAtFullCollection -XX:CMSInitiatingOccupancyFraction=70 $JAVA_OPTS"
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
