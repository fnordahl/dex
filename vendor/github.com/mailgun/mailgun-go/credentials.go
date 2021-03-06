package mailgun

import (
	"fmt"
	"github.com/mbanzon/simplehttp"
	"strconv"
)

// A Credential structure describes a principle allowed to send or receive mail at the domain.
type Credential struct {
	CreatedAt string `json:"created_at"`
	Login     string `json:"login"`
	Password  string `json:"password"`
}

// ErrEmptyParam results occur when a required parameter is missing.
var ErrEmptyParam = fmt.Errorf("empty or illegal parameter")

// GetCredentials returns the (possibly zero-length) list of credentials associated with your domain.
func (mg *MailgunImpl) GetCredentials(limit, skip int) (int, []Credential, error) {
	r := simplehttp.NewHTTPRequest(generateCredentialsUrl(mg, ""))
	if limit != DefaultLimit {
		r.AddParameter("limit", strconv.Itoa(limit))
	}
	if skip != DefaultSkip {
		r.AddParameter("skip", strconv.Itoa(skip))
	}
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	var envelope struct {
		TotalCount int          `json:"total_count"`
		Items      []Credential `json:"items"`
	}
	err := getResponseFromJSON(r, &envelope)
	if err != nil {
		return -1, nil, err
	}
	return envelope.TotalCount, envelope.Items, nil
}

// CreateCredential attempts to create associate a new principle with your domain.
func (mg *MailgunImpl) CreateCredential(login, password string) error {
	if (login == "") || (password == "") {
		return ErrEmptyParam
	}
	r := simplehttp.NewHTTPRequest(generateCredentialsUrl(mg, ""))
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	p := simplehttp.NewUrlEncodedPayload()
	p.AddValue("login", login)
	p.AddValue("password", password)
	_, err := makePostRequest(r, p)
	return err
}

// ChangeCredentialPassword attempts to alter the indicated credential's password.
func (mg *MailgunImpl) ChangeCredentialPassword(id, password string) error {
	if (id == "") || (password == "") {
		return ErrEmptyParam
	}
	r := simplehttp.NewHTTPRequest(generateCredentialsUrl(mg, id))
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	p := simplehttp.NewUrlEncodedPayload()
	p.AddValue("password", password)
	_, err := makePutRequest(r, p)
	return err
}

// DeleteCredential attempts to remove the indicated principle from the domain.
func (mg *MailgunImpl) DeleteCredential(id string) error {
	if id == "" {
		return ErrEmptyParam
	}
	r := simplehttp.NewHTTPRequest(generateCredentialsUrl(mg, id))
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	_, err := makeDeleteRequest(r)
	return err
}
