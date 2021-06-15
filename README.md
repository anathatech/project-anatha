# Project Anatha

The official golang implementation for Project Anatha.

---

For instructions on setting up a validator on the Anatha network, view the guide here: https://app.gitbook.com/@project-anatha/s/validator-guide/

---

## Executables

|   Command   | Description                             |
|--------------------|-----------------------------------------|
| `anathacli`        | The command line Anatha Client tools.  Provides functionality to consume and participate in an existing Anatha Network. |
| `anathad`          | The Anatha daemon for configuring and running a validator node on a network. |
| `anathad-manager`  | A process manager for wrapping the the Anatha daemon, to provide additional functionality. |

Additional details and documentation setup steps (TBD)...

---


## Dependencies
To install required Go dependencies, use:
```
go mod tidy
```

It'll fetch all packages used by Project Anatha and install them accordingly.

And you can rebuild the binaries manually with:
```
make install
```

## Multi-node test network
Local multi-node will run 4 concurrent blockchain validator nodes locally in Docker

Apart from having Go installed, to run the multi-node local network, you'll also have to have Docker and docker-compose installed.

The easiest way to start the local multi-node network is to run the following command:
```
make localnet
```

This command will execute all necessary steps, it'll rebuild the anatha binaries for Linux environment, it'll rebuild the Docker image and start the local multi-node network using docker-compose.

If you only want to restart the nodes without rebuilding the binaries and the Docker image, you can use:
```
make localnet-start
```

If you want more control over the process, you can run these steps manually.

`make build-linux` will rebuild the binaries for Linux environment

`make build-docker` will rebuild the docker images.

---
