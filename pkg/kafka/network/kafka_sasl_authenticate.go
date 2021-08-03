package network

import (
	"github.com/paashzj/kafka_go/pkg/kafka/codec"
	"github.com/paashzj/kafka_go/pkg/kafka/log"
	"github.com/paashzj/kafka_go/pkg/kafka/network/context"
	"github.com/paashzj/kafka_go/pkg/kafka/service"
	"github.com/panjf2000/gnet"
	"k8s.io/klog/v2"
)

func (s *Server) SaslAuthenticate(frame []byte, version int16, context *context.NetworkContext) ([]byte, gnet.Action) {
	if version == 1 || version == 2 {
		return s.ReactSaslHandshakeAuthVersion(frame, version, context)
	}
	klog.Error("unknown handshake auth version ", version)
	return nil, gnet.Close
}

func (s *Server) ReactSaslHandshakeAuthVersion(frame []byte, version int16, context *context.NetworkContext) ([]byte, gnet.Action) {
	req, err := codec.DecodeSaslHandshakeAuthReq(frame, version)
	if err != nil {
		return nil, gnet.Close
	}
	log.Codec().Info("sasl handshake request ", req)
	saslHandshakeResp := codec.NewSaslHandshakeAuthResp(req.CorrelationId)
	authResult, errorCode := service.SaslAuth(s.kafkaImpl, req.Username, req.Password)
	if errorCode != 0 {
		return nil, gnet.Close
	}
	if authResult {
		context.Authed(true)
		return saslHandshakeResp.Bytes(version), gnet.None
	} else {
		return nil, gnet.Close
	}
}
