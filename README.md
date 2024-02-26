# tfvars-atlantis-config

  <p align="center">
  <img src="https://github.com/3bbbeau/tfvars-atlantis-config/assets/72932978/373ea5ca-7c2f-4f52-856d-e4929d424bec" alt="tfvars-atlantis-config logo"/><br><br>
  <b>tfvars pre-workflow hook for Atlantis</b>
</p>


## Quick start
### CLI
	tfvars-atlantis-config generate --automerge --autoplan --parallel --output=atlantis.yaml

### Atlantis Server Side Config
```yaml
repos:
- id: /.*/
  pre_workflow_hooks:
    - run: tfvars-atlantis-config generate --automerge --autoplan --parallel --output=atlantis.yaml
```


## What is `tfvars-atlantis-config`?
Heavily inspired by
[terragrunt-atlantis-config](https://github.com/transcend-io/terragrunt-atlantis-config/),
this tool allows you to dynamically generate your
[Atlantis](https://runatlantis.io) configuration using `tfvars` instead
of environment hiearchies.

This is useful for teams that have simple Terraform
components and don't want the added complexity of using Terragrunt or structured
environment levels.


### Example

#### The following repo structure:
```
my-terraform
├── main.tf
├── dev.tfvars
```

#### Generates the following Atlantis configuration:

```yaml
version: 3
automerge: true
parallel_plan: true
parallel_apply: true
projects:
- name: my-terraform-dev
  dir: my-terraform
  workflow: my-terraform-dev
  autoplan:
    when_modified:
    - '*.tf'
    - dev.tfvars
    enabled: true
- name: my-terraform-prod
  dir: my-terraform
  workflow: my-terraform-prod
  autoplan:
    when_modified:
    - '*.tf'
    - prod.tfvars
    enabled: true
workflows:
  my-terraform-dev:
    plan:
      steps:
      - init
      - plan:
          extra_args:
          - -var-file=dev.tfvars
    apply:
    - apply
  my-terraform-prod:
    plan:
      steps:
      - init
      - plan:
          extra_args:
          - -var-file=prod.tfvars
    apply:
    - apply
```

## Why you should use it?
Dynamically generate your Atlantis configuration based on your Terraform components' `.tfvars` files:
* Auto plan per environment based on the environment `.tfvars` file modified.
* Auto plan all environments when a component's Terraform code is modified (`*.tf`).
* Allow your teams to organize their Terraform components in monorepos/standalone
  repositories how _they_ want, while abstracting the configuration of Atlantis
  away from them, and reducing bloat in your repos.

## Flags

Customize the behavior of this utility through CLI flag values passed in at
runtime.

| Flag Name                     | Description                                                                                                      | Default Value |
| ----------------------------- | ---------------------------------------------------------------------------------------------------------------- | ------------- |
| `--automerge`                 | Enable auto merge.                                                                                               | false         |
| `--autoplan`                  | Enable auto plan.                                                                                                | false         |
| `--default-terraform-version` | Default terraform version to run for Atlantis. Default is determined by the Terraform version constraints.       | ""            |
| `--multienv`                  | Enable injection of environment specific environment variables to each workflow.                                 | false         |
| `--output`                    | Path of the file where configuration will be generated, usually `atlantis.yaml`. Default is to write to `stdout` | `stdout`      |
| `--parallel`                  | Enables plans and applys to happen in parallel.                                                                  | false         |
| `--root`                      | Path to the root directory of the git repo you want to build config for. Default is current dir.                 | `.`           |
| `--use-workspaces`            | Whether to use Terraform workspaces for projects.                                                                | false         |
