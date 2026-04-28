# Google Cloud — Identity Management & Authentication Notes

> Source: [Google Cloud IAM Documentation](https://docs.cloud.google.com/iam/docs/overview) | [Authentication Documentation](https://docs.cloud.google.com/docs/authentication) | [Token Types](https://docs.cloud.google.com/docs/authentication/token-types)

---

## Part 1 — Authentication Basics

### The Core Distinction: Credential ≠ Token

- **Credentials**: a digital object that proves your identity — password, PIN, biometric data, service account key.
- **Tokens**: a digital object that proves you *already provided* valid credentials. Tokens are not credentials.

You present your credentials **once**, receive a token, and use that token for every subsequent API call.

### Principal Types and Their Credentials

| Principal type | Who it represents | How it authenticates |
|----------------|-------------------|----------------------|
| **User account** | Human (developer, admin) | Password + MFA → OAuth 2.0 flow |
| **Service account** | Program / workload | Short-lived credentials or JWT assertion |
| **Federated identity (workforce)** | Human from external IdP | OIDC / SAML 2.0 → STS token exchange |
| **Federated identity (workload)** | External service (AWS, GKE, GitHub…) | OIDC / SAML 2.0 → STS token exchange |

### Application Default Credentials (ADC)

ADC is a strategy used by Google auth libraries to automatically find credentials based on the application environment. It checks known locations in the environment at runtime. This means:

- In **local dev**: ADC uses your user credentials (`gcloud auth application-default login`)
- On a **Compute Engine VM**: ADC uses the attached service account automatically
- Same code, different credentials per environment — no code change needed

> ADC cannot be used to authenticate the gcloud CLI itself, or when using API keys.

---

## Part 2 — Identity Management

### 2.1 — Workforce Identity Federation (humans from external IdP)

Workforce Identity Federation lets users in your external identity provider (IdP) — such as Microsoft Entra ID, Okta, or Active Directory — authenticate to Google Cloud using **single sign-on (SSO)**, without needing a Google Account.

Key properties:
- **Syncless**: no user accounts are stored in Google Cloud. Identities are federated, not synchronized.
- **Protocols supported**: OpenID Connect (OIDC) and SAML 2.0.
- **Token exchange**: follows OAuth 2.0 Token Exchange spec (RFC 8693). You present a credential from your IdP to the **Security Token Service (STS)**, which verifies it and returns a short-lived Google Cloud access token.

#### Workforce Identity Pools

A **workforce identity pool** is the logical container for a group of federated user identities (e.g. `employees`, `partners`). Pools are configured at the **organization level** and let you:
- Group user identities and grant IAM access to an entire pool or a subset
- Federate from one or more IdPs
- Define attribute mappings and conditions

#### Attribute Mapping (technical detail)

After user authentication, the IdP sends attributes (also called *claims* in OIDC or *assertions* in SAML). These are mapped to Google Cloud attributes using Common Expression Language (CEL):

| Google attribute | Required? | Description |
|------------------|-----------|-------------|
| `google.subject` | **Yes** | Unique identifier for the user — often the JWT `sub` claim. Used as the principal in Cloud Audit Logs. Max 127 bytes. |
| `google.groups` | No | Array of groups the user belongs to. Used in IAM policies. Limit: 400 groups max per user. |
| `google.display_name` | No | Display name in the Google Cloud console. Cannot be used in IAM policies. |
| `google.email` | No | Email address mapping for OAuth client integration. Cannot be used in IAM policies. |
| `attribute.KEY` | No | Custom attributes from your IdP (e.g. `attribute.costcenter = "1234"`). Up to 50 custom mappings. |

#### Principal identifiers for IAM policies

| Scope | Identifier format |
|-------|-------------------|
| Single user | `principal://iam.googleapis.com/locations/global/workforcePools/POOL_ID/subject/SUBJECT` |
| All users in a group | `principalSet://iam.googleapis.com/locations/global/workforcePools/POOL_ID/group/GROUP_ID` |
| All users with an attribute | `principalSet://iam.googleapis.com/locations/global/workforcePools/POOL_ID/attribute.KEY/VALUE` |
| All users in a pool | `principalSet://iam.googleapis.com/locations/global/workforcePools/POOL_ID/*` |

### 2.2 — Workload Identity Federation (external services)

Workload Identity Federation lets workloads running **outside** Google Cloud (AWS, Azure, on-premises, GitHub Actions, Kubernetes…) access Google Cloud resources without service account keys.

Same mechanism as Workforce: the external workload presents a token from its own identity provider (e.g. an AWS `GetCallerIdentity` token, a GitHub Actions OIDC token, or a Kubernetes service account JWT), exchanges it with the STS for a short-lived Google Cloud access token.

Supported external token types:
- External JWT (OIDC) — GitHub, Okta, Kubernetes, etc.
- External SAML assertion — Entra ID, AD FS, Okta
- AWS `GetCallerIdentity` token — for AWS workloads
- X.509 certificates

### 2.3 — Google Groups

Google Groups are a principal type in IAM that enables access management at scale.

- Every member of a Google Group **inherits all IAM roles granted to that group**, regardless of their Google Groups role.
- This means you manage access by managing group membership, not by assigning individual IAM bindings.
- Groups are managed by Google Workspace, not IAM directly.

> **Warning**: Deleting a group is irreversible. Always revoke all IAM roles from the group first, then wait at least **7 days** before deletion.

---

## Part 3 — Tokens: Technical Deep Dive

### 3.1 — Three Token Categories

| Category | Path | Purpose |
|----------|------|---------|
| **Access tokens** | Auth server → Client → Google API | Call Google Cloud APIs |
| **Token-granting tokens** | Auth server ↔ Client | Obtain new or different tokens |
| **Identity tokens** | Auth server → Client | Identify the user (not for API calls) |

Access tokens and identity tokens are **bearer tokens**: whoever holds the token can use it. If intercepted over an unencrypted channel, they can be exploited. Always use HTTPS.

---

### 3.2 — Access Tokens

All access tokens share:
- Short-lived (max a few hours)
- Scoped to specific OAuth scopes, endpoints, or resources
- Issued to a specific client

#### User access token

- **Format**: Opaque (proprietary, not directly readable)
- **Lifetime**: 1 hour
- **Revocable**: Yes
- **Introspectable**: Yes — via `https://oauth2.googleapis.com/tokeninfo?access_token=TOKEN`

Introspection response example:
```json
{
  "azp": "0000000000.apps.googleusercontent.com",
  "aud": "0000000000.apps.googleusercontent.com",
  "sub": "00000000000000000000",
  "scope": "openid https://www.googleapis.com/auth/userinfo.email",
  "exp": "1744687132",
  "expires_in": "3568",
  "email": "user@example.com",
  "email_verified": "true"
}
```

| Field | Name | Description |
|-------|------|-------------|
| `sub` | Subject | Unique ID of the authenticated principal |
| `aud` | Audience | The OAuth client this token is for |
| `azp` | Authorized party | The OAuth client that requested the token |
| `scope` | OAuth scopes | Set of APIs the client can access |
| `exp` | Expiry | Unix epoch timestamp |

#### Service account access token

- **Format**: Opaque
- **Lifetime**: Default 1h, configurable from 5 min to 12h via `serviceAccounts.generateAccessToken`
- **Revocable**: **No** — stays valid until expiry
- **Introspectable**: Yes — same `tokeninfo` endpoint

#### Service account JWT (self-signed)

This token is issued by the client itself — no auth server needed. Useful when building custom client libraries or calling Google APIs directly.

Structure (three Base64URL-encoded parts separated by `.`):
```
HEADER.PAYLOAD.SIGNATURE
```

Decoded example:
```json
// Header
{
  "alg": "RS256",
  "kid": "290b7bf588eee0c35d02bf1164f4336229373300",
  "typ": "JWT"
}
// Payload
{
  "iss": "service-account@example.iam.gserviceaccount.com",
  "sub": "service-account@example.iam.gserviceaccount.com",
  "scope": "https://www.googleapis.com/auth/cloud-platform",
  "exp": 1744851267,
  "iat": 1744850967
}
// SIGNATURE (RS256 = RSA + SHA-256, signed with the service account private key)
```

| Field | Description |
|-------|-------------|
| `alg` | Signing algorithm — RS256 (RSA + SHA-256) |
| `kid` | Key ID — identifies which service account key was used |
| `iss` | Issuer — the service account itself |
| `sub` | Subject — the service account itself |
| `scope` | OAuth scopes (or `aud` for a specific API endpoint) |
| `iat` | Issued at — Unix epoch |
| `exp` | Expiry — Unix epoch. Max 1 hour from `iat`. |

- **Lifetime**: 5 min to 1 hour
- **Revocable**: No
- **Format**: JWT (decodable by the client)

#### Federated access token

Issued by the Google Cloud IAM authorization server after a token exchange via STS (RFC 8693). Authenticates a workforce pool principal or workload pool principal.

- **Format**: Opaque
- **Revocable**: No
- **Lifetime**: Workforce → min(session time remaining, 1h). Workload → matches the external token's expiry.

#### Access token summary table

| Token type | Format | Lifetime | Revocable |
|------------|--------|----------|-----------|
| User access token | Opaque | 1 hour | Yes |
| Service account access token | Opaque | 5 min – 12 hours | No |
| Service account JWT | JWT | 5 min – 1 hour | No |
| Federated access token | Opaque | ≤ 1 hour | No |
| Domain-wide delegation token | Opaque | 1 hour | No |

---

### 3.3 — Token-Granting Tokens

These tokens allow obtaining new or different tokens. They are **not** usable to call Google APIs directly.

#### Authorization code

- Generated by the auth server after the user authenticates
- Short-lived: **10 minutes**
- **Single use**: can only be exchanged once
- Opaque

Used as an intermediary in the OAuth 2.0 Authorization Code Flow: the client exchanges it for an access token + refresh token by sending `code + client_id + client_secret` to the token endpoint.

#### Refresh token

- Obtained alongside the access token after the authorization code exchange
- **Long-lived**: remains valid until the user revokes authorization or the Google Cloud session expires
- **Multi-use**: can be used repeatedly to obtain new access tokens
- **Revocable**: Yes
- Tied to a specific client (requires `client_id + client_secret` to use)
- Never sent to the API — stays between the client and the auth server

#### Service account JWT assertion

- Signed by the client itself using the service account private key
- Used to obtain a service account access token or a domain-wide delegation token
- The `aud` field must be set to `https://oauth2.googleapis.com/token`
- Lifetime: up to 1 hour

#### External JWT / SAML assertion

Issued by an external IdP (Entra ID, Okta, GitHub…). Used as a token-granting token in Workforce/Workload Identity Federation: exchanged with the STS for a federated access token.

---

### 3.4 — Identity Tokens (ID Tokens)

Identity tokens identify *who* is making a request. They are **not** for calling Google APIs. Used for:
- Service-to-service authentication (e.g. Cloud Run calling another Cloud Run)
- Authenticating to Identity-Aware Proxy (IAP)

All identity tokens are:
- Formatted as **JWTs** (decodable)
- Short-lived (max 1 hour)
- Not revocable

#### User ID token (OIDC)

Decoded example:
```json
{
  "iss": "https://accounts.google.com",
  "azp": "1234567890-abc.apps.googleusercontent.com",
  "aud": "1234567890-abc.apps.googleusercontent.com",
  "sub": "12345678901234567890",
  "email": "user@example.com",
  "hd": "example.com",
  "iat": 1745361695,
  "exp": 1745365295
}
```

| Field | Description |
|-------|-------------|
| `iss` | Always `https://accounts.google.com` |
| `sub` | Unique ID of the authenticated user |
| `aud` | The OAuth client this token is intended for |
| `hd` | Hosted domain — present only for managed user accounts (Google Workspace / Cloud Identity) |
| `exp` / `iat` | Expiry / issued at, in Unix epoch time |

#### Identity token summary table

| Token type | Lifetime | Signed by |
|------------|----------|-----------|
| User ID token | 1 hour | Google JWKS |
| Service account ID token | 1 hour | Google JWKS |
| IAP assertion | 10 minutes | IAP JWKS |
| SAML assertion | 10 minutes | Cloud Identity key |

---

### 3.5 — The Complete OAuth 2.0 Authorization Code Flow

```
User / App                    Auth Server                    Google API
    |                               |                              |
    |-- credentials (pwd + MFA) --> |                              |
    |                               |                              |
    | <--- authorization_code ------|  (opaque, 10 min, 1 use)    |
    |                               |                              |
    |-- POST /token                 |                              |
    |   {code, client_id, secret} ->|                              |
    |                               |                              |
    | <--- access_token (1h) -------|                              |
    | <--- refresh_token (long) ----|                              |
    |                               |                              |
    |--- Authorization: Bearer access_token ------------------->  |
    | <------------------------------- 200 OK -------------------- |
    |                               |                              |
    |  [ access_token expires ]     |                              |
    |                               |                              |
    |--- Authorization: Bearer expired_token ------------------->  |
    | <------------------------------- 401 Unauthorized ----------- |
    |                               |                              |
    |-- POST /token                 |                              |
    |   {refresh_token, client_id, secret} -> |                   |
    |                               |                              |
    | <--- new access_token (1h) ---|                              |
    |                               |                              |
    |--- Authorization: Bearer new_access_token ---------------->  |
    | <------------------------------- 200 OK -------------------- |
```

### OAuth 2.0 scopes

You can use a **global scope** that authorizes access to all Google Cloud services:
```
https://www.googleapis.com/auth/cloud-platform
```

Or a **limited scope** for a specific service (reduces risk if the token is compromised):
```
https://www.googleapis.com/auth/devstorage.read_only
https://www.googleapis.com/auth/bigquery
```

---

## Summary

```
Credentials  ──────────────────────────────────────────────────────────
(password,                                                             |
 SA key,                                                               ▼
 external JWT)   ──► Auth Server / STS ──► access_token  ──► Google API call
                              │                   (1h)          (Bearer header)
                              │
                              └──► refresh_token  ──► new access_token (on 401)
                              │    (long-lived)
                              │
                              └──► identity_token ──► service-to-service auth
                                   (JWT, 1h)               or IAP
```

> Key rule: **credentials prove identity** once. **Tokens carry that proof** for subsequent calls. The refresh token lets you extend the session without re-authenticating from scratch.