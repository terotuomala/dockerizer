package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Token struct {
	Token string `json:"token"`
}

type Digest struct {
	Manifests []struct {
		Digest   string `json:"digest"`
		Platform Platform
	} `json:"manifests"`
}

type Platform struct {
	Architecture string `json:"architecture"`
	Os           string `json:"os"`
}

func getDockerHubAuthToken(image string) string {
	req, err := http.NewRequest("GET", "https://auth.docker.io/token?scope=repository:"+image+":pull&service=registry.docker.io", nil)
	if err != nil {
		log.Fatalf("Token request failed %v", err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Token fetch failed %v", err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Body read failed %v", err.Error())
	}

	var result Token
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Can not unmarshal JSON %v", err.Error())
	}

	return PrettyfyJson(result.Token)
}

func GetDockerImageDigest(image, tag string) string {
	token := getDockerHubAuthToken(image)
	bearer := "Bearer " + token

	req, err := http.NewRequest("GET", "https://registry-1.docker.io/v2/"+image+"/manifests/"+tag, nil)
	if err != nil {
		log.Fatalf("Digest request failed %v", err.Error())
	}

	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.list.v2+json")
	req.Header.Set("Authorization", bearer)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Digest fetch failed %v", err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Body read failed %v", err.Error())
	}

	var result Digest
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Can not unmarshal JSON %v", err.Error())
	}

	return result.Manifests[0].Digest
}

func PrettyfyJson(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	trim := strings.Trim(string(s), "\"")
	return string(trim)
}
