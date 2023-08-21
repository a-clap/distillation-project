// MIT License
//
// Copyright (c) 2023 a-clap
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package mender_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"mender"
	"mender/mocks"
	"mender/pkg/device"
	"mender/pkg/signer"

	"github.com/carlmjohnson/requests"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/maps"
)

const (
	PrivPem = "-----BEGIN PRIVATE KEY-----\nMIIJQgIBADANBgkqhkiG9w0BAQEFAASCCSwwggkoAgEAAoICAQC5Nm2g46xXKa1U" +
		"\nsoHvOgFfmxq1qRM62ZyG6MJyM8y1FJmp2H094z4wrVfJCTxJNkWdE6dz/d/RfPqz\njMTd9k1B3MKAb1zPsVTpX124VwhpefKnJ869IClQtcR6JkGI6Q/krWAsG6XeV4ay\n22l8LamIeTytPI4BG5ydYlWePrjx9q5UcWt2j9gzdz5zZytDzzdhkQtl8exeuZ05\n7/rcpAN2FZ46MtkWFIv00bMwezkOSfAxfQyhgmoOy17FrUr2ffemkXfOVPKntjcA\n8360NiTCp3kHP1JZPldPT6pqA11KH9+St+/4nC6tl1EA3lCRVM9xYAIGn2wEppt7\nGz0wjfHFb1nv/p1TW1mqKquLqSmMwCMsfuxHj18+GrDXaEFtlWAB95daF93W2YIh\nQv5or0mGSrwE1EczNXqWvlGo4Btr1Va6HYjPkTIbb22RXxKM2u4yccUkarAVJT2E\nHOVKKRdQ5Wk432MhNbmp7RKqGisv/PHytQ6wM6tSfynGisIEnsNKqMgHnHZOGQBm\n2Bs/RmPlnOkMNF1AqDPeleHh2+gue82bu4aA+zByvyphjZcEzfR/VBFaBdyaC82D\nHN1FdESZLBAxXJrlgMQOij70gY+CxwkvqfbyOUVPG0+yS7p/haVFbHron2f6IrBe\nv5C9HxH++TCX5RjjriYU0zrUkk3VfQIDAQABAoIB/y3RL6NcKW9c+XQHFB2RKmAM\n9qGdjf+v1XgyxNKsz6R04drS1KfVRiKkvwXnwi0JMA+VOThwfV0MNlZRCGWT9jrE\nFXL+JjpMRikvK32pcWmxc0Dk+qu+mhp1DG7VVDYBd/zwnaibKyJrEHf9HRzCwUcU\nkJnEw7IeSkPGIr9yFc8jXSFgawVu+RGias7zdXdUB+/n/l3r0MAwZXA/Mdxn57rY\nzFH4wWVHNxVSwUPySSY+KPo045OP4/EJWXbOctGPRlzJsqnRAmkGkQ2edkWFuyyC\nYTUA0VS2drMYHRmda/KOEvBVPy4XwSGl6oPDvHCDmMQNa1sGw/GKVw8ewpNwCpTF\ng3KEL/MKPtnwGmZSqqLmItYnM4vwOLEYHqi6ymiTgr1310exku0rZEan2t9aail3\ntosDcbf1rQydcB9ntaaYD0UU0zGsY5bPmGYto3z+sc7E82gb2/VrSbKwGAu84rHs\nCxzlY5MuV/bCJr/NEDRIk3LndZFBwvIuq+0KN1WyyKntOfvgxfopJIUQqNmzvzqU\nd498fNAo0l4pJVIxOnnvJGUVYGsnPc7RrhU7W8qrzWTf0fGg0txwUo8ETfuHO1ia\nJZ2nnfPuRLXfp/BoAkrkMNA0wUx+DT6+JFvdO+mpSSHCxw+hQAe2kyQ1mEz2gLSF\nuKC61Sss3ebBCelZjbECggEBALoUmx7wOoqn7uoCyDxaUl5qeDerWBPmTZI7YqeI\nsLXpbAyH9CEVmdX0sFtmzobjE7X054yMkA4xthY4KHxc4d6IlpDQKyZ/Vt6cHW9W\nriVN+wtD7GNMdmN3aGi8X5kP6FJmesUhn3NiXWBOXs06UF0iPSnR//w0tQz8CzaV\nx4oDJErLnfR2akfxFz3tqZbiObJGAXCNcO6k1/wnxquj/9Ix7uQD+r1BMzkWqWHr\naN90bAf60i+K4dKSxi2Bm9sRu40SPUV4TdSXgXSpPBanH+VrvLR+jr9yH5BcsHVv\nEpoh1V6lv7nntAgxswGYb4HznMbmyIkGpGV8ihnkSHY6rLMCggEBAP7OVty776X5\n7dTnkw1obOXHruQgDXlL6OeZ59QRyhP2yc0kk+9QE34/ygwFhg8Cnr6cSZitoozb\npa6dzQaE11MeFJfJgZx9jsE1dk7XLfM0X2FHgy/RUeGkt+wwJOPw0bnBS7rTcDB+\nAcsFNHGnPGbtNpdIpEhv+EMqbo9nyVPlACShlRji/6BksFTvAYsLnf2kiH44J5kj\nl4TLd3L9yeZ0iGB31SXmV1Ugvq2gcPpI+fcJIURaSXdh8R5V0KOOMFhHV2pUli+/\ns+EVks4krrvB5POmVDBsc2/MqJQZTtZMdYPnh9B1S4yzaFiruK200OVkS1O056AL\nh4lMOVcU7Q8CggEBAKJXR93+B5TgXfea5caBpkro6Gjo+7agvhxN4wv2nPSX6MQl\n+D7E8alQCGw1jQjxI0kjmL9uAl/fztQjum6FOilDUNiWRI8ZmVgtKyDvpo61Mcfq\nQll/Y+nzSwvVDDIlRrJc5c8GPm4T6xMSTHMP5Pzb2jCaHZKTCUGCwuWkVql5hDgc\n3HliteZ916EXr1ULmPqHSMpBG72X4zcCHLmyIoXnOluDfUWPlHjB5JShJKWOlJGB\nqc4AhHOJyYv5/1doaQ/yUbJB+uT2KOL2oo3A0Hr/O2rc3vz3O8JemzRY6wm9asFg\nKZyvIMnlUh6aNu5Q8v40ac/iE5rWxEVCfFVpazECggEAXlR5J8KTl7iM6ZLJh98u\n4WopPt723f7SPFtnzcTAMN/eGYn/EktszAJFhGnPFWN87Ufinxk44ji2f4x/yHgJ\nVwX9zauVxh9dZ/2ozMswgabT4Kme0WcGjyhxxoiUP6Z5nfEHXiTElc8wTr6giarF\n27zZxuvnlcGOAR+GSqS7jclrYiRHlC0FQZXFCcxpn9YvKSVuOnwDfNgGUe0ZTYLS\n6fQeQMhcKmm5zxQOQyzwZlf54hCJNkrOg9nIb9iJIuOS1jujCwRBjW/E4gEglhxS\na8P+RI1BAaREoBD+H8W2v/MSVkCysOObkn0gliMMfZJA+4tDr7t7PG7IQHXtjGV8\n1QKCAQEAkMZj526SUss2bXTcpa0jOsRRxdBsYDf5Cwtv1qVFsPUiBepeiuPd6yUX\nwDAuMNyxgjZovHJ80jl4pSuhmHe/ByvKf5FDxggUu/lcyiTiyFFZN52WXFjl33UR\nyEt2EUTpO3eTsA1WL1NuH7nRlIPuykEk5n04q37gnv53fZimX1763aUr+nitae9U\nXwG2Jo2XV1bAng036CvQaK9iNT9AYA8oxVuOAYCK27+KzC9v1BGcvKdaGKpoZcDa\nSEsEHD16LrKMZnFgGP6iibbwPXiQCRlRbc52gacflvINcW6twhTFg66dewp/O4G1\n4aDn+w97eJDyF2Q9uSML/Bj3TqaORw==\n-----END PRIVATE KEY-----\n"
)

