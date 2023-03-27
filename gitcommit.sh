#!/bin/bash
function todir() {
  pwd
}

function pull() {
  todir
  echo "git pull"
  git pull
}

function forcepull() {
  todir
  echo "git fetch --all && git reset --hard origin/master && git pull"
  git fetch --all && git reset --hard origin/master && git pull
}


#  shellcheck disable=SC2120
function gitpush() {
  commit=""
  if [ ! -n "$1" ]; then
    commit="$(date '+%Y-%m-%d %H:%M:%S') by ${USER}"
  else
    commit="$1 by ${USER}"
  fi

  echo $commit
  git add .
  git commit -m "$commit"
  #  git push -u origin main
  git push
}

function m() {
    echo "1. 强制更新"
    echo "2. 普通更新"
    echo "3. 提交项目"
    echo "请输入编号:"
    read index

    case "$index" in
    [1]) (forcepull);;
    [2]) (pull);;
    [3]) (gitpush);;
    *) echo "exit" ;;
  esac
}

function bootstrap() {
    case $1 in
    pull) (pull) ;;
    m) (m) ;;
      -f) (forcepull) ;;
       *) ( gitpush $1)  ;;
    esac
}


bootstrap m
