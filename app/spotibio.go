package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type PlayingSong struct {
	Item struct {
		Album struct {
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				Id   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				Uri  string `json:"uri"`
			} `json:"artists"`
		} `json:"album"`
		Artists []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			Id   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			Uri  string `json:"uri"`
		} `json:"artists"`
		Name       string `json:"name"`
		Popularity int    `json:"popularity"`
		PreviewUrl string `json:"preview_url"`
		Uri        string `json:"uri"`
	} `json:"item"`
}

type SpotifyRefreshTokenRes struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}
	loc, _ := time.LoadLocation(os.Getenv("TIMEZONE"))
	logTime := time.Now().In(loc).Format("2006-01-02 15:04:05")

	var listeningTo string
	playingSong, err := getCurrentPlayingSong()
	if err != nil {
		panic(err)
	}

	if playingSong == nil {
		fmt.Println(logTime + " - Not listening to music.")
		if checkIfWasNotListening() {
			fmt.Println("Nothing to update.")
			return
		} else {
			flagAsNotListening()
		}
	} else {
		listeningTo = beautifyListeningTo(*playingSong)
		fmt.Println(logTime + " - Listening to: " + listeningTo)
	}

	description := makeDescription(listeningTo)
	err1 := updateTweeterDescription(description)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
}

func checkIfWasNotListening() bool {
	_, err := os.OpenFile("notListening", os.O_WRONLY, 0600)
	if err != nil {
		return false
	}

	return true
}

func flagAsNotListening() {
	_, err := os.OpenFile("notListening", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	fmt.Println("Flagged as not listening")
}

func beautifyListeningTo(song PlayingSong) string {
	var artistsName []string

	for _, artist := range song.Item.Artists {
		artistsName = append(artistsName, artist.Name)
	}
	listeningTo := song.Item.Name + " | " + strings.Join(artistsName, ", ")

	return listeningTo
}

func makeDescription(listeningTo string) string {
	if listeningTo == "" {
		listeningTo = "Nothing is playing right now 😴\n\n"
	} else {
		listeningTo = fmt.Sprintf("I'm listening to %v On spotify 🎧 \n\n", listeningTo)
	}
	return listeningTo + os.Getenv("DEFAULT_DESCRIPTION")
}

func getCurrentPlayingSong() (*PlayingSong, error) {
	client := &http.Client{}

	accessToken, err0 := getSpotifyAccessToken()
	if err0 != nil {
		return nil, err0
	}

	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	req.Header.Add("Authorization", "Bearer "+*accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 204 {
		return nil, nil
	}
	if resp.StatusCode != 200 {
		err := errors.New(fmt.Sprintf("Http Error code from spotify : %v %v", resp.StatusCode, string(body)))
		return nil, err
	}

	var playingSong PlayingSong
	if err = json.Unmarshal([]byte(body), &playingSong); err != nil {
		return nil, err
	}

	return &playingSong, nil
}

func getSpotifyAccessToken() (*string, error) {
	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", os.Getenv("SPOTIFY_REFRESH_TOKEN"))
	encodedData := form.Encode()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(encodedData))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %v", os.Getenv("SPOTIFY_CLIENT_ID")))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		err := errors.New(fmt.Sprintf("Http Error from spotify refresh access token : %v", string(body)))
		return nil, err
	}
	var spotifyRefreshTokenRes SpotifyRefreshTokenRes
	if err = json.Unmarshal([]byte(body), &spotifyRefreshTokenRes); err != nil {
		return nil, err
	}

	return &spotifyRefreshTokenRes.AccessToken, nil
}

func updateTweeterDescription(description string) error {

	httpClient := createClient()
	path := fmt.Sprintf("https://api.twitter.com/1.1/account/update_profile.json?description=%v", url.QueryEscape(description))
	resp, err0 := httpClient.PostForm(path, nil)
	if err0 != nil {
		return err0
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		err := errors.New(fmt.Sprintf("Http Error code from tweeter : %v %v", resp.StatusCode, string(body)))
		return err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return err
	}

	fmt.Println("Tweeter description updated.")
	return nil
}

func createClient() *http.Client {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN_KEY"), os.Getenv("ACCESS_TOKEN_SECRET"))

	return config.Client(oauth1.NoContext, token)
}
