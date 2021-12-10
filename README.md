# github-actions

This repo builds a Docker image which collects tooling used to setup environments. It serves as documentation of what tools and what versions are in use. Currently it contains primarily:

- Terraform
- Packer

## Using

A [template](./docker/template) for starting a new Terraform IaC repository.

## Building

Build the tooling image:s

```shell
docker build -t dev docker/
```

## Releasing

In order to push a new image to the container registry, you create a new release from the GitHub UI or via the API. The publish release event will trigger a GitHub pipeline that deploys the container from the release tag.

If you need to push a custom image to the registry, you need to go to your GitHub [personal access tokens](https://github.com/settings/tokens) page and create an access token. That token is your password when logging in:

```shell
docker login ghcr.io --username <GITHUB_USERNAME>
docker build -t ghcr.io/xenitab/github-actions/tools:<TAG> ./docker
docker push ghcr.io/xenitab/github-actions/tools:<TAG>
```
