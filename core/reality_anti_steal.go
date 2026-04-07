package core

import (
	"fmt"
	stdnet "net"
	"strconv"
	"strings"

	panel "github.com/wyx2685/v2node/api/v2board"
)

const (
	antiStealLoopbackHost = "127.0.0.1"
	antiStealPortStart    = 20000
	antiStealPortEnd      = 29999
)

func shouldEnableAntiStealReality(info *panel.NodeInfo) bool {
	if info == nil || info.Common == nil {
		return false
	}
	return info.Type == "vless" &&
		info.Security == panel.Reality &&
		info.Common.Network == "tcp" &&
		info.Common.TlsSettings.AntiStealRealityEnabled
}

func getAntiStealDokodemoTag(tag string) (string, bool) {
	if strings.TrimSpace(tag) == "" {
		return "", false
	}
	return tag + "#anti-steal-dokodemo", true
}

func antiStealRealityTargetHost(info *panel.NodeInfo) string {
	if info == nil || info.Common == nil {
		return ""
	}
	return strings.TrimSpace(info.Common.TlsSettings.ServerName)
}

func antiStealRealityTargetPort(info *panel.NodeInfo) int {
	if info == nil || info.Common == nil {
		return 443
	}
	port, err := strconv.Atoi(strings.TrimSpace(info.Common.TlsSettings.ServerPort))
	if err != nil || port <= 0 || port > 65535 {
		return 443
	}
	return port
}

func allocateAntiStealLoopbackPort(nodeID int) (int, error) {
	total := antiStealPortEnd - antiStealPortStart + 1
	start := antiStealPortStart + (nodeID % total)
	for i := 0; i < total; i++ {
		port := antiStealPortStart + ((start - antiStealPortStart + i) % total)
		listener, err := stdnet.Listen("tcp", fmt.Sprintf("%s:%d", antiStealLoopbackHost, port))
		if err != nil {
			continue
		}
		_ = listener.Close()
		return port, nil
	}
	return 0, fmt.Errorf("no available anti-steal loopback port in range %d-%d", antiStealPortStart, antiStealPortEnd)
}
