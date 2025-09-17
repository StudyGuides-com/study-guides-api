package indexing

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	indexingv1 "github.com/studyguides-com/study-guides-api/api/v1/indexing"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	indexingcore "github.com/studyguides-com/study-guides-api/internal/core/indexing"
	"github.com/studyguides-com/study-guides-api/internal/middleware"
	"github.com/studyguides-com/study-guides-api/internal/services"
	"github.com/studyguides-com/study-guides-api/internal/store"
	"github.com/studyguides-com/study-guides-api/internal/store/indexing"
)

// IndexingService provides gRPC interface for indexing operations
// This runs in parallel with the MCP natural language interface
type IndexingService struct {
	indexingv1.UnimplementedIndexingServiceServer
	business *indexingcore.BusinessService
}

// NewIndexingService creates a new IndexingService
func NewIndexingService(store store.Store) *IndexingService {
	return &IndexingService{
		business: indexingcore.NewBusinessService(store),
	}
}

// TriggerIndexing starts a new indexing job
func (s *IndexingService) TriggerIndexing(ctx context.Context, req *indexingv1.TriggerIndexingRequest) (*indexingv1.TriggerIndexingResponse, error) {
	resp, err := services.AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		// Check admin permissions
		if !session.HasRole(sharedpb.UserRole_USER_ROLE_ADMIN) {
			return nil, status.Error(codes.PermissionDenied, "admin access required for indexing operations")
		}

		// Use business service to trigger indexing
		businessReq := indexingcore.TriggerIndexingRequest{
			ObjectType: req.ObjectType,
			Force:      req.Force,
		}
		businessResp, err := s.business.TriggerIndexing(ctx, businessReq)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to trigger indexing: %v", err))
		}

		return &indexingv1.TriggerIndexingResponse{
			JobId:     businessResp.JobID,
			Status:    businessResp.Status,
			Message:   businessResp.Message,
			StartedAt: timestamppb.New(businessResp.StartedAt),
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*indexingv1.TriggerIndexingResponse), nil
}

// GetJobStatus returns the status of a specific indexing job
func (s *IndexingService) GetJobStatus(ctx context.Context, req *indexingv1.GetJobStatusRequest) (*indexingv1.GetJobStatusResponse, error) {
	resp, err := services.AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		// Check admin permissions
		if !session.HasRole(sharedpb.UserRole_USER_ROLE_ADMIN) {
			return nil, status.Error(codes.PermissionDenied, "admin access required for indexing operations")
		}

		// Get job status from business service
		job, err := s.business.GetJobStatus(ctx, req.JobId)
		if err != nil {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("job not found: %v", err))
		}

		// Convert to response format
		response := &indexingv1.GetJobStatusResponse{
			JobId:       job.ID,
			Status:      job.Status,
			Description: job.Description,
		}

		// Set timestamps
		if job.StartedAt != nil {
			response.StartedAt = timestamppb.New(*job.StartedAt)
		}
		if job.CompletedAt != nil {
			response.CompletedAt = timestamppb.New(*job.CompletedAt)
		}
		if job.ErrorMessage != nil {
			response.ErrorMessage = *job.ErrorMessage
		}

		// Convert metadata
		if job.Metadata != nil {
			metadata := &indexingv1.JobMetadata{
				Extra: make(map[string]string),
			}

			// Extract known metadata fields
			if objectType, ok := job.Metadata["objectType"].(string); ok {
				metadata.ObjectType = objectType
			}
			if force, ok := job.Metadata["force"].(bool); ok {
				metadata.Force = force
			}
			if itemsProcessed, ok := job.Metadata["itemsProcessed"].(float64); ok {
				metadata.ItemsProcessed = int64(itemsProcessed)
			}
			if itemsFailed, ok := job.Metadata["itemsFailed"].(float64); ok {
				metadata.ItemsFailed = int64(itemsFailed)
			}

			// Add any extra metadata as strings
			for key, value := range job.Metadata {
				if key != "objectType" && key != "force" && key != "itemsProcessed" && key != "itemsFailed" {
					metadata.Extra[key] = fmt.Sprintf("%v", value)
				}
			}

			response.Metadata = metadata
		}

		return response, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*indexingv1.GetJobStatusResponse), nil
}

// ListRunningJobs returns all currently running indexing jobs
func (s *IndexingService) ListRunningJobs(ctx context.Context, req *indexingv1.ListRunningJobsRequest) (*indexingv1.ListRunningJobsResponse, error) {
	resp, err := services.AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		// Check admin permissions
		if !session.HasRole(sharedpb.UserRole_USER_ROLE_ADMIN) {
			return nil, status.Error(codes.PermissionDenied, "admin access required for indexing operations")
		}

		// Get running jobs from business service
		jobs, err := s.business.ListRunningJobs(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list running jobs: %v", err))
		}

		// Convert to response format
		jobInfos := make([]*indexingv1.JobInfo, 0, len(jobs))
		for _, job := range jobs {
			jobInfo := convertJobStatusToJobInfo(job)
			jobInfos = append(jobInfos, jobInfo)
		}

		return &indexingv1.ListRunningJobsResponse{
			Jobs:       jobInfos,
			TotalCount: int32(len(jobInfos)),
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*indexingv1.ListRunningJobsResponse), nil
}

