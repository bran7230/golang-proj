package server

// this is useful for making sure that the type we are reecieving is a json,
// someone sends an XML to us, and it executes instructions on our server ie:
/**
	Local file disclosure:
		<?xml version="1.0" encoding="ISO-8859-1"?>
		<!DOCTYPE foo [
		<!ELEMENT foo ANY >
		<!ENTITY xxe SYSTEM "file:///etc/passwd" >]>
		<foo>&xxe;</foo>

	Server side request forgery:
		<?xml version="1.0" encoding="ISO-8859-1"?>
		<!DOCTYPE foo [
		<!ELEMENT foo ANY >
		<!ENTITY xxe SYSTEM "http://169.254.169.254/latest/meta-data/iam/security-credentials/" >]>
		<foo>&xxe;</foo>

	Remote code execution:
		<?xml version="1.0" encoding="ISO-8859-1"?>
		<!DOCTYPE foo [
		<!ELEMENT foo ANY >
		<!ENTITY xxe SYSTEM "expect://id" >]>
		<foo>&xxe;</foo>

	Denial of service(billion laughs):
		<?xml version="1.0"?>
		<!DOCTYPE lolz [
		<!ENTITY lol "lol">
		<!ELEMENT lolz (#PCDATA)>
		<!ENTITY lol1 "&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;">
		<!ENTITY lol2 "&lol1;&lol1;&lol1;&lol1;&lol1;&lol1;&lol1;&lol1;&lol1;&lol1;">
		<!ENTITY lol3 "&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;">
		<!ENTITY lol4 "&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;">
		<!ENTITY lol5 "&lol4;&lol4;&lol4;&lol4;&lol4;&lol4;&lol4;&lol4;&lol4;&lol4;">
		]>
		<lolz>&lol5;</lolz>
**/

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func encode[T any](w http.ResponseWriter, status int, v T) error {
	// can change to XML or plain text if needed, and switch to another method of encoding.
	// ie w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

// decode was updated to take in a secret key, hmacSignature to decode and verify that the body has the correct data,
// and the header has the correct signature.
func decode[T any](bodyBytes []byte, secretKey []byte, hmacSignature string) (T, error) {
	var v T

	b, err := verify(bodyBytes, secretKey, hmacSignature)
	if err != nil {
		return v, err
	}

	if !b {
		return v, fmt.Errorf("invalid data, hmac does not match")
	}

	if err := json.Unmarshal(bodyBytes, &v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

// SendError helpful wrapper,
// you can change this as needed ie: you need to pass back the requests response(or par of it) to the client.
func SendError(w http.ResponseWriter, status int, msg string) error {
	payload := map[string]string{"error": msg}

	return encode(w, status, payload)
}
