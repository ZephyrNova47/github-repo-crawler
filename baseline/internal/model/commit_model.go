package model

type CommitResponse struct {
	ID      int64  `json:"id"`
	Hash    string `json:"hash"`
	Message string `json:"message"`
}

type CreateCommitRequest struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
}
