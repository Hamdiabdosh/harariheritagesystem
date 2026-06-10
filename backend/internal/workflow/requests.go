package workflow

import "strings"

type optionalCommentRequest struct {
	CommentText *string `json:"comment_text"`
	Comment     *string `json:"comment"`
}

func (r optionalCommentRequest) value() *string {
	if r.CommentText != nil {
		trimmed := strings.TrimSpace(*r.CommentText)
		if trimmed != "" {
			return &trimmed
		}
	}
	if r.Comment != nil {
		trimmed := strings.TrimSpace(*r.Comment)
		if trimmed != "" {
			return &trimmed
		}
	}
	return nil
}

type requiredCommentRequest struct {
	CommentText string `json:"comment_text"`
	Comment     string `json:"comment"`
}

func (r requiredCommentRequest) value() (string, bool) {
	if trimmed := strings.TrimSpace(r.CommentText); trimmed != "" {
		return trimmed, true
	}
	if trimmed := strings.TrimSpace(r.Comment); trimmed != "" {
		return trimmed, true
	}
	return "", false
}
