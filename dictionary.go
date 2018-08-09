// Package dictionary is the Yandex.Dictionary API client
package dictionary

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

const (
	urlRoot    = "https://dictionary.yandex.net/api/v1/dicservice.json"
	langsPath  = "getLangs"
	lookupPath = "lookup"
)

// Dictionary holds api key and ui lang
type Dictionary struct {
	apiKey string
	ui     string
}

// Params for api request
type Params struct {
	Lang      string
	Text      string
	Family    bool
	Morpho    bool
	PosFilter bool
}

// Entry is the root api response struct
type Entry struct {
	Code    int
	Message string
	Def     []Def
}

// Def is the single definition
type Def struct {
	Text string
	Pos  string
	Ts   string
	Tr   []Tr
}

// Tr is the single translation
type Tr struct {
	Text string
	Pos  string
	Ts   string
	Syn  []Text
	Mean []Text
	Ex   []Ex
}

// Ex is the examples array
type Ex struct {
	Text string
	Tr   []Text
}

// Text encapsulates string explaining the entity (definition, translation, example)
type Text struct {
	Text string
}

// New returns dictionary instance with english ui
func New(apiKey string) *Dictionary {
	return NewUsingLang(apiKey, "en")
}

// NewUsingLang returns dictionary instance with specified ui
func NewUsingLang(apiKey string, ui string) *Dictionary {
	if ui == "" {
		ui = "en"
	}
	return &Dictionary{apiKey: apiKey, ui: ui}
}

// GetLangs returns list of supported languages
func (d *Dictionary) GetLangs() ([]string, error) {
	resp, err := http.PostForm(absURL(langsPath), url.Values{"key": {d.apiKey}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rawResponse interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}

	response, ok := rawResponse.(map[string]interface{})

	// actually "ok" means "error" in this case because the response not an array of languages
	// but a map with error code and message
	if ok {
		return nil, fmt.Errorf("(%v) %v", response["code"], response["message"])
	}

	var langs []string
	for _, v := range rawResponse.([]interface{}) {
		langs = append(langs, v.(string))
	}

	return langs, nil
}

// Lookup returns results of api request wrapped in Entry structs
func (d *Dictionary) Lookup(params *Params) (*Entry, error) {
	errMsg := fmt.Sprintf("can't get definitions for %s", params.Text)

	flagsMask := d.buildFlagsMask(params)
	builtParams := url.Values{"key": {d.apiKey}, "ui": {d.ui}, "lang": {params.Lang}, "text": {params.Text}, "flags": {flagsMask}}
	resp, err := http.PostForm(absURL(lookupPath), builtParams)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	defer resp.Body.Close()

	var entry Entry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, errors.Wrap(err, errMsg)
	}

	if entry.Code != 0 {
		return nil, errors.Errorf("%s: (%d) %s", errMsg, entry.Code, entry.Message)
	}

	if len(entry.Def) == 0 {
		return nil, errors.Errorf("%s: definitions are empty", errMsg)
	}

	return &entry, nil
}

func absURL(route string) string {
	return urlRoot + "/" + route
}

func (d *Dictionary) buildFlagsMask(params *Params) string {
	const (
		family    = 0x0001
		morpho    = 0x0004
		posFilter = 0x0008
	)

	var mask int
	if params.Family {
		mask |= family
	}
	if params.Morpho {
		mask |= morpho
	}
	if params.PosFilter {
		mask |= posFilter
	}

	return strconv.Itoa(mask)
}
