<p align="center">
  <img src="https://raw.githubusercontent.com/filebrowser/logo/master/banner.png" width="550"/>
</p>

![Preview](https://user-images.githubusercontent.com/5447088/50716739-ebd26700-107a-11e9-9817-14230c53efd2.gif)

[![Build](https://github.com/filebrowser/filebrowser/actions/workflows/main.yaml/badge.svg)](https://github.com/filebrowser/filebrowser/actions/workflows/main.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/filebrowser/filebrowser?style=flat-square)](https://goreportcard.com/report/github.com/filebrowser/filebrowser)
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/filebrowser/filebrowser)
[![Version](https://img.shields.io/github/release/filebrowser/filebrowser.svg?style=flat-square)](https://github.com/filebrowser/filebrowser/releases/latest)
[![Chat IRC](https://img.shields.io/badge/freenode-%23filebrowser-blue.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23filebrowser)

filebrowser provides a file managing interface within a specified directory and it can be used to upload, delete, preview, rename and edit your files. It allows the creation of multiple users and each user can have its own directory. It can be used as a standalone app.

## Install

To start this app you have to install:
- [Go lang v1.20.5](https://go.dev/dl/)
- [Node.js v16.19.0](https://nodejs.org/uk/blog/release/v16.19.0)
- [Redis server v7.0](https://redis.io/download/)

## Configuration

[Authentication Method](https://filebrowser.org/configuration/authentication-method) - You can change the way the user authenticates with the filebrowser server

[Command Runner](https://filebrowser.org/configuration/command-runner) - The command runner is a feature that enables you to execute any shell command you want before or after a certain event.

[Custom Branding](https://filebrowser.org/configuration/custom-branding) - You can customize your File Browser installation by change its name to any other you want, by adding a global custom style sheet and by using your own logotype if you want.

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