// Will be used in all tests
var keys = func() *signer.Signer {
	s, err := signer.New(signer.WithPrivKey([]byte(PrivPem)))
	if err != nil {
		panic(err)
	}
	return s
}()

type MenderTestSuite struct {
	suite.Suite
}

func Test_MenderSuite(t *testing.T) {
	suite.Run(t, new(MenderTestSuite))
}

func (ms *MenderTestSuite) TestNew() {
	req := ms.Require()

	keys, err := signer.New(signer.WithPrivKey([]byte(PrivPem)))
	req.Nil(err)
	req.NotNil(keys)

	ctrl := gomock.NewController(ms.T())

	mockLoadSaver := mocks.NewMockLoadSaver(ctrl)
	mockLoadSaver.EXPECT().Load(gomock.Any()).Return(nil)

	defer ctrl.Finish()

	args := []struct {
		name        string
		opts        []mender.Option
		expectedErr bool
		errorsIs    []error
		errorsNotIs []error
	}{
		{
			name:        "no options",
			opts:        nil,
			expectedErr: true,
			errorsIs: []error{
				mender.ErrNeedSignerVerifier, mender.ErrNeedServerURLAndToken, mender.ErrNeedDevice, mender.ErrNeedDownloader, mender.ErrNeedInstaller,
				mender.ErrNeedRebooter, mender.ErrNeedLoadSaver, mender.ErrNeedCallbacks, mender.ErrNeedCommiter,
			},
			errorsNotIs: nil,
		},
		{
			name: "with server",
			opts: []mender.Option{
				mender.WithServer("server", "token"),
			},
			expectedErr: true,
			errorsIs: []error{
				mender.ErrNeedSignerVerifier, mender.ErrNeedDevice, mender.ErrNeedDownloader, mender.ErrNeedInstaller, mender.ErrNeedRebooter,
				mender.ErrNeedLoadSaver, mender.ErrNeedCallbacks, mender.ErrNeedCommiter,
			},
			errorsNotIs: []error{mender.ErrNeedServerURLAndToken},
		},
		{
			name: "with signer verifier",
			opts: []mender.Option{
				mender.WithSigner(keys),
			},
			expectedErr: true,
			errorsIs: []error{
				mender.ErrNeedServerURLAndToken, mender.ErrNeedDevice, mender.ErrNeedDownloader, mender.ErrNeedInstaller, mender.ErrNeedRebooter,
				mender.ErrNeedLoadSaver, mender.ErrNeedCallbacks, mender.ErrNeedCommiter,
			},
			errorsNotIs: []error{mender.ErrNeedSignerVerifier},
		},
		{
			name: "with downloader",
			opts: []mender.Option{
				mender.WithDownloader(mocks.NewMockDownloader(ctrl)),
			},
			expectedErr: true,
			errorsIs: []error{
				mender.ErrNeedServerURLAndToken, mender.ErrNeedDevice, mender.ErrNeedSignerVerifier, mender.ErrNeedInstaller, mender.ErrNeedRebooter,
				mender.ErrNeedLoadSaver, mender.ErrNeedCallbacks, mender.ErrNeedCommiter,
			},
			errorsNotIs: []error{mender.ErrNeedDownloader},
		},
		{
			name: "with installer",
			opts: []mender.Option{
				mender.WithInstaller(mocks.NewMockInstaller(ctrl)),
			},
			expectedErr: true,
			errorsIs: []error{
				mender.ErrNeedServerURLAndToken, mender.ErrNeedDevice, mender.ErrNeedSignerVerifier, mender.ErrNeedDownloader, mender.ErrNeedRebooter,
				mender.ErrNeedLoadSaver, mender.ErrNeedCallbacks, mender.ErrNeedCommiter,
			},
			errorsNotIs: []error{mender.ErrNeedInstaller},
		},
		{
			name: "with rebooter",
			opts: []mender.Option{
				mender.WithRebooter(mocks.NewMockRebooter(ctrl)),
			},
			expectedErr: true,
			errorsIs: []error{
				mender.ErrNeedServerURLAndToken, mender.ErrNeedDevice, mender.ErrNeedSignerVerifier, mender.ErrNeedDownloader, mender.ErrNeedInstaller,
				mender.ErrNeedLoadSaver, mender.ErrNeedCallbacks, mender.ErrNeedCommiter,
			},
			errorsNotIs: []error{mender.ErrNeedRebooter},
		},
		{
			name: "with loadsaver",
			opts: []mender.Option{
				mender.WithLoadSaver(mockLoadSaver),
			},
			expectedErr: true,
			errorsIs: []error{
				mender.ErrNeedServerURLAndToken, mender.ErrNeedDevice, mender.ErrNeedSignerVerifier, mender.ErrNeedDownloader, mender.ErrNeedInstaller,
				mender.ErrNeedRebooter, mender.ErrNeedCallbacks, mender.ErrNeedCommiter,
			},
			errorsNotIs: []error{mender.ErrNeedLoadSaver},
		},
		{
			name: "with callbacks",
			opts: []mender.Option{
				mender.WithCallbacks(mocks.NewMockCallbacks(ctrl)),
			},
			expectedErr: true,
			errorsIs: []error{
				mender.ErrNeedServerURLAndToken, mender.ErrNeedDevice, mender.ErrNeedSignerVerifier, mender.ErrNeedDownloader, mender.ErrNeedInstaller,
				mender.ErrNeedRebooter, mender.ErrNeedLoadSaver, mender.ErrNeedCommiter,
			},
			errorsNotIs: []error{mender.ErrNeedCallbacks},
		},
		{
			name: "with committer",
			opts: []mender.Option{
				mender.WithCommitter(mocks.NewMockCommitter(ctrl)),
			},
			expectedErr: true,
			errorsIs: []error{
				mender.ErrNeedServerURLAndToken, mender.ErrNeedDevice, mender.ErrNeedSignerVerifier, mender.ErrNeedDownloader, mender.ErrNeedInstaller,
				mender.ErrNeedRebooter, mender.ErrNeedLoadSaver, mender.ErrNeedCallbacks,
			},
			errorsNotIs: []error{mender.ErrNeedCommiter},
		},
		{
			name: "all good",
			opts: []mender.Option{
				mender.WithServer("server", "token"),
				mender.WithSigner(keys),
				mender.WithDevice(mocks.NewMockDevice(ctrl)),
				mender.WithDownloader(mocks.NewMockDownloader(ctrl)),
				mender.WithInstaller(mocks.NewMockInstaller(ctrl)),
				mender.WithRebooter(mocks.NewMockRebooter(ctrl)),
				mender.WithLoadSaver(mockLoadSaver),
				mender.WithCallbacks(mocks.NewMockCallbacks(ctrl)),
				mender.WithCommitter(mocks.NewMockCommitter(ctrl)),
			},
			expectedErr: false,
		},
	}
	for _, arg := range args {
		// Create without options result with error
		client, err := mender.New(arg.opts...)
		if !arg.expectedErr {
			req.NotNil(client, arg.name)
			req.Nil(err, arg.name)
			continue
		}

		req.Nil(client, arg.name)
		req.NotNil(err, arg.name)

		for _, errIs := range arg.errorsIs {
			req.ErrorIs(err, errIs, arg.name)
		}

		for _, errNotIs := range arg.errorsNotIs {
			req.NotErrorIs(err, errNotIs, arg.name)
		}
	}
}

