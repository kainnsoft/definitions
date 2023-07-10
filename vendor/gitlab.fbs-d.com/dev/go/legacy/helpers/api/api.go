package api

import (
	"strings"

	"github.com/satori/go.uuid"
	"gitlab.fbs-d.com/dev/go/legacy/exceptions"

	"gitlab.fbs-d.com/dev/go/legacy/helpers"
)

const (
	SysToken = "sys"

	SourceSecured  = "Secured"  // Request from API via OAuth
	SourcePrivate  = "Private"  // Request from API via internal nginx
	SourcePublic   = "Public"   // Request from API without authorization
	SourceInternal = "Internal" // Internal request between services via RPC

	PlatformAndroid = "android_app"
	PlatformIos     = "ios_app"
	PlatformWeb     = "website"
)

type Client struct {
	Ip            string     `json:"ip"`
	Host          string     `json:"host"`
	Domain        string     `json:"domain"`
	UserAgent     string     `json:"userAgent"`
	Language      string     `json:"language"`
	Location      string     `json:"location"`
	UserId        int64      `json:"userId"`
	Platform      string     `json:"platform"`
	Version       int64      `json:"version"`
	Source        string     `json:"source"`
	ApplicationId int64      `json:"applicationId"`
	Analytics     Analytics  `json:"analytics"`
	Role          string     `json:"role"`
	Utm           *Utm       `json:"utm"`
	GoogleIds     *GoogleIds `json:"googleIds"`
	FBIds         *FbIds     `json:"fbIds"`
}

type Analytics struct {
	AppsFlyerId           string `json:"appsFlyerId"`
	FireBaseId            string `json:"fireBaseId"`
	FireBaseAppInstanceId string `json:"fireBaseAppInstanceId"`
	AdvertisingId         string `json:"advertisingId"`
	VendorId              string `json:"vendorId"`
	SessionId             string `json:"sessionId"`
}

type Utm struct {
	Campaign string `json:"campaign"`
	Source   string `json:"source"`
	Medium   string `json:"medium"`
	Term     string `json:"term"`
	Content  string `json:"content"`
}

type GoogleIds struct {
	GclId string `json:"gclid"`
	DclId string `json:"dclid"`
}

type FbIds struct {
	FbClickId   string `json:"fbClickId"`
	FbBrowserId string `json:"fbBrowserId"`
}

type ApiRequest struct {
	Token  string                 `json:"token" validate:"nonzero"`
	Client Client                 `json:"client"`
	Body   map[string]interface{} `json:"body"`
}

// Deprecated: use GenerateTokenV2 instead
func GenerateToken() (token string) {
	uid, _ := uuid.NewV4()
	return strings.Replace(uid.String(), "-", "", -1)
}

func GenerateTokenV2() (token string, ex excp.IException) {
	uid, ex := toolkit.Uuid()
	if ex != nil {
		return
	}
	token = strings.Replace(uid, "-", "", -1)
	return
}

type OptionalInt64 struct {
	Value int64 `json:"value"`
}
