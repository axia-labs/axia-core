# Axia - Axiomatic Trust Graph

An axiomatic trust system for AI agents to build, store, and reference verifiable truth claims.

**Axia Graph is an open protocol for sourcing & rendering Trust relationships**
* **It is a toolkit for building and reading distributed Trust Graphs**
* **An ambitious plan to create interoperability between existing and future Trust Networks**
* **Compatible with existing rating schemes (scores, percentages, star ratings, etc)**
* **Open Source (Apache licensed)**

## Core Concepts

- **Axiomatic Claims**: Foundational truth statements that can be cryptographically verified
- **Trust Network**: A directed graph of trust relationships between AI agents
- **Truth Consensus**: Mechanisms for establishing shared truth through agent cooperation

## Usage

### Create Axiomatic Claims

```
axios claim --help

Usage: axios-claim [options]

  Options:

    -h, --help                   output usage information
    --agent <agent>              DID or URL of AI agent making the claim
    --subject <subject>          DID or URL of claim subject
    --axiom <axiom>             Axiomatic statement being claimed
    --confidence <confidence>    Confidence score in range 0..1
    --tags <tag1, tag2>         Categorical tags for the claim
    --method <method>           Verification method used
    --proof <proof>             Cryptographic proof
```

Example usage:

```
axios claim \
  --agent did:ai:00a65b11-593c-4a46-bf64-8b83f3ef698f \
  --subject did:fact:59f269a0-0847-4f00-8c4c-26d84e6714c4 \
  --axiom 'Sky appears blue due to Rayleigh scattering' \
  --confidence 0.99 \
  --tags 'physics, optics' \
  --method AxiomaticVerification2024 \
  --proof L4mEi7eEdTNNFQEWaa7JhUKAbtHdVvByGAqvpJKC53mfiqunjBjw
```

This creates a signed JSON-LD Verifiable Claim in the following format:

```json
{
    "@context": "https://schema.axios.ai/AxiomaticClaim.jsonld",
    "type": "AxiomaticClaim", 
    "issuer": "did:ai:00a65b11-593c-4a46-bf64-8b83f3ef698f",
    "issued": "2024-01-17T10:05:07Z",
    "claim": {
        "@context": "https://schema.axios.ai/",
        "type": "Axiom",
        "subject": "did:fact:59f269a0-0847-4f00-8c4c-26d84e6714c4",
        "agent": "did:ai:00a65b11-593c-4a46-bf64-8b83f3ef698f",
        "tags": "physics, optics",
        "axiomRating": {
            "@context": "https://schema.axios.ai/",
            "type": "Confidence",
            "maxConfidence": 1,
            "minConfidence": 0, 
            "confidenceValue": "0.99",
            "axiom": "Sky appears blue due to Rayleigh scattering"
        }
    },
    "proof": {
        "type": "AxiomaticVerification2024",
        "created": "2024-01-17T10:05:07Z",
        "verifier": {
            "id": "Axiomatic-key:020d79074ef137d4f338c2e6bef2a49c618109eccf1cd01ccc3286634789baef4b"
        },
        "domain": "axios.ai",
        "proofValue": "IEd/NpCGX7cRe4wc1xh3o4X/y37pY4tOdt8WbYnaGw/Gbr2Oz7GqtkbYE8dxfxjFFYCrISPJGbBNFyaiVBAb6bs="
    }
}
```

### Query Truth Network

Traverse and analyze the network of axiomatic claims.

```
axios truth --help

Usage: axios-truth [options]

  Options:

    -h, --help                output usage information
    --observer <DID>         Observer agent's perspective
    --agent <agent>          Filter by claim-making agent
    --subject <subject>      Filter by claim subject
    --tags <tag1, tag2>      Filter by categorical tags
    --depth <levels>         Search depth in trust network
    --min-confidence <value> Minimum confidence threshold
    --max-confidence <value> Maximum confidence threshold 
    --consensus             Generate consensus analysis
    --decay                 Trust decay with network distance
```

Example query:

```
axios truth \
  --subject did:fact:59f269a0-0847-4f00-8c4c-26d84e6714c4 \
  --tags 'physics, optics' \
  --consensus
```

### Twitter Integration

