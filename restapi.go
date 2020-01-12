package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/julienschmidt/httprouter"
	"github.com/noptics/golog"
	"github.com/noptics/registry/data"
	"github.com/noptics/registry/registrygrpc"
)

const usage = `{"message":"path not found","usage":[{"description":"get data about a channel on a cluster","path":"/:cluster/:channel","method":"GET","example":"GET /nats1/payoutRequests","response":{"cluster":"nats1","channel":"payoutRequests","files":[{"name":"request.proto","data":"<base64 encoded file data>"}]}},{"description":"set channel files (raw protobuf files that will be used to decode messages)","path":"/:cluster/:channel","method":"POST","example":"POST /nats1/payoutReqeusts","body":{"cluster":"nats1 (optional)","channel":"payoutRequests (optional)","files":[{"name":"request.proto","data":"<base64 encoded file data>"}]}},{"description":"set channel message (which proto message should be used when decoding)","path":"/:cluster/:channel/message","method":"POST","example":"POST /nats1/payoutReqeusts/message","body":{"name":"request"}}]}`

// RESTServer wraps the routes and handlers
type RESTServer struct {
	db      data.Store
	l       golog.Logger
	hs      http.Server
	context *Context
	errChan chan error
}

// RESTError is the standard structure for rest api errors
type RESTError struct {
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
	Details     string `json:"details"`
}

// RESTChannelData is the standardized reply from the api
type RESTChannelData struct {
	Cluster string               `json:"cluster"`
	Channel string               `json:"channel"`
	Message string               `json:"message"`
	Files   []*registrygrpc.File `json:"files"`
}

type RESTStatusData struct {
	Context
}

type RouteReply struct {
	Code        int
	Error       *RESTError
	ChannelData *RESTChannelData
	StatusData  *RESTStatusData
	Channels    []string
}

type RouteHandler func(*http.Request, httprouter.Params) *RouteReply

func newReplyError(message, description, details string, status int) *RouteReply {
	return &RouteReply{
		Code: status,
		Error: &RESTError{
			Message:     message,
			Description: description,
			Details:     details,
		},
	}
}

func NewRestServer(db data.Store, port string, errChan chan error, l golog.Logger, context *Context) *RESTServer {
	rs := &RESTServer{
		db:      db,
		l:       l,
		errChan: errChan,
		context: context,
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

func (rs *RESTServer) Status(r *http.Request, ps httprouter.Params) *RouteReply {
	return &RouteReply{Code: 200, StatusData: &RESTStatusData{*rs.context}}
}

// SaveFiles takes the POST body and saves the data to the db.Store. Prefernce is given to
// the channel and cluster parameters provided in the body
// the router params expected are :channel and :cluster
func (rs *RESTServer) SaveFiles(r *http.Request, ps httprouter.Params) *RouteReply {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return newReplyError("request error", "unble to read POST body", err.Error(), 500)
	}

	sfr := &registrygrpc.SaveFilesRequest{}
	err = jsonpb.Unmarshal(bytes.NewBuffer(body), sfr)
	if err != nil {
		return newReplyError("unable to save files", "error parsing POST body", err.Error(), 400)
	}

	if len(sfr.Channel) == 0 {
		sfr.Channel = ps.ByName("channel")
	}

	if len(sfr.ClusterID) == 0 {
		sfr.ClusterID = ps.ByName("cluster")
	}

	if len(sfr.ClusterID) == 0 || len(sfr.Channel) == 0 {
		return newReplyError("unable to save files", "invalid parameters provided", "must provide channel and cluster", 400)
	}

	err = rs.db.SaveFiles(sfr.GetClusterID(), sfr.GetChannel(), sfr.GetFiles())
	if err != nil {
		return newReplyError("unable to save files", "error storing data", err.Error(), 400)
	}

	return &RouteReply{Code: 200}
}

// GetChannel returns the saved channel data
func (rs *RESTServer) GetChannel(r *http.Request, ps httprouter.Params) *RouteReply {
	channel := ps.ByName("channel")
	cluster := ps.ByName("cluster")

	message, fs, err := rs.db.GetChannelData(cluster, channel)
	if err != nil {
		return newReplyError("request error", "trouble reading channel data from store", err.Error(), 500)
	}

	if fs == nil {
		return &RouteReply{Code: 404}
	}

	return &RouteReply{
		Code: 200,
		ChannelData: &RESTChannelData{
			Message: message,
			Channel: channel,
			Cluster: cluster,
			Files:   fs,
		},
	}
}

// GetChannels returns a list of all configured channels for the provided cluster
func (rs *RESTServer) GetChannels(r *http.Request, ps httprouter.Params) *RouteReply {
	cluster := ps.ByName("cluster")

	channels, err := rs.db.GetChannels(cluster)
	if err != nil {
		return newReplyError("unable to get channel list", "trouble reading channel data from store", err.Error(), 500)
	}

	return &RouteReply{Code: 200, Channels: channels}
}

// SetMessage sets the root message that will be sent over a channel.
func (rs *RESTServer) SetMessage(r *http.Request, ps httprouter.Params) *RouteReply {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return newReplyError("request error", "unble to read POST body", err.Error(), 500)
	}

	BodyData := map[string]string{}
	err = json.Unmarshal(body, &BodyData)
	if err != nil {
		return newReplyError("unable to set channel message", "error parsing POST body", err.Error(), 400)
	}

	msg, ok := BodyData["message"]
	if !ok {
		return newReplyError("unable to set channel message", "must provide 'message' field in body", "", 400)
	}

	channel := ps.ByName("channel")
	cluster := ps.ByName("cluster")

	err = rs.db.SetChannelMessage(cluster, channel, msg)
	if err != nil {
		return newReplyError("unable to set channel message", "data store error", err.Error(), 400)
	}

	return &RouteReply{Code: 200}
}