func (ms *MenderTestSuite) TestConnect() {
	req := ms.Require()
	asrt := ms.Assert()

	args := []struct {
		name         string
		statusCode   int
		retBody      string
		expectedErr  error
		teenantToken string
		onID         *struct {
			attribute []device.Attribute
			err       error
		}
	}{
		{
			name:         "proper connection",
			statusCode:   http.StatusOK,
			retBody:      "jwt token",
			teenantToken: "teenant_token",
			onID: &struct {
				attribute []device.Attribute
				err       error
			}{
				attribute: []device.Attribute{{Name: "mac", Value: []string{"01:02:03:04:05:06"}}},
				err:       nil,
			},
		},
		{
			name:         "StatusUnauthorized",
			statusCode:   http.StatusUnauthorized,
			expectedErr:  mender.ErrNeedAuthentication,
			teenantToken: "teenant_token1",
			onID: &struct {
				attribute []device.Attribute
				err       error
			}{
				attribute: []device.Attribute{{Name: "mac", Value: []string{"01:02:03:04:05:06"}}},
				err:       nil,
			},
		},
		{
			name:         "Internal server error",
			statusCode:   http.StatusInternalServerError,
			expectedErr:  requests.ErrValidator,
			teenantToken: "teenant_token2",
			onID: &struct {
				attribute []device.Attribute
				err       error
			}{
				attribute: []device.Attribute{{Name: "mac", Value: []string{"01:02:03:04:05:06"}}},
				err:       nil,
			},
		},
		{
			name:         "Bad request",
			statusCode:   http.StatusBadRequest,
			expectedErr:  requests.ErrValidator,
			teenantToken: "teenant_token3",
			onID: &struct {
				attribute []device.Attribute
				err       error
			}{
				attribute: []device.Attribute{{Name: "mac", Value: []string{"01:02:03:04:05:06"}}},
				err:       nil,
			},
		},
	}
	for _, arg := range args {
		header := make(http.Header)

		var body []byte

		srv := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			// Fetch header and body
			maps.Copy(header, request.Header)

			body, _ = io.ReadAll(request.Body)
			_ = request.Body.Close()

			// Write response
			writer.WriteHeader(arg.statusCode)
			writer.Header().Set("Content-Type", "application/json")
			// And body - if exists
			if len(arg.retBody) > 0 {
				_, err := writer.Write([]byte(arg.retBody))
				asrt.Nil(err, arg.name)
			}
		}))

		ctrl := gomock.NewController(ms.T())
		dev := mocks.NewMockDevice(ctrl)
		mockLoadSaver := mocks.NewMockLoadSaver(ctrl)
		mockLoadSaver.EXPECT().Load(gomock.Any()).Return(nil)

		if arg.onID != nil {
			dev.EXPECT().ID().Return(arg.onID.attribute, arg.onID.err).AnyTimes()
		}

		client, err := mender.New(
			mender.WithServer(srv.URL, arg.teenantToken),
			mender.WithSigner(keys),
			mender.WithDevice(dev),
			mender.WithDownloader(mocks.NewMockDownloader(ctrl)),
			mender.WithInstaller(mocks.NewMockInstaller(ctrl)),
			mender.WithRebooter(mocks.NewMockRebooter(ctrl)),
			mender.WithLoadSaver(mockLoadSaver),
			mender.WithCallbacks(mocks.NewMockCallbacks(ctrl)),
			mender.WithCommitter(mocks.NewMockCommitter(ctrl)),
		)

		// If somehow we didn't create client, fail fast
		if client == nil || err != nil {
			req.Fail("cannot proceed", fmt.Sprintf("%v: %v", arg.name, err))
			srv.Close()
			ctrl.Finish()
		}

		err = client.Connect()
		// Close server and mock
		srv.Close()
		ctrl.Finish()

		// Verify body
		var bodyMap map[string]interface{}
		req.Nil(json.Unmarshal(body, &bodyMap), arg.name)

		pubKey := bodyMap["pubkey"]
		req.Equal(keys.PublicKeyPEM(), pubKey, arg.name)

		tenant_token := bodyMap["tenant_token"]
		req.Equal(arg.teenantToken, tenant_token, arg.name)

		id_data := bodyMap["id_data"]
		expected := fmt.Sprintf(`{"mac":["%v"]}`, arg.onID.attribute[0].Value[0])

		req.Equal(expected, id_data, arg.name)

		// Check, if client generated proper headers
		sigs, ok := header["X-Men-Signature"]
		req.True(ok, arg.name)
		req.Len(sigs, 1, arg.name)

		// Verify sig with body
		req.Nil(client.Verify(body, []byte(sigs[0])))

		if arg.expectedErr != nil {
			req.NotNil(err, arg.name)
			req.ErrorIs(err, arg.expectedErr, arg.name)
			continue
		}

		req.Nil(err, arg.name)
	}
}

