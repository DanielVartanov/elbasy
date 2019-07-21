Elbasy
======

A secure proxy which prevents throttling when talking to remote APIs

Installation and deployment
---------------------------

0. Install mkcert on a local machine
https://github.com/FiloSottile/mkcert

1. Prepare a CA root certificate (it will install the CA into the local machine as well)
`$ mkcert -install"

The certificate and they key will put at `~/.local/share/mkcert`

2. Generate a leaf certificate for every API host you are going to deal with
`$ mkcert amazon.co.uk`
`$ mkcert "*.myshopify.com`

Certificates and keys will be put into the current directory

3. On each production machine:

3.1. Install the given root CA certificate to the system
Note: never leave a key to a root CA certificate on a production machine

3.2. Copy leaf certificates _and their keys_ in a subdirectory `./impostor_certificates` next to `elbasy` binary


Sponsored by [Veeqo](https://github.com/veeqo)
