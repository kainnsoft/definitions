package client

const (
	SourceSecured  = "Secured"  // Request from API via OAuth
	SourcePrivate  = "Private"  // Request from API via internal nginx
	SourcePublic   = "Public"   // Request from API without authorization
	SourceInternal = "Internal" // Internal request between services via RPC
)
