package http_routes

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"code.cloudfoundry.org/routing-acceptance-tests/helpers"

	"testing"
)

func Rtr(args ...string) *Session {
	portString := fmt.Sprintf("%d", routerApiConfig.OAuth.Port)
	oauthUrl := routerApiConfig.OAuth.TokenEndpoint + ":" + portString
	args = append(args, "--api", routerApiConfig.RoutingApiUrl, "--client-id", routerApiConfig.OAuth.ClientName, "--client-secret", routerApiConfig.OAuth.ClientSecret, "--oauth-url", oauthUrl)
	if routerApiConfig.SkipSSLValidation {
		args = append(args, "--skip-tls-verification")
	}
	rtr, exists := os.LookupEnv("ROUTING_API_CLI_BINARY")
	if !exists {
		fmt.Println("You need routing-api-cli (rtr) installed and set ROUTING_API_CLI_BINARY env variable")
		os.Exit(1)
	}
	session, err := Start(exec.Command(rtr, args...), GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}

const (
	DEFAULT_TIMEOUT          = 30 * time.Second
	DEFAULT_POLLING_INTERVAL = 1 * time.Second
	CF_PUSH_TIMEOUT          = 2 * time.Minute
	DEFAULT_MEMORY_LIMIT     = "256M"
)

var routerApiConfig helpers.RoutingConfig

func TestRouting(t *testing.T) {
	RegisterFailHandler(Fail)

	routerApiConfig = helpers.LoadConfig()

	BeforeSuite(func() {
		Expect(routerApiConfig.OAuth.ClientSecret).ToNot(Equal(""), "Must provide a client secret for the routing suite")
	})

	BeforeEach(func() {
		if !routerApiConfig.IncludeHttpRoutes {
			Skip("Skipping this test because Config.IncludeHttpRoutes is set to `false`.")
		}
	})

	RunSpecs(t, "HTTP Routes Suite")
}
