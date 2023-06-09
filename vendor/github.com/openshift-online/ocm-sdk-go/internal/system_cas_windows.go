//go:build windows
// +build windows

/*
Copyright (c) 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file contains the function that returns the trusted CA certificates for Windows. This is
// needed because currently Go doesn't know how to load the Windows trusted CA store. See the
// following issues for more information:
//
//	https://github.com/golang/go/issues/16736
//	https://github.com/golang/go/issues/18609

package internal

import (
	"crypto/x509"
)

// loadSystemCAs loads the certificates of the CAs that we will trust. Currently this uses a fixed
// set of CA certificates, which is obviusly going to break in the future, but there is not much we
// can do (or know to do) till Go learns to read the Windows CA trust store.
func loadSystemCAs() (pool *x509.CertPool, err error) {
	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(ssoCA1)
	pool.AppendCertsFromPEM(ssoCA2)
	pool.AppendCertsFromPEM(apiCA1)
	pool.AppendCertsFromPEM(apiCA2)
	return
}

// The SSO certificates has been obtained with the following command:
//
//	$ openssl s_client -connect sso.redhat.com:443 -showcerts
var ssoCA1 = []byte(`
-----BEGIN CERTIFICATE-----
MIIIJDCCBwygAwIBAgIQCI+RYFfnWP4GuTeagfIMxTANBgkqhkiG9w0BAQsFADB1
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMTQwMgYDVQQDEytEaWdpQ2VydCBTSEEyIEV4dGVuZGVk
IFZhbGlkYXRpb24gU2VydmVyIENBMB4XDTIyMDkxMzAwMDAwMFoXDTIzMDkxMzIz
NTk1OVowgcoxEzARBgsrBgEEAYI3PAIBAxMCVVMxGTAXBgsrBgEEAYI3PAIBAhMI
RGVsYXdhcmUxHTAbBgNVBA8MFFByaXZhdGUgT3JnYW5pemF0aW9uMRAwDgYDVQQF
EwcyOTQ1NDM2MQswCQYDVQQGEwJVUzEXMBUGA1UECBMOTm9ydGggQ2Fyb2xpbmEx
EDAOBgNVBAcTB1JhbGVpZ2gxFjAUBgNVBAoTDVJlZCBIYXQsIEluYy4xFzAVBgNV
BAMTDnNzby5yZWRoYXQuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAxVX7GF5Ca6Hh0GUMCgxz6zZtHz5SV3RaSyQSvaH7BT16XHhBh0PCEVRqvg8u
Dnq6RSx1pyk58fnLOfrNWqiT4+NvFbLLN5GkZl7C34HV+10GjbZyoA2PDQbB6yd5
3fy+XeWbnmlMvBzisW/wcnNb1aD6cmu91EYCsN4Q4lGt2jiPVeS3rplY6vIeq+3E
OZZBVim0KapDhCkIlKvguwS3RUPmdofTGsAoD22c7CjJHsPMfwWBlFfjfjIBRpFK
Hi1p8gR3oqLm9ZKnOMldA1bFpOH3G1gQ5JcdN3EWngY3Mh77aGId8RlEhmF8IlUK
9/iwUMbSS2OJ7tM72skJiQ4H3wIDAQABo4IEWDCCBFQwHwYDVR0jBBgwFoAUPdNQ
pdagre7zSmAKZdMh1Pj41g8wHQYDVR0OBBYEFLku0CpnRgwEXuntzUg6qsAcpwlU
MIIBBAYDVR0RBIH8MIH5gg5zc28ucmVkaGF0LmNvbYITYXV0aC5kZXYucmVkaGF0
LmNvbYIPYXV0aC5yZWRoYXQuY29tghVhdXRoLnN0YWdlLnJlZGhhdC5jb22CE29j
c3AuZGV2LnJlZGhhdC5jb22CF29jc3AucHJlcHJvZC5yZWRoYXQuY29tgg9vY3Nw
LnJlZGhhdC5jb22CFW9jc3Auc3RhZ2UucmVkaGF0LmNvbYIbc3NvLW5ldy1waXBl
bGluZS5yZWRoYXQuY29tgiFzc28tbmV3LXBpcGVsaW5lLnN0YWdlLnJlZGhhdC5j
b22CFHNzby5zdGFnZS5yZWRoYXQuY29tMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUE
FjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwdQYDVR0fBG4wbDA0oDKgMIYuaHR0cDov
L2NybDMuZGlnaWNlcnQuY29tL3NoYTItZXYtc2VydmVyLWczLmNybDA0oDKgMIYu
aHR0cDovL2NybDQuZGlnaWNlcnQuY29tL3NoYTItZXYtc2VydmVyLWczLmNybDBK
BgNVHSAEQzBBMAsGCWCGSAGG/WwCATAyBgVngQwBATApMCcGCCsGAQUFBwIBFhto
dHRwOi8vd3d3LmRpZ2ljZXJ0LmNvbS9DUFMwgYgGCCsGAQUFBwEBBHwwejAkBggr
BgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNlcnQuY29tMFIGCCsGAQUFBzAChkZo
dHRwOi8vY2FjZXJ0cy5kaWdpY2VydC5jb20vRGlnaUNlcnRTSEEyRXh0ZW5kZWRW
YWxpZGF0aW9uU2VydmVyQ0EuY3J0MAkGA1UdEwQCMAAwggGABgorBgEEAdZ5AgQC
BIIBcASCAWwBagB3AOg+0No+9QY1MudXKLyJa8kD08vREWvs62nhd31tBr1uAAAB
gziFhDQAAAQDAEgwRgIhALHB4Lnu/kVO/opahD3Zas42ZxE8FbSYeGDBdeZkgLz9
AiEAsrmKcVzM5VFLQppBvo3mqIg/QEj6VRHoUDCbJJjX/YIAdwA1zxkbv7FsV78P
rUxtQsu7ticgJlHqP+Eq76gDwzvWTAAAAYM4hYSPAAAEAwBIMEYCIQD0SP7fhp43
Hs6IeauXl3yFgoKNFE9+sN7+YnpkZe5p1wIhANvqhGaHvfsABYINVxMGkk1FGShN
ylG8ZZTlf/3Q+qY0AHYAtz77JN+cTbp18jnFulj0bF38Qs96nzXEnh0JgSXttJkA
AAGDOIWEmwAABAMARzBFAiEAljZzvCWMKfTUfbM01ol0ptTEFUBgkBxL6BOWenN/
SeUCICYHgC3Vq27vl3jHvthM721q6tuwqkErXZGJ1RdcXcodMA0GCSqGSIb3DQEB
CwUAA4IBAQBk47U9WGJ3l2CsmiTapQLlt2qUeqYEyJmJ7zETuP64cUeDcbJ/Jsdv
gkoW/BQVDtGU4CKJvHB/WhlohXpbxDGsYuZlYDift1SbFgMnolM63X0JM/jBzVK8
y64LIVLsnIG1TQF/5e1fFVB+bZ+u1gQU3YPmDSyLobpWriHK9rm0oR5njUb/5Tjo
BspZzXr1CpiP8EOQN9W85gaTspZOINYGsV805Gz/vX6F9jIKlirplIepSR+DLWy4
y5McLs8V/KyofxDu7Q92sD96bqnEmDiGhjV1zpXW9IlS0pRvv9EH3DDgqUolfbFm
oK1RA+aDJthxm8EsjbezqLwWrMICIjmj
-----END CERTIFICATE-----
`)

var ssoCA2 = []byte(`
-----BEGIN CERTIFICATE-----
MIIEtjCCA56gAwIBAgIQDHmpRLCMEZUgkmFf4msdgzANBgkqhkiG9w0BAQsFADBs
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSswKQYDVQQDEyJEaWdpQ2VydCBIaWdoIEFzc3VyYW5j
ZSBFViBSb290IENBMB4XDTEzMTAyMjEyMDAwMFoXDTI4MTAyMjEyMDAwMFowdTEL
MAkGA1UEBhMCVVMxFTATBgNVBAoTDERpZ2lDZXJ0IEluYzEZMBcGA1UECxMQd3d3
LmRpZ2ljZXJ0LmNvbTE0MDIGA1UEAxMrRGlnaUNlcnQgU0hBMiBFeHRlbmRlZCBW
YWxpZGF0aW9uIFNlcnZlciBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBANdTpARR+JmmFkhLZyeqk0nQOe0MsLAAh/FnKIaFjI5j2ryxQDji0/XspQUY
uD0+xZkXMuwYjPrxDKZkIYXLBxA0sFKIKx9om9KxjxKws9LniB8f7zh3VFNfgHk/
LhqqqB5LKw2rt2O5Nbd9FLxZS99RStKh4gzikIKHaq7q12TWmFXo/a8aUGxUvBHy
/Urynbt/DvTVvo4WiRJV2MBxNO723C3sxIclho3YIeSwTQyJ3DkmF93215SF2AQh
cJ1vb/9cuhnhRctWVyh+HA1BV6q3uCe7seT6Ku8hI3UarS2bhjWMnHe1c63YlC3k
8wyd7sFOYn4XwHGeLN7x+RAoGTMCAwEAAaOCAUkwggFFMBIGA1UdEwEB/wQIMAYB
Af8CAQAwDgYDVR0PAQH/BAQDAgGGMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEF
BQcDAjA0BggrBgEFBQcBAQQoMCYwJAYIKwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRp
Z2ljZXJ0LmNvbTBLBgNVHR8ERDBCMECgPqA8hjpodHRwOi8vY3JsNC5kaWdpY2Vy
dC5jb20vRGlnaUNlcnRIaWdoQXNzdXJhbmNlRVZSb290Q0EuY3JsMD0GA1UdIAQ2
MDQwMgYEVR0gADAqMCgGCCsGAQUFBwIBFhxodHRwczovL3d3dy5kaWdpY2VydC5j
b20vQ1BTMB0GA1UdDgQWBBQ901Cl1qCt7vNKYApl0yHU+PjWDzAfBgNVHSMEGDAW
gBSxPsNpA/i/RwHUmCYaCALvY2QrwzANBgkqhkiG9w0BAQsFAAOCAQEAnbbQkIbh
hgLtxaDwNBx0wY12zIYKqPBKikLWP8ipTa18CK3mtlC4ohpNiAexKSHc59rGPCHg
4xFJcKx6HQGkyhE6V6t9VypAdP3THYUYUN9XR3WhfVUgLkc3UHKMf4Ib0mKPLQNa
2sPIoc4sUqIAY+tzunHISScjl2SFnjgOrWNoPLpSgVh5oywM395t6zHyuqB8bPEs
1OG9d4Q3A84ytciagRpKkk47RpqF/oOi+Z6Mo8wNXrM9zwR4jxQUezKcxwCmXMS1
oVWNWlZopCJwqjyBcdmdqEU79OX2olHdx3ti6G8MdOu42vi/hw15UJGQmxg7kVkn
8TUoE6smftX3eg==
-----END CERTIFICATE-----
`)

// The API certificates have been obtained with the following command:
//
//	$ openssl s_client -connect api.openshift.com:443 -showcerts
var apiCA1 = []byte(`
-----BEGIN CERTIFICATE-----
MIIFmzCCBSGgAwIBAgIQBQ+NDVzn4QRZqly06Ep3yjAKBggqhkjOPQQDAzBWMQsw
CQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMTAwLgYDVQQDEydEaWdp
Q2VydCBUTFMgSHlicmlkIEVDQyBTSEEzODQgMjAyMCBDQTEwHhcNMjIwNDI5MDAw
MDAwWhcNMjMwNTAzMjM1OTU5WjBsMQswCQYDVQQGEwJVUzEXMBUGA1UECBMOTm9y
dGggQ2Fyb2xpbmExEDAOBgNVBAcTB1JhbGVpZ2gxFjAUBgNVBAoTDVJlZCBIYXQs
IEluYy4xGjAYBgNVBAMTEWFwaS5vcGVuc2hpZnQuY29tMHYwEAYHKoZIzj0CAQYF
K4EEACIDYgAEGdA2EA1g7ynRHLSKRjOlyPtqPnFoEAadzuNqedwTWBN3bA16Q6tr
j18rsF6kERbHmypiNMHDwymxQuYhEEJdlUL9zWnDl//Tt3P97WlJ0yQ96i478ofG
E73IYHzK8C8No4IDnDCCA5gwHwYDVR0jBBgwFoAUCrwIKReMpTlteg7OM8cus+37
w3owHQYDVR0OBBYEFJC5XsZwzzRxFn5ghyFFxkdygoynMDMGA1UdEQQsMCqCEWFw
aS5vcGVuc2hpZnQuY29tghV3d3cuYXBpLm9wZW5zaGlmdC5jb20wDgYDVR0PAQH/
BAQDAgeAMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjCBmwYDVR0fBIGT
MIGQMEagRKBChkBodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vRGlnaUNlcnRUTFNI
eWJyaWRFQ0NTSEEzODQyMDIwQ0ExLTEuY3JsMEagRKBChkBodHRwOi8vY3JsNC5k
aWdpY2VydC5jb20vRGlnaUNlcnRUTFNIeWJyaWRFQ0NTSEEzODQyMDIwQ0ExLTEu
Y3JsMD4GA1UdIAQ3MDUwMwYGZ4EMAQICMCkwJwYIKwYBBQUHAgEWG2h0dHA6Ly93
d3cuZGlnaWNlcnQuY29tL0NQUzCBhQYIKwYBBQUHAQEEeTB3MCQGCCsGAQUFBzAB
hhhodHRwOi8vb2NzcC5kaWdpY2VydC5jb20wTwYIKwYBBQUHMAKGQ2h0dHA6Ly9j
YWNlcnRzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydFRMU0h5YnJpZEVDQ1NIQTM4NDIw
MjBDQTEtMS5jcnQwCQYDVR0TBAIwADCCAX8GCisGAQQB1nkCBAIEggFvBIIBawFp
AHcArfe++nz/EMiLnT2cHj4YarRnKV3PsQwkyoWGNOvcgooAAAGAc3NizAAABAMA
SDBGAiEAwR6POASNS3R+vdTvqP0LpoP0VB8m6JV3P8xn//Z10fYCIQDyV7C7V7yf
XuXzorhEXWg2npekZVkT0fS7jTUSmZJgTgB2ADXPGRu/sWxXvw+tTG1Cy7u2JyAm
Ueo/4SrvqAPDO9ZMAAABgHNzYsoAAAQDAEcwRQIgcTIfCrJMjzjuX3qnyR4lWwI9
e7AIjAm7P+c9B7CpgoECIQDUhSzYogCMx5QdFqFEi18KlX89bjKscSukieVbzc98
DwB2ALNzdwfhhFD4Y4bWBancEQlKeS2xZwwLh9zwAw55NqWaAAABgHNzYvwAAAQD
AEcwRQIhAN3WIvRA3tejaDj1ceCPQgb4tK97AOqWbdzOa4CbA6+oAiAbv0wFXElw
vuMPbYUjraXeyntGPfmTf/8MNlkd0aW20DAKBggqhkjOPQQDAwNoADBlAjEAiho0
712CjvYGMezC6cF6V8+lPWjb1PGRXzTMNYNBvUp0jCla5oznm6hLRpcqw8z3AjAq
sjE4cTA4Jy91GD6dAwd9uhwvEl5IpHIMb8GkxOybDAAHsioc3AGlfXFuls9GBdc=
-----END CERTIFICATE-----
`)

var apiCA2 = []byte(`
-----BEGIN CERTIFICATE-----
MIIEFzCCAv+gAwIBAgIQB/LzXIeod6967+lHmTUlvTANBgkqhkiG9w0BAQwFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBD
QTAeFw0yMTA0MTQwMDAwMDBaFw0zMTA0MTMyMzU5NTlaMFYxCzAJBgNVBAYTAlVT
MRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxMDAuBgNVBAMTJ0RpZ2lDZXJ0IFRMUyBI
eWJyaWQgRUNDIFNIQTM4NCAyMDIwIENBMTB2MBAGByqGSM49AgEGBSuBBAAiA2IA
BMEbxppbmNmkKaDp1AS12+umsmxVwP/tmMZJLwYnUcu/cMEFesOxnYeJuq20ExfJ
qLSDyLiQ0cx0NTY8g3KwtdD3ImnI8YDEe0CPz2iHJlw5ifFNkU3aiYvkA8ND5b8v
c6OCAYIwggF+MBIGA1UdEwEB/wQIMAYBAf8CAQAwHQYDVR0OBBYEFAq8CCkXjKU5
bXoOzjPHLrPt+8N6MB8GA1UdIwQYMBaAFAPeUDVW0Uy7ZvCj4hsbw5eyPdFVMA4G
A1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwdgYI
KwYBBQUHAQEEajBoMCQGCCsGAQUFBzABhhhodHRwOi8vb2NzcC5kaWdpY2VydC5j
b20wQAYIKwYBBQUHMAKGNGh0dHA6Ly9jYWNlcnRzLmRpZ2ljZXJ0LmNvbS9EaWdp
Q2VydEdsb2JhbFJvb3RDQS5jcnQwQgYDVR0fBDswOTA3oDWgM4YxaHR0cDovL2Ny
bDMuZGlnaWNlcnQuY29tL0RpZ2lDZXJ0R2xvYmFsUm9vdENBLmNybDA9BgNVHSAE
NjA0MAsGCWCGSAGG/WwCATAHBgVngQwBATAIBgZngQwBAgEwCAYGZ4EMAQICMAgG
BmeBDAECAzANBgkqhkiG9w0BAQwFAAOCAQEAR1mBf9QbH7Bx9phdGLqYR5iwfnYr
6v8ai6wms0KNMeZK6BnQ79oU59cUkqGS8qcuLa/7Hfb7U7CKP/zYFgrpsC62pQsY
kDUmotr2qLcy/JUjS8ZFucTP5Hzu5sn4kL1y45nDHQsFfGqXbbKrAjbYwrwsAZI/
BKOLdRHHuSm8EdCGupK8JvllyDfNJvaGEwwEqonleLHBTnm8dqMLUeTF0J5q/hos
Vq4GNiejcxwIfZMy0MJEGdqN9A57HSgDKwmKdsp33Id6rHtSJlWncg+d0ohP/rEh
xRqhqjn1VtvChMQ1H3Dau0bwhr9kAMQ+959GG50jBbl9s08PqUU643QwmA==
-----END CERTIFICATE-----
`)
