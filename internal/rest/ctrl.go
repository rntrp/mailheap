package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rntrp/mailheap/internal/config"
	"github.com/rntrp/mailheap/internal/model"
	"github.com/rntrp/mailheap/internal/msg"
	"github.com/rntrp/mailheap/internal/storage"
)

type Controller interface {
	Index(w http.ResponseWriter, r *http.Request)
	GetEml(w http.ResponseWriter, r *http.Request)
	DeleteMails(w http.ResponseWriter, r *http.Request)
	SeekMails(w http.ResponseWriter, r *http.Request)
	UploadMail(w http.ResponseWriter, r *http.Request)
}

func New(s storage.MailStorage, a msg.StoreMailSvc) Controller {
	return &ctrl{storage: s, storeMail: a}
}

type ctrl struct {
	storage   storage.MailStorage
	storeMail msg.StoreMailSvc
}

func (c *ctrl) GetEml(w http.ResponseWriter, r *http.Request) {
	addSecurityHeaders(w.Header())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "numeric ID could not be parsed", http.StatusBadRequest)
		return
	}
	eml, err := c.storage.GetMime(id)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	} else if len(eml) == 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	b := []byte(eml)
	w.Header().Add("Content-Type", "message/rfc822")
	w.Header().Add("Content-Length", strconv.Itoa(len(b)))
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v.eml\"", id))
	w.Write(b)
}

type DeleteMailsResult struct {
	NumDeleted int64
}

func (c *ctrl) DeleteMails(w http.ResponseWriter, r *http.Request) {
	addSecurityHeaders(w.Header())
	var numDeleted int64
	var err error
	if idQuery, ok := r.URL.Query()["id"]; ok {
		if len(idQuery) != 1 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		str := strings.Split(idQuery[0], ",")
		ids := make([]int64, 0, len(str))
		for _, s := range str {
			if id, err := strconv.ParseInt(s, 10, 64); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			} else {
				ids = append(ids, id)
			}
		}
		numDeleted, err = c.storage.DeleteMails(ids...)
	} else {
		numDeleted, err = c.storage.DeleteAllMails()
	}
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	} else if b, err := json.Marshal(DeleteMailsResult{NumDeleted: numDeleted}); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Content-Length", strconv.Itoa(len(b)))
		w.Write(b)
	}
}

type SeekMailsResult struct {
	Id    int64        `json:"id"`
	Total int64        `json:"total"`
	Limit int          `json:"limit"`
	Size  int          `json:"size"`
	Data  []model.Mail `json:"data"`
}

func (c *ctrl) SeekMails(w http.ResponseWriter, r *http.Request) {
	addSecurityHeaders(w.Header())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "numeric ID could not be parsed", http.StatusBadRequest)
		return
	} else if id <= 0 {
		id = math.MaxInt64
	}
	total, err := c.storage.CountMails()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	limit := parseLimit(r.URL.Query())
	mails, err := c.storage.SeekMails(id, limit)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(SeekMailsResult{
		Id:    id,
		Total: total,
		Limit: limit,
		Size:  len(mails),
		Data:  mails})
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}

func parseLimit(query url.Values) int {
	const def = 20
	const min = 10
	const max = 100
	value := query["limit"]
	if len(value) <= 0 {
		return def
	}
	limit, err := strconv.Atoi(value[0])
	switch {
	case err != nil:
		return def
	case limit < min:
		return min
	case limit > max:
		return max
	default:
		return limit
	}
}

func (c *ctrl) UploadMail(w http.ResponseWriter, r *http.Request) {
	addSecurityHeaders(w.Header())
	if !setupFileSizeChecks(w, r) {
		return
	}
	maxMemory := coerceMemoryBufferSize(config.GetHTTPUploadMemoryBufferSize())
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}
	f, _, err := r.FormFile("eml")
	if err != nil {
		http.Error(w, "file 'eml' is missing", http.StatusBadRequest)
		return
	}
	defer f.Close()
	if err := c.storeMail.StoreMail(f); err != nil {
		log.Println(err)
		http.Error(w, "eml could not be stored", http.StatusBadRequest)
	}
}
