#!/bin/bash

set -e -x

if [ ! -d ".git" ];then
    echo "git init ."
    git init .
fi

if [[ -z "$(git config --global user.email)" ]]; then
    git config --global user.email "dice@terminus.io"
fi

if [[ -z "$(git config --global user.name)" ]]; then
    git config --global user.name "dice"
fi

git add .
git commit -m "init templates"

git remote add origin $1
git push origin master
