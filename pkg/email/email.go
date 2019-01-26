package email

import (
	"bytes"
	"context"
	"strconv"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/tasks"
	"github.com/dchenk/mazewire/pkg/util"
)

// A Sender is able to send emails.
// The methods of a Sender do not send emails synchronously: they indicate merely whether an Email is valid and
// was scheduled to be sent out.
type Sender interface {
	Send(*Email) error
	SendMulti([]Email) error
}

// An Email represents all the basic settings we use to send emails from the server.
type Email struct {
	To       Party     `json:"to"`
	From     Party     `json:"from"`
	ReplyTo  Party     `json:"reply_to"`
	Subject  string    `json:"subject"`
	Contents []Content `json:"contents"`
}

// NewEmail constructs an Email, setting only the fields that are absolutely required.
// The mime type of the content is set to text/plain.
func NewEmail(toAddress, message string) *Email {
	return &Email{
		To:       Party{Address: toAddress},
		Contents: []Content{{Type: util.ContentTypeTextPlain, Value: message}},
	}
}

// MarshalJSON implements json.Marshaler. The error is always nil.
func (em *Email) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteString(`{"to":`)
	to, _ := em.To.MarshalJSON()
	b.Write(to)
	b.WriteString(`,"from":`)
	from, _ := em.From.MarshalJSON()
	b.Write(from)
	b.WriteString(`,"reply_to":`)
	replyTo, _ := em.ReplyTo.MarshalJSON()
	b.Write(replyTo)
	b.WriteString(`,"subject":`)
	b.WriteString(strconv.Quote(em.Subject))
	b.WriteString(`,"contents":[`)
	for i := range em.Contents {
		t, _ := em.Contents[i].MarshalJSON()
		b.Write(t)
	}
	b.WriteString("]}")
	return b.Bytes(), nil
}

// A Party is a sender or a receiver of an email.
type Party struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// MarshalJSON implements json.Marshaler. The error is always nil.
func (p *Party) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteString(`{"name":`)
	b.WriteString(strconv.Quote(p.Name))
	b.WriteString(`,"address":`)
	b.WriteString(strconv.Quote(p.Address))
	b.WriteByte('}')
	return b.Bytes(), nil
}

// A Copies contains additional recipients to be copied on an email.
type Copies struct {
	CC  []Party `json:"cc"`
	BCC []Party `json:"bcc"`
}

// MarshalJSON implements json.Marshaler. The error is always nil.
func (c Copies) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteString(`{"cc":[`)
	var t []byte
	for i := range c.CC {
		t, _ = c.CC[i].MarshalJSON()
		b.Write(t)
	}
	b.WriteString(`],"bcc":[`)
	for i := range c.BCC {
		t, _ = c.CC[i].MarshalJSON()
		b.Write(t)
	}
	b.WriteString("]}")
	return b.Bytes(), nil
}

// A Content is piece of the content of an email body. Each email must have at least one Content.
type Content struct {
	// Type is the mime type of the content (for example, "text/plain").
	Type string `json:"type"`

	// Value is the content data.
	Value string `json:"value"`
}

// MarshalJSON implements json.Marshaler. The error is always nil.
func (c *Content) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteString(`{"type":`)
	b.WriteString(strconv.Quote(c.Type))
	b.WriteString(`,"value":`)
	b.WriteString(strconv.Quote(c.Value))
	b.WriteByte('}')
	return b.Bytes(), nil
}

// An Attachment is an email attachment.
type Attachment struct {
	// Content embeds the attachment data. The value must be base64-encoded.
	Content `json:"content"`

	// Filename is the name of the file to show with the attachment.
	Filename string `json:"filename"`

	// Disposition specifies how the attachment is included in the email.
	// A value may be either "inline" or "attachment"; the default, if this field
	// is left blank, is "attachment".
	Disposition string `json:"disposition"`
}

// ScheduleSend schedules an email to be sent asynchronously by the default task queue.
// This function uses the default Sender configured for the website.
func ScheduleSend(ctx context.Context, s *data.Site, m Email) error {
	d, _ := m.MarshalJSON()
	return tasks.New(ctx, tasks.SendEmail, d)
}

// ScheduleSendCustom schedules an email to be sent asynchronously by the default task queue using
// the given Sender.
func ScheduleSendCustom(ctx context.Context, s *data.Site, snd Sender, m Email) error {
	return nil // TODO
}
