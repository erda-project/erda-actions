#!/bin/bash

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

version=$(java -version 2>&1 | awk -F '"' '/version/ {print $2}')
echo version "$version"
IFS=. read major minor extra <<<"$version";
echo major "$major"

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

    if (( major < 11 )); then
      export JAVA_OPTS="-Djava.security.egd=file:/dev/./urandom $JAVA_OPTS"
      # UseContainerSupport: default is true after version JDK8u191
      export JAVA_OPTS="-XX:NewRatio=1 -XX:+UseConcMarkSweepGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSCompactAtFullCollection -XX:CMSInitiatingOccupancyFraction=70 $JAVA_OPTS"
    fi
fi

# spot java agent
if [ -f /opt/action/comp/spot-agent/spot-agent.jar ]; then
    export JAVA_OPTS="$JAVA_OPTS -javaagent:/spot-agent/spot-agent.jar"
fi

echo "JAVA_OPTS="$JAVA_OPTS"" >> $METAFILE
## spot java profiler
#if [ -f /opt/spot/spot-agent/spot-profiler.jar ]; then
#    export JAVA_OPTS="$JAVA_OPTS -javaagent:/opt/spot/spot-agent/spot-profiler.jar"
#fi

# print JAVA_OPTS at first