func (ms *MenderTestSuite) TestUpdateInventory() {
	req := ms.Require()
	asrt := ms.Assert()

	args := []struct {
		name         string
		teenantToken string
		token        string
		statusCode   int
		body         string
		expectedErr  error
		retInfo      device.Info
		retInfoErr   error
		retID        []device.Attribute
		retIDErr     error
		retAttr      []device.Attribute
		retAttrErr   error
	}{
		{
			name:         "basic",
			statusCode:   http.StatusNoContent,
			teenantToken: "blah",
			token:        "jwt token",
			body:         "",
			expectedErr:  nil,
			retInfo:      device.Info{DeviceType: "fake_device", ArtifactName: "artifact_1"},
			retInfoErr:   nil,
			retID:        []device.Attribute{{Name: "mac", Value: []string{"01:02:03:04:05:06"}}},
			retIDErr:     nil,
			retAttr:      []device.Attribute{{Name: "attr", Value: []string{"awesome", "device"}}},
			retAttrErr:   nil,
		},
	}

	for _, arg := range args {
		header := make(http.Header)
		body := make([]byte, 0, 512)

		handle := http.NewServeMux()
		// Return token
		handle.Handle("/api/devices/v1/authentication/auth_requests", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte(arg.token))
		}))

		handle.Handle("/api/devices/v1/inventory/device/attributes", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			// Fetch header and body
			maps.Copy(header, request.Header)

			body, _ = io.ReadAll(request.Body)
			_ = request.Body.Close()

			// Write response
			writer.WriteHeader(arg.statusCode)
			writer.Header().Set("Content-Type", "application/json")
			// And body - if exists
			if len(arg.body) > 0 {
				_, err := writer.Write([]byte(arg.body))
				asrt.Nil(err, arg.name)
			}
		}))

		srv := httptest.NewServer(handle)

		ctrl := gomock.NewController(ms.T())
		dev := mocks.NewMockDevice(ctrl)
		mockLoadSaver := mocks.NewMockLoadSaver(ctrl)

		mockLoadSaver.EXPECT().Load(gomock.Any()).Return(nil)

		dev.EXPECT().Attributes().Return(arg.retAttr, arg.retAttrErr)
		dev.EXPECT().ID().Return(arg.retID, arg.retIDErr)
		dev.EXPECT().Info().Return(arg.retInfo, arg.retInfoErr)

		client, err := mender.New(
			mender.WithServer(srv.URL, arg.teenantToken),
			mender.WithSigner(keys),
			mender.WithDevice(dev),
			mender.WithDownloader(mocks.NewMockDownloader(ctrl)),
			mender.WithInstaller(mocks.NewMockInstaller(ctrl)),
			mender.WithRebooter(mocks.NewMockRebooter(ctrl)),
			mender.WithLoadSaver(mockLoadSaver),
			mender.WithCallbacks(mocks.NewMockCallbacks(ctrl)),
			mender.WithCommitter(mocks.NewMockCommitter(ctrl)),
		)

		if client == nil || err != nil {
			req.Fail(fmt.Sprintln("cannot proceed without client, err is:", err), arg.name)
			srv.Close()
			ctrl.Finish()
		}

		req.Nil(client.Connect())
		// So far so good

		err = client.UpdateInventory()
		// Free resources
		srv.Close()
		ctrl.Finish()

		if arg.expectedErr != nil {
			req.NotNil(err, arg.name)
			req.ErrorIs(err, arg.expectedErr)
			continue
		}

		req.Nil(err)
		// Verify sent attributes
		var attrs []device.Attribute
		req.Nil(json.Unmarshal(body, &attrs))

		// Attributes should contain also info
		arg.retAttr = append(arg.retAttr, []device.Attribute{
			{
				Name:  "device_type",
				Value: []string{arg.retInfo.DeviceType},
			},
			{
				Name:  "artifact_name",
				Value: []string{arg.retInfo.ArtifactName},
			},
		}...)

		req.ElementsMatch(arg.retAttr, attrs, arg.name)

		// Verify sent token
		token := header["Authorization"][0]
		req.EqualValues("Bearer "+arg.token, token, arg.name)
	}
}

