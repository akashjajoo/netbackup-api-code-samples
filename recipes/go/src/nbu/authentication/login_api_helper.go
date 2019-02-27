//This script consists of the helper functions to login to NetBackup using the NetBackup webservices
package authentication

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

//###################
// Global Variables
//###################
const (
    port              = "1556"
    policiesUri       = "config/policies/"
    contentType       = "application/vnd.netbackup+json;version=2.0"
    testPolicyName    = "vmware_test_policy"
    testClientName    = "MEDIA_SERVER"
    testScheduleName  = "vmware_test_schedule"
)

//##############################################################
// Setup the HTTP client to make NetBackup Policy API requests
//##############################################################
func GetHTTPClient() *http.Client {
    tlsConfig := &tls.Config {
        InsecureSkipVerify: true, //for this test, ignore ssl certificate
    }

    tr := &http.Transport{TLSClientConfig: tlsConfig}
    client := &http.Client{Transport: tr}

    return client
}

//#####################################
// Login to the NetBackup webservices
//#####################################
func Login(nbmaster string, httpClient *http.Client, username string, password string, domainName string, domainType string) string {
    fmt.Printf("\nLogin to the NetBackup webservices...\n")

    loginDetails := map[string]string{"userName": username, "password": password}
    if len(domainName) > 0 {
        loginDetails["domainName"] = domainName
    }
    if len(domainType) > 0 {
        loginDetails["domainType"] = domainType
    }
    loginRequest, _ := json.Marshal(loginDetails)

    uri :=  "https://" + nbmaster + ":" + port + "/netbackup/login"

    request, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(loginRequest))
    request.Header.Add("Content-Type", contentType);

    response, err := httpClient.Do(request)

    token := ""
    if err != nil {
        fmt.Printf("The HTTP request failed with error: %s\n", err)
        panic("Unable to login to the NetBackup webservices")
    } else {
        if response.StatusCode == 201 {
            data, _ := ioutil.ReadAll(response.Body)
            var obj interface{}
            json.Unmarshal(data, &obj)
            loginResponse := obj.(map[string]interface{})
            token = loginResponse["token"].(string)
        } else {
            responseBody, _ := ioutil.ReadAll(response.Body)
            fmt.Printf("%s\n", responseBody)
            panic("Unable to login to the NetBackup webservices")
        }
    }

    return token
}

