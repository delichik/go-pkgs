package db

import (
	"context"
)

type CommonOperator interface {
}

type CacheEnabledModel interface {
	GetCacheFormat() string
	GetCacheKey() string
}

func CommonGet[M ModelInterface](ctx context.Context, id int64) (*M, error) {
}

func CommonGetBatch[M ModelInterface](ctx context.Context, ids []int64) ([]*M, error) {
}

func CommonUpdate[M ModelInterface](ctx context.Context, item *M) (*ExecRsp, error) {
}

func CommonUpdateBatch[M ModelInterface](ctx context.Context, items []*M) (*ExecRsp, error) {
}

// CommonUpdateWhere 不支持缓存
func CommonUpdateWhere[M ModelInterface](ctx context.Context, item *M, cond []*Where) (*ExecRsp, error) {
}

// CommonUpsert 不支持缓存
func CommonUpsert[M ModelInterface](ctx context.Context, item *M) (*ExecRsp, error) {
}

// CommonUpsertBatch 不支持缓存
func CommonUpsertBatch[M ModelInterface](ctx context.Context, items []*M) (*ExecRsp, error) {
}

func CommonCreate[M ModelInterface](ctx context.Context, item *M) error {
}

func CommonCreateBatch[M ModelInterface](ctx context.Context, items []*M) error {
}

func CommonDelete[M ModelInterface](ctx context.Context, id int64) (*ExecRsp, error) {
}

func CommonDeleteBatch[M ModelInterface](ctx context.Context, ids []int64) (*ExecRsp, error) {
}

// CommonDeleteWhere 不支持缓存
func CommonDeleteWhere[M ModelInterface](ctx context.Context, cond []*Where) (*ExecRsp, error) {
}

func CommonScan[M ModelInterface](ctx context.Context, req *ScanReq) (*ScanRsp, error) {
}
