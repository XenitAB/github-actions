# Terraform the Xenit way

This repository collects the organization's [infrastructure-as-code](https://en.wikipedia.org/wiki/Infrastructure_as_code). It is encoded using [Hashicorp Terraform](https://www.terraform.io/). It uses an Azure Storage account to keep a shared copy of the Terraform state and also performs extra validation before invoking the Terraform `plan` and `apply` commands.

## Prerequisites

This repository uses [terraform-docker](https://github.com/XenitAB/github-actions/tree/main/docker) to simpliify working with Terraform. In order to work with terraform-docker, you need to install [Docker Engine](https://docs.docker.com/engine/install/) and [make](https://en.wikipedia.org/wiki/Make_(software)) (but not Terraform itself, which is already included).

## Repository organization

The tool assumes that your Terraform is organized into separate subdirectories (called DIR), each of which holds a separate Terraform *configuration*. This repository has two default configurations:

- `governance/`: "Abstract" and administrative resources, typically groups, shared key vaults.
- `core/`: Resources that many other parts interact with, such as peering, shared firewalling, et c.

You can add new subdirectories as necessary to organize your infrastructure.

terraform-docker assumes that you have multiple copies of your configurations deployed so that you can verify any changes before rolling them out in production. Each copy is called an *environment*, or ENV for short. terraform-docker defaults to two environments, called `dev` and `prod`.

## Getting started

Start by familiarizing yourself with the [Terraform workflow](https://www.terraform.io/guides/core-workflow.html). The tool wraps the actual execution of Terraform, but the workflow is the same.

Before you executre `plan` for the first time, you need to initialize your local environment. You need to do this for every ENV/DIR combination.

```shell
make prepare ENV=dev DIR=governance
```

Once you have performed this step, you can start updating your environments.

```shell
make plan ENV=dev DIR=governance
make apply ENV=dev DIR=governance
```