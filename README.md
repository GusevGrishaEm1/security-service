# Security Service

Security service for authentication and authorization.

## Getting Started

This service provides essential security features for authentication and authorization in your application. Below are the steps to build and run the service using Docker.

## Prerequisites

- Docker installed on your system
- Go installed (if you plan to run the service outside of Docker)

## Configuration

Ensure you have a configuration file `dev.yaml` located at `D:\projects\security-service\config\` (or adjust the path as needed). The configuration file should contain necessary settings for the service.

## Building the Docker Image

To build the Docker image, navigate to the root directory of the project and run the following command:

```sh
docker build -t security_service .
