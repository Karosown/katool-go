package web_crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// GenerateRodPageSetupSnippetFromCURL parses a curl command and returns a rod snippet.
func GenerateRodPageSetupSnippetFromCURL(curl string) (string, error) {
	opts := DefaultRodPasteOptions("")
	return GenerateRodPageSetupSnippetFromCURLWithOptions(curl, opts)
}

// GenerateRodPageSetupSnippetFromCURLWithOptions parses a curl command and returns a rod snippet.
func GenerateRodPageSetupSnippetFromCURLWithOptions(curl string, options RodPasteOptions) (string, error) {
	setup, _, err := ParseRodPasteFromCURLWithURL(curl, options)
	if err != nil {
		return "", err
	}
	return setup.ToSnippet("page"), nil
}

// ParseRodPasteFromCURL parses a curl command into a RodPasteResult.
func ParseRodPasteFromCURL(curl string, options RodPasteOptions) (*RodPasteResult, error) {
	setup, _, err := ParseRodPasteFromCURLWithURL(curl, options)
	return setup, err
}

// ParseRodPasteFromCURLWithURL parses a curl command into a RodPasteResult and returns target URL.
func ParseRodPasteFromCURLWithURL(curl string, options RodPasteOptions) (*RodPasteResult, string, error) {
	if strings.TrimSpace(curl) == "" {
		return nil, "", errors.New("curl is empty")
	}
	args := splitShellArgs(curl)
	if len(args) == 0 {
		return nil, "", errors.New("curl has no args")
	}
	if args[0] == "curl" {
		args = args[1:]
	}
	var headers []string
	var cookies []string
	var url string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "" {
			continue
		}
		switch arg {
		case "-H", "--header":
			if i+1 < len(args) {
				i++
				headers = append(headers, args[i])
			}
		case "-b", "--cookie":
			if i+1 < len(args) {
				i++
				cookies = append(cookies, args[i])
			}
		case "-A", "--user-agent":
			if i+1 < len(args) {
				i++
				headers = append(headers, "user-agent: "+args[i])
			}
		case "-e", "--referer", "--referrer":
			if i+1 < len(args) {
				i++
				headers = append(headers, "referer: "+args[i])
			}
		case "--url":
			if i+1 < len(args) {
				i++
				url = args[i]
			}
		case "-X", "--request":
			if i+1 < len(args) {
				i++
			}
		default:
			if strings.HasPrefix(arg, "--header=") {
				headers = append(headers, strings.TrimPrefix(arg, "--header="))
			} else if strings.HasPrefix(arg, "--cookie=") {
				cookies = append(cookies, strings.TrimPrefix(arg, "--cookie="))
			} else if strings.HasPrefix(arg, "--user-agent=") {
				headers = append(headers, "user-agent: "+strings.TrimPrefix(arg, "--user-agent="))
			} else if strings.HasPrefix(arg, "--referer=") || strings.HasPrefix(arg, "--referrer=") {
				val := strings.TrimPrefix(strings.TrimPrefix(arg, "--referer="), "--referrer=")
				headers = append(headers, "referer: "+val)
			} else if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
				if url == "" {
					url = arg
				}
			}
		}
	}

	rawHeaders := strings.Join(headers, "\n")
	rawCookies := strings.Join(cookies, "; ")
	if options.TargetURL == "" {
		options.TargetURL = url
	}
	setup, err := ParseRodPaste(rawHeaders, rawCookies, options)
	if err != nil {
		return nil, "", err
	}
	return setup, options.TargetURL, nil
}

// GenerateRodPageSetupSnippetFromHAR parses a HAR JSON string and returns a rod snippet.
func GenerateRodPageSetupSnippetFromHAR(har string, entryIndex int) (string, error) {
	opts := DefaultRodPasteOptions("")
	return GenerateRodPageSetupSnippetFromHARWithOptions(har, entryIndex, opts)
}

// GenerateRodPageSetupSnippetFromHARWithOptions parses a HAR JSON string and returns a rod snippet.
func GenerateRodPageSetupSnippetFromHARWithOptions(har string, entryIndex int, options RodPasteOptions) (string, error) {
	setup, _, err := ParseRodPasteFromHARWithURL(har, entryIndex, options)
	if err != nil {
		return "", err
	}
	return setup.ToSnippet("page"), nil
}

// ParseRodPasteFromHAR parses a HAR JSON string into a RodPasteResult.
func ParseRodPasteFromHAR(har string, entryIndex int, options RodPasteOptions) (*RodPasteResult, error) {
	setup, _, err := ParseRodPasteFromHARWithURL(har, entryIndex, options)
	return setup, err
}

// ParseRodPasteFromHARWithURL parses a HAR JSON string into a RodPasteResult and returns target URL.
func ParseRodPasteFromHARWithURL(har string, entryIndex int, options RodPasteOptions) (*RodPasteResult, string, error) {
	if strings.TrimSpace(har) == "" {
		return nil, "", errors.New("har is empty")
	}
	var file harFile
	if err := json.Unmarshal([]byte(har), &file); err != nil {
		return nil, "", err
	}
	if len(file.Log.Entries) == 0 {
		return nil, "", errors.New("har has no entries")
	}
	if entryIndex < 0 || entryIndex >= len(file.Log.Entries) {
		return nil, "", fmt.Errorf("entryIndex out of range: %d", entryIndex)
	}
	entry := file.Log.Entries[entryIndex]
	req := entry.Request
	if options.TargetURL == "" {
		options.TargetURL = req.URL
	}

	var headers []string
	for _, h := range req.Headers {
		if strings.TrimSpace(h.Name) == "" {
			continue
		}
		headers = append(headers, h.Name+": "+h.Value)
	}
	rawHeaders := strings.Join(headers, "\n")

	var cookieParts []string
	for _, c := range req.Cookies {
		if strings.TrimSpace(c.Name) == "" {
			continue
		}
		cookieParts = append(cookieParts, c.Name+"="+c.Value)
	}
	rawCookies := strings.Join(cookieParts, "; ")

	setup, err := ParseRodPaste(rawHeaders, rawCookies, options)
	if err != nil {
		return nil, "", err
	}
	return setup, options.TargetURL, nil
}

type harFile struct {
	Log harLog `json:"log"`
}

type harLog struct {
	Entries []harEntry `json:"entries"`
}

type harEntry struct {
	Request harRequest `json:"request"`
}

type harRequest struct {
	Method  string        `json:"method"`
	URL     string        `json:"url"`
	Headers []harNameVal  `json:"headers"`
	Cookies []harCookieNV `json:"cookies"`
}

type harNameVal struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type harCookieNV struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func splitShellArgs(input string) []string {
	var res []string
	var b strings.Builder
	var quote rune
	escaped := false
	flush := func() {
		if b.Len() > 0 {
			res = append(res, b.String())
			b.Reset()
		}
	}
	for _, r := range input {
		if escaped {
			b.WriteRune(r)
			escaped = false
			continue
		}
		switch r {
		case '\\':
			escaped = true
		case '\'', '"':
			if quote == 0 {
				quote = r
			} else if quote == r {
				quote = 0
			} else {
				b.WriteRune(r)
			}
		case ' ', '\t', '\n', '\r':
			if quote != 0 {
				b.WriteRune(r)
			} else {
				flush()
			}
		default:
			b.WriteRune(r)
		}
	}
	flush()
	return res
}
