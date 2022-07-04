# scredenziato

> It's a portmandeu between two italian words:</br>
> "screanzato": one who misbehaves, missing the knowledge of the basic social norms</br>
> "credenziali": credentials, the secrets needed to be recognised

### What is it?

This small tool conflates the functionalities of all the helpers found in
[docker-credential-helpers][d], in a single command line utility.

It tries to find the _correct_ store where *docker* keeps the credentials for connecting
to a registry in an authenticated way. If found, one can then query such store.

### Examples

```
$ screanzato list
registry1.example.org	user1
registry2.example.com	user2

$ screanzato get registry1.example.org
user1
THE_SUPER_SECRET_PASSPHRASE
```

Check the help menu for other options

### Compile

I am on GNU/Linux. This project has 4 supported targets:
 - `GOOS=linux GOARCH=amd64`
 - `GOOS=darwin GOARCH=amd64`
 - `GOOS=darwin GOARCH=arm64`
 - `GOOS=windows GOARCH=amd64`

Depending on the target, the resulting binary pulls in different dependencies. These, in
the case of `GOOS=linux` and `GOOS=darwin` require also to use `cgo`. For this to work
on linux with `darwin` target one needs a working [`osxcross`][o] toolchain (good luck
and follow the steps in the readme) and a windows cross compiler (mingw is the choice on
archlinux). You can build all the targets with

```
$ CC_darwin_amd64=/path/to/o64-clang \
CXX_darwin_amd64=/path/to/o64-clang++ \
CC_darwin_arm64=/path/to/aarch64-apple-darwin21.4-clang \
CXX_darwin_arm64=/path/to/aarch64-apple-darwin21.4-clang++ \
CC_windows_amd64=/path/to/x86_64-w64-mingw32-gcc \
CXX_windows_amd64=/path/to/x86_64-w64-mingw32-g++ \
make build
```

[d]: https://github.com/docker/docker-credential-helpers
[o]: https://github.com/tpoechtrager/osxcross
