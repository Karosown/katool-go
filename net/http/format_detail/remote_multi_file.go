package remote

import (
	"bytes"
	"encoding/json"
	"errors"

	remote "katool/net/http"
)

type RemoteMultiMultiFileDeCodeFormat struct {
	remote.DefaultEnDeCodeFormat
	Req     *remote.Req
	reqBody map[string]any
	Decoder remote.EnDeCodeFormat
	DaoMap  map[string]any
}

func (c *RemoteMultiMultiFileDeCodeFormat) ValidDecode(encode any) (bool, error) {
	if encode == nil {
		return true, nil
	}
	_, ok := encode.(map[string][]string)
	if !ok {
		return false, errors.New("encode is not map[string][]string")
	}
	return true, nil
}
func (e *RemoteMultiMultiFileDeCodeFormat) Encode(obj any) (any, error) {
	marshal, err := json.Marshal(obj)
	if err == nil {
		s := bytes.NewBuffer(marshal).String()
		return &s, nil
	}
	return nil, err
}

func (e *RemoteMultiMultiFileDeCodeFormat) Decode(encode any, back any) (any, error) {
	urlMap := encode.(map[string][]string)
	result := make(map[string][][]any)
	for k, v := range urlMap {
		if dao, ok := e.DaoMap[k]; ok {
			for i, url := range v {
				req := e.Req.Url(url).Format(e.Decoder)
				if body, status := e.reqBody[k]; status {
					req.Data(body)
				}
				build := req.Build(dao)
				if build != nil {
					if result[k] == nil {
						result[k] = make([][]any, len(v))
					}
					array, isArray := build.([]any)
					if isArray {
						result[k][i] = append(result[k][i], array...)
					} else {
						result[k][i] = append(result[k][i], build)
					}
				}
			}
		}
	}

	return result, nil
}