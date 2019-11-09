What to do next
================

* [V] Make the first version runnable locally but fully working
** [V] Make a non-transparent proxy
*** [V] Quickly google how exactly HTTP (without TLS) supports proxying
*** [V] Write the simplest proxy
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

*** [V] Implement leaky-bucket algorithm
**** [V] Make requests wait for quota using a buffered channel and a goroutine which replenishes quota every N milliseconds
**** [V] Make proxy count quotas depending on... (is the quota per shop or per client ID? Well, does not matter, implement the separation of quotas)

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
**** [V] Pre-forge all the certs for supported APIs in advance? Or change them regularly?
*** [V] Try to make it so that you only give the TLS-related code the cert and don't do any other crypto yourself
*** [X] Test TLS server with `openssl` command line tool
*** [V] Second-order server should be one per host. You set them up in advance (on startup), make them load the certs and their listeners' Accept() functions only do `return <-connectionsChannel`, where the channel gets populated by the proxy server as soon as we hijack a connection

* Make it production ready
** [V] Generate Shopify certificate
** [V] Find a proper api key and try to make a request, look at quota headers, see what happens when quota is exceeded
** [V] Before you forget how to do it, write an instruction on how to generate and install a CA certificate and how to generate elbasy certificates for Shopify
** [V] Make the proxy to be a genuine proxy if a request does not match a certicate url
** [V] Make sure clients will *not* go via proxy if they send a http request
** Confirm that Shopify quota is per store, not per client
** [V] Make the throttler apply limits per a store
** [V] Implement the model
** [V] Make it not fail on error
*** [V] (read/find out/ask/google first what is the best practice) Rescue from errors more gracefully (definitely don't log.Fatal/os.Exit every time)
** [V] Make it detect and log quota exceedings
** [V] Check it on a reallt big amount of requests. I have a suspicioun they don't get parallelised
** [V] Present it to fellow devs
*** [V] Prepare a demo with multiple requests to Shopify
** Make it deployable
*** [V] Read the best practices
*** [V] Make it compilable
*** [V] Make a binary
*** [V] Make an automated deployment procedure at Travis
*** Generate and install certificates
**** Figure out for how long a root certificatge and a leaf certificate can last in practice
**** Generate a root certificate with Andrew and leave it with him to store the key very secretly
**** Andrew is to inject the root certificate into the image of every instance (or only to Etsy workers if it is even possible)
** Organise a change of certificates in 1.5 years
** [V] Create an alert for `429 Too Many Requests` in logs
** Think of and make more metrics
*** Difference between calculated quota and quota received in a response Header
** Add support for Etsy.com
*** Test up to its limits
** Cover it with tests
*** Tests for proxying entire requests and responses correctly from client to server and back
**** [V] test when not throttled
***** [V] it copies from client to server
***** [V] it copies from server to client
***** test many requests just go through without delay (meausre time and state it is less than a second, we measure time throttling tests anyway)
Opinion: no need to measure a delay. It is enough to measure time whithin which all requests reach the mockServer and fuzzily compare it to a target time
**** test when throttled
*****  it copies from client to server
*****  it copies from server to client
*****  test throttling per se
*** Tests for closing connections
*** Tests for actual throttling prevention
*** Tests for detecting significant differenence between calculated and actual quota (test it with 7 requests at once)
** [V] Rename ElbasyServer to ImpostorServer, elbasy_certificates to impostor_certificates, entire project to elbasy
** Shutdown gracefully on SIGTERM
** Test it with race detectors https://blog.golang.org/race-detector
*** Include NGinx into the container, make it handle TLS, keep-alive client connections etc. Perhaps it could even work via HTTP/2
** Dockerize it

* Make it fire-and-forgettable
** [V] Update directory structure and switch to go1.13+
** Ignore error "regularConnHandler.handleConnection(): io.Copy(clientConn, serverConn): readfrom tcp 127.0.0.1:8443->127.0.0.1:56470: write tcp 127.0.0.1:8443->127.0.0.1:56470: write: broken pipe"
** Investigate what makes regularproxy dial 0.0.0.0:443
*** If it is not CONNECT, reply with a proper response of redirect to the github page README
** [V] Investigate why elbasy is limited to 1024 file descriptors and fix it
** Measure if file descriptors really leak far beyond 1024
** Fix the file descriptors leak if there is any
sudo lsof | grep elbasy
*** ? Close a connection in case of error in Proxy.handlerFunc()
** Make it easy to re-install from scratch if needed in a year
*** [V] Automate an increase of the file count limit
          # NOTE: Just added `LimitNOFILE=1048576` to the elbasy systemd service definition file

* What to do next
** ! Use HTTP/2 connection to each remote host, so that no matter how many 1000s of incoming request you have all of them are effeciently sent via a single HTTP/2 connection
** There should be a hardcoded list of supported APIs, i.e. those throttlnig properties of which are known and reverse-engineered. A database of API throttling rules.
*** Those properties are: Detection of the API, throttling rules, reading the quota account, reading the current quota, detecting that quota is exceeded
*** Configuration only sets out those which are used in the particular installation
** Print waiting-for-quota times into log and plot a logarithmic graphs from that plot
** Make CI do benchmarks, even if irrelevant to production they will still show sudden or gradual slowdowns
** Write automated tests for that
** Make it not scary to restart, i.e. it does not break exising connections etc
*** React to C-C and another C-C
**** First one causes proxyServer.Shutdown which in its turn causes everything else
*** Use http.Server.Shutdown instead of http.Server.Close everywhere
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
** Log amount of currently-being-handled requests into log every second and (use waitgroup?) and plot a graph of it at monitoring
** When 429 is dicovered (probably due to requests made outside proxy) do an emergency bucket drain
** Implement next throttling models
*** ClaimQuota should return `interface{}` -- a token (even a function if a counter wants so) to identify an exact request
**** I bet it is easier to just accept `action func()` than to deal with tokens
** Make it shutdown gracefully (waiting for the current requests) on receiving TERM signal
** Make it shutdown on `C-c` and `C-c C-c`

** Refactor it
*** Make tests more uniform, apply common patterns etc (do we need LastRequest if we have a function?)
*** Convert all underscore and hyphen based identifiers/names to camel case as the latter apparently _is_ GoLang convention
*** [V] Make proxy & mock server startup process split into synchronous port binding and asyncrhonous connection handling
*** [V] log.Fatal calls os.Exit(1) itself
*** https://github.com/golang/go/wiki/CodeReviewComments
*** Actually learn `testing` package, there are MANY useful functions
https://codesamplez.com/development/golang-unit-testing
*** Read `log` package, what does it provide?
*** errcheck is a program for checking for unchecked errors in go programs. https://github.com/kisielk/errcheck/

* Make it usable by the general public
** Should we make certificates regeneation anyhow automatic?
** https://stackoverflow.com/questions/44929223/why-should-i-use-fork-to-daemonize-my-process/44929497#44929497

* Make the first version presentable at a meetup/conference
** [X] Make Client<->Proxy connection secure as well, it must have a certificate and use ServeTLS
*** Make it optional as Client-Proxy connection in some infras are totally inside a local network
*** Moreover, what to steal there? Server domain names? Everything else is secured anyway
** Write an extensive README/instructions
** Test all the installation instructions in a wild with someone who needs it
** Wait, client should connect to the proxy via its own TLS as well, shouldn't it? Implement it!
** Make mock server stop gracefully (stop it over a quit-channel?)
** Make proxy server stop gracefully (use Server.Shutdown(context) instead of Close)
** Split ListenAndClose so that Listen is called synchronously to avoid for sure that the client is going to connect to a listening port. `server.Serve(listener)` does exist
** Look at implementation style https://github.com/go-httpproxy/httpproxy
** Look at golang github wiki for coding style guide (sort of, that one))
** Ask _around_ what are the best practices to make a software product for sysops: CMD tools best practices, containers, instructions support etc?

* Question/Research: how early will we need HTTP/2 or even /3? Probably it's a long way until it's mandatory, leave it for now
