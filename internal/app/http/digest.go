package http

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/0x131315/hikvision-backup/internal/app/config"
)

type digestChallenge struct {
	qop   string
	realm string
	nonce string
}

type digestContext struct {
	digest digestChallenge
	conf   config.Config
	cnonce string
	nc     int
}

type Digest struct {
	ctx *digestContext
	rnd *rand.Rand
}

func NewDigest(conf config.Config) *Digest {
	rndg := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Digest{ctx: &digestContext{conf: conf}, rnd: rndg}
}

func (d *Digest) updateDigest(digestHeader string) {
	d.ctx.nc = 0
	d.ctx.digest = parseWWWAuthenticate(digestHeader)
}

func (d *Digest) getNextAuthHeader(method, uri string) string {
	d.ctx.nc++
	d.ctx.cnonce = d.randomCnonce()

	return d.buildDigestAuth(method, uri)
}

func (d *Digest) buildDigestAuth(method, uri string) string {
	username := d.ctx.conf.User
	password := d.ctx.conf.Pass
	realm := d.ctx.digest.realm
	nonce := d.ctx.digest.nonce
	qop := d.ctx.digest.qop
	nc := d.ctx.nc
	cnonce := d.ctx.cnonce

	ha1 := md5Hex(fmt.Sprintf("%s:%s:%s", username, realm, password))
	ha2 := md5Hex(fmt.Sprintf("%s:%s", method, uri))
	ncStr := fmt.Sprintf("%08x", nc)
	response := md5Hex(fmt.Sprintf("%s:%s:%s:%s:%s:%s", ha1, nonce, ncStr, cnonce, qop, ha2))

	authHeader := fmt.Sprintf(
		`Digest username="%s", realm="%s", nonce="%s", uri="%s", qop=%s, nc=%s, cnonce="%s", response="%s"`,
		username, realm, nonce, uri, qop, ncStr, cnonce, response,
	)

	return authHeader
}

func parseWWWAuthenticate(header string) digestChallenge {
	if !strings.HasPrefix(header, "Digest ") {
		slog.Error("Expected 'Digest' header")
		os.Exit(1)
	}

	result := digestChallenge{}
	header = strings.TrimPrefix(header, "Digest ")
	for _, part := range strings.Split(header, ",") {
		if kv := strings.SplitN(strings.TrimSpace(part), "=", 2); len(kv) == 2 {
			key := kv[0]
			val := strings.Trim(kv[1], `"`)
			if val == "" {
				slog.Error("Empty digest value", "key", key)
				os.Exit(1)
			}

			switch key {
			case "realm":
				result.realm = val
			case "nonce":
				result.nonce = val
			case "qop":
				result.qop = val
			}
		}
	}
	return result
}

func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func (d *Digest) randomCnonce() string {
	const letters = "abcdef0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[d.rnd.Intn(len(letters))]
	}
	return string(b)
}
