package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	rgSessions "github.com/itslaves/rentalgames-server/common/sessions"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	Vendor       = "vendor"
	UserID       = "userID"
	Nickname     = "nickname"
	ProfileImage = "profileImage"
	Gender       = "gender"
	Email        = "email"
	AccessToken  = "accessToken"
	RefreshToken = "refreshToken"
	Expiry       = "expiry"
	State        = "state"
)

type UserProfile struct {
	UserID       string `json:"userID"`
	Nickname     string `json:"nickname"`
	ProfileImage string `json:"profileImage"`
	Gender       string `json:"gender"`
	Email        string `json:"email"`
}

type UserAuth struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	Expiry       time.Time `json:"expiry"`
}

type SessionValues struct {
	UserProfile
	UserAuth
	Vendor string `json:"vendor"`
	State  string `json:"state,omitempty"`
}

func OAuthVendorUrls(ctx *gin.Context) {
	session := rgSessions.Session(ctx)
	state := randomToken()

	session.Values[State] = state
	if err := session.Save(ctx.Request, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, []map[string]interface{}{
		gin.H{
			"vendor": "kakao",
			"url":    KakaoOAuthConfig().AuthCodeURL(state),
			"image":  "https://developers.kakao.com/assets/img/about/logos/kakaologin/kr/kakao_account_login_btn_large_narrow.png",
		},
		gin.H{
			"vendor": "naver",
			"url":    NaverOAuthConfig().AuthCodeURL(state),
			"image":  "https://static.nid.naver.com/oauth/big_g.PNG",
		},
		gin.H{
			"vendor": "google",
			"url":    GoogleOAuthConfig().AuthCodeURL(state),
			"image":  "https://developers.google.com/identity/images/btn_google_signin_dark_normal_web.png?hl=ko",
		},
	})
}

func randomToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

type UserProfileParser func([]byte) *UserProfile

type callbackHandler struct {
	webAddr           string
	vendor            string
	session           *sessions.Session
	config            *oauth2.Config
	ctx               *gin.Context
	userProfileParser UserProfileParser
	token             *oauth2.Token
	userProfile       *UserProfile
}

func NewCallbackHandler(vendor string, session *sessions.Session, config *oauth2.Config, ctx *gin.Context, userProfileParser UserProfileParser) *callbackHandler {
	return &callbackHandler{
		webAddr:           viper.GetString("web.addr"),
		vendor:            vendor,
		session:           session,
		config:            config,
		ctx:               ctx,
		userProfileParser: userProfileParser,
		token:             nil,
		userProfile:       nil,
	}
}

func (h *callbackHandler) okPage() string {
	return fmt.Sprintf("http://%s/", h.webAddr)
}

func (h *callbackHandler) errPage() string {
	return fmt.Sprintf("http://%s/error", h.webAddr)
}

func (h *callbackHandler) validateState() error {
	sessionState := h.session.Values[State]
	callbackState := h.ctx.Query(State)
	if sessionState != callbackState {
		return errors.New("state mismatched")
	}
	return nil
}

func (h *callbackHandler) loadToken() error {
	token, err := h.config.Exchange(context.TODO(), h.ctx.Query("code"))
	if err != nil {
		return err
	}
	if token == nil || !token.Valid() {
		return errors.New("invalid token")
	}
	h.token = token
	return nil
}

func (h *callbackHandler) loadUserProfile() error {
	client := h.config.Client(context.TODO(), h.token)
	userProfileURL := fmt.Sprintf("oauth.%s.userProfileURL", h.vendor)
	resp, err := client.Get(viper.GetString(userProfileURL))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	h.userProfile = h.userProfileParser(result)
	return nil
}

func (h *callbackHandler) writeSession() error {
	if h.userProfile == nil {
		return errors.New("userProfile does not exist")
	}
	if h.token == nil {
		return errors.New("token does not exist")
	}
	h.session.Values[Vendor] = h.vendor
	h.session.Values[UserID] = h.userProfile.UserID
	h.session.Values[Nickname] = h.userProfile.Nickname
	h.session.Values[ProfileImage] = h.userProfile.ProfileImage
	h.session.Values[Gender] = h.userProfile.Gender
	h.session.Values[Email] = h.userProfile.Email
	h.session.Values[AccessToken] = h.token.AccessToken
	h.session.Values[RefreshToken] = h.token.RefreshToken
	h.session.Values[Expiry] = h.token.Expiry
	if err := h.session.Save(h.ctx.Request, h.ctx.Writer); err != nil {
		return err
	}
	return nil
}

func (h *callbackHandler) Handle() {
	if err := h.validateState(); err != nil {
		fmt.Fprint(os.Stderr, err)
		h.ctx.Redirect(http.StatusPermanentRedirect, h.errPage())
	}
	if err := h.loadToken(); err != nil {
		fmt.Fprint(os.Stderr, err)
		h.ctx.Redirect(http.StatusPermanentRedirect, h.errPage())
	}
	if err := h.loadUserProfile(); err != nil {
		fmt.Fprint(os.Stderr, err)
		h.ctx.Redirect(http.StatusPermanentRedirect, h.errPage())
	}
	if err := h.writeSession(); err != nil {
		fmt.Fprint(os.Stderr, err)
		h.ctx.Redirect(http.StatusPermanentRedirect, h.errPage())
	}
	h.ctx.Redirect(http.StatusPermanentRedirect, h.okPage())
}
