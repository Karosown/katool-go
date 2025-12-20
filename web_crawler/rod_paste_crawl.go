package web_crawler

import (
	"errors"

	"github.com/go-rod/rod"
)

// ReadSourceCodeWithCURL reads source code using headers/cookies from a curl command.
func ReadSourceCodeWithCURL(curl string, renderFunc func(*rod.Page)) SourceCode {
	return DefaultClient.ReadSourceCodeWithCURL(curl, renderFunc)
}

// ReadSourceCodeWithCURL reads source code using headers/cookies from a curl command.
func (c *Client) ReadSourceCodeWithCURL(curl string, renderFunc func(*rod.Page)) SourceCode {
	setup, url, err := ParseRodPasteFromCURLWithURL(curl, DefaultRodPasteOptions(""))
	if err != nil {
		return SourceCode{WebReaderString: "", WebReaderError: WebReaderError{err}}
	}
	if url == "" {
		return SourceCode{WebReaderString: "", WebReaderError: WebReaderError{errors.New("curl has no url")}}
	}
	return c.readSourceCodeWithRodPaste(url, setup, renderFunc)
}

// ReadSourceCodeWithHAR reads source code using headers/cookies from a HAR entry.
func ReadSourceCodeWithHAR(har string, entryIndex int, renderFunc func(*rod.Page)) SourceCode {
	return DefaultClient.ReadSourceCodeWithHAR(har, entryIndex, renderFunc)
}

// ReadSourceCodeWithHAR reads source code using headers/cookies from a HAR entry.
func (c *Client) ReadSourceCodeWithHAR(har string, entryIndex int, renderFunc func(*rod.Page)) SourceCode {
	setup, url, err := ParseRodPasteFromHARWithURL(har, entryIndex, DefaultRodPasteOptions(""))
	if err != nil {
		return SourceCode{WebReaderString: "", WebReaderError: WebReaderError{err}}
	}
	if url == "" {
		return SourceCode{WebReaderString: "", WebReaderError: WebReaderError{errors.New("har entry has no url")}}
	}
	return c.readSourceCodeWithRodPaste(url, setup, renderFunc)
}

func (c *Client) readSourceCodeWithRodPaste(url string, setup *RodPasteResult, renderFunc func(*rod.Page)) SourceCode {
	return c.ReadSourceCode(url, "", func(page *rod.Page) {
		if setup != nil {
			_ = setup.MustApplyToPage(page)
		}
		if renderFunc != nil {
			renderFunc(page)
		}
	})
}
