# CI

<!-- @fadeev why is this doc titled CI? do we need to keep this doc? Are we defining the prerequisites for continuous integration? --> By default, a scaffolded blockchain includes a Github action that builds for AMD64 and ARM64 platforms on Windows, Mac, and Linux.

## Docker Images

For Docker images to build successfully, you must add your Docker Hub credentials as GitHub Actions secrets:

```
https://github.com/{username}/{repository}/settings/secrets/actions
```

Add actions secrets for these environment variables:

```
DOCKERHUB_USERNAME
DOCKERHUB_TOKEN
```

To generate a Docker access token, go to your Docker Hub Account Settings > [security](https://hub.docker.com/settings/security).
