package srv

const (
	CtxKeyUserID      = CtxKey("user_id")
	CtxKeySoiBucketID = CtxKey("soi_bucket_id")
)

type CtxKey string

func (ck CtxKey) String() string {
	return string(ck)
}
