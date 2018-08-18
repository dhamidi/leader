#!/usr/bin/fish

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

bind \\ run-leader
