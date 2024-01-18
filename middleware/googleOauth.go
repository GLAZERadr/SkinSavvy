package middleware

import (
	"fmt"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/InnoFours/skin-savvy/config"
)

type GoogleOauthToken struct {
	Access_token	string 
	Id_token		string
}

type GoogleUserResult struct {
	Id				string 
	Email			string 
	Verified_email	bool
	Name			string 
	Given_name		string 
	Family_name		string 
	Picture			string 
	Locale			string
}

func GetGoogleOauthToken(code string) (*GoogleOauthToken, error) {
	const rootUrl = "https://oauth2.googleapis.com/token"

	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("client_id", config.ConfigGoogleOauthClientId())
	values.Add("client_secret", config.ConfigGoogleOauthClientSecret())
	values.Add("redirect_uri", config.ConfigGoogleOauthRedirectUrl())

	query := values.Encode()

	req, err := http.NewRequest("POST", rootUrl, bytes.NewBufferString(query))
	if err != nil {
		log.Println("error making request on oauth: ", err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("error making requeest on client: ", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		log.Println("can't retrieve any token")
	}

	var resBody bytes.Buffer
	_,err = io.Copy(&resBody, res.Body)
	if err != nil {
		log.Println("can't copy response body of oauth token")
	}

	var GoogleOauthTokenRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleOauthTokenRes); err != nil {
		log.Println("can't unmarshal json body of google oauth token: ", err.Error())
	}

	tokenBody := &GoogleOauthToken{
		Access_token	: GoogleOauthTokenRes["access_token"].(string),
		Id_token		: GoogleOauthTokenRes["id_token"].(string),
	}

	return tokenBody, nil
}

func GetGoogleUser(accessToken string, idToken string) (*GoogleUserResult, error) {
	rootUrl := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", accessToken)

	req, err := http.NewRequest("GET", rootUrl, nil)
	if err != nil {
		log.Println("can't make a request on google api oath: ", err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("BEARER %s", idToken))

	client := http.Client{
		Timeout : time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("can't make an authorize request: ", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		log.Println("error on status code", err.Error())
	}

	var resBody bytes.Buffer
	_,err = io.Copy(&resBody, res.Body)
	if err != nil {
		log.Println("can't copy response body of oauth token")
	}

	var GoogleUserRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleUserRes); err != nil {
		log.Println("can't unmarshal json body of google oauth token: ", err.Error())
	}

	userBody := &GoogleUserResult{
		Id				: GoogleUserRes["id"].(string),
		Email			: GoogleUserRes["email"].(string),
		Verified_email	: GoogleUserRes["verified_email"].(bool),
		Name			: GoogleUserRes["name"].(string),
		Given_name		: GoogleUserRes["given_name"].(string),
		Picture			: GoogleUserRes["picture"].(string),
		Locale			: GoogleUserRes["locale"].(string),
	}
	
	return userBody, nil
}