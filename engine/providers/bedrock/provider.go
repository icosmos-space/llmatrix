package bedrock

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/icosmos-space/llmatrix/engine/aws"
	"github.com/icosmos-space/llmatrix/engine/provider"
)

// Provider implements Amazon Bedrock Converse (LiteLLM bedrock/).
type Provider struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	HTTPClient      *http.Client
	// RuntimeBaseURL overrides the bedrock-runtime host (tests / VPC endpoints).
	RuntimeBaseURL string
}

// Option configures Bedrock.
type Option func(*Provider)

// New creates a Bedrock provider. Credentials default to AWS_* env vars.
func New(opts ...Option) *Provider {
	p := &Provider{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// WithRegion sets AWS region (default us-east-1 / AWS_REGION).
func WithRegion(region string) Option {
	return func(p *Provider) {
		if region != "" {
			p.Region = region
		}
	}
}

// WithCredentials sets static AWS credentials.
func WithCredentials(accessKey, secretKey, sessionToken string) Option {
	return func(p *Provider) {
		p.AccessKeyID = accessKey
		p.SecretAccessKey = secretKey
		p.SessionToken = sessionToken
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(p *Provider) { p.HTTPClient = c }
}

// Info implements provider.Provider.
func (p *Provider) Info() provider.Info {
	return provider.Info{
		Name:        "bedrock",
		DisplayName: "Amazon Bedrock",
		Capabilities: provider.Capabilities{
			Chat: true, Streaming: true, Tools: true,
		},
	}
}

func (p *Provider) httpClient() *http.Client {
	if p.HTTPClient != nil {
		return p.HTTPClient
	}
	return &http.Client{Timeout: 120 * time.Second}
}

func (p *Provider) region() string {
	if p.Region != "" {
		return p.Region
	}
	if v := os.Getenv("AWS_REGION"); v != "" {
		return v
	}
	if v := os.Getenv("AWS_DEFAULT_REGION"); v != "" {
		return v
	}
	return "us-east-1"
}

func (p *Provider) credentials() aws.Credentials {
	creds := aws.Credentials{
		AccessKeyID:     p.AccessKeyID,
		SecretAccessKey: p.SecretAccessKey,
		SessionToken:    p.SessionToken,
	}
	if creds.AccessKeyID == "" {
		creds.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	}
	if creds.SecretAccessKey == "" {
		creds.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}
	if creds.SessionToken == "" {
		creds.SessionToken = os.Getenv("AWS_SESSION_TOKEN")
	}
	return creds
}

func (p *Provider) runtimeEndpoint(modelID, action string) string {
	modelID = strings.TrimPrefix(modelID, "bedrock/")
	// model IDs may contain '/' (e.g. anthropic/… on some gateways); encode path segments.
	escaped := strings.ReplaceAll(urlPathEscapeKeepSlashes(modelID), "%2F", "/")
	base := strings.TrimRight(p.RuntimeBaseURL, "/")
	if base == "" {
		base = "https://bedrock-runtime." + p.region() + ".amazonaws.com"
	}
	return base + "/model/" + escaped + "/" + action
}

func urlPathEscapeKeepSlashes(s string) string {
	parts := strings.Split(s, "/")
	for i, part := range parts {
		parts[i] = strings.ReplaceAll(pathEscape(part), "+", "%20")
	}
	return strings.Join(parts, "/")
}

func pathEscape(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == '.' || c == '~' || c == ':' {
			b.WriteByte(c)
			continue
		}
		b.WriteByte('%')
		const hex = "0123456789ABCDEF"
		b.WriteByte(hex[c>>4])
		b.WriteByte(hex[c&15])
	}
	return b.String()
}
