#!/bin/bash
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
        if ! echo "${JAVA_TOMCAT_OPTS:-}" | grep -q -- "-Xmx"
        then
            export JAVA_TOMCAT_OPTS="-Xmx${max_size}m $JAVA_TOMCAT_OPTS"
        fi
        # no -Xms exist
        if ! echo "${JAVA_TOMCAT_OPTS:-}" | grep -q -- "-Xms"
        then
            export JAVA_TOMCAT_OPTS="-Xms${max_size}m $JAVA_TOMCAT_OPTS"
        fi
    fi

    export JAVA_TOMCAT_OPTS="-XX:+UseContainerSupport -XX:+UnlockExperimentalVMOptions -XX:+UseCGroupMemoryLimitForHeap -XX:+PrintGCDetails -XX:+PrintGCTimeStamps $JAVA_TOMCAT_OPTS"
    echo JAVA_TOMCAT_OPTS=$JAVA_TOMCAT_OPTS
fi

# spot java agent
if [ -f /opt/action/comp/spot-agent/spot-agent.jar ]; then
    export JAVA_TOMCAT_OPTS="$JAVA_TOMCAT_OPTS -javaagent:/spot-agent/spot-agent.jar"
fi

echo "JAVA_TOMCAT_OPTS="$JAVA_TOMCAT_OPTS"" >> $METAFILE
## spot java profiler
#if [ -f /opt/spot/spot-agent/spot-profiler.jar ]; then
#    export JAVA_TOMCAT_OPTS="$JAVA_TOMCAT_OPTS -javaagent:/opt/spot/spot-agent/spot-profiler.jar"
#fi
