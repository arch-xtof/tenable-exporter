package client

import (
	"encoding/json"
	"log"
	"time"

	"github.com/arch-xtof/tenable-exporter/models"
	"github.com/arch-xtof/tenable-exporter/utils"
)

func (client *TenableClient) exportAndDownload(
	optUuid string,
	exporter func() (string, error),
	verifier func(string) (models.ExportStatus, error),
	downloader func(uuid string, chunkId int) (chunk models.Chunk, err error),
) (chunks []models.Chunk, err error) {
	var uuid string

	if optUuid != "" {
		uuid = optUuid
	} else {
		uuid, err = exporter()
		if err != nil {
			return
		}
	}

	log.Printf("%s: Starting export...\n", uuid)
	var status models.ExportStatus
	for i := 0; i <= 120; i++ {
		status, err = verifier(uuid)
		if err != nil {
			return
		}
		log.Printf("%s: %s (%d/%d)\n", uuid, status.Status, status.FinishedChunks, status.TotalChunks)

		if status.Status == "FINISHED" {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	chunksAvailable := status.ChunksAvailable

	chunks, err = client.downloadChunks(uuid, chunksAvailable, downloader)
	if err != nil {
		return
	}
	return
}

func (client *TenableClient) downloadChunks(
	uuid string,
	chunkIds []int,
	downloader func(uuid string, chunkId int) (chunk models.Chunk, err error),
) (chunks []models.Chunk, err error) {
	for _, chunkId := range chunkIds {
		var chunk []byte
		chunk, err = downloader(uuid, chunkId)
		if err != nil {
			return
		}
		chunks = append(chunks, chunk)
	}
	return
}

func (client *TenableClient) GetAssets(optUuid string) (assets []models.Asset, err error) {
	exporter := client.exportAssets
	verifier := client.getAssetsExportStatus
	downloader := client.downloadAssetsChunk

	chunks, err := client.exportAndDownload(optUuid, exporter, verifier, downloader)
	if err != nil {
		return
	}

	for _, chunk := range chunks {
		var assetsBuffer []models.Asset
		err = json.Unmarshal(chunk, &assetsBuffer)
		if err != nil {
			return
		}
		assets = append(assets, assetsBuffer...)
	}

	return
}

func (client *TenableClient) GetVulns(optUuid string) (vulns []models.Vuln, err error) {
	exporter := client.exportVulnerabilities
	verifier := client.getVunlerabilitiesExportStatus
	downloader := client.downloadVulnerabilitiesChunk

	chunks, err := client.exportAndDownload(optUuid, exporter, verifier, downloader)
	if err != nil {
		return
	}

	for _, chunk := range chunks {
		var buffer []models.Vuln
		err = json.Unmarshal(chunk, &buffer)
		if err != nil {
			return
		}
		vulns = append(vulns, buffer...)
	}
	return
}

func MockVulns() (vulns []models.Vuln, err error) {
	bytes := utils.OpenJson("mock/_vulns.json")
	err = json.Unmarshal(bytes, &vulns)
	return
}

func MockAssets() (assets []models.Asset, err error) {
	bytes := utils.OpenJson("mock/_assets.json")
	err = json.Unmarshal(bytes, &assets)
	return
}
