// Package service implements support protobuf for gprc server.
package proto

import (
	context "context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	"github.com/alaleks/shortener/internal/app/logger"
	"github.com/alaleks/shortener/internal/app/service"
	"github.com/alaleks/shortener/internal/app/storage"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
)

const (
	mdAuthName  = "Authorization"
	xrealIPName = "X-Real-IP"
)

// List of typical errors.
var (
	ErrInvalidMetadataAuth = errors.New("invalid authorization token in metadata")
	ErrorEmptyData         = errors.New("request does not contains data")
	ErrEmptyBatch          = errors.New("URL batching error, please check the source data")
	ErrorAccessDenied      = errors.New("access denied")
)

// Define a server struct that implements the grpc server interface.
type (
	Server struct {
		srv            UnsafeShortenerServer
		store          *storage.Store
		log            *logger.AppLogger
		trustedSubnets netip.Prefix
		secret         []byte
	}
)

// New creates a new grpc server.
func New(st *storage.Store, log *logger.AppLogger,
	secret []byte, trustedSubnet string) *Server {
	server := Server{
		srv:    UnimplementedShortenerServer{},
		store:  st,
		secret: secret,
	}

	if network, err := netip.ParsePrefix(trustedSubnet); err == nil {
		server.trustedSubnets = network
	}

	return &server
}