// ListRecentJobs returns recent indexing jobs, optionally filtered by object type
func (s *IndexingService) ListRecentJobs(ctx context.Context, req *indexingv1.ListRecentJobsRequest) (*indexingv1.ListRecentJobsResponse, error) {
	resp, err := services.AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		// Check admin permissions
		if !session.HasRole(sharedpb.UserRole_USER_ROLE_ADMIN) {
			return nil, status.Error(codes.PermissionDenied, "admin access required for indexing operations")
		}

		// Get recent jobs from business service
		businessReq := indexingcore.ListRecentJobsRequest{
			ObjectType: req.ObjectType,
			Limit:      int(req.Limit),
		}
		jobs, err := s.business.ListRecentJobs(ctx, businessReq)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list recent jobs: %v", err))
		}

		// Convert to response format
		jobInfos := make([]*indexingv1.JobInfo, 0, len(jobs))
		for _, job := range jobs {
			jobInfo := convertJobStatusToJobInfo(job)
			jobInfos = append(jobInfos, jobInfo)
		}

		return &indexingv1.ListRecentJobsResponse{
			Jobs: jobInfos,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*indexingv1.ListRecentJobsResponse), nil
}

// TriggerTagIndexing starts a new tag indexing job with filtering
func (s *IndexingService) TriggerTagIndexing(ctx context.Context, req *indexingv1.TriggerTagIndexingRequest) (*indexingv1.TriggerIndexingResponse, error) {
	resp, err := services.AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		// Check admin permissions
		if !session.HasRole(sharedpb.UserRole_USER_ROLE_ADMIN) {
			return nil, status.Error(codes.PermissionDenied, "admin access required for indexing operations")
		}

		// Convert protobuf types to business types
		businessReq := indexingcore.TriggerTagIndexingRequest{
			Force:        req.Force,
			TagTypes:     req.TagTypes,
			ContextTypes: req.ContextTypes,
		}
		businessResp, err := s.business.TriggerTagIndexing(ctx, businessReq)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to trigger tag indexing: %v", err))
		}

		return &indexingv1.TriggerIndexingResponse{
			JobId:     businessResp.JobID,
			Status:    businessResp.Status,
			Message:   businessResp.Message,
			StartedAt: timestamppb.New(businessResp.StartedAt),
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*indexingv1.TriggerIndexingResponse), nil
}

// TriggerSingleIndexing starts a new indexing job for a single specific item
func (s *IndexingService) TriggerSingleIndexing(ctx context.Context, req *indexingv1.TriggerSingleIndexingRequest) (*indexingv1.TriggerIndexingResponse, error) {
	resp, err := services.AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		// Check admin permissions
		if !session.HasRole(sharedpb.UserRole_USER_ROLE_ADMIN) {
			return nil, status.Error(codes.PermissionDenied, "admin access required for indexing operations")
		}

		// Convert protobuf types to business types
		businessReq := indexingcore.TriggerSingleIndexingRequest{
			ObjectType: req.ObjectType,
			ID:         req.Id,
			Force:      req.Force,
		}
		businessResp, err := s.business.TriggerSingleIndexing(ctx, businessReq)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to trigger single indexing: %v", err))
		}

		return &indexingv1.TriggerIndexingResponse{
			JobId:     businessResp.JobID,
			Status:    businessResp.Status,
			Message:   businessResp.Message,
			StartedAt: timestamppb.New(businessResp.StartedAt),
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*indexingv1.TriggerIndexingResponse), nil
}

// convertJobStatusToJobInfo converts internal JobStatus to protobuf JobInfo
func convertJobStatusToJobInfo(job indexing.JobStatus) *indexingv1.JobInfo {
	jobInfo := &indexingv1.JobInfo{
		JobId:       job.ID,
		Status:      job.Status,
		Description: job.Description,
	}

	// Set timestamps
	if job.StartedAt != nil {
		jobInfo.StartedAt = timestamppb.New(*job.StartedAt)
	}
	if job.CompletedAt != nil {
		jobInfo.CompletedAt = timestamppb.New(*job.CompletedAt)
	}
	if job.ErrorMessage != nil {
		jobInfo.ErrorMessage = *job.ErrorMessage
	}

	// Convert metadata
	if job.Metadata != nil {
		metadata := &indexingv1.JobMetadata{
			Extra: make(map[string]string),
		}

		// Extract known metadata fields
		if objectType, ok := job.Metadata["objectType"].(string); ok {
			metadata.ObjectType = objectType
		}
		if force, ok := job.Metadata["force"].(bool); ok {
			metadata.Force = force
		}
		if itemsProcessed, ok := job.Metadata["itemsProcessed"].(float64); ok {
			metadata.ItemsProcessed = int64(itemsProcessed)
		}
		if itemsFailed, ok := job.Metadata["itemsFailed"].(float64); ok {
			metadata.ItemsFailed = int64(itemsFailed)
		}

		// Add any extra metadata as strings
		for key, value := range job.Metadata {
			if key != "objectType" && key != "force" && key != "itemsProcessed" && key != "itemsFailed" {
				metadata.Extra[key] = fmt.Sprintf("%v", value)
			}
		}

		jobInfo.Metadata = metadata
	}

	return jobInfo
}