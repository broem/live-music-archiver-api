package igscraper

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
)

var following []string

const (
	// Instagram User-Agent for iOS
	InstagramUserAgent = "Instagram 123.0.0.21.114 (iPhone; CPU iPhone OS 11_4 like Mac OS X; en_US; en-US; scale=2.00; 750x1334) AppleWebKit/605.1.15"
)

type Scraper struct {
	c         *http.Client
	csrfToken string
	loggedOn  bool
}

func NewScraper() *Scraper {
	cj1, _ := cookiejar.New(nil)
	return &Scraper{
		c: &http.Client{
			Jar: cj1,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
		loggedOn: false,
	}
}

func (s *Scraper) IsLoggedOn() bool {
	if s.loggedOn {
		return true
	}

	u, _ := url.Parse("https://www.instagram.com/")
	cookies, err := readCookiesFromDisk()
	if err != nil {
		log.Println(err)
		return false
	}

	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, cookies)
	s.c.Jar = jar

	req, err := http.NewRequest("GET", "https://www.instagram.com/", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := s.c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return len(body) > 0
}

func (s *Scraper) DoLogin(user, pass string) {
	req1, err := http.NewRequest("GET", "https://www.instagram.com/", nil)
	if err != nil {
		log.Fatal("wtf")
	}

	req1.Header.Set("Referer", InstagramUserAgent)
	bb, err := s.c.Do(req1)
	if err != nil {
		log.Fatal(err)
	}
	defer bb.Body.Close()
	_, err = ioutil.ReadAll(bb.Body)
	if err != nil {
		log.Fatal(err)
	}

	u, _ := url.Parse("https://www.instagram.com/")
	for _, value := range s.c.Jar.Cookies(u) {
		if strings.Contains(value.Name, "csrftoken") {
			s.csrfToken = value.Value
		}
	}

	data, _ := json.Marshal(map[string]string{"username": user, "password": pass, "_csrftoken": s.csrfToken})

	vs := url.Values{}
	bf := bytes.NewBuffer([]byte{})

	for k, v := range generateSignature(b2s([]byte(data))) {
		vs.Add(k, v)
	}

	bf.WriteString(vs.Encode())

	req, err := http.NewRequest("POST", "https://www.instagram.com/accounts/login/ajax/", bf)
	if err != nil {
		log.Fatal("somethings wrong")
	}
	req.Header.Add("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Add("user-agent", InstagramUserAgent)

	resp, err := s.c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(body)
	for _, value := range s.c.Jar.Cookies(u) {
		if strings.Contains(value.Name, "csrftoken") {
			s.csrfToken = value.Value
		}
	}

	saveHeadersToDisk(req.Header)
	saveCookiesToDisk(s.c.Jar.Cookies(u))

	s.loggedOn = true
}

func readCookiesFromDisk() ([]*http.Cookie, error) {
	file := "cookies.txt"
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cookies := []*http.Cookie{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var cookie *http.Cookie
		if cookie, err = CookieFromString(scanner.Text()); err != nil {
			return nil, err
		}
		cookies = append(cookies, cookie)
	}
	return cookies, nil
}

func CookieFromString(s string) (*http.Cookie, error) {
	var c http.Cookie
	parts := strings.Split(s, ";")
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid cookie: %q", part)
		}
		switch kv[0] {
		case "Domain":
			c.Domain = kv[1]
		case "Expires":
			c.Expires = time.Now()
		case "Max-Age":
			secs, err := strconv.Atoi(kv[1])
			if err != nil {
				return nil, err
			}
			c.MaxAge = secs
		case "Path":
			c.Path = kv[1]
		case "Secure":
			c.Secure = true
		case "HttpOnly":
			c.HttpOnly = true
		case "SameSite":
			ss, err := SameSiteFromString(kv[1])
			if err != nil {
				return nil, err
			}
			c.SameSite = http.SameSite(ss)
		case "Raw":
			c.Raw = kv[1]
		case "Unparsed":
			c.Unparsed = append(c.Unparsed, kv[1])
		default:
			c.Name = kv[0]
			c.Value = kv[1]
		}
	}
	return &c, nil
}

func SameSiteFromString(s string) (http.SameSite, error) {
	switch s {
	case "Strict":
		return http.SameSiteStrictMode, nil
	case "Lax":
		return http.SameSiteLaxMode, nil
	case "None":
		return http.SameSiteNoneMode, nil
	}
	return 0, fmt.Errorf("invalid SameSite value: %q", s)
}

func saveCookiesToDisk(cookies []*http.Cookie) {
	file := "cookies.txt"
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for _, value := range cookies {
		f.WriteString(value.Name + "=" + value.Value + "\n")
	}
}

func readHeadersFromDisk() (http.Header, error) {
	file := "headers.txt"
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	headers := make(http.Header)
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				headers.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}
	return headers, nil
}

