# CI

By default, a scaffolded blockchain includes a Github action that builds for amd64 and arm64 on Windows, Mac, and Linux.

## Docker Images And Pi Images

In order for Docker images and Raspberry Pi images to build successfully, please add your docker hub credentials as secrets: https://github.com/{username}/{repository}/settings/secrets/actions

Add these:

```
DOCKERHUB_USERNAME
DOCKERHUB_TOKEN
```

You can get the token [here](https://hub.docker.com/settings/security).