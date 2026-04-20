# Security Policy

## Supported Versions

This project is pre-alpha (0.x). Security fixes are applied to the latest version only.

## Reporting a Vulnerability

Please do **not** open a public GitHub issue for security vulnerabilities.

Instead, email `mahimairaja@mahimai.ai` with:
- A description of the vulnerability
- Steps to reproduce
- Potential impact
- Your suggested fix (if any)

You will receive a response within 72 hours. If the vulnerability is confirmed, a fix will be released as soon as possible and you will be credited in the changelog (unless you prefer to remain anonymous).

## Scope

This project benchmarks self-hosted AI memory systems. The primary security concerns are:

1. **Config file injection**: Run configs may reference environment variables. Never commit secrets in config files.
2. **Dataset integrity**: Dataset loaders verify SHA-256 hashes. Do not override the expected hash constant.
3. **Arbitrary code execution**: Do not run benchmark configs from untrusted sources.
