// Package webteleport is a client library for creating webteleport connections that is easy to use:
//
//   - Call [Listen] the same way you use [net.Listen] to get a [Listener]
//   - Call [Listener.Accept] to create new [net.Conn]
//
// The URL you pass to [Listen] should be a [WebTeleport Server](github.com/webteleport/server)
//
// With webteleport, you can easily serve on a public address, even if you are behind a firewall.
//
// For possible use cases, check https://github.com/webteleport/ufo
package webteleport
