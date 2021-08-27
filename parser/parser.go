// Copyright © 2021 Alibaba Group Holding Ltd.

package parser

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/alibaba/sealer/version"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/alibaba/sealer/logger"
	v1 "github.com/alibaba/sealer/types/api/v1"
)

var validLayer = []string{"FROM", "COPY", "RUN", "CMD"}

type Interface interface {
	Parse(kubeFile []byte) *v1.Image
}

type Parser struct{}

func NewParse() Interface {
	return &Parser{}
}

func (p *Parser) Parse(kubeFile []byte) *v1.Image {
	image := &v1.Image{
		TypeMeta: metaV1.TypeMeta{APIVersion: "", Kind: "Image"},
		Spec:     v1.ImageSpec{SealerVersion: version.Get().GitVersion},
		Status:   v1.ImageStatus{},
	}
	scanner := bufio.NewScanner(strings.NewReader(string(kubeFile)))
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.Trim(text, " \t\n")
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		layerType, layerValue, err := decodeLine(text)
		if err != nil {
			logger.Warn("decode kubeFile line failed, err: %v", err)
			return nil
		}
		if layerType == "" {
			continue
		}

		//TODO count layer hash
		image.Spec.Layers = append(image.Spec.Layers, v1.Layer{
			ID:    "",
			Type:  layerType,
			Value: layerValue,
		})
	}
	return image
}

func decodeLine(line string) (string, string, error) {
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", nil
	}
	//line = strings.TrimPrefix(line, " ")
	ss := strings.SplitN(line, " ", 2)
	if len(ss) != 2 {
		return "", "", fmt.Errorf("unknown line %s", line)
	}
	var flag bool
	for _, v := range validLayer {
		if ss[0] == v {
			flag = true
		}
	}
	if !flag {
		return "", "", fmt.Errorf("invalid command %s %s", ss[0], line)
	}

	return ss[0], strings.TrimSpace(ss[1]), nil
}
