package models

import (
	"fmt"
)

type TagMap map[string]string

type AssetMap map[string]TagMap

func CreateAssetMap(assets []Asset) (assetMap AssetMap) {
	assetMap = make(AssetMap)
	for _, asset := range assets {
		tags := make(TagMap)
		for _, tag := range asset.Tags {
			if tags[tag.Key] == "" {
				tags[tag.Key] = tag.Value
				continue
			}
			tags[tag.Key] = fmt.Sprintf("%s,%s", tags[tag.Key], tag.Value)
			assetMap[asset.Id] = tags
		}
		assetMap[asset.Id] = tags
	}
	return
}

type VulnAssetTag struct {
	Count int
	Tags  []Tag
}

type AssetVuln struct {
	Count    int
	Severity string
	Hostname string
	Labels   []string
}
