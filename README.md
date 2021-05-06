# drone-git

Drone plugin to clone `git` repositories.

## Build

Build the Docker image with the following commands:

```
docker build --rm -f docker/Dockerfile.linux.amd64 -t drone/git .
```

## Usage

Clone a commit:

```
docker run --rm \
  -e DRONE_WORKSPACE=/drone \
  -e DRONE_REMOTE_URL=https://github.com/drone/envsubst.git \
  -e DRONE_BUILD_EVENT=push \
  -e DRONE_COMMIT_SHA=15e3f9b7e16332eee3bbdff9ef31f95d23c5da2c \
  -e DRONE_COMMIT_BRANCH=master \
  drone/git
```
