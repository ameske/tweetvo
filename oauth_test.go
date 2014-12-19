package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

var expectedHashKey = "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw&LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE"
var expectedSigString = "POST&https%3A%2F%2Fapi.twitter.com%2F1%2Fstatuses%2Fupdate.json&include_entities%3Dtrue%26oauth_consumer_key%3Dxvz1evFS4wEEPTGEFPHBog%26oauth_nonce%3DkYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1318622958%26oauth_token%3D370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb%26oauth_version%3D1.0%26status%3DHello%2520Ladies%2520%252B%2520Gentlemen%252C%2520a%2520signed%2520OAuth%2520request%2521"
var expectedSignature = "tnnArxj06cWHq44gCs1OSKk/jLY="

var testRequest *http.Request

var testConfig = OAuthConfig{
	consumerKey:       "xvz1evFS4wEEPTGEFPHBog",
	consumerSecret:    "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
	accessToken:       "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
	accessTokenSecret: "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
	version:           "1.0",
}

var testHeader = oauthHeader{
	ConsumerKey:     "xvz1evFS4wEEPTGEFPHBog",
	Nonce:           "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
	SignatureMethod: "HMAC-SHA1",
	Timestamp:       "1318622958",
	AccessToken:     "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
	Version:         "1.0",
}

var testBody = url.Values{
	"include_entities": []string{"true"},
	"status":           []string{"Hello Ladies + Gentlemen, a signed OAuth request!"},
}

func init() {
}

func TestConfigHashKey(t *testing.T) {
	key := testConfig.HashKey()
	if key != expectedHashKey {
		t.Errorf("Expected: %s\tGot: %s", expectedHashKey, key)
	}
}

func TestOAuthSigning(t *testing.T) {
	req, err := http.NewRequest("POST", "https://api.twitter.com/1/statuses/update.json", bytes.NewBuffer([]byte(testBody.Encode())))
	if err != nil {
		t.Fatal(err.Error())
	}

	params, err := sortedOauthParameters(req, &testHeader)
	if err != nil {
		t.Fatal(err.Error())
	}

	sigString := oauthSignatureString(req, params)

	if sigString != expectedSigString {
		t.Errorf("\n\tExpected: %s\n\tGot: %s", expectedSigString, sigString)
	}

	signature := oauthSign(sigString, testConfig.HashKey())

	if signature != expectedSignature {
		t.Errorf("Expected: %s\tGot: %s", expectedSignature, signature)
	}
}
