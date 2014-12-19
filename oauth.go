package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type OAuthConfig struct {
	consumerKey       string `yaml:"CONSUMER_KEY"`
	consumerSecret    string `yaml:"CONSUMER_SECRET"`
	accessToken       string `yaml:"ACCESS_TOKEN"`
	accessTokenSecret string `yaml:"ACCESS_TOKEN_SECRET"`
	version           string `yaml:"VERSION"`
}

func (config OAuthConfig) header() *oauthHeader {
	return &oauthHeader{
		ConsumerKey:     config.consumerKey,
		Nonce:           oauthNonce(),
		SignatureMethod: "HMAC-SHA1",
		Timestamp:       strconv.Itoa(int(time.Now().Unix())),
		AccessToken:     config.accessToken,
		Version:         "1.0",
	}
}

func (config OAuthConfig) HashKey() string {
	return config.consumerSecret + "&" + config.accessTokenSecret
}

var headerFormat = `OAuth %s="%s", %s="%s", %s="%s", %s="%s", %s="%s", %s="%s", %s=%s"`

type oauthHeader struct {
	ConsumerKey     string
	Nonce           string
	Signature       string
	SignatureMethod string
	Timestamp       string
	AccessToken     string
	Version         string
}

func (h *oauthHeader) Map() map[string]string {
	return map[string]string{
		"oauth_consumer_key":     h.ConsumerKey,
		"oauth_nonce":            h.Nonce,
		"oauth_signature_method": h.SignatureMethod,
		"oauth_timestamp":        h.Timestamp,
		"oauth_token":            h.AccessToken,
		"oauth_version":          h.Version,
	}
}

func (h *oauthHeader) String() string {
	return fmt.Sprintf(headerFormat,
		url.QueryEscape(h.ConsumerKey),
		url.QueryEscape(h.Nonce),
		url.QueryEscape(h.Signature),
		url.QueryEscape(h.SignatureMethod),
		url.QueryEscape(h.Timestamp),
		url.QueryEscape(h.AccessToken),
		url.QueryEscape(h.Version))
}

func (config OAuthConfig) Sign(req *http.Request, header *oauthHeader) error {
	params, err := sortedOauthParameters(req, header)
	if err != nil {
		return err
	}

	s := oauthSignatureString(req, params)

	header.Signature = oauthSign(s, config.HashKey())
	req.Header.Add("Authorization", header.String())

	return nil
}

func oauthSign(s string, key string) string {
	sha1Hmac := hmac.New(sha1.New, []byte(key))
	sha1Hmac.Write([]byte(s))

	hash := sha1Hmac.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash)
}

func oauthNonce() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Int31()

	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, &n)
	return base64.URLEncoding.EncodeToString(b.Bytes())
}

func oauthSignatureString(req *http.Request, params []oauthParam) string {
	var signatureString bytes.Buffer

	// Method
	signatureString.WriteString(req.Method)
	signatureString.WriteByte('&')

	// URL
	signatureString.WriteString(percentEncode(req.URL.String()))
	signatureString.WriteByte('&')

	// Parm String
	var buf bytes.Buffer
	for i, p := range params {
		buf.WriteString(p.key)
		buf.WriteString("=")
		buf.WriteString(p.value)
		if i+1 < len(params) {
			buf.WriteByte('&')
		}
	}
	signatureString.WriteString(percentEncode(buf.String()))

	return signatureString.String()
}

type oauthParam struct {
	key   string
	value string
}

type paramList []oauthParam

func (p paramList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p paramList) Len() int { return len(p) }

func (p paramList) Less(i, j int) bool { return p[i].key < p[j].key }

func sortedOauthParameters(req *http.Request, h *oauthHeader) ([]oauthParam, error) {
	var bodyBytes bytes.Buffer
	bodyBytes.ReadFrom(req.Body)

	// Gather the needed components for the signing
	header := h.Map()
	query := req.URL.Query()
	body, err := url.ParseQuery(bodyBytes.String())
	if err != nil {
		return nil, err
	}

	params := make([]oauthParam, 0, len(header)+len(query)+len(body))
	for k, v := range query {
		params = append(params, oauthParam{key: percentEncode(k), value: percentEncode(v[0])})
	}
	for k, v := range body {
		params = append(params, oauthParam{key: percentEncode(k), value: percentEncode(v[0])})
	}
	for k, v := range header {
		params = append(params, oauthParam{key: percentEncode(k), value: percentEncode(v)})
	}

	sort.Sort(paramList(params))

	return params, nil
}
