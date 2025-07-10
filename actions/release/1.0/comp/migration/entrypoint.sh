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

# If not default limit_in_bytes in cgroup
if [ "${memory_unlimited}" -ne 1 ]
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

    # UseContainerSupport: default is true after version JDK8u191
    export JAVA_OPTS="-XX:+PrintGCDetails -XX:+PrintGCTimeStamps $JAVA_OPTS"
    echo JAVA_OPTS=$JAVA_OPTS
fi

# spot java agent
if [ -f /opt/spot/java-agent/java-agent.jar ]; then
    export JAVA_OPTS="$JAVA_OPTS -javaagent:/opt/spot/java-agent/java-agent.jar"
fi

exec java $JAVA_OPTS -Djava.security.egd=file:/dev/./urandom -jar \
/app.jar --spring.profiles.active=${SPRING_PROFILES_ACTIVE}