package handler

// import (
// 	"bytes"
// 	"net/http"
// 	"net/http/httptest"
// 	"reflect"
// 	"testing"

// 	"github.com/balugcath/william/pkg/types"
// )

// var (
// 	pckt1 = `{
// 	"NAS-IP-Address":"9.9.9.9",
// 	"Acct-Status-Type":"Start",
// 	"Acct-Session-Id":"00123553-94463159-AV1",
// 	"Cisco-AVPair":"xpgk-record-id=00123553-94463159",
// 	"Cisco-AVPair":"xpgk-local-src-signaling-address=9.9.9.9",
// 	"Cisco-AVPair":"xpgk-src-number-out=1234567890",
// 	"NAS-Port-Type":"Async",
// 	"Cisco-AVPair":"h323-remote-id=MVTS-II",
// 	"h323-conf-id":"h323-conf-id=14392814 011511EB ABB6D8D3 85E1BACC",
// 	"Service-Type":"Login-User",
// 	"Cisco-AVPair":"h323-incoming-call-id=72E4CA60 1643DACF 29B62268 0AF88524",
// 	"Cisco-AVPair":"h323-incoming-conf-id=145D3402 011511EB ABB6D8D3 85E1BACC",
// 	"Cisco-AVPair":"h323-call-id=14392814 011511EB ABB6D8D3 85E1BACC",
// 	"h323-call-origin":"h323-call-origin=answer",
// 	"Cisco-AVPair":"xpgk-dst-number-in=1234567890",
// 	"Cisco-AVPair":"xpgk-dst-number-out=001234567890",
// 	"Acct-Delay-Time":"0",
// 	"Framed-IP-Address":"1.2.5.2",
// 	"h323-remote-address":"h323-remote-address=9.9.9.9",
// 	"h323-gw-id":"h323-gw-id=DeltaDigital_gw",
// 	"Cisco-AVPair":"h323-gw-address=1.0.7.1",
// 	"h323-connect-time":"h323-connect-time=23:00:21.017 UTC Sun Sep 27 2020",
// 	"h323-setup-time":"h323-setup-time=22:59:21.252 UTC Sun Sep 27 2020",
// 	"h323-call-type":"h323-call-type=VoIP",
// 	"Calling-Station-Id":"1234567890",
// 	"User-Name":"1234567890",
// 	"Called-Station-Id":"1234567890",
// 	"Cisco-AVPair":"xpgk-route-retries=0",
// 	"Cisco-AVPair":"xpgk-src-number-in=1234567890",
// 	"Event-Timestamp":"Sep 28 2020 02:00:21 MSK",
// 	"Tmp-String-9":"ai:",
// 	"Acct-Input-Octets":"806140",
// 	"Acct-Output-Octets":"3268464",
// 	"Acct-Unique-Session-Id":"3d90e499e4a1ff416e79691f40e7b064",
// 	"Timestamp":1601247621,
// 	"h323-billing-model": "fmc_megafon",
// 	"NAS-Port-Id" : "0/0/1/115"
// 	}	
// `
// 	r1 = types.RadiusAccounting{
// 		AcctStatusType:  types.AcctStatusTypeStart,
// 		NASIPAddress:    "9.9.9.9",
// 		UserName:        "1234567890",
// 		AcctSessionTime: "0",

// 		FramedIPAddress:  "1.2.5.2",
// 		AcctInputOctets:  "806140",
// 		AcctOutputOctets: "3268464",
// 		AcctSessionID:    "00123553-94463159-AV1",
// 		NASPortID:        "0/0/1/115",

// 		CalledStationID:  "1234567890",
// 		CallingStationID: "1234567890",
// 		H323CallOrigin:   "h323-call-origin=answer",
// 		H323ConfID:       "h323-conf-id=14392814 011511EB ABB6D8D3 85E1BACC",
// 		H323SetupTime:    "h323-setup-time=22:59:21.252 UTC Sun Sep 27 2020",
// 		H323BillingModel: "fmc_megafon",
// 	}
// )

// func TestHTTPHandler_radAcct(t *testing.T) {
// 	type args struct {
// 		body   string
// 		token  string
// 		header string
// 	}
// 	type fields struct {
// 		token  string
// 		queuer *queueMock
// 	}
// 	tests := []struct {
// 		name       string
// 		fields     fields
// 		args       args
// 		wantStatus int
// 		wantStruct types.RadiusAccounting
// 	}{
// 		{
// 			name:       "test1",
// 			fields:     fields{token: "123", queuer: &queueMock{}},
// 			args:       args{body: "123", token: "123", header: types.RadiusHTTPAuthHeader},
// 			wantStatus: http.StatusBadRequest,
// 			wantStruct: types.RadiusAccounting{},
// 		},
// 		{
// 			name:       "test2",
// 			fields:     fields{token: "123", queuer: &queueMock{}},
// 			args:       args{body: "123", token: "123", header: types.RadiusHTTPAuthHeader + "s"},
// 			wantStatus: http.StatusNotFound,
// 			wantStruct: types.RadiusAccounting{},
// 		},
// 		{
// 			name:       "test3",
// 			fields:     fields{token: "123", queuer: &queueMock{}},
// 			args:       args{body: pckt1, token: "123", header: types.RadiusHTTPAuthHeader},
// 			wantStatus: http.StatusNoContent,
// 			wantStruct: r1,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &HTTPHandler{
// 				token:  tt.fields.token,
// 				queuer: tt.fields.queuer,
// 			}

// 			w := httptest.NewRecorder()
// 			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.args.body))
// 			req.Header = map[string][]string{
// 				tt.args.header: {tt.args.token},
// 			}

// 			s.radAcct(w, req)

// 			if tt.wantStatus != w.Result().StatusCode {
// 				t.Errorf("want %+v got %+v", tt.wantStatus, w.Result().StatusCode)
// 			}

// 			if reflect.DeepEqual(tt.fields.queuer.res, tt.wantStruct) {
// 				t.Errorf("want %+v got %+v", tt.wantStruct, tt.fields.queuer.res)
// 			}

// 		})
// 	}
// }
