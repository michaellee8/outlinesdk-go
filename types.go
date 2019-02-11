package outlinesdk

import "time"

type NameType struct {
	Name string `json:"name"`
}

type MetricsSetting struct {
	MetricsEnabled bool `json:"metricsEnabled"`
}

type AccessKey struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Port      int    `json:"port"`
	Method    string `json:"method"`
	AccessUrl string `json:"accessUrl"`
}

type AccessKeyList map[string]AccessKey

type UsageInfo struct {
	BytesTransferredByUserId map[string]int64 `json:"bytesTransferredByUserId"`
}

type ServerInfo struct {
	Name                 string
	ServerId             string
	MetricsEnabled       string
	CreatedTimestampMs   int64
	PortForNewAccessKeys int
}

func (i *ServerInfo) CreatedTimestamp() time.Time {
	return time.Unix(0, i.CreatedTimestampMs*1000000)
}
