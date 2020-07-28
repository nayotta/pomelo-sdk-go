package pomelo

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"path"

	"github.com/levigross/grequests"
)

type PomeloSDKOption struct {
	BaseURL            string
	InsecureSkipVerify bool
}

type PomeloSDK interface {
	SendSMS(ctx context.Context, id string, numbers []string, args map[string]string) error
}

func NewPomeloSDK(opt *PomeloSDKOption) (PomeloSDK, error) {
	return &PomeloSDKImpl{
		opt: opt,
	}, nil
}

type PomeloSDKImpl struct {
	opt *PomeloSDKOption
}

func (sdk *PomeloSDKImpl) SendSMS(ctx context.Context, id string, numbers []string, args map[string]string) error {
	u, err := url.Parse(sdk.opt.BaseURL)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "sms", id)
	body := map[string]interface{}{
		"phoneNumberSet": numbers,
		"arguments":      args,
	}
	res, err := grequests.Post(u.String(), &grequests.RequestOptions{
		JSON:               body,
		InsecureSkipVerify: sdk.opt.InsecureSkipVerify,
	})
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(string(res.Bytes()))
	}

	return nil
}
