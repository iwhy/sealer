// Copyright © 2021 Alibaba Group Holding Ltd.

package distributionutil

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	dockerRegistryClient "github.com/docker/distribution/registry/client"
	dockerAuth "github.com/docker/distribution/registry/client/auth"
	dockerTransport "github.com/docker/distribution/registry/client/transport"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/dockerversion"
	dockerRegistry "github.com/docker/docker/registry"
	"github.com/docker/go-connections/tlsconfig"
)

func NewRepository(ctx context.Context, authConfig types.AuthConfig, repoName string, config registryConfig, actions ...string) (distribution.Repository, error) {
	tlsConfig := tlsconfig.ServerDefault()
	tlsConfig.InsecureSkipVerify = config.Insecure

	rurlStr := strings.TrimSuffix(config.Domain, "/")
	if !strings.HasPrefix(rurlStr, "https://") && !strings.HasPrefix(rurlStr, "http://") {
		if !config.NonSSL {
			rurlStr = "https://" + rurlStr
		} else {
			rurlStr = "http://" + rurlStr
		}
	}

	rurl, err := url.Parse(rurlStr)
	if err != nil {
		return nil, err
	}

	direct := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	// TODO(dmcgowan): Call close idle connections when complete, use keep alive
	base := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         direct.DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     tlsConfig,
		// TODO(dmcgowan): Call close idle connections when complete and use keep alive
		DisableKeepAlives: true,
	}

	modifiers := dockerRegistry.Headers(dockerversion.DockerUserAgent(ctx), nil)
	authTransport := dockerTransport.NewTransport(base, modifiers...)

	challengeManager, _, err := dockerRegistry.PingV2Registry(rurl, authTransport)
	if err != nil {
		return nil, err
	}
	// typically this filed would be empty
	if authConfig.RegistryToken != "" {
		passThruTokenHandler := &existingTokenHandler{token: authConfig.RegistryToken}
		modifiers = append(modifiers, dockerAuth.NewAuthorizer(challengeManager, passThruTokenHandler))
	} else {
		scope := dockerAuth.RepositoryScope{
			Repository: repoName,
			Actions:    actions,
			Class:      "image",
		}

		creds := dockerRegistry.NewStaticCredentialStore(&authConfig)
		tokenHandlerOptions := dockerAuth.TokenHandlerOptions{
			Transport:   authTransport,
			Credentials: creds,
			Scopes:      []dockerAuth.Scope{scope},
			ClientID:    dockerRegistry.AuthClientID,
		}
		tokenHandler := dockerAuth.NewTokenHandlerWithOptions(tokenHandlerOptions)
		basicHandler := dockerAuth.NewBasicHandler(creds)
		modifiers = append(modifiers, dockerAuth.NewAuthorizer(challengeManager, tokenHandler, basicHandler))
	}

	tr := dockerTransport.NewTransport(base, modifiers...)
	repoNameRef, err := reference.WithName(repoName)
	if err != nil {
		return nil, err
	}

	return dockerRegistryClient.NewRepository(repoNameRef, rurl.String(), tr)
}

type existingTokenHandler struct {
	token string
}

func (th *existingTokenHandler) Scheme() string {
	return "bearer"
}

func (th *existingTokenHandler) AuthorizeRequest(req *http.Request, params map[string]string) error {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", th.token))
	return nil
}
