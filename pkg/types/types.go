package types

import (
	"errors"
	"strconv"
)

const (
	RadAcctMetricName  = "radius_acct"
	RadAcctMetricHelp  = "radius accounting packet"
	UserIDMetricName   = "user_id_req"
	UserIDMetricHelp   = "user_id request"
	QueueLenMetricName = "queue_len"
	QueueLenMetricHelp = "length queue"
)

var (
	// ErrSQLListen ...
	ErrSQLListen = errors.New("SQL listen")
	// ErrSQLListenTimeout ...
	ErrSQLListenTimeout = errors.New("SQL listen i/o timeout")
	// ErrHTTPHandler ...
	ErrHTTPHandler = errors.New("HTTP handler")
	// ErrSQLProcess ...
	ErrSQLProcess = errors.New("SQL process")
)

const (
	// AcctStatusTypeStart ...
	AcctStatusTypeStart = "Start"
	// AcctStatusTypeInterimUpdate ...
	AcctStatusTypeInterimUpdate = "Interim-Update"
	// AcctStatusTypeStop ...
	AcctStatusTypeStop = "Stop"
)

// RadiusAccounting ...
type RadiusAccounting struct {
	AcctStatusType  string `json:"Acct-Status-Type"`
	NASIPAddress    string `json:"NAS-IP-Address"`
	UserName        string `json:"User-Name"`
	AcctSessionTime string `json:"Acct-Session-Time"`

	FramedIPAddress  string `json:"Framed-IP-Address"`
	AcctInputOctets  string `json:"Acct-Input-Octets"`
	AcctOutputOctets string `json:"Acct-Output-Octets"`
	AcctSessionID    string `json:"Acct-Session-Id"`
	NASPortID        string `json:"NAS-Port-Id"`

	CalledStationID  string `json:"Called-Station-Id"`
	CallingStationID string `json:"Calling-Station-Id"`
	H323CallOrigin   string `json:"h323-call-origin"`
	H323ConfID       string `json:"h323-conf-id"`
	H323SetupTime    string `json:"h323-setup-time"`
	H323BillingModel string `json:"h323-billing-model"`

	// ?
	LoginHost string `json:"Login-Host"`

	// ????
	AsteriskSrc       string `json:"Asterisk-Src"`
	AsteriskDst       string `json:"Asterisk-Dst"`
	AsteriskStartTime string `json:"Asterisk-Start-Time"`
	AsteriskBillSec   string `json:"Asterisk-Bill-Sec"`
}

func (s RadiusAccounting) String() string {
	return s.AcctSessionID
}

// UserID ...
type UserID struct {
	UserID int `json:"user_id"`
}

func (s UserID) String() string {
	return strconv.Itoa(s.UserID)
}

// Config ...
type Config struct {
	NodeName string `required:"true"`

	DBURI           string `required:"true"`
	RadAcctSQLQuery string `required:"true"`
	UserIDSQLQuery  string `required:"true"`
	SQLListenChan   string `required:"true"`

	RadHTTPListen     string `default:":1234"`
	RadHTTPToken      string `required:"true"`
	RadHTTPHeader     string `default:"X-Auth"`
	RadiusHTTPAcctURL string `default:"/acct"`

	AcctCntWorker   int `default:"4"`
	UserIDCntWorker int `default:"4"`

	NodeCheckURL string `default:"/hb"`

	PrometheusListen string `default:":9090"`
	PrometheusPath   string `default:"/metrics"`

	DebugLevel bool `default:"false"`
}
