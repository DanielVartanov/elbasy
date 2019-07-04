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
**** https://github.com/golang/go/wiki/CodeReviewComments
**** Actually learn `testing` package, there are MANY useful functions
**** Read `log` package, what does it provide?
**** errcheck is a program for checking for unchecked errors in go programs. https://github.com/kisielk/errcheck/

// An idea:
        // ConnState specifies an optional callback function that is
        // called when a client connection changes state. See the
        // ConnState type and associated constants for details.
        http.Server.ConnState func(net.Conn, ConnState)

// bufio.Scanner provides a convenient interface for reading data such as
// a file of newline-delimited lines of text. Successive calls to
// the Scan method will step through the 'tokens' of a file, skipping
// the bytes between the tokens.

*** [V] Implement leaky-bucket algorithm
**** [V] Make requests wait for quota using a buffered channel and a goroutine which replenishes quota every N milliseconds
**** Make proxy count quotas depending on... (is the quota per shop or per client ID? Well, does not matter, implement the separation of quotas)

** [ ] Make it an HTTPS proxy
*** [V] Read that medium post about hijacking an HTTP(S?) connection via proxy server written in Go
*** [V] Make a genuine proxy, have a TCP connection and copy everything back and forth,
**** [V] Check it works with curl
**** [V] Read actual code of bufio.Reader and Writer, how are they different from the buffer I made myself?
Conclusion
! Buffered writer is bad (or you'll have implement a version of Copy/WriteTo manually and do Flush every time yourself. Buffer is not flushed if not pushed by more data to the same Writer)
! io.Copy seems to be working (but with unbuffered readers/writers!)
We are reading in large chunks anyway (not in tiny ones, my god, io.Copy buffer is 32kilobytes!) so buffers won't actually reduce amount of actual reads and writes to and from a physical device! Buffers are useless here.
perhaps EXACTLY because of this buffering many of our attempts have failed (including WriteTo ones)

! Yes, io.Copy on UNBUFFERED  connection readers/writers will work. It actually implements what we have implemented already, but with all thre precautions.
  Well, at least we know how to debug it. We can just copy its implementation and see in between what's in there. But I believe everything will be fine.
**** [X] bufio.Reader.WriteTo should actually work
**** [X] Convert everything to Buffered R/W
**** [ ] Rescue from errors gracefully
**** [ ] Check it works with browser
***** [ ] Make it work with any Host
**** [V] Read actual source code of http.Server to see how it detects and reads HTTP request (response not needed for now) from the connection
**** [V] Find HTTP parser inside builtin Go

* Make the implementation clean and responsible
** Act responsibly on server shutdown
// Shutdown does not attempt to close nor wait for hijacked
// connections such as WebSockets. The caller of Shutdown should
// separately notify such long-lived connections of shutdown and wait
// for them to close, if desired. See RegisterOnShutdown for a way to
// register shutdown notification functions.
** Close clientConn and serverConn
// After hijacking it becomes the caller's responsibility to manage
// and close the connection.
// When do we do it? When the request is over or when server is shutdown?
// Or we have to read and respect Keep-Alive?

* Try it in production
** Find an API we use (PB or Channels or whatever) that supports HTTP
** (read/find out/ask/google first what is the best practice) Rescue from errors more gracefully (definitely don't log.Fatal/os.Exit every time)
** Write a letter to devs
*** It must not become yet another thing to maintain, it must be optimised for fire-and-forget. We (and support) should not know that such thing even exists

* Make the first version presentable at a meetup/conference
** Make mock server stop gracefully (stop it over a quit-channel?)
** Make proxy server stop gracefully (use Server.Shutdown(context) instead of Close)
** Split ListenAndClose so that Listen is called synchronously to avoid for sure that the client is going to connect to a listening port. `server.Serve(listener)` does exist
** Look at implementation style https://github.com/go-httpproxy/httpproxy
** Look at golang github wiki for coding style guide (sort of, that one))


* Question/Research: how early will we need HTTP/2 or even /3? Probably it's a long way until it's mandatory, leave it for now
