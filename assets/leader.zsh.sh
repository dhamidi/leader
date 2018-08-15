#!/bin/zsh

() {
  leader_widget() {
    local leader_exit leader_next
    leader_next=$(SHELL=/bin/zsh BUFFER=$BUFFER CURSOR=$CURSOR leader print)
    leader_exit=$?
    if [ $leader_exit -eq 3 ]; then
        BUFFER="${BUFFER}${KEYS}"
        CURSOR=$((CURSOR + $#KEYS))
        return "$leader_exit"
    fi
    eval "$leader_next"
    stty sane
  }

  zle -N leader_widget
  bindkey '\\' leader_widget
}
