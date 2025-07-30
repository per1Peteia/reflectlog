package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostJSON(t *testing.T) {

	t.Run("happy end to end in memory", func(t *testing.T) {
		var srvPayload struct {
			A string `json:"a"`
			B string `json:"b"`
			C string `json:"c"`
		}

		tres := "hi, this is the test server"
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "POST", r.Method)

			data, err := io.ReadAll(r.Body)
			assert.NoError(t, err)

			err = json.Unmarshal(data, &srvPayload)
			assert.NoError(t, err)

			fmt.Fprint(w, tres)
		}))
		defer ts.Close()

		payload := struct {
			A string `json:"a"`
			B string `json:"b"`
			C string `json:"c"`
		}{
			A: "s1",
			B: "s2",
			C: "s3",
		}

		res, err := postJSON(ts.URL, payload)
		assert.NoError(t, err)

		s, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		res.Body.Close()

		assert.Equal(t, tres, string(s))

		assert.Equal(t, "s1", srvPayload.A)
		assert.Equal(t, "s2", srvPayload.B)
		assert.Equal(t, "s3", srvPayload.C)
	})
}
