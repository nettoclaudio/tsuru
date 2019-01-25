package vault

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/api"
	"github.com/tsuru/config"
	"github.com/tsuru/tsuru/app/bind"
)

type vaultKeyValueStorager struct {
	client  *api.Client
	prefix  string
	appName string
}

type vaultEnvVar struct {
	bind.EnvVar
	Deleted bool `json:"deleted"`
}

func (v *vaultKeyValueStorager) Get(envs ...string) (map[string]bind.EnvVar, error) {
	path := fmt.Sprintf("secret/data/%s/%s", v.prefix, v.appName)
	secret, err := v.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	data := secret.Data["data"].(map[string]interface{})

	result := make(map[string]bind.EnvVar)

	if len(envs) > 0 {
		for _, env := range envs {
			var envVar bind.EnvVar
			json.Unmarshal([]byte(data[env].(string)), &envVar)
			result[env] = envVar
		}

		return result, nil
	}

	for key, rawEnvVar := range data {
		var env bind.EnvVar
		json.Unmarshal([]byte(rawEnvVar.(string)), &env)
		result[key] = env
	}

	return result, nil
}

func (v *vaultKeyValueStorager) Set(envs []bind.EnvVar) error {
	data := make(map[string]interface{}, 1)

	result := make(map[string]interface{})

	for _, env := range envs {

		vaultEnvVar := vaultEnvVar{env, Deleted: false}

		rawEnv, _ := json.Marshal(env)
		result[env.Name] = string(rawEnv)
	}

	data["data"] = result

	path := fmt.Sprintf("secret/data/%s/%s", v.prefix, v.appName)
	_, err := v.client.Logical().Write(path, data)
	return err
}

func (v *vaultKeyValueStorager) Unset(envs ...string) error {
	return nil
}

func NewVaultKeyValueStorager(appName string) *vaultKeyValueStorager {
	client, _ := api.NewClient(nil)
	address, _ := config.GetString("vault:address")

	certficatePath, _ := config.GetString("vault:cert-file")
	keyPath, _ := config.GetString("vault:key-file")

	certificate, err := tls.LoadX509KeyPair(certficatePath, keyPath)

	fmt.Printf("%#v\n", err)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	tlsConfig.BuildNameToCertificate()

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	httpClient := &http.Client{Transport: transport}

	namespace, _ := config.GetString("vault:namespace")

	url := fmt.Sprintf("%s/v1/auth/cert/login", address)

	body, _ := json.Marshal(map[string]string{"name": namespace})

	response, _ := httpClient.Post(url, "application/json", bytes.NewBuffer(body))

	defer response.Body.Close()

	type vaultResponse struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}

	var vr vaultResponse

	json.NewDecoder(response.Body).Decode(&vr)

	client.SetAddress(address)
	client.SetToken(vr.Auth.ClientToken)

	return &vaultKeyValueStorager{
		client:  client,
		prefix:  "lab/tsuru",
		appName: appName,
	}
}
