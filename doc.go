//Package gowaf is a lightweight MVC framework for building fast, scalable and robust web applications
//
// gowaf is built with focus of composability. This means, you can reap all the
// benefit of using gowaf while maintaining harmony with existing libraries and
// frameworks.
//
// Core feaures of gowaf
//
// * Session management
//
// * Middleware support
//
// * Data acess( Relational database)
//
// * Logging
//
// * Templates (views)
//
// * Configuration
//
// * Static content serving
//
// Why gowaf
//
// There are many frameworks out there, you might be wondering why do we need
// yet another framework?. We probably don't.. gowaf is just a summary of the
// tools, and techniques I use to develop web  applications in Go.
//
// This includes the best libraries, and the best organization of the code base.
// gowaf has one of the very handy Controller( you will see more details in the
// controller section)
//
// These are the common libraries I use
//
// * gorilla/mux: for http routing.
//
// * gorilla/session: for session management
//
// * gorm: for object relationl mapping.
//
// In all of my web application I noticed that, keeping global state brings a
// lot of complexities to the application, so uttron avoids this. gowaf
// applicatio is a simple struct that you can play around with in whatever way
// that you find suits you.
package gowaf
