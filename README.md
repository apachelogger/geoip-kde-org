# Description

GeoIP service based on Maxmind's GeoLite2 data. Supports multiple API endpoints
for different output formats for different installers. The underlying GeoLite2
data is automatically updated on every start iff the existing data were last
updated more than a week ago. The service automatically terminates after 14 days
of uptime (thus triggering the aforementioned update). Because of this the
listening sockets are managed through systemd so no connections are lost during
this restart dance.

# Requirements

Needs Go 1.8

# Deployment

Database it automatically downloaded and managed in PWD, so PWD should be suitable.
systemd/* contains example socket and service.

# Documentation

Documentation uses apidocjs.com. Run `make doc` to generate it (requires npm).
