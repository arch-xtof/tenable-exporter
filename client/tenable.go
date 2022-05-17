package client

import (
	"net/http"

	"github.com/arch-xtof/tenable-exporter/models"
)

type Auth struct {
	AccessKey string
	SecretKey string
}

type TenableClient struct {
	httpClient *http.Client
	auth       Auth
}

func NewTenableClient(httpClient *http.Client, auth Auth) (client *TenableClient, err error) {
	client = &TenableClient{
		httpClient: httpClient,
		auth:       auth,
	}
	return
}

func (client *TenableClient) GetServerVulnerabilityCount(vulns []models.Vuln, assets models.AssetMap) (svc models.ServerVulnerabilityCount) {
	svc = make(models.ServerVulnerabilityCount)
	for _, vuln := range vulns {
		if vuln.SeverityModificationType == "ACCEPTED" {
			continue
		}
		tags := assets[vuln.Asset.Uuid]
		if tags["Type"] != "Servers" {
			continue
		}

		if svc[vuln.Asset.Hostname] == nil {
			svc[vuln.Asset.Hostname] = make(map[string]*struct {
				Region string
				Group  string
				Count  int
			})
		}
		if svc[vuln.Asset.Hostname][vuln.Severity] == nil {
			svc[vuln.Asset.Hostname][vuln.Severity] = &struct {
				Region string
				Group  string
				Count  int
			}{
				Region: tags["Region"],
				Group:  tags["Servers"],
				Count:  1,
			}
			continue
		}
		svc[vuln.Asset.Hostname][vuln.Severity].Count++
	}
	return
}

func (client *TenableClient) GetWorkstationVulnerabilityCount(vulns []models.Vuln, assets models.AssetMap) (wvc models.WorkstationVulnerabilityCount) {
	wvc = make(models.WorkstationVulnerabilityCount)
	for _, vuln := range vulns {
		if vuln.SeverityModificationType == "ACCEPTED" {
			continue
		}
		tags := assets[vuln.Asset.Uuid]
		if tags["Type"] != "Workstations" {
			continue
		}

		if wvc[vuln.Asset.Hostname] == nil {
			wvc[vuln.Asset.Hostname] = make(map[string]*struct {
				OS    string
				App   string
				Count int
			})
		}

		if wvc[vuln.Asset.Hostname][vuln.Severity] == nil {
			wvc[vuln.Asset.Hostname][vuln.Severity] = &struct {
				OS    string
				App   string
				Count int
			}{
				OS:    tags["Workstations"],
				App:   cpeGuess(vuln.Output, vuln.Plugin.Cpe),
				Count: 1,
			}
			continue
		}
		wvc[vuln.Asset.Hostname][vuln.Severity].Count++
	}
	return
}
