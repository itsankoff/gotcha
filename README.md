# gotcha
Chat server and client library and ready applications written in Golang.
* [Intro](#intro)
* [Install](#install)
* [Usage](#usage)
* [Troubleshooting](#troubleshooting)
* [License](#license)

<a name="intro">
## Intro
This repository provides both library which you can integrate in your project and read Goland server and client.

<a name="install">
## Install
go get

<a name="usage">
## Usage
**1. Simple usage of client and server applications:**  

    ./server
    ./client

**2. More sophisticated usage:**  

    ./server --host=<ip|host> --port=<port> --key_path=<private key path> --cert_path=<cert path>
    ./client --host=<ip|host> --port=<port>

**3. For full information about all possible options for both client and server run:**  

    ./server --help 
    ./client --help

**4. For library API reference check the documentation and [examples](https://github.com/itsankoff/gotcha/blob/master/examples/):**

<a name="troubleshooting">
## Troubleshooting

<a name="license">
## License
[MIT License](https://github.com/itsankoff/gotcha/blob/master/LICENSE)
