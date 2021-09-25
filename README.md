# Traefik Cookie Handler Plugin

[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)

A middleware plugin for Traefik that retrieves any `Set-Cookie` headers from a
custom HTTP/HTTPS request and applies them to the `Cookie` header of the request
that is forwarded to the next Traefik middleware or Traefik service.
