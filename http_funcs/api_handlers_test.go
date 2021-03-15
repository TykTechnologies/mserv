package http_funcs_test

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io"
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

func setupServerAndTempDir(t *testing.T) *http_funcs.HttpServ {
	t.Helper()
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
	cfg.Mserv.RetainUploads = false // would default to false anyway, but we set it explicitly for these tests

	fileCountPath := filepath.Join(cfg.Mserv.MiddlewarePath, "plugins")
	is.NoErr(os.MkdirAll(fileCountPath, 0o700)) // could not prepare upload directory

	cfgBytes, err := json.Marshal(cfg)
	is.NoErr(err)                                            // could not marshal config struct
	is.NoErr(ioutil.WriteFile(cfgFilePath, cfgBytes, 0o600)) // could not write config out to file

	// Make sure the config in use is current and aligned with what was just established here
	config.Reload()

	// Return a new server
	return http_funcs.NewServer("http://mserv.io", &mock.Storage{})
}

func prepareRequest(t *testing.T, getAttachment func(*testing.T) []byte) *http.Request {
	t.Helper()
	is := is.New(t)

	// Get the attachment ready
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	part, err := writer.CreateFormFile(http_funcs.UploadFormField, "attachment.zip")
	is.NoErr(err) // could not create part for file being uploaded

	size, err := part.Write(getAttachment(t))
	is.NoErr(err) // could not write test attachment file bytes to multipart

	is.NoErr(writer.Close()) // could not close multipart writer cleanly
	t.Logf("attachment size is '%d' bytes", size)

	// Get the request ready with the attachment
	req := httptest.NewRequest(http.MethodPost, "http://mserv.io/api/mw", reqBody)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	is.NoErr(req.ParseForm()) // could not parse form on HTTP request

	return req
}

func getCompressed(t *testing.T, flat bool, files ...string) []byte {
	t.Helper()
	is := is.New(t)

	buf := bytes.Buffer{}
	w := zip.NewWriter(&buf)

	for _, file := range files {
		var (
			f   io.Writer
			err error
		)

		if flat {
			f, err = w.Create(filepath.Base(file))
		} else {
			f, err = w.Create(file)
		}
		is.NoErr(err) // could not create file in zip.Writer

		_, err = f.Write(getUncompressed(t, file))
		is.NoErr(err) // could not write body of file into zip.Writer
	}

	is.NoErr(w.Close()) // could not close zip.Writer

	return buf.Bytes()
}

func getUncompressed(t *testing.T, files ...string) []byte {
	t.Helper()
	is := is.New(t)

	buf := bytes.Buffer{}

	for _, file := range files {
		body, err := ioutil.ReadFile(file) //nolint:gosec // File paths are hard-coded to these test helpers
		is.NoErr(err)                      // could not read file

		_, err = buf.Write(body)
		is.NoErr(err) // could not write body of file into buffer
	}

	return buf.Bytes()
}

func compressedTestData(t *testing.T) []byte {
	return getCompressed(t, true, "testdata/uncompressed/manifest.json", "testdata/uncompressed/middleware.py")
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
