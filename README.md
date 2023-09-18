<p align="center">
  <img src="https://raw.githubusercontent.com/filebrowser/logo/master/banner.png" width="550"/>
</p>


filebrowser provides a file managing interface within a specified directory and it can be used to upload, delete, preview, rename and edit your files. It allows the creation of multiple users and each user can have its own directory. It can be used as a standalone app.

## Install

## Prerequisites

Before starting the installation process, make sure you have the following dependencies installed:

- [Go Lang v1.20.5](https://go.dev/dl/)
- [Node.js v16.19.0](https://nodejs.org/uk/blog/release/v16.19.0)
- [Redis Server v7.0](https://redis.io/download/)

## Frontend Setup
### 1. Navigate to the Frontend Directory
Move to the "frontend" directory using your terminal.

### 2. Install Frontend Dependencies
Run the following command to install frontend dependencies:
```shell
npm install
```

### 3. Build the Frontend
To build the frontend for production, use:
```shell
npm run build
```
For development, use:
```shell
npm run watch
```

## Backend Setup

### 1. Navigate to the Backend Directory
Move to the "filebrowser" directory using your terminal.

### 2. Download Go Modules
Run the following command to download Go modules:
```shell
go mod download
```

### 3. Build the Backend
To build the backend for production, use:
```shell
go build
```

For development, use:
```shell
go build -tags dev
```

## Environment Configuration

Create and populate a .filebrowser.env file based on the example provided in .filebrowser.example.env.

## Running the Application

To start the application, execute the following command from the terminal:
```shell
./filebrowser
```
or this command for development:
```shell
go run main.go
```
Your application will be served on the port specified in the .filebrowser.env file.

## Technologies

BackEnd Language - [GO](https://go.dev/)
Database - [BoltDB](https://github.com/boltdb/bolt)
Database ORM - [Storm](https://github.com/asdine/storm)
File system framework - [Afero](https://github.com/spf13/afero)
Env Config - [Viper](https://github.com/spf13/viper)
CLI Interface - [Cobra](https://github.com/spf13/cobra)
API Routing = [gorilla/mux](https://github.com/gorilla/mux)


FrontEnd Language - [JS](https://developer.mozilla.org/en-US/docs/Web/JavaScript)
FrontEnd Framework - [VUE 2.6](https://v2.vuejs.org/)