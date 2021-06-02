#!/bin/bash

git clone https://github.com/intel-secl/secure-docker-daemon 2>/dev/null

cd secure-docker-daemon
git fetch
git checkout v3.6.0
git pull

#Build secure docker daemon

make >/dev/null

if [ $? -ne 0 ]; then
  echo "could not build secure docker daemon"
  exit 1
fi

echo "Successfully built secure docker daemon"
