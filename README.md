# gotcha
Chat server and client library and ready applications written in Golang.
* [Intro](https://github.com/itsankoff/gotcha/blob/master/README.md#Intro)
* [Installtion](https://github.com/itsankoff/gotcha/blob/master/README.md#Installtion)
* [Usage](https://github.com/itsankoff/gotcha/blob/master/README.md#Usage)
* [Troubleshooting](https://github.com/itsankoff/gotcha/blob/master/README.md#Troubleshooting)
* [License](https://github.com/itsankoff/gotcha/blob/master/README.md#License)

## Intro
This repository provides both library which you can integrate in your project and read Goland server and client.

## Installation
go get

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

## Troubleshooting


## License
[MIT License](https://github.com/itsankoff/gotcha/blob/master/LICENSE)
