package main

import (
	pbmsg "boardwallfloor/chatbox_message_module/pb/messaging/v1"
	"context"
	"database/sql"
)

type sqliteDB struct {
	pbmsg.UnimplementedMessagingServiceServer
	Db *sql.DB
}

func (sl *sqliteDB) SendMessage(ctx context.Context, in *pbmsg.SendMessageRequest) (*pbmsg.SendMessageResponse, error)
func (sl *sqliteDB) GetMessages(ctx context.Context, in *pbmsg.GetMessagesRequest) (*pbmsg.GetMessagesResponse, error)
func (sl *sqliteDB) MarkMessageAsDelivered(ctx context.Context, in *pbmsg.MarkMessageAsDeliveredRequest) (*pbmsg.MarkMessageAsDeliveredResponse, error)
