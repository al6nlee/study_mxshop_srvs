package handler

import "study_mxshop_srvs/userop_srv/proto"

type UserOpServer struct {
	proto.UnimplementedAddressServer
	proto.UnimplementedUserFavServer
	proto.UnimplementedMessageServer
}
