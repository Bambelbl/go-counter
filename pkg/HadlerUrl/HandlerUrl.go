package HadlerUrl

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type UrlSource struct {
	Url string
}

func (u UrlSource) Handler(stringPattern *regexp.Regexp) (currentCount uint64, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, u.Url, nil)
	if err != nil {
		return 0, fmt.Errorf("error new request to %s: %s", u.Url, err.Error())
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error do request to %s: %s", u.Url, err.Error())
	}

	defer func() {
		err = res.Body.Close()
		if err != nil {
			err = fmt.Errorf("error in close body: %s", err.Error())
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("error read body in %s: %s", u.Url, err.Error())
	}

	count := len(stringPattern.FindAll(body, -1))
	return uint64(count), nil
}
