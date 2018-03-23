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
	switch ctx.String("provider") {
	case "aws":
		return &AmazonWebServices{Region: ctx.String("region")}
	}
	return &AmazonWebServices{}
}
