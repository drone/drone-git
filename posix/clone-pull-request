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

git fetch ${FLAGS} origin +refs/heads/${DRONE_COMMIT_BRANCH}:
git checkout ${DRONE_COMMIT_BRANCH}

git fetch origin ${DRONE_COMMIT_REF}:
git merge ${DRONE_COMMIT_SHA}