func saveHeadersToDisk(h http.Header) {
	file := "headers.txt"
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for k, v := range h {
		f.WriteString(k + ": " + v[0] + "\n")
	}
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func generateSignature(data string) map[string]string {
	m := make(map[string]string)
	m["ig_sig_key_version"] = "4"
	m["signed_body"] = fmt.Sprintf(
		"%s.%s", generateHMAC(data, "iGuessThisMEansNothing"), data,
	)
	return m
}

func generateHMAC(text, key string) string {
	hasher := hmac.New(sha256.New, []byte(key))
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s *Scraper) GetUserInfo(username string) (string, error) {
	UrlUserInfo := []string{"https://www.instagram.com/", username, "/"}
	completeUrl := strings.Join(UrlUserInfo, "")
	u, _ := url.Parse(completeUrl)
	cookies, err := readCookiesFromDisk()
	if err != nil {
		println(err)
	}

	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, cookies)
	s.c.Jar = jar

	req, err := http.NewRequest("GET", completeUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Add("user-agent", InstagramUserAgent)
	req.Header.Set("Referer", InstagramUserAgent)
	resp, err := s.c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	decode := string(body)
	//fmt.Println(decode)
	var getUserId []string
	for {
		getUserId = s.multiBetween(decode, `"page_id":"profilePage_`, `","profile_id":`)
		if getUserId[0] != "52578078225" {
			break //this is just to prevent it from picking me.
		}
		decode = getUserId[1]
	}

	userInfo := getUserId[0]
	return userInfo, nil
}

func (s *Scraper) GetFollowing(userInfo string) []string {
	rand.Seed(time.Now().UnixNano())
	max := 4
	min := 1
	delayScrape := rand.Intn(max-min) + min
	time.Sleep(time.Duration(delayScrape) * time.Second)

	makeUrl := []string{"https://i.instagram.com/api/v1/friendships/", userInfo, "/following/?count=12"}
	completeUrl := strings.Join(makeUrl, "")
	log.Println(completeUrl)

	u, _ := url.Parse(completeUrl)
	cookies, err := readCookiesFromDisk()
	if err != nil {
		println(err)
	}

	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, cookies)
	s.c.Jar = jar

	req, err := http.NewRequest("GET", completeUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Add("user-agent", InstagramUserAgent)
	req.Header.Set("Referer", InstagramUserAgent)
	resp, err := s.c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	decode := string(body)
	fmt.Println(decode)

	if err != nil {
		log.Fatal(err)
	}
	var nextFollowing string
	var repeatNameCheck = false
	var startTrack = 0
	var firstName = ""
	for {
		tempDecode := decode
		time.Sleep(time.Duration(delayScrape) * time.Second)
		for {
			nameAndVal := s.multiBetween(tempDecode, `"username":"`, `","full_name":`)
			if nameAndVal[0] != "" {
				if firstName == nameAndVal[0] {
					repeatNameCheck = true
				}
				if startTrack == 0 {
					firstName = nameAndVal[0]
				}
				if startTrack == 0 {
					startTrack++
				}
				following = append(following, nameAndVal[0])
			}
			tempDecode = nameAndVal[1]
			if nameAndVal[0] == "" {
				break
			}
		}
		if repeatNameCheck {
			break
		}
		nextFollowing = s.between(decode, `"next_max_id":"`, `","status":"ok"`)
		decode = ""
		makeUrlNext := []string{completeUrl, "&max_id=", nextFollowing}
		urlNext := strings.Join(makeUrlNext, "")
		log.Println(urlNext)
		req2, err := http.NewRequest("GET", urlNext, nil)
		if err != nil {
			log.Fatal(err)
		}
		req2.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req2.Header.Set("Accept-Language", "en-US")
		req2.Header.Add("user-agent", InstagramUserAgent)
		req2.Header.Set("Referer", InstagramUserAgent)
		resp2, err := s.c.Do(req2)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body2, err := ioutil.ReadAll(resp2.Body)
		if err != nil {
			log.Fatal(err)
		}
		decode = string(body2)
		fmt.Println(decode)
	}

	return following
}

func (s *Scraper) GetPosts(userInfo string) []string {
	rand.Seed(time.Now().UnixNano())
	max := 50
	min := 20
	delayScrape := rand.Intn(max-min) + min
	time.Sleep(time.Duration(delayScrape) * time.Second)

	makeUrl := []string{"https://i.instagram.com/api/v1/feed/user/", userInfo, ""}
	completeUrl := strings.Join(makeUrl, "")
	log.Println(completeUrl)

	u, _ := url.Parse(completeUrl)
	cookies, err := readCookiesFromDisk()
	if err != nil {
		println(err)
	}

	jar, _ := cookiejar.New(nil)
	jar.SetCookies(u, cookies)
	s.c.Jar = jar

	req, err := http.NewRequest("GET", completeUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Add("user-agent", InstagramUserAgent)
	req.Header.Set("Referer", InstagramUserAgent)
	resp, err := s.c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	decode := string(body)
	fmt.Println(decode)

	var feedMedia *FeedMedia
	err = json.Unmarshal([]byte(decode), &feedMedia)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	// turn FeedMedia into a slice of strings
	var posts []string
	for _, item := range feedMedia.Items {
		// turn the item into a string
		b, err := json.Marshal(item)
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		posts = append(posts, string(b))
	}

	return posts
}

func (s *Scraper) ScrapeFollowing(following []string) {

	for x := 0; x < len(following); x++ {
		rand.Seed(time.Now().UnixNano())
		max := 30
		min := 10
		delayScrape := rand.Intn(max-min) + min
		time.Sleep(time.Duration(delayScrape) * time.Second)
		UrlUserInfo := []string{"https://www.instagram.com/", following[x], "/"}
		completeUrl := strings.Join(UrlUserInfo, "")
		u, _ := url.Parse(completeUrl)
		cookies, err := readCookiesFromDisk()
		if err != nil {
			println(err)
		}

		jar, _ := cookiejar.New(nil)
		jar.SetCookies(u, cookies)
		s.c.Jar = jar

		req, err := http.NewRequest("GET", completeUrl, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("Accept-Language", "en-US")
		req.Header.Add("user-agent", InstagramUserAgent)
		req.Header.Set("Referer", InstagramUserAgent)
		resp, err := s.c.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		decode := string(body)
		fmt.Println(decode)
	}

}

func (s *Scraper) between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func (s *Scraper) multiBetween(value string, a string, b string) []string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return []string{"", ""}
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return []string{"", ""}
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return []string{"", ""}
	}

	posLastAdjust := len(b)
	newVal := value[posLast+posLastAdjust:]
	return []string{value[posFirstAdjusted:posLast], newVal}
}
