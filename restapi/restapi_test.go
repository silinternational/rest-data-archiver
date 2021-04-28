package restapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRestAPI_httpRequest(t *testing.T) {
	server := getTestServer()
	endpoints := getFakeEndpoints()

	tests := []struct {
		name    string
		restAPI RestAPI
		verb    string
		url     string
		body    string
		headers map[string]string
		want    []byte
		wantErr string
	}{
		{
			name: "basic auth",
			restAPI: RestAPI{
				AuthType: endpoints[EndpointListWorkday].authType,
				Username: endpoints[EndpointListWorkday].username,
				Password: endpoints[EndpointListWorkday].password,
			},
			verb: endpoints[EndpointListWorkday].method,
			url:  server.URL + endpoints[EndpointListWorkday].path,
			want: []byte(endpoints[EndpointListWorkday].responseBody),
		},
		{
			name: "basic auth fail",
			restAPI: RestAPI{
				AuthType: endpoints[EndpointListWorkday].authType,
				Username: endpoints[EndpointListWorkday].username,
				Password: "bad password",
			},
			verb:    endpoints[EndpointListWorkday].method,
			url:     server.URL + endpoints[EndpointListWorkday].path,
			wantErr: "401 Unauthorized",
		},
		{
			name: "bearer token",
			restAPI: RestAPI{
				AuthType: endpoints[EndpointListOther].authType,
				Password: endpoints[EndpointListOther].password,
			},
			verb: endpoints[EndpointListOther].method,
			url:  server.URL + endpoints[EndpointListOther].path,
			want: []byte(endpoints[EndpointListOther].responseBody),
		},
		{
			name: "bearer token fail",
			restAPI: RestAPI{
				AuthType: endpoints[EndpointListOther].authType,
				Password: "bad token",
			},
			verb:    endpoints[EndpointListOther].method,
			url:     server.URL + endpoints[EndpointListOther].path,
			wantErr: "401 Unauthorized",
		},
		{
			name: "salesforce",
			restAPI: RestAPI{
				AuthType: endpoints[EndpointListSalesforce].authType,
				Password: endpoints[EndpointListSalesforce].password,
			},
			verb: endpoints[EndpointListSalesforce].method,
			url:  server.URL + endpoints[EndpointListSalesforce].path,
			want: []byte(endpoints[EndpointListSalesforce].responseBody),
		},
		{
			name: "salesforce fail",
			restAPI: RestAPI{
				AuthType: endpoints[EndpointListSalesforce].authType,
				Password: "bad token",
			},
			verb:    endpoints[EndpointListSalesforce].method,
			url:     server.URL + endpoints[EndpointListSalesforce].path,
			wantErr: "401 Unauthorized",
		},
		{
			name: "bearer create",
			restAPI: RestAPI{
				AuthType: endpoints[EndpointCreateOther].authType,
				Password: endpoints[EndpointCreateOther].password,
			},
			verb: endpoints[EndpointCreateOther].method,
			url:  server.URL + endpoints[EndpointCreateOther].path,
			body: `{"email":"test@example.com","id":"1234"}`,
			want: []byte(endpoints[EndpointCreateOther].responseBody),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.restAPI.httpRequest(tt.verb, tt.url, tt.body, tt.headers)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr, "httpRequest() incorrect error")
				return
			}
			require.Equal(t, tt.want, got, "httpRequest() got = %v, want %v", got, tt.want)
		})
	}
}
