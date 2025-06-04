package http

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/0x131315/hikvision-backup/internal/app/config"
	"github.com/0x131315/hikvision-backup/internal/app/util"
	"math/rand"
	"strings"
	"time"
)

type digestChallenge struct {
	qop   string
	realm string
	nonce string
}

type context struct {
	digest digestChallenge
	conf   config.Config
	cnonce string
	nc     int
}

var ctx *context
var rnd *rand.Rand

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	ctx = &context{}
	ctx.conf = config.Get()
}

func updateDigest(digestHeader string) {
	ctx.nc = 0
	ctx.digest = parseWWWAuthenticate(digestHeader)
}

func getNextAuthHeader(method, uri string) string {
	ctx.nc++
	ctx.cnonce = randomCnonce()

	return buildDigestAuth(ctx, method, uri)
}

func buildDigestAuth(ctx *context, method, uri string) string {
	username := ctx.conf.User
	password := ctx.conf.Pass
	realm := ctx.digest.realm
	nonce := ctx.digest.nonce
	qop := ctx.digest.qop
	nc := ctx.nc
	cnonce := ctx.cnonce

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
		util.FatalError("Expected 'Digest' header")
	}

	result := digestChallenge{}
	header = strings.TrimPrefix(header, "Digest ")
	for _, part := range strings.Split(header, ",") {
		if kv := strings.SplitN(strings.TrimSpace(part), "=", 2); len(kv) == 2 {
			key := kv[0]
			val := strings.Trim(kv[1], `"`)
			if val == "" {
				util.FatalError(fmt.Sprintf("Empty digest key: %s", key))
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

func randomCnonce() string {
	const letters = "abcdef0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[rnd.Intn(len(letters))]
	}
	return string(b)
}
