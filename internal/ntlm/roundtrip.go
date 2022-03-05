package ntlm

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type NTLMNegotiator struct {
	http.RoundTripper
	User, Pass, Domain, Workstation string
	IsHash                          bool
	Debug                           bool
}

func listContainsPrefix(list []string, prefix string) bool {
	for _, v := range list {
		if strings.HasPrefix(v, prefix) {
			return true
		}
	}
	return false
}

func (n *NTLMNegotiator) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := n.RoundTripper
	if rt == nil {
		rt = http.DefaultTransport
	}

	body := bytes.Buffer{}
	if req.Body != nil {
		_, err := body.ReadFrom(req.Body)
		if err != nil {
			if n.Debug {
				log.Printf("[DEBUG] [Read Request Body] %s\n", err.Error())
			}
			return nil, err
		}

		req.Body.Close()
		req.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))
	}

	req.Header.Del("Authorization")
	res, err := rt.RoundTrip(req)
	if err != nil {
		if n.Debug {
			log.Printf("[DEBUG] [Initial RT Error] %s\n", err.Error())
		}
		return nil, err
	}
	if res.StatusCode != http.StatusUnauthorized {
		return res, nil
	}
	resauth := []string(res.Header.Values("Www-Authenticate"))
	if listContainsPrefix(resauth, "NTLM") || listContainsPrefix(resauth, "Negotiate") {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()

		// Negotiate
		negotiateMessage, err := n.newNegotiateMessage(n.Domain, n.Workstation)
		if err != nil {
			if n.Debug {
				log.Printf("[DEBUG] [Negotiation Message Creation] %s\n", err.Error())
			}
			return nil, err
		}

		if listContainsPrefix(resauth, "NTLM") {
			req.Header.Set("Authorization", "NTLM "+base64.StdEncoding.EncodeToString(negotiateMessage))
		} else {
			req.Header.Set("Authorization", "Negotiate "+base64.StdEncoding.EncodeToString(negotiateMessage))
		}

		req.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))

		res, err = rt.RoundTrip(req)
		if err != nil {
			if n.Debug {
				log.Printf("[DEBUG] [Negotiation RT Error] %s\n", err.Error())
			}
			return nil, err
		}

		// Process challenge
		resauth = []string(res.Header.Values("Www-Authenticate"))
		challengeMessage, err := GetData(resauth)
		if err != nil {
			if n.Debug {
				log.Printf("[DEBUG] [Challenge Message Creation] %s\n", err.Error())
			}
			return nil, err
		}

		if !(listContainsPrefix(resauth, "NTLM") || listContainsPrefix(resauth, "Negotiate")) || len(challengeMessage) == 0 {
			// Negotiation failed.
			return res, nil
		}
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()

		// Create authentication message
		authMessage, err := n.craftResponse(challengeMessage, n.User, n.Pass, n.IsHash)
		if err != nil {
			if n.Debug {
				log.Printf("[DEBUG] [Authentication Message Creation] %s\n", err.Error())
			}
			return nil, err
		}

		if listContainsPrefix(resauth, "NTLM") {
			req.Header.Set("Authorization", "NTLM "+base64.StdEncoding.EncodeToString(authMessage))
		} else {
			req.Header.Set("Authorization", "Negotiate "+base64.StdEncoding.EncodeToString(authMessage))
		}

		req.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))
		return rt.RoundTrip(req)
	}

	return res, nil
}

func GetData(h []string) ([]byte, error) {
	for _, s := range h {
		if strings.HasPrefix(string(s), "NTLM") || strings.HasPrefix(string(s), "Negotiate") || strings.HasPrefix(string(s), "Basic ") {
			p := strings.Split(string(s), " ")
			if len(p) < 2 {
				return nil, nil
			}
			return base64.StdEncoding.DecodeString(string(p[1]))
		}
	}
	return nil, nil
}