// ShortenURL implements URL shortening.
func (s *Server) ShortenURL(ctx context.Context, in *ShortenRequest) (*ShortenResponse, error) {
	if err := service.IsURL(in.Url); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, err := s.auth(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	shortURL, err := s.store.St.Add(in.Url, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ShortenResponse{
		Result:  shortURL,
		Success: true,
	}, nil
}

// GetStatAPI implements getting statistics on the use of a short URL.
func (s *Server) GetStat(ctx context.Context, in *StatRequest) (*StatResponse, error) {
	stat, err := s.store.St.Stat(in.Shortuid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &StatResponse{
		Shorturl:  stat.ShortURL,
		Longurl:   stat.LongURL,
		CreatedAt: stat.CreatedAt,
		Usage:     uint64(stat.Usage),
	}, nil
}

// GetUsersURL returns all shortened URLs for current user.
func (s *Server) GetUsersURL(ctx context.Context, in *Empty) (*UsersURL, error) {
	userID, err := s.definitionUser(ctx)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	out, err := s.store.St.GetUrlsUser(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	userURLS := UsersURL{
		Urls: make([]*UserURL, 0, len(out)),
	}

	for _, v := range out {
		userURLS.Urls = append(userURLS.Urls, &UserURL{
			LongUrl:  v.LongURL,
			ShortUrl: v.ShortUID,
		})
	}

	return &userURLS, nil
}

// ShortenURLBatch implements url batch shortening.
func (s *Server) ShortenURLBatch(ctx context.Context, in *ShortenBatchRequest) (*ShortenBatchResponse, error) {
	if len(in.Urls) == 0 {
		return nil, status.Error(codes.NotFound, ErrorEmptyData.Error())
	}

	userID, err := s.definitionUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	output, err := s.processingURLBatch(userID, in)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return output, nil
}

// ShortenDelete performs deletion all shortened URLs
func (s *Server) ShortenDelete(ctx context.Context, in *ShortenDeleteRequest) (*Empty, error) {
	if len(in.Urls) == 0 {
		return nil, status.Error(codes.NotFound, ErrorEmptyData.Error())
	}

	userID, err := s.definitionUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	var data = struct {
		userID         string
		shortUIDForDel []string
	}{
		userID:         userID,
		shortUIDForDel: in.Urls,
	}

	s.store.Pool.AddTask(func() error {
		err := s.store.St.DelUrls(data.userID,
			checkShortUID(data.shortUIDForDel...)...)
		if err != nil {
			return fmt.Errorf("deletion error: %w", err)
		}

		return nil
	})

	return &Empty{}, nil
}

// StatsInternal implement getting data about the number of shortened URLs
func (s *Server) StatsInternal(ctx context.Context, in *Empty) (*StatsInternalReponse, error) {
	var (
		xrealIP string
	)

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if values := md.Get(xrealIPName); len(values) > 0 {
			xrealIP = values[0]
		}
	}

	if len(xrealIP) == 0 {
		s.log.LZ.Error(ErrorAccessDenied)
		return nil, status.Error(codes.Unauthenticated, ErrorAccessDenied.Error())
	}

	realIP, err := netip.ParseAddr(xrealIP)
	if err != nil {
		s.log.LZ.Error(ErrorAccessDenied)
		return nil, status.Error(codes.Unauthenticated, ErrorAccessDenied.Error())
	}

	if !s.trustedSubnets.Contains(realIP) {
		s.log.LZ.Error(ErrorAccessDenied)
		return nil, status.Error(codes.Unauthenticated, ErrorAccessDenied.Error())
	}

	stat, err := s.store.St.GetInternalStats()
	if err != nil {
		s.log.LZ.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &StatsInternalReponse{
		Urls:  int64(stat.UrlsSize),
		Users: int64(stat.Users),
	}, nil
}

// auth performs the authorization this grpc server implements.
func (s *Server) auth(ctx context.Context) (string, error) {
	var (
		userID    string
		mdAuthVal string
	)

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if values := md.Get(mdAuthName); len(values) > 0 {
			mdAuthVal = values[0]
		}
	}

	if len(mdAuthVal) == 0 {
		uid := s.store.St.Create()

		header := metadata.Pairs("Authorization", s.encryptAuth(uid))
		grpc.SendHeader(ctx, header)

		return strconv.Itoa(int(uid)), nil
	}

	uid, err := s.decryptAuth(mdAuthVal)
	if err != nil {
		return userID, ErrInvalidMetadataAuth
	}

	header := metadata.Pairs("Authorization", mdAuthVal)
	grpc.SendHeader(ctx, header)

	return strconv.Itoa(int(uid)), nil
}

// decryptAuth perfoms decryption metadata value
// and returns user ID and error value.
func (s *Server) decryptAuth(mdAuth string) (uint, error) {
	signedVal, err := base64.URLEncoding.DecodeString(mdAuth)
	if err != nil {
		return 0, fmt.Errorf("metadata Authorization decoding error: %w", err)
	}

	signature := signedVal[:sha256.Size]
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(signedVal[sha256.Size:])
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal(signature, expectedSignature) {
		return 0, ErrInvalidMetadataAuth
	}

	return uint(int64(binary.LittleEndian.Uint64(signedVal[sha256.Size:]))), nil
}

// encryptAuth perfoms encryption user ID.
//
// Encryption is performed according to the sha256 algorithm.
// The secret key and user ID are used for encryption.
func (s *Server) encryptAuth(uid uint) string {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(uid))

	mac := hmac.New(sha256.New, s.secret)
	mac.Write(b)
	signature := mac.Sum(nil)
	signature = append(signature, b...)

	return base64.URLEncoding.EncodeToString(signature)
}

// definitionUser performs checking user ID from metadata.
func (s *Server) definitionUser(ctx context.Context) (string, error) {
	var (
		userID    string
		mdAuthVal string
	)

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if values := md.Get(mdAuthName); len(values) > 0 {
			mdAuthVal = values[0]
		}
	}

	if len(mdAuthVal) == 0 {
		return userID, ErrInvalidMetadataAuth
	}

	uid, err := s.decryptAuth(mdAuthVal)
	if err != nil {
		return userID, ErrInvalidMetadataAuth
	}

	header := metadata.Pairs("Authorization", mdAuthVal)
	grpc.SendHeader(ctx, header)

	return strconv.Itoa(int(uid)), nil
}

// processingURLBatch performs processing of batch of URLs.
func (s *Server) processingURLBatch(userID string, in *ShortenBatchRequest) (*ShortenBatchResponse, error) {
	out := ShortenBatchResponse{
		Urls: make([]*ShortenBatchResponseItem, 0, len(in.Urls)),
	}

	for _, item := range in.Urls {
		err := service.IsURL(item.OriginalUrl)

		if err == nil {
			shortURL := s.store.St.AddBatch(item.OriginalUrl, userID, item.CorrelationId)
			out.Urls = append(out.Urls, &ShortenBatchResponseItem{
				CorrelationId: item.CorrelationId,
				ShortUrl:      shortURL,
			})
		} else {
			out.Urls = append(out.Urls, &ShortenBatchResponseItem{
				CorrelationId: item.CorrelationId,
				Error:         err.Error(),
			})
		}
	}

	if len(out.Urls) == 0 {
		return nil, ErrEmptyBatch
	}

	return &out, nil
}

// checkShortUID perform checking shortUID before deleting.
func checkShortUID(shortUID ...string) []string {
	var correctShortUID []string

	for _, v := range shortUID {
		if !strings.Contains(v, "/") && v != "" {
			correctShortUID = append(correctShortUID, v)
		} else if sUID := v[strings.LastIndex(v, "/")+1:]; sUID != "" {
			correctShortUID = append(correctShortUID, sUID)
		}
	}

	return correctShortUID
}

// mustEmbedUnimplementedShortenerServer implements interface UnsafeShortenerServer.
func (s *Server) mustEmbedUnimplementedShortenerServer() {}
