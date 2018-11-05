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
