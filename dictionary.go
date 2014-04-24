package yandex_dictionary

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	URL_ROOT    = "https://dictionary.yandex.net/api/v1/dicservice.json"
	LANGS_PATH  = "getLangs"
	LOOKUP_PATH = "lookup"
)

type YandexDictionary struct {
	apiKey string
	ui     string
}

type Params struct {
	Lang      string
	Text      string
	Family    bool
	Morpho    bool
	PosFilter bool
}

type Entry struct {
	Code    int
	Message string
	Def     []Def
}

type Def struct {
	Text string
	Pos  string
	Ts   string
	Tr   []Tr
}

type Tr struct {
	Text string
	Pos  string
	Ts   string
	Syn  []Text
	Mean []Text
	Ex   []Ex
}

type Ex struct {
	Text string
	Tr   []Text
}

type Text struct {
	Text string
}

func New(apiKey string) *YandexDictionary {
	return NewUsingLang(apiKey, "en")
}

func NewUsingLang(apiKey string, ui string) *YandexDictionary {
	if ui == "" {
		ui = "en"
	}
	return &YandexDictionary{apiKey: apiKey, ui: ui}
}

func (d *YandexDictionary) GetLangs() ([]string, error) {
	resp, err := http.PostForm(absUrl(LANGS_PATH), url.Values{"key": {d.apiKey}})
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

func (d *YandexDictionary) Lookup(params *Params) (*Entry, error) {
	flagsMask := d.buildFlagsMask(params)
	builtParams := url.Values{"key": {d.apiKey}, "ui": {d.ui}, "lang": {params.Lang}, "text": {params.Text}, "flags": {flagsMask}}
	resp, err := http.PostForm(absUrl(LOOKUP_PATH), builtParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var entry Entry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, err
	}

	if entry.Code != 0 {
		return nil, fmt.Errorf("(%v) %v", entry.Code, entry.Message)
	}

	return &entry, nil
}

func absUrl(route string) string {
	return URL_ROOT + "/" + route
}

func (d *YandexDictionary) buildFlagsMask(params *Params) string {
	const (
		FAMILY     = 0x0001
		MORPHO     = 0x0004
		POS_FILTER = 0x0008
	)

	var mask int
	if params.Family {
		mask |= FAMILY
	}
	if params.Morpho {
		mask |= MORPHO
	}
	if params.PosFilter {
		mask |= POS_FILTER
	}

	return strconv.Itoa(mask)
}
