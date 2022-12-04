package data

type BrowserConfig struct {
	UserAgent string

	SecChUa, SecChUaPlatform, SecChUaMobile string
	AcceptLanguage                          string
}

func NewConfig(regex ...string) *BrowserConfig {
	userAgent := fetchUserAgent()

	return &BrowserConfig{
		UserAgent: userAgent,

		SecChUa:         parseSecChUa(userAgent),
		SecChUaPlatform: parseSecChUaPlatform(userAgent),
		SecChUaMobile:   parseSecChUaMobile(userAgent),

		AcceptLanguage: fetchAcceptLanguage(),
	}
}
