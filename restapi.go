package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/noptics/golog"
	"github.com/sethjback/noptics/registry/data"
	"github.com/sethjback/noptics/registry/registrygrpc"
)

const usage = `{"message":"path not found","usage":[{"description":"get data about a channel on a cluster","path":"/:cluster/:channel","method":"GET","example":"GET /nats1/payoutRequests","response":{"cluster":"nats1","channel":"payoutRequests","files":[{"name":"request.proto","data":"<base64 encoded file data>"}]}},{"description":"set channel files (raw protobuf files that will be used to decode messages)","path":"/:cluster/:channel","method":"POST","example":"POST /nats1/payoutReqeusts","body":{"cluster":"nats1 (optional)","channel":"payoutRequests (optional)","files":[{"name":"request.proto","data":"<base64 encoded file data>"}]}},{"description":"set channel message (which proto message should be used when decoding)","path":"/:cluster/:channel/message","method":"POST","example":"POST /nats1/payoutReqeusts/message","body":{"name":"request"}}]}`

type RESTServer struct {
	db      data.Store
	l       golog.Logger
	hs      http.Server
	errChan chan error
}

type RESTError struct {
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
	Details     string `json:"details"`
}

type RESTChannelData struct {
	Cluster string               `json:"cluster"`
	Channel string               `json:"channel"`
	Message string               `json:"message"`
	Files   []*registrygrpc.File `json:"files"`
}

func NewRestServer(db data.Store, port string, errChan chan error, l golog.Logger) *RESTServer {
	rs := &RESTServer{
		db:      db,
		l:       l,
		errChan: errChan,
	}

	rs.hs = http.Server{Addr: ":" + port, Handler: rs.Router()}

	l.Infow("starting rest server", "port", port)
	go func() {
		if err := rs.hs.ListenAndServe(); err != nil {
			rs.errChan <- err
		}
	}()

	return rs
}

func (rs *RESTServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return rs.hs.Shutdown(ctx)
}

func (rs *RESTServer) SaveFiles(cluster, channel string, body []byte) (int, error) {
	sfr := &registrygrpc.SaveFilesRequest{}
	err := jsonpb.Unmarshal(bytes.NewBuffer(body), sfr)
	if err != nil {
		es, _ := json.Marshal(&RESTError{Message: "unable to save files", Description: "error parsing POST body", Details: err.Error()})
		return 400, fmt.Errorf("%s", es)
	}

	err = rs.db.SaveFiles(sfr.GetCluster(), sfr.GetChannel(), sfr.GetFiles())
	if err != nil {
		es, _ := json.Marshal(&RESTError{Message: "unable to save files", Description: "error storing data", Details: err.Error()})
		return 400, fmt.Errorf("%s", es)
	}

	return 200, nil
}

func (rs *RESTServer) GetChannel(cluster, channel string) (string, *registrygrpc.SaveFilesRequest, int, error) {
	message, fs, err := rs.db.GetChannelData(cluster, channel)
	if err != nil {
		es, _ := json.Marshal(&RESTError{Message: "unable to get file", Details: err.Error()})
		return "", nil, 400, fmt.Errorf("%s", es)
	}

	if fs == nil {
		return "", nil, 404, nil
	}

	return message, &registrygrpc.SaveFilesRequest{Cluster: cluster, Channel: channel, Files: fs}, 200, nil
}

func (rs *RESTServer) SetMessage(cluster, channel string, bdy []byte) (int, error) {
	BodyData := map[string]string{}
	err := json.Unmarshal(bdy, &BodyData)
	if err != nil {
		es, _ := json.Marshal(&RESTError{Message: "unable to set channel message", Description: "error parsing POST body", Details: err.Error()})
		return 400, fmt.Errorf("%s", es)
	}

	msg, ok := BodyData["message"]
	if !ok {
		es, _ := json.Marshal(&RESTError{Message: "unable to set channel message", Description: "must provide 'message' field in body"})
		return 400, fmt.Errorf("%s", es)
	}

	err = rs.db.SetChannelMessage(cluster, channel, msg)
	if err != nil {
		es, _ := json.Marshal(&RESTError{Message: "unable to set message for channel", Details: err.Error()})
		return 400, fmt.Errorf("%s", es)
	}

	return 200, nil
}

func (rs *RESTServer) wrapRoute(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method == "GET" {
		message, files, status, err := rs.GetChannel(ps.ByName("cluster"), ps.ByName("channel"))
		var body []byte
		if err != nil {
			body = []byte(err.Error())
		} else {
			if len(message) == 0 && files == nil {
				status = 404
			} else {
				rcd := RESTChannelData{
					Message: message,
					Channel: ps.ByName("channel"),
					Cluster: ps.ByName("cluster"),
				}

				if files != nil {
					rcd.Files = files.Files
				}
				body, _ = json.Marshal(rcd)
			}
		}

		writeReply(status, body, w)
	}

	if r.Method == "POST" {
		pathString := strings.Split(r.URL.Path, "/")[1:]

		bdy, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeReply(400, []byte(fmt.Sprintf(`{"message": "error reading post body", "details": "%s"}`, err.Error)), w)
			return
		}

		var status int
		if pathString[2] == "files" {
			status, err = rs.SaveFiles(pathString[0], pathString[1], bdy)
		} else {
			status, err = rs.SetMessage(pathString[0], pathString[1], bdy)
		}

		var body []byte
		if err != nil {
			body = []byte(err.Error())
		}
		writeReply(status, body, w)
	}

}

func writeReply(status int, body []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if len(body) != 0 {
		w.Write(body)
	}
}

func (rs *RESTServer) Router() *httprouter.Router {
	r := httprouter.New()
	r.POST("/:cluster/:channel/files", rs.wrapRoute)
	r.POST("/:cluster/:channel/message", rs.wrapRoute)
	r.GET("/:cluster/:channel", rs.wrapRoute)

	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write([]byte(usage))
	})

	return r
}
