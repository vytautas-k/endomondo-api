package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// User ...
type User struct {
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
	Remember bool   `json:"remember" form:"remember" query:"remember"`
}

func main() {
	e := echo.New()

	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		fmt.Printf("%s\n", resBody)
	}))

	// e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
	// 	TokenLookup: "header:" + echo.HeaderXCSRFToken,
	// 	CookieName:  "CSRF_TOKEN",
	// }))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/users/:id", func(c echo.Context) error {
		return c.String(http.StatusOK, "/users/:id")
	})

	// Handler
	session := func(c echo.Context) (err error) {
		data := url.Values{}
		data.Set("email", "weecka@gmail.com")
		data.Set("password", "freebsd")
		data.Set("remember", "true")

		mapD := map[string]string{"email": "weecka@gmail.com", "password": "freebds", "remember": "true"}
		mapB, _ := json.Marshal(mapD)

		var csrfToken = "-"

		// csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
		// fmt.Println(strings.NewReader(data.Encode()))

		// s := "{\"email\": \"weecka@gmail.com\", \"password\":\"freebds\", \"remember\": \"true\"}"
		// fmt.Println(strings.NewReader(s))

		fmt.Println(strings.NewReader(string(mapB)))
		fmt.Println(bytes.NewBufferString(string(mapB)))

		var jsonStr = []byte(`{"email": "weecka@gmail.com", "password": "freebds", "remember": "true"}`)

		// Build the request
		req, err := http.NewRequest(echo.POST, "https://www.endomondo.com/rest/session", bytes.NewBuffer(jsonStr))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Cookie", "CSRF_TOKEN="+csrfToken)
		req.Header.Add("X-CSRF-TOKEN", csrfToken)

		if err != nil {
			fmt.Println("Error is req: ", err)
		}

		cookieJar, err := cookiejar.New(nil)

		if err != nil {
			fmt.Println("Cookie Error:", err)
		}

		// create a Client
		client := &http.Client{
			Jar: cookieJar,
		}

		// Do sends an HTTP request and
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error in send req: ", err)
		}

		// Defer the closing of the body
		defer resp.Body.Close()

		// body, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// fmt.Println("error parsing data: ", err)
		// }
		// fmt.Println(body)

		buf, bodyErr := ioutil.ReadAll(req.Body)
		if bodyErr != nil {
			fmt.Println("bodyErr ", bodyErr.Error())
			return
		}

		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
		fmt.Println("BODY: ", rdr1)
		req.Body = rdr2

		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(requestDump))

		responseDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(responseDump))

		// Use json.Decode for reading streams of JSON data
		// if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		// 	fmt.Println(err)
		// }

		return c.JSON(http.StatusOK, "OK")
	}

	e.GET("api/session", session)

	e.GET("api/test", func(c echo.Context) (err error) {
		test()
		return c.JSON(http.StatusOK, "TEST")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func test() {
	url := "https://www.endomondo.com/rest/session"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"email": "weecka@gmail.com", "password": "freebds", "remember": "true"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "CSRF_TOKEN=-")
	req.Header.Add("X-CSRF-TOKEN", "-")
	req.Header.Add("Cache-Control", "no-cache")

	reqBody, _ := ioutil.ReadAll(req.Body)
	fmt.Println("Request Body:", string(reqBody))

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	cookieJar, err := cookiejar.New(&options)

	if err != nil {
		fmt.Println("Cookie Error:", err)
	}

	client := &http.Client{
		Jar: cookieJar,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(requestDump))

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
