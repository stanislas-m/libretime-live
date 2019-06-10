package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/pkg/errors"
)

const apiURL = "/api/live-info-v2"

// APIFetcher gets the meta data to build live info.
type APIFetcher interface {
	PollAPI() error
	Live() liveInfo
}

type apiFetcher struct {
	client *http.Client
	host   string
	moot   sync.RWMutex
	live   liveInfo
}

func (a *apiFetcher) Live() liveInfo {
	a.moot.RLock()
	defer a.moot.RUnlock()
	return a.live
}

func (a *apiFetcher) PollAPI() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", a.host, apiURL), nil)
	if err != nil {
		return errors.Wrap(err, "could not fetch live info")
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := a.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "could not make request for live info")
	}
	if res.StatusCode != 200 {
		return errors.Errorf("error while fetching API: %d", res.StatusCode)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "could not read live info")
	}

	live := liveInfo{}
	if err := json.Unmarshal(body, &live); err != nil {
		return errors.Wrap(err, "could not decode live info")
	}

	a.moot.Lock()
	a.live = live
	a.moot.Unlock()

	return nil
}