func (ms *MenderTestSuite) TestSendStatus() {
	req := ms.Require()

	args := []struct {
		name           string
		statusCode     int
		deployID       string
		sendStatus     mender.DeploymentStatus
		expectedStatus string
		err            error
	}{
		{
			name:           "StatusBadRequest",
			statusCode:     http.StatusBadRequest,
			sendStatus:     mender.Downloading,
			expectedStatus: "downloading",
			deployID:       "1233",
			err:            fmt.Errorf("%v", http.StatusBadRequest),
		},
		{
			name:           "StatusNotFound",
			statusCode:     http.StatusNotFound,
			sendStatus:     mender.PauseBeforeInstalling,
			expectedStatus: "pause_before_installing",
			deployID:       "12345",
			err:            fmt.Errorf("%v", http.StatusNotFound),
		},
		{
			name:           "StatusConflict",
			statusCode:     http.StatusConflict,
			sendStatus:     mender.Installing,
			expectedStatus: "installing",
			deployID:       "123456",
			err:            fmt.Errorf("%v", http.StatusConflict),
		},
		{
			name:           "StatusInternalServerError",
			statusCode:     http.StatusInternalServerError,
			sendStatus:     mender.PauseBeforeRebooting,
			expectedStatus: "pause_before_rebooting",
			deployID:       "1234567",
			err:            fmt.Errorf("%v", http.StatusInternalServerError),
		},
		{
			name:           "NoContent",
			statusCode:     http.StatusNoContent,
			sendStatus:     mender.PauseBeforeCommitting,
			expectedStatus: "pause_before_committing",
			deployID:       "1233",
			err:            nil,
		},
		{
			name:           "NoContent",
			statusCode:     http.StatusNoContent,
			sendStatus:     mender.Success,
			expectedStatus: "success",
			deployID:       "1233",
			err:            nil,
		},
		{
			name:           "NoContent",
			statusCode:     http.StatusNoContent,
			sendStatus:     mender.Failure,
			expectedStatus: "failure",
			deployID:       "1233",
			err:            nil,
		},
		{
			name:           "NoContent",
			statusCode:     http.StatusNoContent,
			sendStatus:     mender.AlreadyInstalled,
			expectedStatus: "already-installed",
			deployID:       "1233",
			err:            nil,
		},
	}

	for _, arg := range args {
		handle := http.NewServeMux()
		// Return token
		handle.Handle("/api/devices/v1/authentication/auth_requests", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte("token"))
		}))

		handle.Handle("/api/devices/v1/deployments/device/deployments/next", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			body, _ := json.Marshal(mender.DeploymentInstructions{
				ID: arg.deployID,
				Artifact: mender.DeploymentArtifact{
					Name: "my-app-0.1",
					Source: mender.DeploymentSource{
						URI:    "https://aws.my_update_bucket.com/image_123",
						Expire: "2016-03-11T13:03:17.063493443Z",
					},
					Compatible: []string{
						"device",
						"rspi2",
						"rspi0",
					},
				},
			})
			_, _ = writer.Write(body)
		}))

		var body []byte
		deployURL := fmt.Sprintf("/api/devices/v1/deployments/device/deployments/%v/status", arg.deployID)
		handle.Handle(deployURL, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(arg.statusCode)
			body, _ = io.ReadAll(request.Body)
			_ = request.Body.Close()
		}))

		srv := httptest.NewServer(handle)

		ctrl := gomock.NewController(ms.T())
		dev := mocks.NewMockDevice(ctrl)
		mockLoadSaver := mocks.NewMockLoadSaver(ctrl)

		mockLoadSaver.EXPECT().Load(gomock.Any()).Return(nil)
		mockLoadSaver.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

		dev.EXPECT().ID().Return([]device.Attribute{{Name: "id", Value: []string{"id"}}}, nil)
		dev.EXPECT().Info().Return(device.Info{DeviceType: "device", ArtifactName: "artifact"}, nil)

		client, err := mender.New(
			mender.WithServer(srv.URL, "teenant token"),
			mender.WithSigner(keys),
			mender.WithDevice(dev),
			mender.WithDownloader(mocks.NewMockDownloader(ctrl)),
			mender.WithInstaller(mocks.NewMockInstaller(ctrl)),
			mender.WithRebooter(mocks.NewMockRebooter(ctrl)),
			mender.WithLoadSaver(mockLoadSaver),
			mender.WithCallbacks(mocks.NewMockCallbacks(ctrl)),
			mender.WithCommitter(mocks.NewMockCommitter(ctrl)),
		)

		if client == nil || err != nil {
			srv.Close()
			ctrl.Finish()
			req.Fail("cannot proceed", arg.name)
		}

		// Fetch token
		req.Nil(client.Connect())

		// Get deploy instructions
		newRelease, err := client.PullReleases()
		req.True(newRelease, arg.name)
		req.Contains(client.AvailableReleases(), "my-app-0.1", arg.name)
		req.Nil(err, arg.name)

		err = client.NotifyServer(arg.sendStatus, "my-app-0.1")

		// Free resources
		ctrl.Finish()
		srv.Close()

		var bodyMap map[string]interface{}
		req.Nil(json.Unmarshal(body, &bodyMap), arg.name)

		req.EqualValues(arg.expectedStatus, bodyMap["status"], arg.name)
		if arg.err != nil {
			req.NotNil(err, arg.name)
			req.ErrorContains(err, arg.err.Error(), arg.name)
			continue
		}
		req.Nil(err, arg.name)
	}
}

