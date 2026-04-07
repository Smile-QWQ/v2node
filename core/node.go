package core

import (
	"fmt"

	panel "github.com/wyx2685/v2node/api/v2board"
)

func (v *V2Core) AddNode(tag string, info *panel.NodeInfo) error {
	inBoundConfigs, err := buildInbounds(info, tag)
	if err != nil {
		return fmt.Errorf("build inbound error: %s", err)
	}
	for _, inBoundConfig := range inBoundConfigs {
		err = v.addInbound(inBoundConfig)
		if err != nil {
			for _, addedInbound := range inBoundConfigs {
				if addedInbound == nil || addedInbound.Tag == inBoundConfig.Tag {
					break
				}
				_ = v.removeInbound(addedInbound.Tag)
			}
			return fmt.Errorf("add inbound error: %s", err)
		}
	}
	return nil
}

func (v *V2Core) DelNode(tag string, info *panel.NodeInfo) error {
	err := v.removeInbound(tag)
	if err != nil {
		return fmt.Errorf("remove in error: %s", err)
	}
	if shouldEnableAntiStealReality(info) {
		antiStealDokodemoTag, ok := getAntiStealDokodemoTag(tag)
		if !ok {
			return fmt.Errorf("anti-steal dokodemo tag is invalid")
		}
		err = v.removeInbound(antiStealDokodemoTag)
		if err != nil {
			return fmt.Errorf("remove anti steal inbound error: %s", err)
		}
	}
	return nil
}
