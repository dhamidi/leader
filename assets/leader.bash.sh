#!/bin/bash

leader_widget() {
  local leader_exit leader_next
  leader_next=$(leader "$@")
  leader_exit=$?
  if [ "$leader_exit" = 3 ]; then
      READLINE_LINE="${READLINE_LINE}\\"
      READLINE_POINT=$((READLINE_POINT + 1))
      return $leader_exit
  fi
  eval "$leader_next"
}

bind -x '"\\":leader_widget print'
