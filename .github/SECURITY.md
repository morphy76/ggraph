# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| main    | :white_check_mark: |
| < 1.0   | :x:                |

**Note:** As this project is currently in early development (pre-1.0), we recommend always using the latest version from the `main` branch.

## Reporting a Vulnerability

We take the security of ggraph seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### How to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

1. **GitHub Security Advisories** (Preferred)
   - Navigate to the [Security tab](../../security/advisories) of this repository
   - Click "Report a vulnerability"
   - Fill out the form with details about the vulnerability

2. **Direct Contact**
   - Open a private issue and request direct contact with maintainers
   - We will provide a secure communication channel

### What to Include

Please include the following information in your report:

- Type of issue (e.g., buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

### What to Expect

- **Acknowledgment**: We will acknowledge receipt of your vulnerability report within 48 hours
- **Updates**: We will send you regular updates about our progress
- **Timeline**: We aim to:
  - Confirm the vulnerability within 7 days
  - Develop and test a fix within 30 days
  - Release a patched version and public disclosure within 90 days

### Disclosure Policy

- We ask that you give us reasonable time to address the vulnerability before any public disclosure
- We will credit you for discovering the vulnerability (unless you prefer to remain anonymous)
- We will coordinate the disclosure with you once a fix is available

## Security Best Practices for Users

When using ggraph in your projects:

1. **Keep Updated**: Always use the latest version of ggraph
2. **Dependency Management**: Regularly update dependencies using `go get -u` and `go mod tidy`
3. **Input Validation**: Validate all inputs to your graph nodes and agents
4. **Sensitive Data**: Avoid storing sensitive information in graph state without proper encryption
5. **API Keys**: Never commit API keys or secrets to version control
6. **Environment Variables**: Use environment variables or secure vaults for sensitive configuration

## Known Security Considerations

### AI Integration

- **Prompt Injection**: Be cautious of user inputs that may manipulate AI behavior
- **Data Leakage**: Ensure sensitive data is not inadvertently sent to AI providers
- **Rate Limiting**: Implement appropriate rate limiting to prevent abuse

### State Management

- **Persistence**: If using persistent state, ensure proper access controls
- **Serialization**: Be aware of potential deserialization vulnerabilities when loading state

### Dependencies

We use Dependabot to monitor our dependencies for known vulnerabilities. You can view the dependency graph and security alerts in the [Security tab](../../security).

## Security Updates

Security updates will be:

- Released as soon as possible after confirmation
- Announced in the [release notes](../../releases)
- Documented in the CHANGELOG with a `[SECURITY]` prefix

## Comments on This Policy

If you have suggestions on how this process could be improved, please submit a pull request or open an issue.

---

Thank you for helping keep ggraph and its users safe!
