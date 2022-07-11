package gapi

import (
	"context"
	"database/sql"

	db "github.com/simplebank/db/sqlc"
	"github.com/simplebank/pb"
	"github.com/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "username not found:%s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get user:%s", err)

	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password:%s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create token:%s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh_token:%s", err)
	}
	mtdt := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session:%s", err)
	}

	rsp := &pb.LoginUserResponse{
		SessionId:            session.ID.String(),
		AccessToken:          accessToken,
		AccessTokenExpireAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:         refreshToken,
		RefreshTokenExpireAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                 convertUser(user),
	}

	return rsp, nil
}
