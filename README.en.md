# ConfigWeb

ConfigWeb is a tool for mapping the configuration of a real infrastructure to the configuration of an abstract application package.

## Summary

The problem ConfigWeb aims to solve is simple.
Let's say you have separated multiple packages that utilize a Redis Cluster based on their roles and responsibilities.
In this case, your application project will have multiple packages (or classes) that serve as Redis Cluster clients with different names.
These packages will have their own configurations, and these configurations must match the settings of the actual Redis Cluster they connect to.

If there are as many Redis Clusters as there are packages, you can simply input the Addresses and Password of each Redis Cluster into each configuration.
However, if there is an M:N relationship between packages and the actual configured Redis Clusters, this configuration becomes very inefficient.
Furthermore, it becomes difficult to determine which package is connected to which Redis Cluster infrastructure later on, and it's also hard to modify this relationship.
Therefore, there is a need for a tool that, when you define the configuration for each Redis Cluster infrastructure and specify which infrastructure each package connects to, maps the environment settings according to that configuration.

To solve this problem, ConfigWeb was created to map the configuration of a real infrastructure to the configuration of an abstract application package.

## Installation

```bash
go install github.com/snowmerak/configweb@latest
```

## Usage

1. Project Creation

```bash
mkdir example
configweb n -p .
```

First, open the `config-set.yaml` file and modify it as follows:

```yaml
providers:
    - name: SHARED_VALKEY_CLUSTER
      type: yaml
      location: ./infra/shared_valkey_cluster.yaml
```

2. Create Infrastructure Configuration File

Create and modify the previously registered `shared_valkey_cluster.yaml` file.

```bash
configweb n -y shared_valkey_cluster
```

When executed, the `./infra/shared_valkey_cluster.yaml` file is created, and modify it as follows:

```yaml
ADDRESSES:
    - 10.4.3.5
    - 10.4.3.7
    - 10.4.3.9
PASSWORD: z1mp8qow5jd2jj2
```

3. Create Package Configuration File

```bash
configweb n -k auth
```

When executed, the `./package/auth.yaml` file is created, and modify it as follows:

```yaml
APP_NAME: auth
LISTEN_PORT: 3030
LOG_LEVEL: info
VALKEY_SERVER: $SHARED_VALKEY_CLUSTER
```

4. Generate Configuration File

```bash
configweb j auth.yaml
```

When executed, the `./config/auth.yaml` file is created.

```yaml
APP_NAME: auth
LISTEN_PORT: 3030
LOG_LEVEL: info
VALKEY_SERVER:
    ADDRESSES:
        - 10.4.3.5
        - 10.4.3.7
        - 10.4.3.9
    PASSWORD: z1mp8qow5jd2jj2
```
```
