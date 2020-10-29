#!/bin/sh

FLAGS=""
if [[ ! -z "${PLUGIN_DEPTH}" ]]; then
	FLAGS="--depth=${PLUGIN_DEPTH}"
fi

if [ ! -d .git ]; then
	git init
	git remote add origin ${DRONE_REMOTE_URL}
fi

# the branch may be empty for certain event types,
# such as github deployment events. If the branch
# is empty we checkout the sha directly. Note that
# we intentially omit depth flags to avoid failed
# clones due to lack of history.
if [[ -z "${DRONE_COMMIT_BRANCH}" ]]; then
	set -e
	set -x
	git fetch origin
	git checkout -qf ${DRONE_COMMIT_SHA}
	exit 0
fi

# the commit sha may be empty for builds that are
# manually triggered in Harness CI Enterprise. If
# the commit is empty we clone the branch.
if [[ -z "${DRONE_COMMIT_SHA}" ]]; then
	set -e
	set -x
	git fetch ${FLAGS} origin +refs/heads/${DRONE_COMMIT_BRANCH}:
	git checkout -b ${DRONE_COMMIT_BRANCH} origin/${DRONE_COMMIT_BRANCH}
	exit 0
fi

set -e
set -x

git fetch ${FLAGS} origin +refs/heads/${DRONE_COMMIT_BRANCH}:
git checkout ${DRONE_COMMIT_SHA} -b ${DRONE_COMMIT_BRANCH}
