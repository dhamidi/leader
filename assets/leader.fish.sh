#!/usr/bin/fish

function do-nothing
end

function run-leader
   set --export FISH_INPUT (commandline)
   set --export FISH_POINT (commandline -C)
   set next_command (leader print)
   if test $status -eq 3
      commandline --insert "\\"
   else
      echo $next_command | source
   end
end

function fish_user_key_bindings
  bind \\ run-leader
  bind -M literal --sets-mode default '' self-insert
  bind -M default --sets-mode literal \cv do-nothing
end
