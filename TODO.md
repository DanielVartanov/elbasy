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
**** [V] Make it local-usage-ready
***** [V] Make the main func actually run proxy server
***** [V] Make proxy server dump requests and responses in a log
***** [V] Try it with an actual curl (or even browser?)

** Refactor it
*** Make tests more uniform, apply common patterns etc (do we need LastRequest if we have a function?)
*** Make proxy & mock server startup process split into synchronous port binding and asyncrhonous connection handling
*** [V] log.Fatal calls os.Exit(1) itself
*** https://github.com/golang/go/wiki/CodeReviewComments
*** Actually learn `testing` package, there are MANY useful functions
*** Read `log` package, what does it provide?
*** errcheck is a program for checking for unchecked errors in go programs. https://github.com/kisielk/errcheck/

*** [V] Implement leaky-bucket algorithm
**** [V] Make requests wait for quota using a buffered channel and a goroutine which replenishes quota every N milliseconds
**** Make proxy count quotas depending on... (is the quota per shop or per client ID? Well, does not matter, implement the separation of quotas)

** [V] Make it an HTTPS proxy
*** [V] Read that medium post about hijacking an HTTP(S?) connection via proxy server written in Go
*** [V] Make a genuine proxy, have a TCP connection and copy everything back and forth,
**** [V] Check it works with curl
**** [V] Read actual code of bufio.Reader and Writer, how are they different from the buffer I made myself?
**** [X] bufio.Reader.WriteTo should actually work
**** [X] Convert everything to Buffered R/W
**** [V] Check it works with browser
***** [V] Make it work with any Host
**** [V] Read actual source code of http.Server to see how it detects and reads HTTP request (response not needed for now) from the connection
**** [V] Find HTTP parser inside builtin Go

** [V] Impersonate remote server to the client over TLS connection
*** [V] Learn how TLS works
*** [V] Look at TLS implementation code in current Go stdlib
*** [V] Make the simplest TLS connection get established
**** [V] Test it with `curl --insecure`
*** [V] Generate elbasyCertificate upon the certificate received from the remote server https://github.com/FiloSottile/mkcert
**** [X] Only public key is to be changes (possibly an algorithm too). Make sure it is clear for a human that the certificate is forged
**** Pre-forge all the certs for supported APIs in advance? Or change them regularly?
*** [V] Try to make it so that you only give the TLS-related code the cert and don't do any other crypto yourself
*** [X] Test TLS server with `openssl` command line tool
*** [V] Second-order server should be one per host. You set them up in advance (on startup), make them load the certs and their listeners' Accept() functions only do `return <-connectionsChannnnel`, where the channel gets populated by the proxy server as soon as we hijack a connection

* Make it minimally production ready
** [V] Generate Shopify certificate
** [V] Find a proper api key and try to make a request, look at quota headers, see what happens when quota is exceeded
** [V] Before you forget how to do it, write an instruction on how to generate and install a CA certificate and how to generate elbasy certificates for Shopify
** [ ] Make the proxy to be a genuine proxy if a request does not match a certicate url
** [ ] Make sure clients will *not* go via proxy if they send a http request
** [ ] Find out whether Shopify quota is per client or per store
*** [ ] Make the throttler act upon that
** [V] Implement the model
** Write automated tests for that
** (read/find out/ask/google first what is the best practice) Rescue from errors more gracefully (definitely don't log.Fatal/os.Exit every time)
** Write a letter to devs
*** It must not become yet another thing to maintain, it must be optimised for fire-and-forget. We (and support) should not know that such thing even exists
** [ ] Make the implementation clean and responsible (so that it is not scary to restart it, i.e. it does not break exising connections etc)
*** Act responsibly on server shutdown
// Shutdown does not attempt to close nor wait for hijacked
// connections such as WebSockets. The caller of Shutdown should
// separately notify such long-lived connections of shutdown and wait
// for them to close, if desired. See RegisterOnShutdown for a way to
// register shutdown notification functions.
*** Close clientConn and serverConn
// After hijacking it becomes the caller's responsibility to manage
// and close the connection.
// When do we do it? When the request is over or when server is shutdown?
// Or we have to read and respect Keep-Alive?
** [ ] Make it not fail on error

* Make the first version presentable at a meetup/conference
** Wait, client should connect to the proxy via its own TLS as well, shouldn't it? Implement it!
** Make mock server stop gracefully (stop it over a quit-channel?)
** Make proxy server stop gracefully (use Server.Shutdown(context) instead of Close)
** Split ListenAndClose so that Listen is called synchronously to avoid for sure that the client is going to connect to a listening port. `server.Serve(listener)` does exist
** Look at implementation style https://github.com/go-httpproxy/httpproxy
** Look at golang github wiki for coding style guide (sort of, that one))
** Ask _around_ what are the best practices to make a software product for sysops: CMD tools best practices, containers, instructions support etc?


* Question/Research: how early will we need HTTP/2 or even /3? Probably it's a long way until it's mandatory, leave it for now
