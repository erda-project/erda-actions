package manualReview

import (
	"time"

	"github.com/erda-project/erda/apistructs"
)

type ErrorResponse struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Ctx  interface{} `json:"ctx"`
}

type Header struct {
	Success bool          `json:"success" `
	Error   ErrorResponse `json:"err"`
}

type ManualReview struct {
	Header
	Data *ManualReviewData `json:"data"`
}
type ManualReviewData struct {
	Total          int64  `json:"total"`
	Approvalstatus string `json:"Approvalstatus"`
	Id             uint64 `json:"id"`
}
type CreateReviewRequest struct {
	BuildId         uint64    `json:"buildId"`
	ProjectId       uint64    `json:"projectId"`
	ApplicationId   uint64    `json:"applicationId"`
	ApplicationName string    `json:"applicationName"`
	SponsorId       string    `json:"sponsorId"`
	CommitID        string    `json:"commitID"`
	OrgId           uint64    `json:"orgId"`
	TaskId          uint64    `json:"taskId"`
	ProjectName     string    `json:"projectName"`
	BranchName      string    `json:"branchName"`
	ApprovalStatus  string    `json:"approvalStatus"`
	CreatedAt       time.Time `json:"createdAt"`
}
type CreateReviewResponse struct {
	Header
	Data int64 `json:"data"`
}
type CreateReviewUserRequest struct {
	Operator string `json:"operator"`
	OrgId    uint64 `json:"orgId"`
	TaskId   uint64 `json:"taskId"`
}

type CreateReviewUserResponse struct {
	Header
	Data *apistructs.CreateReviewUserResponse `json:"data"`
}
