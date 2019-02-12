package outlinesdk

import "time"

// NameType is a struct with only a name field
type NameType struct {
	Name string `json:"name"`
}

// MetricsSetting is a struct that represents the server's metrics settings
type MetricsSetting struct {
	MetricsEnabled bool `json:"metricsEnabled"`
}

// AccessKey is a struct that represents an access key provided by the server
type AccessKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Port      int    `json:"port"`
	Method    string `json:"method"`
	AccessURL string `json:"accessUrl"`
}

// AccessKeyList is a map with access key's id as its key
type AccessKeyList map[string]AccessKey

// UsageInfo is a map with access key's id as its key and bytes transferred per user as its value
type UsageInfo struct {
	BytesTransferredByUserID map[string]int64 `json:"bytesTransferredByUserId"`
}

// ServerInfo is a struct that provides server's information
type ServerInfo struct {
	Name                 string
	ServerID             string
	MetricsEnabled       string
	CreatedTimestampMs   int64
	PortForNewAccessKeys int
}

// CreatedTime parses CreatedTimestampMs and returns a time.Time
func (i *ServerInfo) CreatedTime() time.Time {
	return time.Unix(0, i.CreatedTimestampMs*1000000)
}
