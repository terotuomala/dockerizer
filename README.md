# Dockerizer
> NB. This project is just a get to know for golang cli applications for me and not meant to offer "all encompassing" features.

An interactive CLI for creating initial Dockerfile for Go and Node.js applications using good practises such as:

- In order to keep the Docker image size optimal a multi-stage builds is used
- Only the layers from the `release` stage are pushed when the Docker image is build
- SHA256 digest pinning is used to achieve reliable and reproducable builds

<p align="center"><img src="./dockerizer-usage-flow.gif?raw=true"/></p>

## Usage

> I like to use [go-task](https://taskfile.dev/installation/) in order to make the usage a bit easier:

### Install depedencies
```sh
task dependecies
```

### Run the application
```sh
task run
```

### Build
```sh
task build
```