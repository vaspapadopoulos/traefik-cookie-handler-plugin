# Traefik Cookie Handler Plugin

[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)

A middleware plugin for [Traefik](https://doc.traefik.io/traefik/) that
retrieves `Set-Cookie` headers from a custom HTTP/HTTPS request and assigns them
to the `Cookie` header of the request which is forwarded to the next Traefik
middleware or Traefik service.

## Configuration

```yaml
testData:
  url: https://a-domain.com/a-path
  method: POST
  responseCookies:
    - COOKIE-A
    - COOKIE-B
```

- `url` the target of the HTTP/HTTPS request
- `method` the HTTP method
- `responseCookies` the variable names of the `Set-Cookie` header values in the
  HTTP response that will be appended to the `Cookie` header of the original
  HTTP request.

Based on the above `testData`, the request, which will be forwarded to the next
Traefik middleware or Traefik service, will contain a `Cookie` header with the
value `COOKIE-A=<x>; COOKIE-B=<y>` where `x` and `y` are the values retrieved by
the HTTP request defined in the plugin's configuration.

## SonarQube use case

The reason for developing this plugin for Traefik was one of the limitations of
SonarQube Community Edition 8.9.2 LTS. The limitation was that
[SonarQube's Web API](https://docs.sonarqube.org/latest/extend/web-api/) was
restricting unauthorized users to hit endpoints, like
`api/project_badges/quality_gate`, and retrieve info for private projects. Here,
it should be mentioned that this make sense since a project is private. However,
there was a need by the SonarQube community to retrieve project badges. In
short, this is the story:

- Community Topic #1:
  [Badges on private projects](https://community.sonarsource.com/t/badges-on-private-projects/4894/46)
- Feature implemented for SonarCloud, but not for SonarQube:
  [\[MMF-1178\] Allow usage of project badges on private projects](https://jira.sonarsource.com/browse/MMF-1178)
- Community Topic #2:
  [Badges for private projects on self hosted SonarQube](https://community.sonarsource.com/t/badges-for-private-projects-on-self-hosted-sonarqube/35783)
- Feature pending for SonarQube:
  [\[MMF-1942\] Allow usage of project badges on private projects](https://jira.sonarsource.com/browse/MMF-1942)
- [Feature planned for SonarQube 9.2](https://portal.productboard.com/sonarsource/3-sonarqube/c/129-project-badges-for-private-projects)

Instead of waiting, I decided to solve this issue by deploying a Traefik
instance in front of my SonarQube instance and developing a Traefik plugin that
implicitly authenticates any requests to the Web API of SonarQube using a
"read-only" user.

### Example configuration for SonarQube

You will need to modify or extend the below snippets according to your setup.
Notice that this config will authenticate to SonarQube via endpoint
`/api/authentication/login` using the relevant query parameters as defined by
SonarQube's Web API (Community Edition 8.9.2 LTS).

- Static configuration

```yaml
pilot:
  token: <your-token>

experimental:
  plugins:
    cookie-handler:
      moduleName: github.com/vaspapadopoulos/traefik-cookie-handler-plugin
      version: v0.1.0
```

- Dynamic configuration

```yaml
http:
  routers:
    sonarqube-router:
      service: sonarqube
      rule: Host(`my-sonarqube.domain.com`)
    sonarqube-router-badges:
      service: sonarqube
      middlewares:
        - sonarQubeBadgesAuth
      rule: >
        Host(`my-sonarqube.domain.com`) &&
        (Path(`/api/project_badges/measure`) ||
          Path(`/api/project_badges/quality_gate`)
        ) &&
        Method(`GET`)
  middlewares:
    sonarQubeBadgesAuth:
      plugin:
        cookie-handler:
          url: http://x.y.z.w:p/api/authentication/login?login=some_user&password=some_password
          method: POST
          responseCookies:
            - JWT-SESSION
            - XSRF-TOKEN
  services:
    sonarqube:
      loadBalancer:
        servers:
          - url: http://x.y.z.w:p/
```
