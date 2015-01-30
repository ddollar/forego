// Handle updating new releases from a godist server
package dist

import (
	"bitbucket.org/kardianos/osext"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/ddollar/go-update"
	"github.com/kr/binarydist"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var DigicertHighAssuranceCert = `-----BEGIN CERTIFICATE-----
MIIDxTCCAq2gAwIBAgIQAqxcJmoLQJuPC3nyrkYldzANBgkqhkiG9w0BAQUFADBs
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSswKQYDVQQDEyJEaWdpQ2VydCBIaWdoIEFzc3VyYW5j
ZSBFViBSb290IENBMB4XDTA2MTExMDAwMDAwMFoXDTMxMTExMDAwMDAwMFowbDEL
MAkGA1UEBhMCVVMxFTATBgNVBAoTDERpZ2lDZXJ0IEluYzEZMBcGA1UECxMQd3d3
LmRpZ2ljZXJ0LmNvbTErMCkGA1UEAxMiRGlnaUNlcnQgSGlnaCBBc3N1cmFuY2Ug
RVYgUm9vdCBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMbM5XPm
+9S75S0tMqbf5YE/yc0lSbZxKsPVlDRnogocsF9ppkCxxLeyj9CYpKlBWTrT3JTW
PNt0OKRKzE0lgvdKpVMSOO7zSW1xkX5jtqumX8OkhPhPYlG++MXs2ziS4wblCJEM
xChBVfvLWokVfnHoNb9Ncgk9vjo4UFt3MRuNs8ckRZqnrG0AFFoEt7oT61EKmEFB
Ik5lYYeBQVCmeVyJ3hlKV9Uu5l0cUyx+mM0aBhakaHPQNAQTXKFx01p8VdteZOE3
hzBWBOURtCmAEvF5OYiiAhF8J2a3iLd48soKqDirCmTCv2ZdlYTBoSUeh10aUAsg
EsxBu24LUTi4S8sCAwEAAaNjMGEwDgYDVR0PAQH/BAQDAgGGMA8GA1UdEwEB/wQF
MAMBAf8wHQYDVR0OBBYEFLE+w2kD+L9HAdSYJhoIAu9jZCvDMB8GA1UdIwQYMBaA
FLE+w2kD+L9HAdSYJhoIAu9jZCvDMA0GCSqGSIb3DQEBBQUAA4IBAQAcGgaX3Nec
nzyIZgYIVyHbIUf4KmeqvxgydkAQV8GK83rZEWWONfqe/EW1ntlMMUu4kehDLI6z
eM7b41N5cdblIZQB2lWHmiRk9opmzN6cN82oNLFpmyPInngiK3BD41VHMWEZ71jF
hS9OMPagMRYjyOfiZRYzy78aG6A9+MpeizGLYAiJLQwGXFK3xPkKmNEVX58Svnw2
Yzi9RKR/5CYrCsSXaQ3pjOLAEFe4yHYSkVXySGnYvCoCWw9E1CAx2/S6cCZdkGCe
vEsXCS+0yx5DaMkHJ8HSXPfqIbloEpw8nL+e/IBcm2PN7EeqJSdnoDfzAIJ9VNep
+OkuE6N36B9K
-----END CERTIFICATE-----`

type Dist struct {
	Host    string
	Name    string
	Project string
	Version string
}

type Release struct {
	Version string
	Url     string
}

// Initialize a new godist client, speciying a project name
// e.g. "ddollar/forego"
func NewDist(project string, version string) (d *Dist) {
	d = new(Dist)
	d.Host = "https://godist.herokuapp.com"
	d.Name = strings.Split(project, "/")[1]
	d.Project = project
	d.Version = version
	return
}

// Update the currently running binary to the latest version
func (d *Dist) Update() (to string, err error) {
	releases, err := d.fetchReleases()
	if len(releases) < 1 {
		return "", errors.New("no releases")
	}
	to = releases[0].Version
	return to, d.UpdateTo(to)
}

// Update the currently running binary to a specific version
func (d *Dist) UpdateTo(to string) (err error) {
	if d.Version == to {
		return errors.New("nothing to update")
	}
	binary, _ := osext.Executable()
	reader, err := os.Open(binary)
	if err != nil {
		return err
	}
	defer reader.Close()
	url := fmt.Sprintf("%s/projects/%s/diff/%s/%s/%s-%s", d.Host, d.Project, d.Version, to, runtime.GOOS, runtime.GOARCH)
	patch, err := d.httpGet(url)
	if err != nil {
		return err
	}
	writer := new(bytes.Buffer)
	err = binarydist.Patch(reader, writer, bytes.NewReader(patch))
	if err != nil {
		return err
	}
	reader.Close()
	err, _ = update.FromStream(writer)
	return
}

// Update to a specific version regardless of the starting version
func (d *Dist) FullUpdate(to string) (err error) {
	url := fmt.Sprintf("%s/projects/%s/releases/%s/%s-%s/%s", d.Host, d.Project, to, runtime.GOOS, runtime.GOARCH, d.Name)
	reader, err := d.httpGet(url)
	if err != nil {
		return err
	}
	err, _ = update.FromStream(bytes.NewReader(reader))
	return
}

func (d *Dist) fetchReleases() (releases []Release, err error) {
	body, err := d.httpGet(fmt.Sprintf("%s/projects/%s/releases/%s-%s", d.Host, d.Project, runtime.GOOS, runtime.GOARCH))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &releases)
	return
}

func (d *Dist) httpClient() (client *http.Client) {
	chain := d.rootCertificate()
	config := tls.Config{}
	config.RootCAs = x509.NewCertPool()
	for _, cert := range chain.Certificate {
		x509Cert, err := x509.ParseCertificate(cert)
		if err != nil {
			panic(err)
		}
		config.RootCAs.AddCert(x509Cert)
	}
	config.BuildNameToCertificate()
	tr := http.Transport{TLSClientConfig: &config}
	client = &http.Client{Transport: &tr}
	return
}

func (d *Dist) httpGet(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", fmt.Sprintf("%s/%s dist/%s (%s-%s)", d.Name, d.Version, Version, runtime.GOOS, runtime.GOARCH))
	res, err := d.httpClient().Do(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(body))
	}
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}

func (d *Dist) updateFromUrl(url string) (err error) {
	client := d.httpClient()
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err, _ = update.FromStream(res.Body)
	return
}

func (d *Dist) rootCertificate() (cert tls.Certificate) {
	certPEMBlock := []byte(DigicertHighAssuranceCert)
	var certDERBlock *pem.Block
	for {
		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			break
		}
		if certDERBlock.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
		}
	}
	return
}
