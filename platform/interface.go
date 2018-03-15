package platform

// Platform ...
type Platform interface {
	Validate() error
}

// Context ...
type Context interface {
	String(string) string
}

// Get ...
func Get(ctx Context) Platform {
	switch ctx.String("provider") {
	case "aws":
		return &AmazonWebServices{Region: ctx.String("region")}
	}
	return &Virtualbox{}
}
