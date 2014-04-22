package dictionary

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

const (
	LANGS_URL  = "https://dictionary.yandex.net/api/v1/dicservice.json/getLangs"
	LOOKUP_URL = "https://dictionary.yandex.net/api/v1/dicservice.json/lookup"
)

type YandexDictionary struct {
	apiKey string
	ui     string
}

type YD struct {
	Lang      string
	Text      string
	Family    bool
	Morpho    bool
	PosFilter bool
}

type YandexDictionaryEntry struct {
	Code    int
	Message string
	Def     []YandexDictionaryDef
}

type YandexDictionaryDef struct {
	Text string
	Pos  string
	Ts   string
	Tr   []YandexDictionaryTr
}

type YandexDictionaryTr struct {
	Text string
	Pos  string
	Ts   string
	Syn  []YandexDictionaryText
	Mean []YandexDictionaryText
	Ex   []YandexDictionaryEx
}

type YandexDictionaryEx struct {
	Text string
	Tr   []YandexDictionaryText
}

type YandexDictionaryText struct {
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
	resp, err := http.PostForm(LANGS_URL, url.Values{"key": {d.apiKey}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var langs []string
	if err := json.NewDecoder(resp.Body).Decode(&langs); err != nil {
		return nil, err
	}

	return langs, nil
}

func (d *YandexDictionary) Lookup(params *YD) (*YandexDictionaryEntry, error) {
	flagsMask := d.buildFlagsMask(params)
	builtParams := url.Values{"key": {d.apiKey}, "ui": {d.ui}, "lang": {params.Lang}, "text": {params.Text}, "flags": {flagsMask}}
	resp, err := http.PostForm(LOOKUP_URL, builtParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var entry YandexDictionaryEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, err
	}

	if entry.Code != 0 {
		return nil, errors.New(entry.Message)
	}

	return &entry, nil
}

func (d *YandexDictionary) buildFlagsMask(params *YD) string {
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
