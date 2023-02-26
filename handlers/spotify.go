package handlers

import (
	"file-api/structs"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

const SPOTIFY_CLIENT_ID = "f89f5d1747434a36b690f12445a4c77b"
const SPOTIFY_CLIENT_SECRET = "8305d8edc9d24fb8944d8b312e68980b"
const REDIRECT_URL = "https://pichsereyvattanachan.vercel.app/home/"

func LoginSpotify(c echo.Context) error {
	//"https://accounts.spotify.com/authorize?response_type=code&client_id=" + encodeURI("f89f5d1747434a36b690f12445a4c77b") + "&redirect_uri=" + encodeURI("http://localhost:3000/" + "&scope=" + encodeURI("user-read-currently-playing") + "&client_secret=" + encodeURI("8305d8edc9d24fb8944d8b312e68980b")
	// c.Redirect(200, fmt.Sprintf(`https://accounts.spotify.com/authorize?response_type=code&client_id=%s&client_secret=%s&redirect_uri=%s&scope=user-read-currently-playing`,  SPOTIFY_CLIENT_ID, SPOTIFY_CLIENT_SECRET, REDIRECT_URL))

	// c.Redirect(http.StatusFound, fmt.Sprintf(`https://accounts.spotify.com/authorize?response_type=token&client_id=%s&client_secret=%s&redirect_uri=%s&scope=user-read-currently-playing&state=aaaaaaaaa&show_dialog=false`,  SPOTIFY_CLIENT_ID, SPOTIFY_CLIENT_SECRET, REDIRECT_URL))
	// fmt.Println(fmt.Sprintf(`https://accounts.spotify.com/authorize?response_type=code&client_id=%s&client_secret=%s&redirect_uri=%s&scope=user-read-currently-playing&state=aaaaaaaaa&show_dialog=false`,  SPOTIFY_CLIENT_ID, SPOTIFY_CLIENT_SECRET, REDIRECT_URL))
	request, _ := http.NewRequest("GET", fmt.Sprintf(`https://accounts.spotify.com/authorize?response_type=token&client_id=%s&client_secret=%s&redirect_uri=%s&scope=user-read-currently-playing&show_dialog=false`,  SPOTIFY_CLIENT_ID, SPOTIFY_CLIENT_SECRET, REDIRECT_URL), nil)
    client := &http.Client{}
    resp, err := client.Do(request)
    if err != nil {
        panic(err)
    }
	//https://accounts.spotify.com/authorize?scope=user-read-currently-playing&response_type=token&client_secret=8305d8edc9d24fb8944d8b312e68980b&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2F&state=aaaaaaaaa&client_id=f89f5d1747434a36b690f12445a4c77b&show_dialog=false
	// c.Redirect(http.StatusFound, "https://accounts.spotify.com/authorize?scope=user-read-currently-playing&response_type=token&client_secret=8305d8edc9d24fb8944d8b312e68980b&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2F&state=aaaaaaaaa&client_id=f89f5d1747434a36b690f12445a4c77b&show_dialog=false")


	
	return c.JSON(200, structs.Message{Message: resp.Request.URL.String(), Code: 200})
}
