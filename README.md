# Prow job generator for ci-operator

The purpose of this tool is to reduce an amount of boilerplate that component
owners need to write when they use
[ci-operator](https://github.com/openshift/ci-operator) to set up CI for their
component. The generator is able to entirely generate the necessary Prow job
configuration from the ci-operator configuration file.

## Use

To use the generator, you need to build it:

```
$ make build
```

Alternatively, you may obtain a containerized version from the registry on
`api.ci.openshift.org`:

```
$ docker pull registry.svc.ci.openshift.org/ci/ci-operator-prowgen:latest
```

### Generate Prow jobs for new ci-operator config file

The generator can use the naming conventions and directory structure of the
[openshift/release](https://github.com/openshift/release) repository. Provided
you placed your `ci-operator` configuration file to the correct place in
[ci-operator/config](https://github.com/openshift/release/tree/master/ci-operator/config),
you may run the following (`$REPO is a path to `openshift/release` working
copy):

```
$ ./ci-operator-prowgen --from-file $REPO/ci-operator/config/org/component/branch.json \
 --to-dir $REPO/ci-operator/jobs
```

This extracts the `org` and `component` from the configuration file path, reads
the `branch.json` file and generates new Prow job configuration files in the
`(...)/ci-operator/jobs/` directory, creating the necessary directory structure
and files if needed. If the target files already exist and contain Prow job
configuration, newly generated jobs will be merged with the old ones (jobs are
matched by name).

### Generate Prow jobs for multiple ci-operator config files

The generator may take a directory as an input. In this case, the generator
walks the directory structure under the given directory, finds all JSON files
there and generates jobs for all of them.

You can generate jobs for a certain component, organization, or everything:

```
$ ./ci-operator-prowgen --from-dir $REPO/ci-operator/config/org/component --to-dir $REPO/ci-operator/jobs
$ ./ci-operator-prowgen --from-dir $REPO/ci-operator/config/org --to-dir $REPO/ci-operator/jobs
$ ./ci-operator-prowgen --from-dir $REPO/ci-operator/config --to-dir $REPO/ci-operator/jobs
```

If you have cloned `openshift/release` with `go get` and you have `$GOPATH` set
correctly, the generator can derive the paths for the input/output directories.
These invocations are equivalent:

```
$ ./ci-operator-prowgen --from-release-repo --to-release-repo
$ ./ci-operator-prowgen --from-dir $GOPATH/src/github.com/openshift/release/ci-operator/config \
 --to-dir $GOPATH/src/github.com/openshift/release/ci-operator/jobs
```

## What does the generator create

The generator creates one presubmit job for each test specified in the
ci-operator config file (in `tests` list):

```yaml
presubmits:
  ORG/REPO:
  - agent: kubernetes
    always_run: true
    branches:
    - master
    context: ci/prow/TEST
    decorate: true
    name: pull-ci-ORG-REPO-BRANCH-TEST
    rerun_command: /test TEST
    skip_cloning: true
    spec:
      containers:
      - args:
        - --artifact-dir=$(ARTIFACTS)
        - --target=TEST
        command:
        - ci-operator
        env:
        - name: CONFIG_SPEC
          valueFrom:
            configMapKeyRef:
              key: BRANCH.json
              name: ci-operator-ORG-REPO
        image: ci-operator:latest
        name: ""
        resources: {}
      serviceAccountName: ci-operator
    trigger: ((?m)^/test( all| TEST),?(\\s+|$))
```

Also, if the configuration file has a non-empty `images` list, one additional
presubmit and postsubmit job is generated with `--target=[images]` option passed
to `ci-operator` to attempt to build the component images. This postsubmit job
also uses the `--promote` option to promote the component images built in this
way.

## Develop

To build the generator, run:

```
$ make build
```

To run unit-tests, run:

```
$ make test
```
