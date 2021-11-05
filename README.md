# tx2-k8s-gpu-plugin
This repository contains source code of gpu device plugin of Jetson tx2 board of Kubernetes container orchestration.

**NOTE:** This repo is not yet finished, while still contains codes that works but requires configuration before using it.

## TODO
1. Device plugin containerization
2. Deployment orchestration

## Prerequisites
1. A working k8s cluster
2. A Jetson Tx2 board joined into k8s cluster
3. Docker runtime set default to nvidia-docker
4. `golang` version greater than 1.16

# Usage
Running `go run .` on the installation directory on tx2 board.
