// Copyright 2012 Jack Pearkes. All rights reserved.
// License information and location goes here.
/*
Package up is a simple web service for monitoring the uptime
of other web services, and optionally sending alerts.

In order to use up, you must interact with it via a client.

A list of clients can be found in the readme located at
http://github.com/pearkes/up-server.

up is designed to be a simple way to keep tabs on your web
applications. It is not expected to be highly reliable, but rather
cheap and simple to set-up. I recommend using up for side projects
and other non-critical applications.

The up-server intends to start without much hassle, so it
tries to infer some information about your environment.

    SECRET

        description:

            a secret key used to authorize clients with your up
            service. up will panic and crash if this is not set.

        example:

            $ export SECRET='my_really_secret_key'

    PORT

        description:

            the port that the up service will listen on.

        example:

            $ export PORT=4321

        default:

            4747

    DATABASE_URL

        description:

            the url of the database you are backing your up
            service with.

        example:

            $ export DATBASE_URL='postgres://jack@localhost/up?sslmode=disable'

        default:

            postgres://$USER@localhost/up?sslmode=disable
*/

package up

// main runs the up service
func main() {
	// Initialize the orm
	initOrm()
	// Initalize the server
	initServer()
}
