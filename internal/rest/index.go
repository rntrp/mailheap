package rest

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/base64"
	"hash/fnv"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

//go:embed index.html
var indexHtml []byte

var indexHtmlEtag string

var indexHtmlGz []byte

var indexHtmlGzEtag string

func init() {
	hasher := fnv.New128a()
	if _, err := hasher.Write(indexHtml); err != nil {
		log.Fatal(err)
	}
	indexHtmlEtag = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	buf := new(bytes.Buffer)
	if gz, err := gzip.NewWriterLevel(buf, gzip.BestCompression); err != nil {
		log.Fatal(err)
	} else if _, err := gz.Write(indexHtml); err != nil {
		log.Fatal(err)
	} else {
		gz.Close()
	}
	indexHtmlGz = buf.Bytes()
	slog.Info("📄 index.html successfully gzipped & cached",
		"before", strconv.Itoa(len(indexHtml)),
		"after", strconv.Itoa(len(indexHtmlGz)))
	hasher.Reset()
	if _, err := hasher.Write(indexHtmlGz); err != nil {
		log.Fatal(err)
	}
	indexHtmlGzEtag = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (c *ctrl) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	acceptEncoding := r.Header.Get("Accept-Encoding")
	acceptEncoding = strings.ToLower(acceptEncoding)
	addSecurityHeaders(w.Header())
	if strings.Contains(acceptEncoding, "gzip") {
		w.Header().Add("Content-Encoding", "gzip")
		w.Header().Add("Content-Length", strconv.Itoa(len(indexHtmlGz)))
		w.Header().Add("ETag", indexHtmlGzEtag)
		if r.Header.Get("If-None-Match") == indexHtmlGzEtag {
			w.WriteHeader(http.StatusNotModified)
		} else {
			w.Write(indexHtmlGz)
		}
	} else {
		w.Header().Add("Content-Length", strconv.Itoa(len(indexHtml)))
		w.Header().Add("ETag", indexHtmlEtag)
		if r.Header.Get("If-None-Match") == indexHtmlEtag {
			w.WriteHeader(http.StatusNotModified)
		} else {
			w.Write(indexHtml)
		}
	}
}