The system processes trust claims from tweets using a standardized command syntax. Each command creates a cryptographically signed axiomatic claim that gets added to the trust network.

#### Command Structure
```
@axia_terminal <command> <parameters> [description]
```

#### Supported Commands

1. `#rug` - Report malicious/fraudulent activity
```bash
@axia_terminal #rug @<subject> $<project> [details]
```
Parameters:
- `@<subject>`: Twitter handle of the reported entity
- `$<project>`: Project identifier (token/protocol)
- `[details]`: Optional context (market cap, amount, etc)

Example:
```bash
@axia_terminal #rug @malicious_dev $token_xyz rugpulled at 50m mc
```

Generated Claim:
```json
{
    "type": "AxiomaticClaim",
    "issuer": "twitter:reporter_id",
    "subject": "project:token_xyz",
    "claim": {
        "type": "SecurityReport",
        "action": "rugpull",
        "target": "@malicious_dev",
        "context": {
            "marketCap": "50m",
            "reportType": "rug"
        }
    },
    "confidence": 0.85,
    "tags": ["security", "rugpull", "fraud"]
}
```

2. `#rate` - Submit reputation rating
```bash
@axia_terminal #rate <score> @<subject> [context]
```
Parameters:
- `<score>`: Rating value (0.0-5.0)
- `@<subject>`: Twitter handle being rated
- `[context]`: Optional context about the interaction

Example:
```bash
@axia_terminal #rate 4.5 @trusted_dev great code audit work
```

Generated Claim:
```json
{
    "type": "AxiomaticClaim",
    "issuer": "twitter:rater_id",
    "subject": "identity:trusted_dev",
    "claim": {
        "type": "ReputationRating",
        "score": 4.5,
        "normalizedScore": 0.9,
        "context": "code audit work"
    },
    "confidence": 0.95,
    "tags": ["reputation", "audit", "positive"]
}
```

#### Confidence Scoring

The system automatically generates confidence scores based on various factors:

1. Command-specific base confidence:
   - `#rug`: 0.85 base (serious allegation)
   - `#rate`: 0.95 base (direct experience)

2. Confidence modifiers:
   - +0.05: Detailed context provided
   - +0.10: Verifiable amounts/figures included
   - -0.15: New/unverified reporter account
   - +0.10: Reporter has high trust score

#### State Machine

Each report goes through a state machine:
```
RECEIVED -> VALIDATED -> PROCESSED -> [VERIFIED|DISPUTED]
```

1. `RECEIVED`: Initial tweet ingestion
2. `VALIDATED`: Command syntax and parameters verified
3. `PROCESSED`: Claim created and added to network
4. `VERIFIED`: Additional sources confirmed claim
5. `DISPUTED`: Conflicting claims exist

#### Implementation

To process tweets:

1. Configure webhook:
```bash
axios server --port 8080 --twitter-webhook /webhook/twitter
```

2. Set up authentication:
```bash
# Required headers
X-Axia-Key: your-secret-key
X-Twitter-Webhook-Token: twitter-verification-token
```

3. Monitor state changes:
```bash
axios monitor --command rug --state verified
```

4. Query trust impact:
```bash
axios trust query --subject @malicious_dev --depth 2
```

#### Database Schema

Reports are stored in the `twitter_reports` table:
```sql
CREATE TABLE twitter_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tweet_id VARCHAR(255) NOT NULL UNIQUE,
    author_id VARCHAR(255) NOT NULL,
    subject_handle VARCHAR(255) NOT NULL,
    action VARCHAR(100) NOT NULL,
    project VARCHAR(100) NOT NULL,
    amount VARCHAR(50),
    claim_id UUID NOT NULL REFERENCES claims(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

#### Error Handling

1. Invalid command format:
```bash
@axia_terminal #rug @xyz  # Missing project
> Error: Invalid command format. Usage: #rug @<subject> $<project> [details]
```

2. Rate limiting:
- Maximum 5 reports per hour per user
- Cooldown period of 24h for #rug commands
- Rate limit headers included in responses

3. Duplicate detection:
- Same subject/project combination within 24h
- Similar text content detection

### IPFS Storage

Store and retrieve trust graphs using IPFS via Tatum:

```bash
# Upload entire trust graph
axios ipfs upload

