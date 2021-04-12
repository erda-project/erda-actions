#!/bin/bash

set -eo pipefail

sonarqubeHome="${SONARQUBE_HOME}"
cd "${sonarqubeHome}"

(
  # 等待 h2 db 启动, 60s is enough
  sleep 60s
  # 获取用户密码
  adminPass="${SONAR_ADMIN_PASSWORD}"
  if [[ "${adminPass}" == "" ]]; then
    echo "no password passed in, use default."
    adminPass="terminus/dice/sonar"
  fi

  # 使用 htpasswd 生成新密码
  bcryptedPw="$(htpasswd -nbBC 12 "" "${adminPass}" | tr -d ":\n")"
  # 将 2y 转换为 2a (sonar 使用 2a version 的 checkpw)
  # 将 $ 转义为 \$，以便在 h2-cli 中作为命令行参数使用
  #bcryptedPw="$(echo ${bcryptedPw} | sed 's/\$2y\$/\$2a\$/g' | sed 's/\$/\\\$/g')"
  bcryptedPw="$(echo ${bcryptedPw} | sed 's/\$2y\$/\$2a\$/g')"
  # 生成用于修改密码的 h2 sql
  reinstateSql="$(printf "update users set crypted_password = '%s', salt=null, hash_method='BCRYPT' where login = 'admin'" "${bcryptedPw}")"
  # 修改密码
  java -cp lib/jdbc/h2/h2-1.4.199.jar org.h2.tools.Shell -url jdbc:h2:tcp://localhost:9092/sonar -sql "${reinstateSql}"
  printf "=====================================\n  set new password for admin success!\n=====================================\n"

  # 插入 user token

) &

(
  # 等待 h2 db 启动, 60s is enough
  sleep 60s
  # 获取用户 token
  adminToken="${SONAR_ADMIN_TOKEN}"
  if [[ "${adminToken}" == "" ]]; then
    echo "no token passed in, use default."
    adminToken="terminus-dice-sonar-action-token"
  fi

  tokenName="auto-generated-sonar-action-token"

  # 删除已有 token
  deleteTokenSql="$(printf "delete from user_tokens where name='%s'" "${tokenName}")}"
  java -cp lib/jdbc/h2/h2-1.4.199.jar org.h2.tools.Shell -url jdbc:h2:tcp://localhost:9092/sonar -sql "${deleteTokenSql}"

  # sha384hex 处理 token
  tokenSha384hexed="$(echo -n "${adminToken}" | sha384sum | tr -d " \-\n")"

  # 插入新 token
  insertTokenSql="$(printf "insert into user_tokens
  (user_uuid, name, token_hash, last_connection_date, created_at, uuid) values ((select uuid from users where login='admin'), '%s', '%s', null, '%s', '%s')" \
  "${tokenName}" "${tokenSha384hexed}" "$(date +%s%3N)" "${tokenName}")"
  java -cp lib/jdbc/h2/h2-1.4.199.jar org.h2.tools.Shell -url jdbc:h2:tcp://localhost:9092/sonar -sql "${insertTokenSql}"
  printf "=====================================\n  set new token for admin success!\n=====================================\n"
) &

# 启动服务
exec bin/run.sh bin/sonar.sh
