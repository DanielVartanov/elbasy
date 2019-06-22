What to do next
================

* Make the first version runnable locally but fully working
** Make a transparent proxy
*** Quickly google how exactly HTTP (without TLS) supports proxying

*** Write the simplest proxy
**** [V] Bind to a TCP port (which one by convention?)
**** [V] Accept client connections
**** [V] Write tests for it (faster development further! No need to run curls etc)
***** [V] Google whether there are special tools for testing HTTP
***** [V] In a separate goroutine run a mock HTTP
**** In tests make mock server add all reached requests into an output channel
**** [ ] Make the code nice and pleasant to work with
**** [ ] Initiate a connection from a proxy to a server
**** [ ] Read from client, write to server
**** [ ] Read from server, write to client
**** Detect server connection closure, close client connection if so (it's not our job to act on Keep-Alive)

*** Make proxy parse the requests
**** Read out entire request (how did we detect the request is fully read? Content-Length?)
**** Open a server connection and put a request there only when entire request is fully read

*** Implement leaky-bucket algorithm
**** Make requests wait for quota using a buffered channel and a goroutine which replenishes quota every N milliseconds
**** Make proxy count quotas depending on... (is the quota per shop or per client ID? Well, does not matter, implement the separation of quotas)

* Make the first version presentable at a meetup/conference
** [ ] Make mock server stop gracefully (stop it over a quit-channel?)
** [ ] Make proxy server stop gracefully ("use of closed connection" is annoying)

* Make the first version usable in production

* How early will we need HTTP/2 or even /3? Probably it's a long way until it's mandatory, leave it for now