func (ms *MenderTestSuite) TestCheckDeployment() {
	mustMarshal := func(v any) []byte {
		ret, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		return ret
	}

	req := ms.Require()

	args := []struct {
		name        string
		devInfo     device.Info
		statusCode  int
		body        []byte
		newRelease  bool
		releaseName string
		releaseErr  error
	}{
		{
			name:        "StatusInternalServerError",
			devInfo:     device.Info{DeviceType: "StatusInternalServerError", ArtifactName: "0.1.0"},
			statusCode:  http.StatusInternalServerError,
			body:        nil,
			newRelease:  false,
			releaseName: "",
			releaseErr:  fmt.Errorf("%v", http.StatusInternalServerError),
		},
		{
			name:        "StatusBadRequest",
			devInfo:     device.Info{DeviceType: "StatusBadRequest", ArtifactName: "0.1.31"},
			statusCode:  http.StatusBadRequest,
			body:        nil,
			newRelease:  false,
			releaseName: "",
			releaseErr:  fmt.Errorf("%v", http.StatusBadRequest),
		},
		{
			name:        "StatusNotFound",
			devInfo:     device.Info{DeviceType: "StatusNotFound", ArtifactName: "01.1.31"},
			statusCode:  http.StatusNotFound,
			body:        nil,
			newRelease:  false,
			releaseName: "",
			releaseErr:  fmt.Errorf("%v", http.StatusNotFound),
		},
		{
			name:        "StatusConflict",
			devInfo:     device.Info{DeviceType: "StatusConflict", ArtifactName: "0.21.31"},
			statusCode:  http.StatusConflict,
			body:        nil,
			newRelease:  false,
			releaseName: "",
			releaseErr:  fmt.Errorf("%v", http.StatusConflict),
		},
		{
			name:        "StatusNoContent",
			devInfo:     device.Info{DeviceType: "fake_device2", ArtifactName: "11.1.31"},
			statusCode:  http.StatusNoContent,
			body:        nil,
			newRelease:  false,
			releaseName: "",
			releaseErr:  nil,
		},
		{
			name:       "StatusOK",
			devInfo:    device.Info{DeviceType: "fake_device2", ArtifactName: "11.1.31"},
			statusCode: http.StatusOK,
			body: mustMarshal(mender.DeploymentInstructions{
				ID: "w81s4fae-7dec-11d0-a765-00a0c91e6bf6",
				Artifact: mender.DeploymentArtifact{
					Name: "my-app-0.1",
					Source: mender.DeploymentSource{
						URI:    "https://aws.my_update_bucket.com/image_123",
						Expire: "2016-03-11T13:03:17.063493443Z",
					},
					Compatible: []string{
						"fake_device2",
						"rspi2",
						"rspi0",
					},
				},
			}),
			newRelease:  true,
			releaseName: "my-app-0.1",
			releaseErr:  nil,
		},
		{
			name:       "Incompatible device",
			devInfo:    device.Info{DeviceType: "fake_device2", ArtifactName: "11.1.31"},
			statusCode: http.StatusOK,
			body: mustMarshal(mender.DeploymentInstructions{
				ID: "w81s4fae-7dec-11d0-a765-00a0c91e6bf6",
				Artifact: mender.DeploymentArtifact{
					Name: "my-app-0.1",
					Source: mender.DeploymentSource{
						URI:    "https://aws.my_update_bucket.com/image_123",
						Expire: "2016-03-11T13:03:17.063493443Z",
					},
					Compatible: []string{
						"rspi1",
						"rspi2",
						"rspi0",
					},
				},
			}),
			newRelease:  false,
			releaseName: "",
			releaseErr:  nil,
		},
	}

	for _, arg := range args {
		header := make(http.Header)
		values := make(url.Values)

		handle := http.NewServeMux()
		// Return token
		handle.Handle("/api/devices/v1/authentication/auth_requests", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write([]byte("token"))
		}))

		handle.Handle("/api/devices/v1/deployments/device/deployments/next", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			// Fetch header
			values = request.URL.Query()
			maps.Copy(header, request.Header)

			writer.WriteHeader(arg.statusCode)
			if arg.body != nil {
				_, _ = writer.Write(arg.body)
			}
		}))

		srv := httptest.NewServer(handle)

		ctrl := gomock.NewController(ms.T())
		dev := mocks.NewMockDevice(ctrl)

		mockLoadSaver := mocks.NewMockLoadSaver(ctrl)
		mockLoadSaver.EXPECT().Load(gomock.Any()).Return(nil)
		if arg.releaseName != "" {
			mockLoadSaver.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
		}

		dev.EXPECT().ID().Return([]device.Attribute{{Name: "id", Value: []string{"id"}}}, nil)
		dev.EXPECT().Info().Return(arg.devInfo, nil)

		client, err := mender.New(
			mender.WithServer(srv.URL, "teenant token"),
			mender.WithSigner(keys),
			mender.WithDevice(dev),
			mender.WithDownloader(mocks.NewMockDownloader(ctrl)),
			mender.WithInstaller(mocks.NewMockInstaller(ctrl)),
			mender.WithRebooter(mocks.NewMockRebooter(ctrl)),
			mender.WithLoadSaver(mockLoadSaver),
			mender.WithCallbacks(mocks.NewMockCallbacks(ctrl)),
			mender.WithCommitter(mocks.NewMockCommitter(ctrl)),
		)

		if client == nil || err != nil {
			req.Fail("cannot proceed", fmt.Sprintf("%v: %v", arg.name, err))
			srv.Close()
			ctrl.Finish()
		}

		// Fetch token
		req.Nil(client.Connect())

		// Check deployment
		newRelease, err := client.PullReleases()
		// Free resources
		ctrl.Finish()
		srv.Close()

		// Verify params
		artifact_name := values["artifact_name"][0]
		req.Equal(arg.devInfo.ArtifactName, artifact_name, arg.name)

		device_type := values["device_type"][0]
		req.Equal(arg.devInfo.DeviceType, device_type, arg.name)

		if arg.releaseErr != nil {
			req.NotNil(err, arg.name)
			req.ErrorContains(err, arg.releaseErr.Error(), arg.name)
			continue
		}

		req.Nil(err, arg.name)
		// Verify response
		req.Equal(arg.newRelease, newRelease, arg.name)
		releases := client.AvailableReleases()
		if arg.releaseName == "" {
			req.Len(releases, 0, arg.name)
		} else {
			req.Contains(releases, arg.releaseName, arg.name)
		}
	}
}

