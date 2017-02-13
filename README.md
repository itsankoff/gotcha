# gotcha
Chat server and client library and ready applications written in Golang.
* [Intro](#intro)
* [Install](#install)
* [Usage](#usage)
* [Troubleshooting](#troubleshooting)
* [License](#license)

<a name="intro">
## Intro
This repository provides both library which you can integrate in your project and ready Golang server and client.

<a name="install">
## Install

    go get github.com/itsankoff/gotcha
    go install github.com/itsankoff/gotcha/cmd/server
    go install github.com/itsankoff/gotcha/cmd/client

<a name="usage">
## Usage
**1. Simple usage of client and server applications:**  

    /path/to/your/workspace/bin/server
    /path/to/your/workspace/bin/client

**2. More sophisticated usage:**  

    /path/to/your/workspace/bin/server --host=<ip|host> --port=<port> --key_path=<private key path> --cert_path=<cert path>
    /path/to/your/workspace/bin/client --host=<ip|host> --port=<port>

**3. For full information about all possible options for both client and server run:**  

    /path/to/your/workspace/bin/server --help
    /path/to/your/workspace/bin/client --help

**4. For library API reference check the documentation and
    [client examples files](https://github.com/itsankoff/gotcha/blob/master/src/client/) and
    [server examples files](https://github.com/itsankoff/gotcha/blob/master/src/server/):**

<a name="troubleshooting">
## Troubleshooting

<a name="license">
## License
[MIT License](https://github.com/itsankoff/gotcha/blob/master/LICENSE)
