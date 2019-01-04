/*
Copyright 2015 All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedirectToAuthorizationUnauthorized(t *testing.T) {
	requests := []fakeRequest{
		{URI: "/admin", ExpectedCode: http.StatusUnauthorized},
	}
	newFakeProxy(nil).RunTests(t, requests)
}

func TestRedirectToAuthorization(t *testing.T) {
	requests := []fakeRequest{
		{
			URI:              "/admin",
			Redirects:        true,
			ExpectedLocation: "/oauth/authorize?state",
			ExpectedCode:     http.StatusTemporaryRedirect,
		},
	}
	newFakeProxy(nil).RunTests(t, requests)
}

func TestRedirectToAuthorizationWith303Enabled(t *testing.T) {
	cfg := newFakeKeycloakConfig()
	cfg.InvalidAuthRedirectsWith303 = true

	requests := []fakeRequest{
		{
			URI:              "/admin",
			Redirects:        true,
			ExpectedLocation: "/oauth/authorize?state",
			ExpectedCode:     http.StatusSeeOther,
		},
	}
	newFakeProxy(cfg).RunTests(t, requests)
}

func TestRedirectToAuthorizationSkipToken(t *testing.T) {
	requests := []fakeRequest{
		{URI: "/admin", ExpectedCode: http.StatusUnauthorized},
	}
	c := newFakeKeycloakConfig()
	c.SkipTokenVerification = true
	newFakeProxy(c).RunTests(t, requests)
}

func assertAlmostEquals(t *testing.T, expected time.Duration, actual time.Duration) {
	delta := expected - actual
	if delta < 0 {
		delta = -delta
	}
	assert.True(t, delta < time.Duration(1)*time.Minute, "Diff should be less than a minute but delta is %s", delta)
}

func TestGetAccessCookieExpiration_NoExp(t *testing.T) {
	token := newTestToken("foo").getToken()
	refreshToken := token.Encode()
	c := newFakeKeycloakConfig()
	c.AccessTokenDuration = time.Duration(1) * time.Hour
	proxy := newFakeProxy(c).proxy
	duration := proxy.getAccessCookieExpiration(token, refreshToken)
	assertAlmostEquals(t, c.AccessTokenDuration, duration)
}

func TestGetAccessCookieExpiration_ZeroExp(t *testing.T) {
	ft := newTestToken("foo")
	ft.setExpiration(time.Unix(0, 0))
	token := ft.getToken()
	refreshToken := token.Encode()
	c := newFakeKeycloakConfig()
	c.AccessTokenDuration = time.Duration(1) * time.Hour
	proxy := newFakeProxy(c).proxy
	duration := proxy.getAccessCookieExpiration(token, refreshToken)
	assert.True(t, duration > 0, "duration should be positive")
	assertAlmostEquals(t, c.AccessTokenDuration, duration)
}

func TestGetAccessCookieExpiration_PastExp(t *testing.T) {
	ft := newTestToken("foo")
	ft.setExpiration(time.Now().AddDate(-1, 0, 0))
	token := ft.getToken()
	refreshToken := token.Encode()
	c := newFakeKeycloakConfig()
	c.AccessTokenDuration = time.Duration(1) * time.Hour
	proxy := newFakeProxy(c).proxy
	duration := proxy.getAccessCookieExpiration(token, refreshToken)
	assertAlmostEquals(t, c.AccessTokenDuration, duration)
}

func TestGetAccessCookieExpiration_ValidExp(t *testing.T) {
	ft := newTestToken("foo")
	token := ft.getToken()
	refreshToken := token.Encode()
	c := newFakeKeycloakConfig()
	c.AccessTokenDuration = time.Duration(1) * time.Hour
	proxy := newFakeProxy(c).proxy
	duration := proxy.getAccessCookieExpiration(token, refreshToken)
	val, ok, _ := ft.claims.TimeClaim("exp")
	assert.True(t, ok)
	expectedDuration := time.Until(val)
	assertAlmostEquals(t, expectedDuration, duration)
}
