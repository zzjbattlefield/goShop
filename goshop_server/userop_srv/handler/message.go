package handler

import (
	"context"
	"goshop/userop_srv/global"
	"goshop/userop_srv/model"
	"goshop/userop_srv/proto"
)

//获取留言信息列表 带上用户id则只获取当前用户,不带用户id获取所有信息列表
func (*UseropServer) MessageList(ctx context.Context, req *proto.MessageRequest) (*proto.MessageListResponse, error) {
	messageModel := make([]model.LeavingMessages, 0)
	messageRsp := &proto.MessageListResponse{}
	res := global.DB.Where(&model.LeavingMessages{User: req.UserId}).Find(&messageModel)
	messageRsp.Total = int32(res.RowsAffected)
	messageList := make([]*proto.MessageResponse, 0)
	for _, info := range messageModel {
		messageList = append(messageList, &proto.MessageResponse{
			Id:          info.ID,
			UserId:      info.User,
			MessageType: info.MessageType,
			Subject:     info.Subject,
			Message:     info.Message,
			File:        info.File,
		})
	}
	messageRsp.Data = messageList
	return messageRsp, nil
}

func (*UseropServer) CreateMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	messageModel := model.LeavingMessages{
		User:        req.UserId,
		MessageType: req.MessageType,
		Subject:     req.Subject,
		Message:     req.Message,
		File:        req.File,
	}
	global.DB.Save(&messageModel)
	return &proto.MessageResponse{Id: messageModel.ID}, nil
}
