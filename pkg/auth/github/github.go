package github

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	ghoauth2 "golang.org/x/oauth2/github"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	// TODO: randomize it
	oauthStateString = "pseudo-random"
)

type SessionCreator interface {
	NewSession(http.ResponseWriter)
}

type GitHub struct {
	oauth2.Config
	sessionCreator  SessionCreator
	AuthorizedUsers []string
}

type userInfo struct {
	Login string `json:"login"`
}

// New constructs a new GitHub struct from a config file.
func New(filename string, s SessionCreator) *GitHub {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	github := &GitHub{}
	github.Endpoint = ghoauth2.Endpoint
	err = json.NewDecoder(f).Decode(github)
	if err != nil {
		panic(err)
	}
	github.sessionCreator = s
	return github
}

func (g *GitHub) Login(w http.ResponseWriter, r *http.Request) {
	url := g.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (g *GitHub) Callback(w http.ResponseWriter, r *http.Request) {
	content, err := g.getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var u userInfo
	err = json.Unmarshal(content, &u)
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
		return
	}
	if g.isAuthorized(u.Login) {
		g.sessionCreator.NewSession(w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		fmt.Fprintf(w, "Access Denied.")
	}
}

func (g *GitHub) isAuthorized(username string) bool {
	for _, authorizedUser := range g.AuthorizedUsers {
		if username == authorizedUser {
			return true
		}
	}
	return false
}

func (g *GitHub) getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := g.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://api.github.com/user?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}
