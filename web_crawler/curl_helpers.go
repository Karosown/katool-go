package web_crawler

import (
	"errors"
	"strings"

	"github.com/go-rod/rod"
)

// parseCurlSetup 解析 curl，返回 rod 配置与目标 URL。
func parseCurlSetup(curl string) (*RodPasteResult, string, error) {
	setup, url, err := ParseRodPasteFromCURLWithURL(curl, DefaultRodPasteOptions(""))
	if err != nil {
		return nil, "", err
	}
	if strings.TrimSpace(url) == "" {
		return nil, "", errors.New("curl has no url")
	}
	return setup, url, nil
}

// combineRenderFunc 合并 rod 配置与渲染回调。
func combineRenderFunc(setup *RodPasteResult, renderFunc func(*rod.Page)) func(*rod.Page) {
	if setup == nil {
		return renderFunc
	}
	return func(page *rod.Page) {
		_ = setup.MustApplyToPage(page)
		if renderFunc != nil {
			renderFunc(page)
		}
	}
}

// combineRequestModifiers 合并 rod 配置导出的请求修饰与用户提供的修饰器。
func combineRequestModifiers(setup *RodPasteResult, requestModifiers []RequestWith) []RequestWith {
	if setup == nil {
		return requestModifiers
	}
	mod := setup.RequestModifier()
	if mod == nil {
		return requestModifiers
	}
	if len(requestModifiers) == 0 {
		return []RequestWith{mod}
	}
	res := make([]RequestWith, 0, len(requestModifiers)+1)
	res = append(res, mod)
	res = append(res, requestModifiers...)
	return res
}
