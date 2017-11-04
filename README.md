# gotcha

Chat server and client library and ready applications written in Golang.
* [Intro](#intro)
* [Install](#install)
* [Usage](#usage)
* [License](#license)

---

<a name="intro">

## Intro

This repository provides both library which you can integrate in your project and ready Golang server and client.

<a name="install">

## Install

```
go get github.com/itsankoff/gotcha/cmd/server
go get github.com/itsankoff/gotcha/cmd/client
```

<a name="usage">

## Usage

**1. Simple usage of client and server applications**  
```
/path/to/your/workspace/bin/server
/path/to/your/workspace/bin/client
```

**2. More sophisticated usage**  
```
/path/to/your/workspace/bin/server --host=<ip:port> --key_path=<private key path> --cert_path=<cert path>
/path/to/your/workspace/bin/client --host=<ip:port>
```

**3. For full information about all possible options for both client and server run**  
```
/path/to/your/workspace/bin/server --help
/path/to/your/workspace/bin/client --help
```

**4. For library API reference check the documentation and
    [client examples files](https://github.com/itsankoff/gotcha/blob/master/client/client_test.go) and
    [server examples files](https://github.com/itsankoff/gotcha/blob/master/server/server_test.go)**


<a name="license">

## License
[MIT License](https://github.com/itsankoff/gotcha/blob/master/LICENSE)
