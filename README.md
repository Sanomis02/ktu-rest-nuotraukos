# KTU REST API FAMILY PHOTO STORAGE

This repository is an example showcasing features of deploying services via `Docker-compose`, reverse proxy services using `nginx` server and backend logic featuring `Golang` controllers and http routers.


# Project structure

For future development/debugging reference to given project structure:

    .
    ├── backend						# Consists of main.go file where handlers are registered
    │   ├── Dockerfile
    │   ├── go.mod
    │   ├── go.sum
    │   ├── handlers				# List of handlers that are referenced by main.go
    │   │   ├── auth_middleware.go
    │   │   ├── image_delete.go
    │   │   ├── images_list.go
    │   │   ├── image_upload.go
    │   │   ├── login.go
    │   │   ├── register.go
    │   │   └── users.go
    │   ├── main.go
    │   └── uploads					# Mount point for images that are served by http file server
    ├── certs						# Directory to store certificate and key file
    ├── docker-compose.yml			# Configuration for services
    ├── Makefile					# Makefile for dependency installation and deployment
    ├── nginx
    │   ├── certs
    │   ├── conf.d
    │   │   └── default.conf
    │   └── html
    │       ├── index.html
    │       ├── login.html
    │       └── upload.html
    ├── README.md
    ├── scripts						# Scripts that act as wrappers for POST/GET requests
    │   ├── get_request.sh
    │   ├── image_delete.sh
    │   ├── image_upload.sh
    │   ├── login.sh
    │   └── register.sh
    └── uploads

## Deployment

There are multiple ways of deploying project. Reference to `Makefile` content. Main method of deploying project is running command as super-user `make build-services`
