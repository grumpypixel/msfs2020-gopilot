package app

import (
	"sync"

	"github.com/grumpypixel/msfs2020-simconnect-go/simconnect"
)

type Var struct {
	Name, Moniker string
}

type Request struct {
	ClientID string
	Meta     string
	Vars     map[simconnect.DWord]*Var
}

func NewRequest(clientID string, meta string) *Request {
	return &Request{
		ClientID: clientID,
		Meta:     meta,
		Vars:     make(map[simconnect.DWord]*Var),
	}
}

func (req *Request) Add(defineID simconnect.DWord, name, moniker string) bool {
	if len(name) == 0 {
		return false
	}
	if len(moniker) == 0 {
		moniker = name
	}
	req.Vars[defineID] = &Var{name, moniker}
	return true
}

func (req *Request) GetVar(defineID simconnect.DWord) (string, string, bool) {
	if value, exists := req.Vars[defineID]; exists {
		return value.Name, value.Moniker, true
	}
	return "", "", false
}

type RequestManager struct {
	Requests []*Request
	mutex    sync.Mutex
}

func NewRequestManager() *RequestManager {
	mgr := &RequestManager{
		Requests: make([]*Request, 0),
	}
	return mgr
}

func (mgr *RequestManager) RequestCount() int {
	return len(mgr.Requests)
}

func (mgr *RequestManager) AddRequest(request *Request) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	mgr.Requests = append(mgr.Requests, request)
}

func (mgr *RequestManager) RefCount(simVarName string) int {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	count := 0
	for _, request := range mgr.Requests {
		for _, v := range request.Vars {
			if v.Name == simVarName {
				count++
			}
		}
	}
	return count
}
