package jwt_test

// Example HTTP auth using asymmetric crypto/RSA keys
// This is based on a (now outdated) example at https://gist.github.com/cryptix/45c33ecf0ae54828e63b

import (
	"crypto/rsa"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
)

// location of the files used for signing and verification
const (
	privKeyPath = "test/sample_key"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = "test/sample_key.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

var (
	verifyKey  *rsa.PublicKey
	signKey    *rsa.PrivateKey
	serverPort int
)

// read the key files before starting http handlers
func init() {
	signBytes, err := os.ReadFile(privKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := os.ReadFile(pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)

	http.HandleFunc("/authenticate", authHandler)
	http.HandleFunc("/restricted", restrictedHandler)

	// Setup listener
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{})
	fatal(err)
	serverPort = listener.Addr().(*net.TCPAddr).Port

	log.Println("Listening...")
	go func() {
		fatal(http.Serve(listener, nil))
	}()
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Define some custom types were going to use within our tokens
type CustomerInfo struct {
	Name string
	Kind string
}

type CustomClaimsExample struct {
	jwt.RegisteredClaims
	TokenType string
	CustomerInfo
}

func Example_getTokenViaHTTP() {
	// See func authHandler for an example auth handler that produces a token
	res, err := http.PostForm(fmt.Sprintf("http://localhost:%v/authenticate", serverPort), url.Values{
		"user": {"test"},
		"pass": {"known"},
	})
	fatal(err)

	if res.StatusCode != 200 {
		fmt.Println("Unexpected status code", res.StatusCode)
	}

	// Read the token out of the response body
	buf, err := io.ReadAll(res.Body)
	fatal(err)
	_ = res.Body.Close()
	tokenString := strings.TrimSpace(string(buf))

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaimsExample{}, func(token *jwt.Token) (any, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return verifyKey, nil
	})
	fatal(err)

	claims := token.Claims.(*CustomClaimsExample)
	fmt.Println(claims.Name)

	// Output: test
}

func Example_useTokenViaHTTP() {
	// Make a sample token
	// In a real world situation, this token will have been acquired from
	// some other API call (see Example_getTokenViaHTTP)
	token, err := createToken("foo")
	fatal(err)

	// Make request.  See func restrictedHandler for example request processor
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%v/restricted", serverPort), nil)
	fatal(err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	res, err := http.DefaultClient.Do(req)
	fatal(err)

	// Read the response body
	buf, err := io.ReadAll(res.Body)
	fatal(err)
	_ = res.Body.Close()
	fmt.Printf("%s", buf)

	// Output: Welcome, foo
}

func createToken(user string) (string, error) {
	// create a signer for rsa 256
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	// set our claims
	t.Claims = &CustomClaimsExample{
		jwt.RegisteredClaims{
			// set the expire time
			// see https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.4
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)),
		},
		"level1",
		CustomerInfo{user, "human"},
	}

	// Creat token string
	return t.SignedString(signKey)
}

// reads the form values, checks them and creates the token
func authHandler(w http.ResponseWriter, r *http.Request) {
	// make sure its post
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintln(w, "No POST", r.Method)
		return
	}

	user := r.FormValue("user")
	pass := r.FormValue("pass")

	log.Printf("Authenticate: user[%s] pass[%s]\n", user, pass)

	// check values
	if user != "test" || pass != "known" {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintln(w, "Wrong info")
		return
	}

	tokenString, err := createToken(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintln(w, "Sorry, error while Signing Token!")
		log.Printf("Token Signing error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/jwt")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(w, tokenString)
}

// only accessible with a valid token
func restrictedHandler(w http.ResponseWriter, r *http.Request) {
	// Get token from request
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (any, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return verifyKey, nil
	}, request.WithClaims(&CustomClaimsExample{}))

	// If the token is missing or invalid, return error
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintln(w, "Invalid token:", err)
		return
	}

	// Token is valid
	_, _ = fmt.Fprintln(w, "Welcome,", token.Claims.(*CustomClaimsExample).Name)
}
