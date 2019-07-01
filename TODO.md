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
**** [V] Make the code nice and pleasant to work with
**** [V] Initiate a connection from a proxy to a server
**** [V] Read from client, write to server, make a test for it
***** [V] Dump requests just to be sure everything is the same https://godoc.org/net/http/httputil#DumpRequest
***** [V] Write a test that all parts of request get passed to server (just pass entire request except for requestURL)
****** [V] Write a test that request URL becomes request path before being sent to server
**** [V] Read from server, write to client
***** [V] Write a test that all parts of response get passed to client (just pass entire response)
**** Make it local-usage-ready
***** [V] Make the main func actually run proxy server
***** [V] Make proxy server dump requests and responses in a log
***** Try it with an actual curl (or even browser?)
*** Refactor it
**** Make tests more uniform, apply common patterns etc (do we need LastRequest if we have a function?)
**** Make proxy & mock server startup process split into synchronous port binding and asyncrhonous connection handling
**** log.Fatal calls os.Exit(1) itself

// An idea:
        // ConnState specifies an optional callback function that is
        // called when a client connection changes state. See the
        // ConnState type and associated constants for details.
        http.Server.ConnState func(net.Conn, ConnState)

*** [V] Implement leaky-bucket algorithm
**** [V] Make requests wait for quota using a buffered channel and a goroutine which replenishes quota every N milliseconds
**** Make proxy count quotas depending on... (is the quota per shop or per client ID? Well, does not matter, implement the separation of quotas)

** [ ] Make it an HTTPS proxy
*** Read that medium post about hijacking an HTTP(S?) connection via proxy server written in Go
*** [ ] Make a genuine proxy, have a TCP connection and copy everything back and forth, check it works with curl AND BROWSER
**** [ ] Try http.Server.Server(net.Listener)
**** ! [ ] Read actual source code of http.Server to see how it detects and reads HTTP request (response not needed for now) from the connection
**** [ ] Find HTTP parser inside builtin Go

* Try it in production
** Find an API we use (PB or Channels or whatever) that supports HTTP
** Write a letter to devs
*** It must not become yet another thing to maintain, it must be optimised for fire-and-forget. We (and support) should not know that such thing even exists

* Make the first version presentable at a meetup/conference
** Make mock server stop gracefully (stop it over a quit-channel?)
** Make proxy server stop gracefully (use Server.Shutdown(context) instead of Close)
** Split ListenAndClose so that Listen is called synchronously to avoid for sure that the client is going to connect to a listening port. `server.Serve(listener)` does exist
** Look at implementation style https://github.com/go-httpproxy/httpproxy

* Question/Research: how early will we need HTTP/2 or even /3? Probably it's a long way until it's mandatory, leave it for now
