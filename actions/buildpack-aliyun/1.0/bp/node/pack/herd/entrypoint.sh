#!/bin/sh

START_CMD="${@:-`which npm` run start}"

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
    # default young generation memory is 32 MB
    max_old_space_size=$(expr $limit_in_megabytes - 32)
    echo MAX_OLD_SPACE_SIZE=$max_old_space_size
fi

exec node --max-old-space-size=$max_old_space_size $START_CMD
