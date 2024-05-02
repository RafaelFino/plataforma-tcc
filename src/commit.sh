#!/bin/bash
git add . 
git commit -a -m "$(curl -sL whatthecommit.com/index.txt)" 
git push 