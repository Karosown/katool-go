package web_crawler

import (
	"errors"
	"fmt"
	nurl "net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// RodPasteOptions controls how to generate rod page setup code.
type RodPasteOptions struct {
	TargetURL        string
	ExcludeHeaders   []string
	IncludeUserAgent bool
	IncludeHeaders   bool
	IncludeCookies   bool
}

// DefaultRodPasteOptions returns default options for generating rod page setup code.
func DefaultRodPasteOptions(targetURL string) RodPasteOptions {
	return RodPasteOptions{
		TargetURL:        targetURL,
		ExcludeHeaders:   []string{"cookie", "content-length"},
		IncludeUserAgent: true,
		IncludeHeaders:   true,
		IncludeCookies:   true,
	}
}

// GenerateRodPageSetupSnippet converts raw headers/cookies into rod page setup code.
func GenerateRodPageSetupSnippet(rawHeaders, rawCookies, targetURL string) (string, error) {
	opts := DefaultRodPasteOptions(targetURL)
	return GenerateRodPageSetupSnippetWithOptions(rawHeaders, rawCookies, opts)
}

// GenerateRodPageSetupSnippetWithOptions converts raw headers/cookies into rod page setup code.
func GenerateRodPageSetupSnippetWithOptions(rawHeaders, rawCookies string, options RodPasteOptions) (string, error) {
	setup, err := ParseRodPaste(rawHeaders, rawCookies, options)
	if err != nil {
		return "", err
	}
	return setup.ToSnippet("page"), nil
}

// RodPasteResult holds parsed headers/cookies and can emit code snippets.
type RodPasteResult struct {
	Headers        map[string]string
	Cookies        []*proto.NetworkCookieParam
	UserAgent      string
	AcceptLanguage string
}

// ParseRodPaste parses raw headers/cookies into a reusable structure.
func ParseRodPaste(rawHeaders, rawCookies string, options RodPasteOptions) (*RodPasteResult, error) {
	headers := make(map[string]string)
	cookies := make(map[string]string)
	var ua string
	var acceptLang string

	if rawHeaders != "" {
		parseHeaders(rawHeaders, headers, cookies, &ua, &acceptLang)
	}
	if rawCookies != "" {
		for k, v := range parseCookieString(rawCookies) {
			cookies[k] = v
		}
	}

	for _, key := range options.ExcludeHeaders {
		delete(headers, strings.ToLower(key))
	}

	var cookieParams []*proto.NetworkCookieParam
	if options.IncludeCookies && len(cookies) > 0 {
		if options.TargetURL == "" {
			return nil, errors.New("targetURL is required for cookies")
		}
		if _, err := nurl.Parse(options.TargetURL); err != nil {
			return nil, fmt.Errorf("invalid targetURL: %w", err)
		}
		names := make([]string, 0, len(cookies))
		for name := range cookies {
			if name != "" {
				names = append(names, name)
			}
		}
		sort.Strings(names)
		for _, name := range names {
			cookieParams = append(cookieParams, &proto.NetworkCookieParam{
				Name:  name,
				Value: cookies[name],
				URL:   options.TargetURL,
			})
		}
	}

	res := &RodPasteResult{
		Headers: headers,
		Cookies: cookieParams,
	}
	if options.IncludeUserAgent {
		res.UserAgent = ua
		res.AcceptLanguage = acceptLang
	}
	if !options.IncludeHeaders {
		res.Headers = map[string]string{}
	}
	return res, nil
}

// ToSnippet renders a Go snippet that sets headers/cookies on a rod.Page.
func (r *RodPasteResult) ToSnippet(pageVar string) string {
	var b strings.Builder
	if pageVar == "" {
		pageVar = "page"
	}
	if r.UserAgent != "" {
		b.WriteString(pageVar)
		b.WriteString(".MustSetUserAgent(&proto.NetworkSetUserAgentOverride{")
		b.WriteString("UserAgent: ")
		b.WriteString(strconv.Quote(r.UserAgent))
		if r.AcceptLanguage != "" {
			b.WriteString(", AcceptLanguage: ")
			b.WriteString(strconv.Quote(r.AcceptLanguage))
		}
		b.WriteString("})\n")
	}

	if len(r.Headers) > 0 {
		keys := make([]string, 0, len(r.Headers))
		for k := range r.Headers {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		b.WriteString("cleanup := ")
		b.WriteString(pageVar)
		b.WriteString(".MustSetExtraHeaders(\n")
		for _, k := range keys {
			b.WriteString("\t")
			b.WriteString(strconv.Quote(k))
			b.WriteString(", ")
			b.WriteString(strconv.Quote(r.Headers[k]))
			b.WriteString(",\n")
		}
		b.WriteString(")\n_ = cleanup\n")
	}

	if len(r.Cookies) > 0 {
		b.WriteString(pageVar)
		b.WriteString(".MustSetCookies(\n")
		for _, c := range r.Cookies {
			b.WriteString("\t&proto.NetworkCookieParam{")
			b.WriteString("Name: ")
			b.WriteString(strconv.Quote(c.Name))
			b.WriteString(", Value: ")
			b.WriteString(strconv.Quote(c.Value))
			if c.URL != "" {
				b.WriteString(", URL: ")
				b.WriteString(strconv.Quote(c.URL))
			}
			b.WriteString("},\n")
		}
		b.WriteString(")\n")
	}
	return b.String()
}

// ApplyToPage applies headers/cookies/user-agent to a rod page.
func (r *RodPasteResult) ApplyToPage(page *rod.Page) (cleanup func(), err error) {
	if page == nil {
		return nil, errors.New("page is nil")
	}
	if r == nil {
		return nil, errors.New("rod paste result is nil")
	}
	if r.UserAgent != "" {
		if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent:      r.UserAgent,
			AcceptLanguage: r.AcceptLanguage,
		}); err != nil {
			return nil, err
		}
	}
	if len(r.Headers) > 0 {
		dict := make([]string, 0, len(r.Headers)*2)
		for k, v := range r.Headers {
			dict = append(dict, k, v)
		}
		cleanup, err = page.SetExtraHeaders(dict)
		if err != nil {
			return cleanup, err
		}
	}
	if len(r.Cookies) > 0 {
		if err := page.SetCookies(r.Cookies); err != nil {
			return cleanup, err
		}
	}
	return cleanup, nil
}

