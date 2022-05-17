package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arch-xtof/tenable-exporter/models"
)

func (client *TenableClient) exportVulnerabilities() (uuid string, err error) {
	method := http.MethodPost
	path := "vulns/export"
	params := "{\"num_assets\":200,\"filters\":{\"severity\":[\"low\",\"medium\",\"high\",\"critical\"],\"state\":[\"open\",\"reopened\"]}}"

	body, err := client.Query(method, path, params)
	if err != nil {
		return
	}

	var bodyMap map[string]string
	_ = json.Unmarshal(body, &bodyMap)
	uuid = bodyMap["export_uuid"]
	return
}

func (client *TenableClient) getVunlerabilitiesExportStatus(uuid string) (status models.ExportStatus, err error) {
	method := http.MethodGet
	path := fmt.Sprintf("vulns/export/%s/status", uuid)
	params := ""

	body, err := client.Query(method, path, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &status)
	return
}

func (client *TenableClient) downloadVulnerabilitiesChunk(uuid string, chunkId int) (chunk models.Chunk, err error) {
	method := http.MethodGet
	path := fmt.Sprintf("vulns/export/%s/chunks/%d", uuid, chunkId)
	params := ""

	chunk, err = client.Query(method, path, params)
	return
}

func (client *TenableClient) exportAssets() (uuid string, err error) {
	method := http.MethodPost
	path := "assets/export"
	params := "{\"chunk_size\":200}"

	body, err := client.Query(method, path, params)
	if err != nil {
		return
	}

	var bodyMap map[string]string
	_ = json.Unmarshal(body, &bodyMap)
	uuid = bodyMap["export_uuid"]
	return
}

func (client *TenableClient) getAssetsExportStatus(uuid string) (status models.ExportStatus, err error) {
	method := http.MethodGet
	path := fmt.Sprintf("assets/export/%s/status", uuid)
	params := ""

	body, err := client.Query(method, path, params)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &status)
	return
}

func (client *TenableClient) downloadAssetsChunk(uuid string, chunkId int) (chunk models.Chunk, err error) {
	method := http.MethodGet
	path := fmt.Sprintf("assets/export/%s/chunks/%d", uuid, chunkId)
	params := ""

	chunk, err = client.Query(method, path, params)
	return
}
