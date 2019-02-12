package outlinesdk

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	client, err := NewClient(os.Getenv("API_URL"), os.Getenv("CERT_SHA_256"))
	if err != nil {
		t.Error(err)
	}

	t.Log("set metrics setting to true")
	err = client.SetMetricsSetting(true)
	if err != nil {
		t.Error(err)
	}
	metricSetting, err := client.GetMetricsSetting()
	if err != nil {
		t.Error(err)
	}
	t.Logf("GetMetricsSetting: %v", *metricSetting)

	t.Log("set metrics setting to false")
	err = client.SetMetricsSetting(false)
	if err != nil {
		t.Error(err)
	}
	metricSetting, err = client.GetMetricsSetting()
	if err != nil {
		t.Error(err)
	}
	t.Logf("GetMetricsSetting: %v", *metricSetting)

	serverInfo, err := client.GetServerInfo()
	if err != nil {
		t.Error(err)
	}
	originalName := serverInfo.Name
	t.Logf("Original server name is %s", serverInfo.Name)
	randomName := randomString()
	t.Logf("set server name to %s", randomName)
	err = client.RenameServer(randomName)
	if err != nil {
		t.Error(err)
	}
	serverInfo, err = client.GetServerInfo()
	if err != nil {
		t.Error(err)
	}
	if randomName != serverInfo.Name {
		t.Error("server name not changed")
	}
	err = client.RenameServer(originalName)
	if err != nil {
		t.Error(err)
	}
	serverInfo, err = client.GetServerInfo()
	if err != nil {
		t.Error(err)
	}
	t.Logf("GetServerInfo: %v", *serverInfo)
	t.Logf("server created in %v", serverInfo.CreatedTime())

	t.Logf("creates a new access key")
	accessKey, err := client.CreateAccessKey()
	if err != nil {
		t.Error(err)
	}
	keyID := accessKey.ID
	t.Logf("keyId: %s", keyID)
	list, err := client.GetAccessKeys()
	if err != nil {
		t.Error(err)
	}
	t.Logf("GetAccessKeys: %v", *list)
	if getKeyWithID(list, keyID) == nil {
		t.Error("such key does not exists")
	}
	originalName = getKeyWithID(list, keyID).Name
	t.Logf("original key name is %s", originalName)
	randomName = randomString()
	err = client.RenameAccessKey(keyID, randomName)
	if err != nil {
		t.Error(err)
	}
	changedList, err := client.GetAccessKeys()
	if err != nil {
		t.Error(err)
	}
	if getKeyWithID(changedList, keyID).Name != randomName {
		t.Error("key does not match new name")
	}
	err = client.DeleteAccessKey(keyID)
	if err != nil {
		t.Error(err)
	}
	newList, err := client.GetAccessKeys()
	if err != nil {
		t.Error(err)
	}
	if getKeyWithID(newList, keyID) != nil {
		t.Error("key is not deleted")
	}
	t.Logf("Latest GetAccessKeys: %v", *newList)
	usageInfo, err := client.GetUsageMetrics()
	if err != nil {
		t.Error(err)
	}
	t.Logf("GetUsageMetrics: %v", *usageInfo)
}

func getKeyWithID(l *AccessKeyList, id string) *AccessKey {
	for _, k := range *l {
		if k.ID == id {
			return &k
		}
	}
	return nil
}

func randomString() string {
	return fmt.Sprintf("%d", rand.Int())
}
