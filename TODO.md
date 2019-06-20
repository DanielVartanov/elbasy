What to do next
================

* Make the first version runnable locally but fully working
** Make a transparent proxy
*** Quickly google how exactly HTTP (without TLS) supports proxying

*** Write the simplest proxy
**** [V] Bind to a TCP port (which one by convention?)
**** [V] Accept client connections
**** Initiate a connection to a server
**** Read from client, write to server
**** Read from server, write to client
**** Detect server connection closure, close client connection if so (it's not our job to act on Keep-Alive)
**** Write tests for it
***** Google whether there are special tools for testing HTTP
***** In a separate goroutine run a mock HTTP server which adds all reached requests into an output channel
****** Make it stop and start gracefully (stop it over a quit-channel?)

*** Make proxy parse the requests
**** Read out entire request (how did we detect the request is fully read? Content-Length?)
**** Open a server connection and put a request there only when entire request is fully read

*** Implement leaky-bucket algorithm
**** Make requests wait for quota using a buffered channel and a goroutine which replenishes quota every N milliseconds

* Make the first version presentable at a meetup/conference

* Make the first version usable in production

* How early will we need HTTP/2 or even /3? Probably it's a long way until it's mandatory, leave it for now
