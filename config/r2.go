package config

import "os"

type R2 struct {
	PublicBucketURL string
	PublicBucket    string
	PrivateBucket   string
	AccountID       string
	AccessKeyID     string
	AccessKeySecret string
}

func (r *R2) Parse() error {
	r.PublicBucketURL = os.Getenv("R2_PUBLIC_BUCKET_URL")
	r.PublicBucket = os.Getenv("R2_PUBLIC_BUCKET")
	r.PrivateBucket = os.Getenv("R2_PRIVATE_BUCKET")
	r.AccountID = os.Getenv("R2_ACCOUNT_ID")
	r.AccessKeyID = os.Getenv("R2_ACCESS_KEY_ID")
	r.AccessKeySecret = os.Getenv("R2_ACCESS_KEY_SECRET")
	return nil
}
