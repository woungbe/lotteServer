package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

var JSESSIONID string

func main() {

	userId := ""
	userPw := ""

	getClientjsessionid() // 처음 접근하면 jsessionid 를 주는데 이게 있어야됨

	// HTTP 클라이언트 생성
	client, err := getClientWithCookie(JSESSIONID)
	if err != nil {
		fmt.Println("HTTP 클라이언트 생성 실패:", err)
		return
	}

	// 로그인 API 요청을 위한 데이터 설정 (body)
	body := url.Values{}
	body.Set("method", "login")
	body.Set("userId", userId)
	body.Set("password", userPw)

	// 로그인 API 요청 보내기
	resp, err := client.PostForm("https://www.dhlottery.co.kr/userSsl.do", body)
	if err != nil {
		fmt.Println("로그인 API 요청 실패:", err)
		return
	}
	defer resp.Body.Close()

	// 응답 바디 읽기
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("응답 바디 읽기 실패:", err)
		return
	}

	// 응답 바디를 통해 로그인 성공 여부 확인
	setCookieHeader := resp.Header.Get("Set-Cookie")
	fmt.Println("setCookieHeader : ", setCookieHeader)

	/*
		uid := extractUIDFromResponseHeader(resp.Header)
		if uid != "" {
			fmt.Println("로그인 성공!")
			fmt.Println("UID:", uid)
			// 여기서 uid를 사용하여 추가 작업을 수행할 수 있습니다.
		} else {
			fmt.Println("UID를 찾을 수 없습니다.")
		}
	*/

}

func getClientjsessionid() bool {
	url := "https://dhlottery.co.kr/common.do?method=main"

	// GET 요청 보내기
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("GET 요청 실패:", err)
		return false
	}
	defer resp.Body.Close()

	// JSESSIONID 쿠키 가져오기
	jsessionid := ""
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "JSESSIONID" {
			jsessionid = cookie.Value
			break
		}
	}

	if jsessionid != "" {
		fmt.Println("JSESSIONID:", jsessionid)
		JSESSIONID = jsessionid
		return true
	} else {
		fmt.Println("JSESSIONID를 찾을 수 없습니다.")
	}

	return false
}

// 아이디 가져오기
func getID(jsessionID, userid, password string) {

	// POST 요청에 필요한 데이터 설정 (body)
	body := fmt.Sprintf("method=login&userId=%s&password=%s", userid, password)

	// POST 요청을 보낼 URL 설정
	url := "https://www.dhlottery.co.kr/userSsl.do"

	// HTTP 클라이언트 생성
	client, err := getClientWithCookie(jsessionID)
	if err != nil {
		fmt.Println("HTTP 클라이언트 생성 실패:", err)
		return
	}

	// POST 요청 보내기
	resp, err := client.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(body))
	if err != nil {
		fmt.Println("POST 요청 실패:", err)
		return
	}
	defer resp.Body.Close()

	// 응답 바디 읽기
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("응답 바디 읽기 실패:", err)
		return
	}

	// fmt.Println("responseBody : ", string(responseBody))

	// fmt.Println("resp.Header : ", resp.Header)

	// UID 추출
	/*
		uid := extractUIDFromResponseHeader(resp.Header)
		if uid != "" {
			fmt.Println("UID:", uid)
			// 여기서 uid를 사용하여 추가 작업을 수행할 수 있습니다.
		} else {
			fmt.Println("로그인에 실패하였습니다.")
		}
	*/

}

// HTTP 클라이언트 생성 함수
func getClient() (*http.Client, error) {
	// Cookie 저장을 위한 Cookie Jar 생성
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	// HTTP 클라이언트 생성 및 Cookie Jar 설정
	client := &http.Client{Jar: cookieJar}
	return client, nil
}

// JSESSIONID를 쿠키로부터 추출하는 함수
func getJSESSIONIDFromCookie(cookies []*http.Cookie) string {
	jsessionID := ""
	for _, cookie := range cookies {
		if cookie.Name == "JSESSIONID" {
			jsessionID = cookie.Value
			break
		}
	}
	return jsessionID
}

// HTTP 클라이언트 생성 및 쿠키 설정 함수
func getClientWithCookie(jsessionID string) (*http.Client, error) {
	// Cookie 저장을 위한 Cookie Jar 생성
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	// JSESSIONID를 사용하여 쿠키 생성
	cookieURL, _ := url.Parse("https://www.dhlottery.co.kr")
	cookie := &http.Cookie{
		Name:  "JSESSIONID",
		Value: jsessionID,
	}
	cookieJar.SetCookies(cookieURL, []*http.Cookie{cookie})

	// HTTP 클라이언트 생성 및 Cookie Jar 설정
	client := &http.Client{Jar: cookieJar}
	return client, nil
}

// 응답 헤더에서 UID 추출하는 함수
func extractUIDFromResponseHeader(header http.Header) string {
	uid := ""
	setCookieHeader := header.Get("Set-Cookie")
	regex := regexp.MustCompile(`UID=([^;]+)`)
	matches := regex.FindStringSubmatch(setCookieHeader)
	if len(matches) > 1 {
		uid = matches[1]
	}
	return uid
}