# Upload filtered subset
axios ipfs upload --filter subject=project:example --filter tag=security

# Retrieve graph from IPFS
axios ipfs get QmX...
```

## Installation & Setup

### Prerequisites
- Go 1.19 or higher
- PostgreSQL 13 or higher
- Tatum API key for IPFS storage

### Installation

1. Clone the repository:
```bash
git clone https://github.com/your-org/axia.git
cd axia
```

2. Install dependencies:
```bash
make dev-deps
```

3. Generate secret key:
```bash
make generate-key
```

4. Set up environment variables:
```bash
# Required
export AXIA_SECRET_KEY=$(cat .env)
export DB_HOST=localhost
export DB_USER=your_user
export DB_PASSWORD=your_password
export DB_NAME=axia_trust

# Optional
export TATUM_API_KEY=your_tatum_api_key  # Required for IPFS storage
```

5. Initialize the database:
```bash
axios migrate
```

## Security

### Authentication

All operations require the `AXIA_SECRET_KEY`. This can be provided in several ways:

1. CLI operations use environment variable:
```bash
export AXIA_SECRET_KEY=your-secret-key
axios claim ...
```

2. HTTP requests require header:
```bash
curl -H "X-Axia-Key: your-secret-key" http://localhost:8080/webhook/twitter
```

3. Programmatic usage requires context:
```go
ctx := context.WithValue(context.Background(), "secret_key", os.Getenv("AXIA_SECRET_KEY"))
claim, err := manager.CreateClaim(ctx, ...)
```

## Database Schema

### Claims
```sql
CREATE TABLE claims (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issuer VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    axiom_text TEXT NOT NULL,
    confidence DECIMAL(4,3) NOT NULL CHECK (confidence >= 0 AND confidence <= 1),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    proof_type VARCHAR(100) NOT NULL,
    proof_value TEXT NOT NULL,
    proof_created_at TIMESTAMP WITH TIME ZONE NOT NULL
);
```

### Tags
```sql
CREATE TABLE claim_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    claim_id UUID NOT NULL REFERENCES claims(id) ON DELETE CASCADE,
    tag VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Trust Network
```sql
CREATE TABLE trust_edges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_node UUID NOT NULL REFERENCES claims(id),
    to_node UUID NOT NULL REFERENCES claims(id),
    weight DECIMAL(4,3) NOT NULL CHECK (weight >= 0 AND weight <= 1),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### IPFS Records
```sql
CREATE TABLE ipfs_records (
    id UUID PRIMARY KEY,
    ipfs_id VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);
```

## Development

### Project Structure
```
.
├── cmd/
│   └── trust/           # CLI implementation
├── internal/
│   ├── actions/         # Core business logic
│   ├── auth/            # Authentication system
│   ├── axiom/           # Axiomatic claim management
│   ├── database/        # PostgreSQL integration
│   ├── graph/           # Trust graph implementation
│   ├── logging/         # Structured logging
│   ├── social/          # Social media integrations
│   │   └── twitter/     # Twitter webhook handler
│   ├── state/          # State machine management
│   └── storage/        # External storage (IPFS)
└── doc/                # Documentation
```

### Adding New Features

1. Create new claims:
```go
claim, err := manager.CreateClaim(ctx,
    "agent-id",
    "subject-id",
    "axiom statement",
    0.95,
    []string{"tag1", "tag2"},
)
```

2. Query the network:
```go
results, err := network.Query(ctx, trust.QueryOptions{
    Subject:       "project:example",
    MinConfidence: 0.8,
    Depth:         3,
    UseConsensus:  true,
})
```

3. Store on IPFS:
```go
ipfsClient := ipfs.NewTatumClient(os.Getenv("TATUM_API_KEY"), logger)
ipfsID, err := ipfsClient.UploadGraph(ctx, graphData)
```

### Testing
```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/axiom/...

# Run with coverage
go test -cover ./...
```

### Linting
```bash
make lint
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License - see the [LICENSE](LICENSE) file for details.

## Support

For support:
- Open an issue in the GitHub repository
- Join our Discord community
- Contact the maintainers at support@axios.ai
