package mc

import (
	"github.com/delongchen/cdl-rcon/pkg/rcon"
	"log"
)

type Handler struct {
	s    *rcon.RCONSession
	pMap rcon.CMDTaskMap
}

func CoverSession(s *rcon.RCONSession) *Handler {
	return &Handler{
		s:    s,
		pMap: make(rcon.CMDTaskMap),
	}
}

func (handler *Handler) ExecCMD(cmd string) *rcon.BasicPacket {
	id := handler.s.ID
	handler.s.ID += 1

	ch := make(chan *rcon.BasicPacket)
	handler.pMap[id] = &rcon.CMDTask{
		Ch:  ch,
		CMD: cmd,
		Vec: make([]*rcon.BasicPacket, 0),
	}

	handler.s.In <- &rcon.BasicPacket{
		ID:   id,
		Body: cmd,
		Type: rcon.SERVERDATA_EXECCOMMAND,
	}

	return <-ch
}

func (handler *Handler) Start(quite chan struct{}) {
	if handler.s == nil {
		log.Fatal("not init session!")
	}

	err := handler.s.StartLoop()
	if err != nil {
		log.Fatal("session start fail!")
	}

	for p := range handler.s.Out {
		if p.ID != rcon.PING_ID && p.Type == rcon.SERVERDATA_RESPONSE_VALUE {
			task, exist := handler.pMap[p.ID]
			if !exist {
				continue
			}

			task.Ch <- p
			close(task.Ch)
			delete(handler.pMap, p.ID)
		}
	}

	quite <- struct{}{}
}
