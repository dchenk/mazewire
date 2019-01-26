package sendgrid

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dchenk/mazewire/pkg/email"
	"github.com/dchenk/mazewire/pkg/log"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Endpoint is the SendGrid API endpoint this package uses to send emails.
const Endpoint = "/v3/mail/send"

// Sender implements the email.Sender interface to send emails using the SendGrid API.
type Sender struct {
	apiKey  string
	headers map[string]string
}

// NewSender creates a new Sender given an API key.
func NewSender(APIKey string) *Sender {
	APIKey = strings.TrimSpace(APIKey)
	if APIKey == "" {
		msg := "sendgrid: blank API key given"
		log.Critical(nil, msg, fmt.Errorf("invalid arguments"))
		panic(msg)
	}
	return &Sender{apiKey: APIKey}
}

// AddHeader adds a header to use in the emails.
func (s *Sender) AddHeader(key, val string) {
	if s.headers == nil {
		s.headers = make(map[string]string, 2)
	}
	s.headers[key] = val
}

// newSendReq constructs a POST request initialized for the SendGrid email sending API.
func newSendReq(APIKey string) rest.Request {
	r := sendgrid.GetRequest(APIKey, Endpoint, "")
	r.Method = "POST"
	return r
}

// rateLimitSleep is the default sleep time when retrying after a failure to send due to a rate limit.
const rateLimitSleep = 600 * time.Millisecond

// sgMailFromEmail constructs a SGMailV3 out of an email.Email.
func sgMailFromEmail(e *email.Email) *mail.SGMailV3 {
	m := mail.SGMailV3{
		From:    &mail.Email{Name: e.From.Name, Address: e.From.Address},
		Subject: e.Subject,
		// TODO
	}
	return &m
}

// makeRequestRetry makes a synchronous request but retries up to three times in the event of a rate limit.
// This function copies the basic logic of sendgrid.MakeRequestRetry but uses our http.Client from the
// context for each request.
func makeRequestRetry(ctx context.Context, request *rest.Request, email *mail.SGMailV3) (*rest.Response, error) {

	client := sendgrid.Client{Request: *request}

	var retry uint8
	for {

		response, err := client.Send(email)
		if err != nil {
			return nil, err
		}

		if response.StatusCode != http.StatusTooManyRequests {
			return response, nil
		}

		if retry > 3 {
			return nil, errors.New("rate limit retry exceeded")
		}
		retry++

		now := time.Now()
		resetTime := now.Add(rateLimitSleep)

		reset, ok := response.Headers["X-RateLimit-Reset"]
		if ok && len(reset) > 0 {
			t, err := strconv.ParseInt(reset[0], 10, 64)
			if err == nil {
				resetTime = time.Unix(t, 0)
			}
		}
		time.Sleep(resetTime.Sub(now))

	}

}

// According the SendGrid docs, either a 200 or a 202 status indicates that the request succeeded,
// thought the 200 status is used only in sandbox mode.
// For more than just sending email, consider also 201 and 204:
// https://sendgrid.com/docs/API_Reference/Web_API_v3/How_To_Use_The_Web_API_v3/responses.html
func statusOK(code int) bool {
	return code == http.StatusOK || code == http.StatusAccepted
}
