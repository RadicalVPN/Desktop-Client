# RadicalVPN Desktop Client (Windows/MacOS)

[![RadicalVPN Desktop Build](https://github.com/RadicalVPN/Desktop-Client/actions/workflows/ci.yml/badge.svg)](https://github.com/RadicalVPN/Desktop-Client/actions/workflows/ci.yml)
[![Nightly Build](https://github.com/RadicalVPN/Desktop-Client/actions/workflows/nightly.yml/badge.svg)](https://github.com/RadicalVPN/Desktop-Client/actions/workflows/nightly.yml)
[![RadicalVPN Desktop Build](https://github.com/RadicalVPN/Desktop-Client/actions/workflows/release.yml/badge.svg)](https://github.com/RadicalVPN/Desktop-Client/actions/workflows/release.yml)

**Official desktop client for RadicalVPN.**

<img src="https://radicalvpn.com/logo_dark.svg" alt="RadicalLogo" width="200"/>

- [About RadicalVPN Desktop](#about)
- [Project structure](#project-structure)
- [Requirements](#requirements)
  - [Windows](#requirements-windows)
  - [MacOS](#requirements-macos)
- [Building for production](#building)
  - [MacOS](#building-macos)
  - [Windows](#building-windows)
- [Local development](#local-development)
  - [Running the daemon](#local-development-daemon)
  - [Running the GUI](#local-development-gui)
- [Security](#security)

<a name="about"></a>

## About RadicalVPN Desktop

This repository contains the source code for the official desktop client for RadicalVPN. It is based on [Electron](https://electronjs.org/), [Vue](https://vuejs.org/) and [Go](https://golang.org/).

<a name="project-structure"></a>

## Project structure

The project is split into two parts:

- **Daemon** - The Go daemon that handles the VPN connection and the communication with the RadicalVPN API.
- **GUI** - The User Interface, built with Vue and Electron.

<a name="requirements"></a>

## Requirements

<a name="requirements-windows"></a>

### Windows

- TBD

<a name="requirements-macos"></a>

### MacOS

- [Go Lang](https://go.dev) 1.19
- [Node.JS](https://nodejs.org) 20.x
- [Homebrew](https://brew.sh)

To install all build dependencies, run the following command:

```bash
brew install autoconf automake libtool
```

<a name="security"></a>

<a name="building"></a>

## Building production binaries

<a name="building-macos"></a>

### MacOS

To build the MacOS binaries, run the following command:

```bash
bash build/darwin/build.sh -c {TEAM_ID}
```

<a name="building-windows"></a>

### Windows

- TBD

<a name="local-development"></a>

## Local development

<a name="local-development-daemon"></a>

### Running the daemon

```bash
cd radicalvpnd
RADICALVPND_SECRET="test" RADICALVPND_PORT="80" sudo --preserve-env=RADICALVPND_SECRET --preserve-env=RADICALVPND_PORT go run -tags debug  .
```

<a name="local-development-gui"></a>

### Running the GUI

```bash
cd gui
npm install
npm run dev
```

<a name="security"></a>

## Security Policy

If you discover any security related issues, please visit our [Security File](/SECURITY.md)
