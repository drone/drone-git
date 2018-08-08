#!/bin/sh

FLAGS=""
if [[ ! -z "${PLUGIN_DEPTH}" ]]; then
	FLAGS="--depth=${PLUGIN_DEPTH}"
fi

if [ ! -d .git ]; then
	git init
	git remote add origin ${DRONE_REMOTE_URL}
fi

set -e
set -x

git fetch ${FLAGS} origin +refs/tags/${DRONE_TAG}:
git checkout -qf FETCH_HEAD