func (ms *MenderTestSuite) TestUpdate() {
	req := ms.Require()

	const (
		deployID     = "1234"
		artifactName = "my-app-0.1"
	)

	deployArtifact := mender.DeploymentInstructions{
		ID: deployID,
		Artifact: mender.DeploymentArtifact{
			Name: artifactName,
			Source: mender.DeploymentSource{
				URI:    "https://aws.my_update_bucket.com/image_123",
				Expire: "2016-03-11T13:03:17.063493443Z",
			},
			Compatible: []string{"device"},
		},
	}

	// Prepare server
	handle := http.NewServeMux()
	// Return token
	handle.Handle("/api/devices/v1/authentication/auth_requests", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte("token"))
	}))

	handle.Handle("/api/devices/v1/deployments/device/deployments/next", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		body, _ := json.Marshal(deployArtifact)
		_, _ = writer.Write(body)
	}))

	type jsonStatus map[string]string
	// Make chan big enough to not block
	serverReceivedStatuses := make(chan jsonStatus, 10)

	deployURL := fmt.Sprintf("/api/devices/v1/deployments/device/deployments/%v/status", deployID)
	handle.Handle(deployURL, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Proper response
		writer.WriteHeader(http.StatusOK)
		// Get body
		body, _ := io.ReadAll(request.Body)
		_ = request.Body.Close()
		// Parse body to json
		js := make(jsonStatus)
		// Make sure it is correct
		req.Nil(json.Unmarshal(body, &js))

		serverReceivedStatuses <- js
	}))

	srv := httptest.NewServer(handle)

	ctrl := gomock.NewController(ms.T())
	mockDevice := mocks.NewMockDevice(ctrl)

	mockDevice.EXPECT().ID().Return([]device.Attribute{{Name: "id", Value: []string{"id"}}}, nil)
	mockDevice.EXPECT().Info().Return(device.Info{DeviceType: "device", ArtifactName: "artifact"}, nil)

	mockDownloader := mocks.NewMockDownloader(ctrl)
	mockInstaller := mocks.NewMockInstaller(ctrl)
	mockRebooter := mocks.NewMockRebooter(ctrl)
	mockCallbacks := mocks.NewMockCallbacks(ctrl)
	mockLoadSaver := mocks.NewMockLoadSaver(ctrl)

	mockLoadSaver.EXPECT().Load(gomock.Any()).Return(nil)
	// First call on init
	mockLoadSaver.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

	client, err := mender.New(
		mender.WithServer(srv.URL, "teenant token"),
		mender.WithSigner(keys),
		mender.WithDevice(mockDevice),
		mender.WithDownloader(mockDownloader),
		mender.WithInstaller(mockInstaller),
		mender.WithRebooter(mockRebooter),
		mender.WithLoadSaver(mockLoadSaver),
		mender.WithCallbacks(mockCallbacks),
		mender.WithCommitter(mocks.NewMockCommitter(ctrl)),
	)

	r := ms.Require()

	r.Nil(err)
	r.NotNil(client)

	r.Nil(client.Connect())
	got, err := client.PullReleases()
	r.True(got)

	r.Contains(client.AvailableReleases(), artifactName)
	r.Nil(err)

	// First step is downloading
	errChan := make(chan error, 1)
	progressChan := make(chan int, 100)
	for i := 1; i <= 100; i++ {
		mockCallbacks.EXPECT().Update(mender.Downloading, i)
		progressChan <- i
	}

	downloadFinished := make(chan struct{})
	startInstall := make(chan struct{})
	mockCallbacks.EXPECT().Update(mender.PauseBeforeInstalling, 100).Times(1)
	mockCallbacks.EXPECT().NextState(mender.Installing).Return(true).Do(func(any) {
		downloadFinished <- struct{}{}
		<-startInstall
	})

	mockDownloader.EXPECT().Download(gomock.Any(), gomock.Any()).Return(progressChan, errChan, nil)

	err = client.Update(artifactName)
	r.Nil(err)

	select {
	case <-downloadFinished:
	case <-time.After(1 * time.Millisecond):
		ms.FailNow("failed at downloading")
	}

	// We expect to receive statuses as follows:
	for _, expectedStatus := range []string{"downloading", "pause_before_installing"} {
		select {
		case status := <-serverReceivedStatuses:
			req.Equal(expectedStatus, status["status"])
		case <-time.After(1 * time.Millisecond):
			req.Fail("failed to receive proper status:", expectedStatus)
		}
	}

	// Now we should expect some install calls
	progressChan = make(chan int, 100)
	for i := 1; i <= 100; i++ {
		mockCallbacks.EXPECT().Update(mender.Installing, i)
		progressChan <- i
	}

	errChan = make(chan error, 1)
	mockInstaller.EXPECT().Install(gomock.Any()).Return(progressChan, errChan, nil)

	installFinished := make(chan struct{})
	startReboot := make(chan struct{})

	mockCallbacks.EXPECT().Update(mender.PauseBeforeRebooting, 100).Times(1)
	mockCallbacks.EXPECT().NextState(mender.Rebooting).Return(true).Do(func(any) {
		installFinished <- struct{}{}
		<-startReboot
	})

	// Start installl
	close(startInstall)

	select {
	case <-installFinished:
	case <-time.After(1 * time.Millisecond):
		ms.FailNow("failed at installing")
	}

	// We expect to receive statuses as follows:
	for _, expectedStatus := range []string{"installing", "pause_before_rebooting"} {
		select {
		case status := <-serverReceivedStatuses:
			req.Equal(expectedStatus, status["status"])
		case <-time.After(1 * time.Millisecond):
			req.Fail("failed to receive proper status:", expectedStatus)
		}
	}

	rebootFinished := make(chan struct{})
	mockCallbacks.EXPECT().Update(mender.Rebooting, 1).Times(1)

	mockRebooter.EXPECT().Reboot().Return(nil).Do(func() {
		rebootFinished <- struct{}{}
	})

	// Client should store this data
	expectedSaveData := mender.Artifacts{
		Current: &mender.CurrentDeployment{
			State:                  mender.PauseBeforeCommitting,
			DeploymentInstructions: &deployArtifact,
		},
		Archive: []mender.DeploymentInstructions{deployArtifact},
	}

	mockLoadSaver.EXPECT().Save(gomock.Any(), expectedSaveData).Return(nil)

	close(startReboot)

	select {
	case <-rebootFinished:
	case <-time.After(1 * time.Millisecond):
		ms.FailNow("failed at rebooting")
	}
	// We expect to receive statuses as follows:
	for _, expectedStatus := range []string{"rebooting"} {
		select {
		case status := <-serverReceivedStatuses:
			req.Equal(expectedStatus, status["status"])
		case <-time.After(1 * time.Millisecond):
			req.Fail("failed to receive proper status:", expectedStatus)
		}
	}
}

