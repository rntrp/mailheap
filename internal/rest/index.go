package rest

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/base64"
	"hash/fnv"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//go:embed index.html
var indexHtml []byte

var indexHtmlEtag string

var indexHtmlGz []byte

var indexHtmlGzEtag string

//go:embed favicon.ico
var faviconIco []byte

var faviconIcoEtag string

//go:embed favicon.svg
var faviconSvg []byte

var faviconSvgEtag string

var faviconSvgGz []byte

var faviconSvgGzEtag string

//go:embed index.css
var indexCss []byte

var indexCssEtag string

var indexCssGz []byte

var indexCssGzEtag string

//go:embed index.js
var indexJs []byte

var indexJsEtag string

var indexJsGz []byte

var indexJsGzEtag string

//go:embed index.jsmimeparser.min.js
var jsmimeparser []byte

var jsmimeparserEtag string

var jsmimeparserGz []byte

var jsmimeparserGzEtag string

func InitIndex() {
	indexHtmlEtag = etag(indexHtml)
	indexHtmlGz = gz(indexHtml)
	indexHtmlGzEtag = etag(indexHtmlGz)

	faviconIcoEtag = etag(faviconIco)

	faviconSvgEtag = etag(faviconSvg)
	faviconSvgGz = gz(faviconSvg)
	faviconSvgGzEtag = etag(faviconSvgGz)

	indexCssEtag = etag(indexCss)
	indexCssGz = gz(indexCss)
	indexCssGzEtag = etag(indexCssGz)

	indexJsEtag = etag(indexJs)
	indexJsGz = gz(indexJs)
	indexJsGzEtag = etag(indexJsGz)

	jsmimeparserEtag = etag(jsmimeparser)
	jsmimeparserGz = gz(jsmimeparser)
	jsmimeparserGzEtag = etag(jsmimeparserGz)
}

func (c *ctrl) Index(w http.ResponseWriter, r *http.Request) {
	addHeaders(w.Header(), "text/html; charset=utf-8")
	addSecurityHeaders(w.Header())
	gzBroker(w, r, indexHtml, indexHtmlGz, indexHtmlEtag, indexHtmlGzEtag)
}

func (c *ctrl) IndexFaviconIco(w http.ResponseWriter, r *http.Request) {
	addHeaders(w.Header(), "image/x-icon")
	addSecurityHeaders(w.Header())
	etagBroker(w, r, faviconIco, faviconIcoEtag)
}

func (c *ctrl) IndexFaviconSvg(w http.ResponseWriter, r *http.Request) {
	addHeaders(w.Header(), "image/svg+xml")
	addSecurityHeaders(w.Header())
	gzBroker(w, r, faviconSvg, faviconSvgGz, faviconSvgEtag, faviconSvgGzEtag)
}

func (c *ctrl) IndexCss(w http.ResponseWriter, r *http.Request) {
	addHeaders(w.Header(), "text/css; charset=utf-8")
	addSecurityHeaders(w.Header())
	gzBroker(w, r, indexCss, indexCssGz, indexCssEtag, indexCssGzEtag)
}

func (c *ctrl) IndexJs(w http.ResponseWriter, r *http.Request) {
	addHeaders(w.Header(), "text/javascript; charset=utf-8")
	addSecurityHeaders(w.Header())
	gzBroker(w, r, indexJs, indexJsGz, indexJsEtag, indexJsGzEtag)
}

func (c *ctrl) IndexJsMimeParser(w http.ResponseWriter, r *http.Request) {
	addHeaders(w.Header(), "text/javascript; charset=utf-8")
	addSecurityHeaders(w.Header())
	gzBroker(w, r, jsmimeparser, jsmimeparserGz, jsmimeparserEtag, jsmimeparserGzEtag)
}

func addHeaders(hdr http.Header, contentType string) {
	hdr.Add("Content-Type", contentType)
	hdr.Add("Cache-Control", "max-age=3600, must-revalidate")
}

func gz(src []byte) []byte {
	buf := new(bytes.Buffer)
	if gz, err := gzip.NewWriterLevel(buf, gzip.BestCompression); err != nil {
		log.Fatal(err)
	} else if _, err := gz.Write(src); err != nil {
		log.Fatal(err)
	} else {
		gz.Close()
	}
	return buf.Bytes()
}

func etag(src []byte) string {
	hasher := fnv.New128a()
	if _, err := hasher.Write(indexHtml); err != nil {
		log.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func etagBroker(w http.ResponseWriter, r *http.Request, src []byte, srcEtag string) {
	w.Header().Add("Content-Length", strconv.Itoa(len(src)))
	w.Header().Add("ETag", srcEtag)
	if r.Header.Get("If-None-Match") == srcEtag {
		w.WriteHeader(http.StatusNotModified)
	} else {
		w.Write(src)
	}
}

func gzBroker(w http.ResponseWriter, r *http.Request, src, gz []byte, srcEtag, gzEtag string) {
	acceptEncoding := r.Header.Get("Accept-Encoding")
	acceptEncoding = strings.ToLower(acceptEncoding)
	if strings.Contains(acceptEncoding, "gzip") {
		w.Header().Add("Content-Encoding", "gzip")
		w.Header().Add("Content-Length", strconv.Itoa(len(gz)))
		w.Header().Add("ETag", gzEtag)
		if r.Header.Get("If-None-Match") == gzEtag {
			w.WriteHeader(http.StatusNotModified)
		} else {
			w.Write(gz)
		}
	} else {
		etagBroker(w, r, src, srcEtag)
	}
}
