package platform

// Platform ...
type Platform interface {
	Validate() error
}

// Context ...
type Context interface {
	String(string) string
	Int(string) int
}

// Get ...
func Get(ctx Context) Platform {
	switch Provider(ctx.String("provider")) {
	case AWS:
		return &AmazonWebServices{Region: ctx.String("aws-region")}
	}
	return &AmazonWebServices{}
}
