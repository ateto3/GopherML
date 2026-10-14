# 🔒 Security

## Supported Versions

| Version | Supported |
| :--- | :--- |
| Latest version on `main` | ✅ |
| Any older version | ❌ |

This project is distributed without warranties (see [LICENSE](LICENSE)). Only the main branch receives vulnerability fixes.

## Reporting a Vulnerability

Please **do not open a public issue** to report security vulnerabilities.

Use GitHub's [Security Advisories](https://github.com/ateto3/gopherml/security/advisories/new) feature to report privately. If you cannot use it, contact the maintainer directly through GitHub ([@ateto3](https://github.com/ateto3)).

Include in the report:

1. **Description** of the problem and its potential impact.
2. **Steps to reproduce** (pseudocode, snippet or minimal test).
3. **Version / commit** affected.
4. **Known mitigation**, if any.

## What to expect

- **Acknowledgment** within a reasonable time (target: 72 hours).
- **Assessment** of the impact and confirmation of the vulnerability.
- **Fix and release** before any coordinated disclosure.
- Credit in the `CHANGELOG.md` if you wish (or anonymity if you prefer).

## Scope

In scope:

- Data validation (`utils.CheckFinite`) and security issues in input handling.
- Persistence (`persistence.Load`): file corruption, model integrity validation.
- Any memory leak or race condition detectable with `go test -race`.

Out of scope:

- Performance issues without security impact.
- Vulnerabilities in exchanges or third-party APIs used alongside GopherML.
- `go-nlp/nltk/`: reference code imported from NLTK (Python) for algorithmic comparison; it is not part of the Go runtime.
