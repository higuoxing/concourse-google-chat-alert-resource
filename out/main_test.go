package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/higuoxing/concourse-google-chat-alert-resource/concourse"
)

func TestOut(t *testing.T) {
	ok := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ok.Close()
	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer bad.Close()

	env := map[string]string{
		"ATC_EXTERNAL_URL":    "https://ci.example.com",
		"BUILD_TEAM_NAME":     "main",
		"BUILD_PIPELINE_NAME": "demo",
		"BUILD_JOB_NAME":      "test",
		"BUILD_NAME":          "2",
	}

	cases := map[string]struct {
		outRequest *concourse.OutRequest
		want       *concourse.OutResponse
		env        map[string]string
		err        bool
	}{
		"default alert": {
			outRequest: &concourse.OutRequest{
				Source: concourse.Source{URL: ok.URL},
			},
			want: &concourse.OutResponse{
				Version: concourse.Version{"ver": "static"},
				Metadata: []concourse.Metadata{
					{Name: "type", Value: "default"},
					{Name: "alerted", Value: "true"},
				},
			},
			env: env,
		},
		"success alert": {
			outRequest: &concourse.OutRequest{
				Source: concourse.Source{URL: ok.URL},
				Params: concourse.OutParams{AlertType: "success"},
			},
			want: &concourse.OutResponse{
				Version: concourse.Version{"ver": "static"},
				Metadata: []concourse.Metadata{
					{Name: "type", Value: "success"},
					{Name: "alerted", Value: "true"},
				},
			},
			env: env,
		},
		"failed alert": {
			outRequest: &concourse.OutRequest{
				Source: concourse.Source{URL: ok.URL},
				Params: concourse.OutParams{AlertType: "failed"},
			},
			want: &concourse.OutResponse{
				Version: concourse.Version{"ver": "static"},
				Metadata: []concourse.Metadata{
					{Name: "type", Value: "failed"},
					{Name: "alerted", Value: "true"},
				},
			},
			env: env,
		},
		"started alert": {
			outRequest: &concourse.OutRequest{
				Source: concourse.Source{URL: ok.URL},
				Params: concourse.OutParams{AlertType: "started"},
			},
			want: &concourse.OutResponse{
				Version: concourse.Version{"ver": "static"},
				Metadata: []concourse.Metadata{
					{Name: "type", Value: "started"},
					{Name: "alerted", Value: "true"},
				},
			},
			env: env,
		},
		"aborted alert": {
			outRequest: &concourse.OutRequest{
				Source: concourse.Source{URL: ok.URL},
				Params: concourse.OutParams{AlertType: "aborted"},
			},
			want: &concourse.OutResponse{
				Version: concourse.Version{"ver": "static"},
				Metadata: []concourse.Metadata{
					{Name: "type", Value: "aborted"},
					{Name: "alerted", Value: "true"},
				},
			},
			env: env,
		},
		"custom alert": {
			outRequest: &concourse.OutRequest{
				Source: concourse.Source{URL: ok.URL},
				Params: concourse.OutParams{
					AlertType: "non-existent-type",
					Message:   "Deploying",
					Color:     "#ffffff",
				},
			},
			want: &concourse.OutResponse{
				Version: concourse.Version{"ver": "static"},
				Metadata: []concourse.Metadata{
					{Name: "type", Value: "default"},
					{Name: "alerted", Value: "true"},
				},
			},
			env: env,
		},
		"override channel at Source": {
			outRequest: &concourse.OutRequest{
				Source: concourse.Source{URL: ok.URL, Channel: "#source"},
			},
			want: &concourse.OutResponse{
				Version: concourse.Version{"ver": "static"},
				Metadata: []concourse.Metadata{
					{Name: "type", Value: "default"},
					{Name: "alerted", Value: "true"},
				},
			},
			env: env,
		},
		"override channel at Params": {
			outRequest: &concourse.OutRequest{
				Source: concourse.Source{URL: ok.URL, Channel: "#source"},
				Params: concourse.OutParams{Channel: "#params"},
			},
			want: &concourse.OutResponse{
				Version: concourse.Version{"ver": "static"},
				Metadata: []concourse.Metadata{
					{Name: "type", Value: "default"},
					{Name: "alerted", Value: "true"},
				},
			},
			env: env,
		},
		"disable alert": {
			outRequest: &concourse.OutRequest{
				Source: concourse.Source{URL: bad.URL},
				Params: concourse.OutParams{Disable: true},
			},
			want: &concourse.OutResponse{
				Version: concourse.Version{"ver": "static"},
				Metadata: []concourse.Metadata{
					{Name: "type", Value: "default"},
					{Name: "alerted", Value: "false"},
				},
			},
			env: env,
		},
		"error without Google Chat URL": {
			outRequest: &concourse.OutRequest{
				Source: concourse.Source{URL: ""},
			},
			env: env,
			err: true,
		},
		"error with bad request": {
			outRequest: &concourse.OutRequest{
				Source: concourse.Source{URL: bad.URL},
			},
			env: env,
			err: true,
		},
		"error without basic auth for fixed type": {
			outRequest: &concourse.OutRequest{
				Source: concourse.Source{URL: ok.URL, Username: "", Password: ""},
				Params: concourse.OutParams{AlertType: "fixed"},
			},
			env: env,
			err: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			for k, v := range c.env {
				os.Setenv(k, v)
			}

			got, err := out(c.outRequest, "")
			if err != nil && !c.err {
				t.Fatalf("unexpected error from out:\n\t(ERR): %s", err)
			} else if err == nil && c.err {
				t.Fatalf("expected an error from out:\n\t(GOT): nil")
			} else if !cmp.Equal(got, c.want) {
				t.Fatalf("unexpected concourse.OutResponse value from out:\n\t(GOT): %#v\n\t(WNT): %#v\n\t(DIFF): %v", got, c.want, cmp.Diff(got, c.want))
			}
		})
	}
}

func TestPreviousBuildName(t *testing.T) {
	cases := map[string]struct {
		build string
		want  string

		err bool
	}{
		"standard": {
			build: "6",
			want:  "6",
		},
		"rerun 1": {
			build: "6.1",
			want:  "6",
		},
		"rerun x": {
			build: "6.2",
			want:  "6.1",
		},
		"error 1": {
			build: "X",
			err:   true,
		},
		"error x": {
			build: "6.X",
			err:   true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := previousBuildName(c.build)
			if err != nil && !c.err {
				t.Fatalf("unexpected error from previousBuildName:\n\t(ERR): %s", err)
			} else if err == nil && c.err {
				t.Fatalf("expected an error from previousBuildName:\n\t(GOT): nil")
			} else if err != nil && c.err {
				return
			}

			if err != nil {
				t.Fatalf("unexpected value from previousBuildName:\n\t(GOT): %#v\n\t(WNT): %#v", got, c.want)
			}
		})
	}
}
