package middleware

import (
	"context"
	"strings"

	"github.com/studyguides-com/study-guides-api/internal/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorUnaryInterceptor converts application errors to friendly gRPC status errors
func ErrorUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		
		if err != nil {
			// Convert specific application errors to friendly gRPC status errors
			switch err {
			case errors.ErrToolNotFound:
				return nil, status.Error(codes.InvalidArgument, "I couldn't understand how to help with that request. Could you please rephrase or ask about something else?")
			case errors.ErrNotFound:
				return nil, status.Error(codes.NotFound, "The requested resource was not found.")
			case errors.ErrSystemPromptEmpty:
				return nil, status.Error(codes.Internal, "System configuration error. Please try again later.")
			case errors.ErrUserPromptEmpty:
				return nil, status.Error(codes.InvalidArgument, "Please provide a message or question.")
			case errors.ErrNoCompletionChoicesReturned:
				return nil, status.Error(codes.Internal, "I'm having trouble processing your request right now. Please try again.")
			case errors.ErrFailedToCreateChatCompletionWithTools:
				return nil, status.Error(codes.Internal, "I'm experiencing technical difficulties. Please try again later.")
			default:
				// Check for specific error patterns in the error message
				errMsg := err.Error()
				if strings.Contains(errMsg, "AI did not call any tools") {
					return nil, status.Error(codes.InvalidArgument, "I couldn't understand how to help with that request. Try asking me to 'list tags', 'count tags', or 'show root tags'.")
				}
				if strings.Contains(errMsg, "AI returned no choices") {
					return nil, status.Error(codes.Internal, "I'm having trouble processing your request right now. Please try again.")
				}
				if strings.Contains(errMsg, "failed to parse AI response") {
					return nil, status.Error(codes.Internal, "I'm experiencing technical difficulties. Please try again later.")
				}
				// For any other errors, return them as-is
				return nil, err
			}
		}
		
		return resp, nil
	}
} 