func (ms *MenderTestSuite) TestContinueUpdateAfterReboot() {
	req := ms.Require()

	const (
		deployID     = "1234"
		artifactName = "my-app-0.1"
	)

	// Prepare server
	handle := http.NewServeMux()
	// Return token
	handle.Handle("/api/devices/v1/authentication/auth_requests", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte("token"))
	}))

	handle.Handle("/api/devices/v1/deployments/device/deployments/next", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	}))

	type jsonStatus map[string]string
	bodyStatus := make(chan jsonStatus, 3)
	// bodyStatus := make([]jsonStatus, 0)

	deployURL := fmt.Sprintf("/api/devices/v1/deployments/device/deployments/%v/status", deployID)
	handle.Handle(deployURL, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// Proper response
		writer.WriteHeader(http.StatusOK)
		// Get body
		body, _ := io.ReadAll(request.Body)
		_ = request.Body.Close()
		// Parse body to json
		js := make(jsonStatus)
		// Make sure it is correct
		req.Nil(json.Unmarshal(body, &js))
		// bodyStatus = append(bodyStatus, js)
		bodyStatus <- js
	}))

	srv := httptest.NewServer(handle)

	ctrl := gomock.NewController(ms.T())
	mockDevice := mocks.NewMockDevice(ctrl)
	mockDevice.EXPECT().ID().Return([]device.Attribute{{Name: "id", Value: []string{"id"}}}, nil)

	mockDownloader := mocks.NewMockDownloader(ctrl)
	mockInstaller := mocks.NewMockInstaller(ctrl)
	mockRebooter := mocks.NewMockRebooter(ctrl)
	mockCallbacks := mocks.NewMockCallbacks(ctrl)
	mockLoadSaver := mocks.NewMockLoadSaver(ctrl)

	arti := mender.Artifacts{
		Current: &mender.CurrentDeployment{
			State: mender.PauseBeforeCommitting,
			DeploymentInstructions: &mender.DeploymentInstructions{
				ID: deployID,
				Artifact: mender.DeploymentArtifact{
					Name: artifactName,
					Source: mender.DeploymentSource{
						URI:    "https://aws.my_update_bucket.com/image_123",
						Expire: "2016-03-11T13:03:17.063493443Z",
					},
					Compatible: []string{"device"},
				},
			},
		},
		Archive: nil,
	}

	mockLoadSaver.EXPECT().Load(gomock.Any()).Return(arti)
	// Save should get called with new args
	saveCalled := make(chan struct{})
	mockLoadSaver.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Do(func(any, any) {
		close(saveCalled)
	})

	nextStateCalled := make(chan struct{})
	// Should continue update
	mockCallbacks.EXPECT().NextState(mender.Success).Return(true).Do(func(any) {
		close(nextStateCalled)
	})

	client, err := mender.New(
		mender.WithServer(srv.URL, "teenant token"),
		mender.WithSigner(keys),
		mender.WithDevice(mockDevice),
		mender.WithDownloader(mockDownloader),
		mender.WithInstaller(mockInstaller),
		mender.WithRebooter(mockRebooter),
		mender.WithLoadSaver(mockLoadSaver),
		mender.WithCallbacks(mockCallbacks),
		mender.WithCommitter(mocks.NewMockCommitter(ctrl)),
	)

	req.Nil(err)
	req.NotNil(client)

	// Wait for NextState
	select {
	case <-nextStateCalled:
	case <-time.After(1 * time.Millisecond):
		req.Fail("client should have already called NextState")
	}

	// Wait for Status
	select {
	case st := <-bodyStatus:
		req.Equal("success", st["status"])
	case <-time.After(1 * time.Millisecond):
		req.Fail("didn't receive status")
	}

	// Wait for Save
	select {
	case <-saveCalled:
	case <-time.After(1 * time.Millisecond):
		req.Fail("client should have called Save already")
	}
}