// SetChannelData saves all relevent data related to a given channel
func (rs *RESTServer) SetChannelData(r *http.Request, ps httprouter.Params) *RouteReply {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return newReplyError("request error", "unble to read POST body", err.Error(), 500)
	}

	type SetChannelData struct {
		Message string               `json:"message"`
		Files   []*registrygrpc.File `json:"files"`
	}

	BodyData := &SetChannelData{}
	err = json.Unmarshal(body, &BodyData)
	if err != nil {
		return newReplyError("unable to set channel data", "error parsing POST body", err.Error(), 400)
	}

	err = rs.db.SaveChannelData(ps.ByName("cluster"), ps.ByName("channel"), BodyData.Message, BodyData.Files)
	if err != nil {
		return newReplyError("unable to set channel message", "data store error", err.Error(), 400)
	}

	return &RouteReply{Code: 200}
}

func wrapRoute(rh RouteHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		response := rh(r, ps)

		setHeaders(w, r.Header)
		var body []byte
		var err error

		// TODO: bad
		if response.Error != nil {
			body, err = json.Marshal(response.Error)
		} else if response.ChannelData != nil {
			body, err = json.Marshal(response.ChannelData)
		} else if response.Channels != nil {
			body, err = json.Marshal(response.Channels)
		} else if response.StatusData != nil {
			body, err = json.Marshal(response.StatusData)
		}

		if err != nil {
			body = []byte(fmt.Sprintf(`{"message": "error marshalling response", "description":"unable to json encode data", "details":"%s"`, err.Error()))
			response.Code = 500
		}

		w.WriteHeader(response.Code)
		if len(body) != 0 {
			w.Write(body)
		}
	}
}

func setHeaders(w http.ResponseWriter, rh http.Header) {
	allowMethod := "GET, POST, PUT, DELETE, OPTIONS"
	allowHeaders := "Content-Type"
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Allow", allowMethod)
	w.Header().Set("Access-Control-Allow-Methods", allowMethod)
	w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

	o := rh.Get("Origin")
	if o == "" {
		o = "*"
	}
	w.Header().Set("Access-Control-Allow-Origin", o)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
}

func writeReply(status int, body []byte, w http.ResponseWriter) {
	w.WriteHeader(status)
	if len(body) != 0 {
		w.Write(body)
	}
}

func (rs *RESTServer) Router() *httprouter.Router {
	r := httprouter.New()

	r.HandleOPTIONS = true
	r.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w, r.Header)
		w.WriteHeader(http.StatusNoContent)
	})

	r.POST("/:cluster/:channel/files", wrapRoute(rs.SaveFiles))
	r.POST("/:cluster/:channel/message", wrapRoute(rs.SetMessage))
	r.GET("/:cluster/:channel", wrapRoute(rs.GetChannel))
	r.POST("/:cluster/:channel", wrapRoute(rs.SetChannelData))
	r.GET("/:cluster", wrapRoute(rs.GetChannels))
	r.GET("/", wrapRoute(rs.Status))

	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write([]byte(usage))
	})

	return r
}
