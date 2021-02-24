package sessions

import (
	"encoding/base32"
	"encoding/gob"
	"fmt"
	"net/http"
	"strings"
	"time"

	rgRedis "github.com/skyoo2003/rentalgames-server/internal/third_party/redis"
	"github.com/spf13/viper"

	"github.com/go-redis/redis"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func init() {
	// gob encoder가 time.Time 타입을 인식할 수 있도록 등록
	gob.Register(time.Time{})
}

type RedisStore struct {
	Codecs    []securecookie.Codec
	Options   *sessions.Options
	keyPrefix string
	client    redis.UniversalClient
}

func NewRedisStore() *RedisStore {
	hashKey := []byte(viper.GetString("session_hash_key"))
	blockKey := []byte(viper.GetString("session_block_key"))
	maxAge := viper.GetInt("session.maxAge")
	domain := viper.GetString("session.domain")

	rs := &RedisStore{
		Codecs: securecookie.CodecsFromPairs(hashKey, blockKey),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: maxAge,
			Domain: domain,
		},
		keyPrefix: "session",
		client:    rgRedis.Client(),
	}
	rs.MaxAge(rs.Options.MaxAge)
	rs.SetSerializer(securecookie.GobEncoder{})
	return rs
}

func (rs *RedisStore) sessionKey(s *sessions.Session) string {
	return fmt.Sprintf("%s.%s", rs.keyPrefix, s.ID)
}

func (rs *RedisStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(rs, name)
}

func (rs *RedisStore) New(r *http.Request, name string) (*sessions.Session, error) {
	s := sessions.NewSession(rs, name)
	opts := *rs.Options
	s.Options = &opts
	s.IsNew = true

	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = securecookie.DecodeMulti(name, c.Value, &s.ID, rs.Codecs...)
		if err == nil {
			err = rs.load(s)
			if err == nil {
				s.IsNew = false
			}
		}
	}
	return s, err
}

func (rs *RedisStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	// Delete if max-age is <= 0
	if s.Options.MaxAge <= 0 {
		if err := rs.erase(s); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(s.Name(), "", s.Options))
		return nil
	}

	if s.ID == "" {
		// Because the ID is used in the redis, encode it to use alphanumeric characters only.
		s.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
	}

	if err := rs.save(s); err != nil {
		return err
	}

	encoded, err := securecookie.EncodeMulti(s.Name(), s.ID, rs.Codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(s.Name(), encoded, s.Options))
	return nil
}

func (rs *RedisStore) MaxAge(age int) {
	rs.Options.MaxAge = age
	for _, codec := range rs.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

func (rs *RedisStore) SetSerializer(serializer securecookie.Serializer) {
	for _, codec := range rs.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.SetSerializer(serializer)
		}
	}
}

func (rs *RedisStore) save(s *sessions.Session) error {
	encoded, err := securecookie.EncodeMulti(s.Name(), s.Values, rs.Codecs...)
	if err != nil {
		return err
	}
	age := time.Duration(s.Options.MaxAge) * time.Second
	if err := rs.client.Set(rs.sessionKey(s), encoded, age).Err(); err != nil {
		return err
	}
	return nil
}

func (rs *RedisStore) load(s *sessions.Session) error {
	r := rs.client.Get(rs.sessionKey(s))
	if r.Err() != nil {
		return r.Err()
	}
	if data, err := r.Result(); err == nil {
		err = securecookie.DecodeMulti(s.Name(), data, &s.Values, rs.Codecs...)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (rs *RedisStore) erase(s *sessions.Session) error {
	if err := rs.client.Del(rs.sessionKey(s)).Err(); err != nil {
		return err
	}
	return nil
}
