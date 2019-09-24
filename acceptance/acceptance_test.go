package acceptance

import (
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var (
	serverPort     int
	hostName       string
	apiServicePath string
	err            error
	session        *gexec.Session
	client         *http.Client
	args           apiCommandArgs
)

type apiCommandArgs struct {
	serverPort int
	hostName   string
}

func (args apiCommandArgs) toStringSlice() []string {
	stringArgs := []string{
		"--server-port", strconv.Itoa(args.serverPort),
		"--server-hostname", args.hostName,
	}
	return stringArgs
}

var _ = BeforeSuite(func() {
	var err error
	apiServicePath, err = gexec.Build("github.com/schedulerapi/cmd/")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	os.Remove(apiServicePath)
})

var _ = BeforeEach(func() {
	serverPort = 4567 + GinkgoParallelNode() + rand.Int()%20
	hostName = "localhost"
	args = apiCommandArgs{
		serverPort: serverPort,
		hostName:   hostName,
	}
	client = &http.Client{}

})

// Build and run the api service
var _ = JustBeforeEach(func() {
	apiServiceCMD := exec.Command(apiServicePath, args.toStringSlice()...)

	var err error
	session, err = gexec.Start(apiServiceCMD, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	Eventually(session.Out, "3s", "20ms").Should(gbytes.Say("Scheduler API Service listening on port:%d", serverPort))
})

var _ = AfterEach(func() {
	session.Interrupt()
})
