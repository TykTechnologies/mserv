package http_funcs_test

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"

	config "github.com/TykTechnologies/mserv/conf"
	"github.com/TykTechnologies/mserv/http_funcs"
	"github.com/TykTechnologies/mserv/util/storage/mock"
)

type addMWTestCase struct {
	testBodyBytes  func(*testing.T) []byte
	expectedStatus int
}

var addMWTestCases = map[string]addMWTestCase{
	"Compressed (ZIP) upload is OK": {
		testBodyBytes:  compressedTestData,
		expectedStatus: http.StatusOK,
	},
	"Uncompressed upload is not OK": {
		testBodyBytes:  uncompressedTestData,
		expectedStatus: http.StatusUnprocessableEntity,
	},
	"Uncompressed manifest by itself (no Python) is not OK": {
		testBodyBytes:  uncompressedJSONTestData,
		expectedStatus: http.StatusUnprocessableEntity,
	},
	"Uncompressed Python by itself (no manifest JSON) is not OK": {
		testBodyBytes:  uncompressedPythonTestData,
		expectedStatus: http.StatusUnprocessableEntity,
	},
	"Random byte stream can not be classified/detected": {
		testBodyBytes:  randomByteStream,
		expectedStatus: http.StatusUnsupportedMediaType,
	},
}

func TestAddMW(t *testing.T) {
	is := is.New(t)

	// Prepare the config file and the plugin uploads directory
	testTemp := t.TempDir()
	t.Logf("operating out of '%s' directory", testTemp)

	cfgFilePath := filepath.Join(testTemp, "mserv.conf")
	is.NoErr(os.Setenv("TYK_MSERV_CONFIG", cfgFilePath)) // could not set config file location in environment

	cfg := config.Config{}
	cfg.Mserv.FileStore = &config.FileStorage{}
	cfg.Mserv.FileStore.Kind = "local"
	cfg.Mserv.FileStore.Local = &config.LocalStore{}
	cfg.Mserv.FileStore.Local.ConfigKeyPath = filepath.Join(testTemp, "files")
	cfg.Mserv.MiddlewarePath = filepath.Join(testTemp, "middleware")
	cfg.Mserv.PluginDir = filepath.Join(testTemp, "plugins")

	cfgBytes, err := json.Marshal(cfg)
	is.NoErr(err)                                           // could not marshal config struct
	is.NoErr(ioutil.WriteFile(cfgFilePath, cfgBytes, 0600)) // could not write config out to file

	// Create a new server
	srv := http_funcs.NewServer("http://mserv.io", &mock.Storage{})

	for name, tc := range addMWTestCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			// Get the attachment ready
			reqBody := &bytes.Buffer{}
			writer := multipart.NewWriter(reqBody)
			part, err := writer.CreateFormFile(http_funcs.UploadFormField, "attachment.zip")
			is.NoErr(err) // could not create part for file being uploaded

			size, err := part.Write(tc.testBodyBytes(t))
			is.NoErr(err) // could not write compressed test file data to multipart

			is.NoErr(writer.Close()) // could not close multipart writer cleanly
			t.Logf("attachment size is '%d' bytes", size)

			// Get the request ready with the attachment
			req := httptest.NewRequest(http.MethodPost, "http://mserv.io/api/mw", reqBody)
			req.Header.Add("Content-Type", writer.FormDataContentType())

			is.NoErr(req.ParseForm()) // could not parse form on HTTP request
			req.Form.Add("store_only", "true")

			// Handle the request
			w := httptest.NewRecorder()
			srv.AddMW(w, req)

			// Parse the response
			resp := w.Result()
			defer is.NoErr(resp.Body.Close()) // could not close response body cleanly

			is.Equal(tc.expectedStatus, resp.StatusCode) // expected response status does not equal actual

			bod, errRead := ioutil.ReadAll(resp.Body)
			is.NoErr(errRead) // could not read response body
			t.Logf("response: %s %s", resp.Status, bod)
		})
	}
}

func compressedTestData(t *testing.T) []byte {
	t.Helper()

	buf := bytes.Buffer{}
	w := zip.NewWriter(&buf)
	files := []string{"testdata/uncompressed/manifest.json", "testdata/uncompressed/middleware.py"}

	for _, file := range files {
		f, err := w.Create(file)
		if err != nil {
			t.Fatalf("could not create file '%s' in zip.Writer: %v", file, err)
		}

		_, err = f.Write(getUncompressed(t, file))
		if err != nil {
			t.Fatalf("could not write body of file '%s' into zip.Writer: %v", file, err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatalf("could not close zip.Writer: %v", err)
	}

	return buf.Bytes()
}

func getUncompressed(t *testing.T, files ...string) []byte {
	t.Helper()

	buf := bytes.Buffer{}

	for _, file := range files {
		body, err := ioutil.ReadFile(file) //nolint:gosec // File paths are hard-coded to these test helpers
		if err != nil {
			t.Fatalf("could not read file '%s': %v", file, err)
		}

		if _, err := buf.Write(body); err != nil {
			t.Fatalf("could not write body of file '%s' into buffer: %v", file, err)
		}
	}

	return buf.Bytes()
}

func uncompressedTestData(t *testing.T) []byte {
	return getUncompressed(t, "testdata/uncompressed/manifest.json", "testdata/uncompressed/middleware.py")
}

func uncompressedJSONTestData(t *testing.T) []byte {
	return getUncompressed(t, "testdata/uncompressed/manifest.json")
}

func uncompressedPythonTestData(t *testing.T) []byte {
	return getUncompressed(t, "testdata/uncompressed/middleware.py")
}

func randomByteStream(t *testing.T) []byte {
	t.Helper()
	is := is.New(t)

	bytes := make([]byte, 1024)
	count, err := rand.Read(bytes)
	is.NoErr(err)
	is.Equal(1024, count)

	return bytes
}
