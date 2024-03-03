A simple and straight-forward TLS proxy, that is planned to work as a MITM and break the TLS session.
It can be used for:
  1. Monitoring traffic (to and from the server)
  2. Modifing traffic (to and from the server)

To configure the TLS one should provide a x.509 certificate (instructions are embeeded in the code as comments) and to change the upstream TLS server (hardcoded, but very easy to change)
Defulat port is 4443, but can be easily changed if needed

There is a simple example of traffic modifcation in the code, it can be easily extended to what ever data you would like to change and even other formats (defaul encoding is utf8, but can be extended to whatever..)

