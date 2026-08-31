# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| 1.x     | ✅ Active support |
| < 1.0   | ❌ No longer supported |

## Reporting a Vulnerability

The Somidax security team takes all vulnerabilities seriously.

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, report them privately via one of these channels:

- **Email:** security@somidax.net
- **GitHub Private Advisory:** [Report a vulnerability](https://github.com/somidaxAI/go-sdk/security/advisories/new)

### What to include in your report

Please provide as much of the following as possible:

- Type of vulnerability (e.g. authentication bypass, injection, data exposure)
- File path(s) and line number(s) of the affected source code
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact assessment — what an attacker could achieve

### What to expect

| Timeline | Action |
|----------|--------|
| Within 48 hours | Acknowledgement of your report |
| Within 7 days | Initial assessment and severity rating |
| Within 30 days | Patch developed and tested |
| Within 45 days | Patch released and advisory published |

We follow [Coordinated Vulnerability Disclosure](https://en.wikipedia.org/wiki/Coordinated_vulnerability_disclosure). We will credit researchers in the release notes unless you prefer to remain anonymous.

## Security Best Practices for SDK Users

- **Never hardcode API keys** — use environment variables: `os.Getenv("SOMIDAX_API_KEY")`
- **Always verify webhook signatures** using `somidax.VerifySignature()` before processing events
- **Use HTTPS only** — never override `WithBaseURL` with an HTTP endpoint in production
- **Rotate API keys** immediately if you suspect compromise — contact support@somidax.net
- **Keep the SDK updated** — run `go get github.com/somidaxAI/go-sdk@latest` regularly

## Token Contract Addresses

If you discover a vulnerability related to the $SMDX token contracts, report to security@somidax.net immediately.

| Chain | Contract |
|-------|----------|
| Ethereum ERC-20 | `0x7e8539D1E5cB91d63E46B8e188403b3f262a949B` |
| BNB Chain BEP-20 | `0xea8c5b9c537f3ebbcc8f2df0573f2d084e9e2bdb` |

## Scope

**In scope:**
- Authentication and authorisation flaws in the SDK
- Webhook signature bypass vulnerabilities
- Data leakage through API responses
- Dependency vulnerabilities

**Out of scope:**
- Vulnerabilities in third-party dependencies (report to the dependency maintainer)
- Issues in the Somidax web application (report to security@somidax.net separately)
- Social engineering attacks

---

Thank you for helping keep Somidax and its users safe.
