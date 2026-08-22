# OIDC authentication

This kit expects the API behind it to validate every request using
**OpenID Connect** rather than a home-rolled JWT scheme. This is a
deliberate step up from a hand-built `Authorization: Bearer <jwt>`
check: token issuance, rotation, and validation are delegated to a
dedicated identity provider instead of being reimplemented per project.

## Why not just hand-rolled JWT?

A custom JWT implementation (signing, expiry, refresh) is easy to get
subtly wrong — weak secrets, missing expiry checks, no revocation path.
OIDC pushes that complexity to a provider whose only job is to get it
right, and gives you standard, auditable flows (authorization code,
client credentials, etc.).

## Local setup (no external provider needed)

For local development and testing, run [Keycloak](https://www.keycloak.org/)
in a container — it speaks OIDC natively and needs no cloud account:

```bash
docker run -d --name keycloak \
  -p 8081:8080 \
  -e KEYCLOAK_ADMIN=admin \
  -e KEYCLOAK_ADMIN_PASSWORD=admin \
  quay.io/keycloak/keycloak:latest start-dev
```

Then:
1. Open `http://localhost:8081`, log in with `admin` / `admin`.
2. Create a realm (e.g. `secure-deploy-kit`).
3. Create a client for your API (`app-backend`), with a client secret.
4. Note the discovery document URL:
   `http://localhost:8081/realms/secure-deploy-kit/.well-known/openid-configuration`

## Wiring it into a Spring Boot API

Add the OAuth2 resource server starter and point it at the discovery URL:

```yaml
# application.yml
spring:
  security:
    oauth2:
      resourceserver:
        jwt:
          issuer-uri: http://localhost:8081/realms/secure-deploy-kit
```

Spring Security then validates incoming bearer tokens against the
provider automatically — no manual signature checking needed.

See [`example-config.yaml`](example-config.yaml) for a fuller example,
including role-based authorization mapping.

## Going to production

Swap Keycloak for a managed OIDC provider if you don't want to operate
your own identity server (e.g. Auth0, AWS Cognito, or a self-hosted
Keycloak on real infrastructure). Only the `issuer-uri` changes —
the application code stays the same.
