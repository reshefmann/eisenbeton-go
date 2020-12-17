# eisenbeton-go

**Experimental**

A web server with the following goals:
- Make the transition from sync HTTP calls to async flow inside the system as close as possible to the edge.
- Take care of the use case of high traffic of small, independent messages, so not support for file uploads, static file serving etc.
- Support two modes: 
  1) "Publish" only - Just receive the message, return a constant response to the caller and handle internally.
  2) "Request/Reply" - Receive a message, send asynchronously to a handler and publish the response back to the server so it will send the answer to the caller

The implementation is based on [nats.io](http://nats.io) for the communication between the web server and the handler.

Since nats.io is used for the messaging the system is **polyglot**.

```
                                                        +-------------+
                                                        |             |
                                                 +----->+   Handler   |
                                                 |      |             |
   +--------------+        +-------------+       |      --------------+
   |              |        |             |       |
   | Server       |        |   nats.io   |       |
   | (eisenbeton) +<------>+   (Server   +<------+
   |              |        |   or        +<------+      +-------------+
   +--------------+        |   Cluster)  |       |      |             |
                           |             |       +----->+   Handler   |
                           +-------------+              |             |
                                                        +-------------+



```


### Terminology
- *System* - The whole setup of server, nats and handlers.
- *Server* - The HTTP server that take care of handling the incoming requests.
- *Handler* - The internal handler that performs the actual work. Connect to the system with client lib.
- *nats* - well, the nats.io server/cluster that takes care of the communication.
- *Eisenbeton handler lib* - A library embedded in the handler to communicate with nats.

### Messaging 
For the messages *eisenbeton* uses flatbuffers.

### Client libraries
I'm working also on client libraries in several languages (first will be Clojure) and will publish them here when ready.

### Other implementations
This is the golang implementation. I plan to implement the server in more languages - Clojure (WIP), FSharp and Rust. Why? For comparison and fun.

### The name
Copyrights on the name go to [Shlomi Izikovich](https://github.com/shlomii) :smile:

### Sources of inspiration
[http://mongrel2.org/manual/book-final.html](http://mongrel2.org/manual/book-final.html) - Not sure how much info is still online for this deprecated server but the design is really innovative even now.  
[https://resgate.io/docs/get-started/introduction/#how-does-it-work](https://resgate.io/docs/get-started/introduction/#how-does-it-work) - I like the concept here although they aim for a different use case  
[https://doc.traefik.io/traefik/middlewares/circuitbreaker/](https://doc.traefik.io/traefik/middlewares/circuitbreaker/)
