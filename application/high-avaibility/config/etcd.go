package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Struktur data member Patroni
type PatroniMember struct {
	ConnURL      string `json:"conn_url"`
	APIURL       string `json:"api_url"`
	Role         string `json:"role"`
	State        string `json:"state"`
	Version      string `json:"version"`
	XlogLocation int    `json:"xlog_location"`
	Timeline     int    `json:"timeline"`
}

func GetConnectionPrimary() (string, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.1.15:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		return "", err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Ambil Nama Leader
	resp, err := cli.Get(ctx, "/service/postgres-ha/leader")
	if err != nil || len(resp.Kvs) == 0 {
		return "", fmt.Errorf("leader not found")
	}
	leaderName := string(resp.Kvs[0].Value)
	fmt.Printf("📡 Etcd Leader Name: %s\n", leaderName)

	// 2. Ambil Detail Member (IP)
	memberKey := fmt.Sprintf("/service/postgres-ha/members/%s", leaderName)
	mResp, err := cli.Get(ctx, memberKey)
	if err != nil || len(mResp.Kvs) == 0 {
		return "", fmt.Errorf("member detail not found")
	}

	// 3. Decode JSON Patroni
	var member PatroniMember
	if err := json.Unmarshal(mResp.Kvs[0].Value, &member); err != nil {
		return "", err
	}

	u, err := url.Parse(member.ConnURL)
	if err != nil {
		return "", err
	}
	host := u.Hostname()

	user := "postgres"
	pass := "super_password"
	dbName := "postgres"

	finalDSN := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		host, user, pass, dbName)

	return finalDSN, nil
}
