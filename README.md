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

[d]: https://github.com/docker/docker-credential-helpers
