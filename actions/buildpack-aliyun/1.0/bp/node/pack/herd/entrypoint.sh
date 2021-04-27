#!/bin/sh

START_CMD="${@:-`which npm` run start}"

limit_in_bytes=$(cat /sys/fs/cgroup/memory/memory.limit_in_bytes)

# If not default limit_in_bytes in cgroup
if [ "$limit_in_bytes" -ne "9223372036854771712" ]
then
    limit_in_megabytes=$(expr $limit_in_bytes \/ 1048576)
    # default young generation memory is 32 MB
    max_old_space_size=$(expr $limit_in_megabytes - 32)
    echo MAX_OLD_SPACE_SIZE=$max_old_space_size
fi

exec node --max-old-space-size=$max_old_space_size $START_CMD
