package server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Arkadiyche/TP_proxy/database"
	"github.com/Arkadiyche/TP_proxy/models"
	"github.com/Arkadiyche/TP_proxy/utils"
	"github.com/jackc/pgx"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)



func NewServer(port string, db *pgx.ConnPool) *http.Server {
	return &http.Server{
		Addr:         port,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				checkPattern := `^/check/[0-9]+$`
				requestsPattern := `^/requests$`
				requestPattern := `^/request/[0-9]+$`
				if match, _ := regexp.Match(requestsPattern, []byte(r.URL.String())); match {
					RequestList(w, r, db)
				} else if match, _ := regexp.Match(requestPattern, []byte(r.URL.String())); match {
					RepeatRequest(r.URL.String(), w, r, db)
				} else if match, _ := regexp.Match(checkPattern, []byte(r.URL.String())); match {
					CheckWithParamMiner(r.URL.String(), w, r, db)
				} else {
					//fmt.Println(r.URL)
					handleHTTP(w, r, db)
				}
			}
		}),
	}
}

func handleHTTP(w http.ResponseWriter, r *http.Request,  db *pgx.ConnPool) {
	var resp *http.Response
	var err error
	err = database.LogRequest(r, db)
	if err != nil {
		//fmt.Println(err)
		return
	}
	switch r.Method {
	case "GET":
		resp, err = http.DefaultTransport.RoundTrip(r)
	case "POST":
		resp, err = http.Post(r.URL.String(), r.Header.Get("Content-Type"), r.Body)
	default:
		resp, err = http.Get(r.URL.String())
	}

	if err != nil {
		return
	}
	defer resp.Body.Close()
	for mime, val := range resp.Header {
		if mime == "Proxy-Connection" {
			continue
		}
		w.Header().Set(mime, val[0])
	}
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type")+"; charset=utf8")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	//fmt.Println("writer", w)
	return
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	//fmt.Println(dest_conn)
	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func RequestList(w http.ResponseWriter, r *http.Request, db *pgx.ConnPool) {
	result := database.GetAllRequests(db)
	answer, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(answer)
}

func RepeatRequest(url string, w http.ResponseWriter, r *http.Request, db *pgx.ConnPool) {
	buffer := strings.Split(url, "/")
	id, err := strconv.Atoi(buffer[2])
	if err != nil {
		return
	}
	request := database.GetRequest(id, db)
	r = &request
	//fmt.Println(request)
	http.Redirect(w, &request, request.URL.String(), 301)
	return
}

func CheckWithParamMiner(url string, w http.ResponseWriter, r *http.Request, db *pgx.ConnPool)  {
	flag := false
	buffer := strings.Split(url, "/")
	id, err := strconv.Atoi(buffer[2])
	if err != nil {
		return
	}
	request := database.GetRequest(id, db)

	if request.Method == "" {
		fmt.Println("request doesn't exist")
		return
	}
	//fmt.Println(request.URL)
	for _, val := range models.Params {
		randomString := utils.RandStringRunes()
		request.URL.RawQuery = val+"="+randomString
		//fmt.Println("!!!", request.URL, key)
		resp, err := http.DefaultTransport.RoundTrip(&request)
		if err != nil {
			fmt.Println("error with round trip")
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error with Read all")
			return
		}
		if strings.Contains(string(body), randomString) {
			w.Write([]byte(val + "-найден скрытый гет параметр\n"))
			flag = true
		}
	}
	if flag == false {
		w.Write([]byte("скрытые гет параметры не найдены\n"))
	}
}
