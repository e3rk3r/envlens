# envlens

> Diff, audit, and sanitize `.env` files across environments — with secret detection and schema validation.

---

## Installation

```bash
go install github.com/yourusername/envlens@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envlens.git && cd envlens && go build ./...
```

---

## Usage

```bash
# Diff two .env files
envlens diff .env.development .env.production

# Audit a .env file for exposed secrets
envlens audit .env

# Validate against a schema file
envlens validate .env --schema .env.schema

# Sanitize a .env file by redacting secret values
envlens sanitize .env --output .env.sanitized
```

**Example output:**

```
[DIFF] Missing in production: DEBUG, LOG_LEVEL
[DIFF] Changed keys: DATABASE_URL, API_BASE_URL
[AUDIT] Potential secret detected: AWS_SECRET_KEY (line 12)
[VALIDATE] Missing required keys: SMTP_HOST, REDIS_URL
```

---

## Features

- 🔍 **Diff** — Compare `.env` files across environments and surface missing or changed keys
- 🔐 **Secret Detection** — Flag values that look like passwords, tokens, or API keys
- ✅ **Schema Validation** — Ensure required keys are present and correctly formatted
- 🧹 **Sanitize** — Redact sensitive values for safe sharing or logging

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you'd like to change.

---

## License

[MIT](LICENSE)