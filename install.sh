#!/bin/bash

which brew &> /dev/null

if [ $? -eq 1 ]
then
	/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
fi
which go &> /dev/null

if [ $? -eq 1 ]
then
	brew install go
fi






