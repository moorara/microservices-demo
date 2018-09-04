# Docker Swarm

## Commands

| Command         | Description                                             |
|-----------------|---------------------------------------------------------|
| `make up`       | Brings up a local Swarm environment using `vagrant` vms |
| `make down`     | Destroys the Swarm environment `vagrant` vms            |
| `make clean`    | Removes created directories and files                   |

## Accessing Swarm

| Command                    | Example                  | Description                                        |
|----------------------------|--------------------------|----------------------------------------------------|
| `./connect.sh <node-name>` | `./connect.sh manager-1` | Opens a tunnel to docker running on specified node |

## Documentation

  - https://docs.docker.com/engine/reference/commandline/swarm_init
  - https://docs.docker.com/engine/reference/commandline/swarm_update
  - https://docs.docker.com/engine/reference/commandline/swarm_join-token
  - https://docs.docker.com/engine/reference/commandline/node

  * https://docs.docker.com/engine/swarm/ingress
  * https://docs.docker.com/engine/swarm/configs
  * https://docs.docker.com/engine/swarm/secrets
  * https://docs.docker.com/engine/swarm/stack-deploy
  * https://docs.docker.com/engine/swarm/how-swarm-mode-works/nodes
  * https://docs.docker.com/engine/swarm/how-swarm-mode-works/services