// MustApplyToPage applies headers/cookies/user-agent and panics on error.
func (r *RodPasteResult) MustApplyToPage(page *rod.Page) (cleanup func()) {
	cleanup, err := r.ApplyToPage(page)
	if err != nil {
		panic(err)
	}
	return cleanup
}

// RenderFunc returns a render function that applies the rod paste settings.
func (r *RodPasteResult) RenderFunc() func(*rod.Page) {
	return func(page *rod.Page) {
		_ = r.MustApplyToPage(page)
	}
}

func parseHeaders(raw string, headers map[string]string, cookies map[string]string, ua *string, acceptLang *string) {
	lines := splitHeaderLines(raw)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		key, val, ok := splitHeader(line)
		if !ok {
			continue
		}
		lower := strings.ToLower(key)
		switch lower {
		case "cookie":
			for k, v := range parseCookieString(val) {
				cookies[k] = v
			}
		case "user-agent":
			*ua = val
		case "accept-language":
			*acceptLang = val
			headers[lower] = val
		default:
			headers[lower] = val
		}
	}
}

func splitHeaderLines(raw string) []string {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	return strings.Split(raw, "\n")
}

func splitHeader(line string) (key, val string, ok bool) {
	idx := strings.Index(line, ":")
	if idx <= 0 {
		return "", "", false
	}
	key = strings.TrimSpace(line[:idx])
	val = strings.TrimSpace(line[idx+1:])
	return key, val, true
}

func parseCookieString(raw string) map[string]string {
	res := make(map[string]string)
	for _, part := range strings.Split(raw, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		idx := strings.Index(part, "=")
		if idx <= 0 {
			continue
		}
		name := strings.TrimSpace(part[:idx])
		val := strings.TrimSpace(part[idx+1:])
		if name == "" {
			continue
		}
		res[name] = val
	}
	return res
}
