package utility

import (
	"os"
	"testing"
)

const (
	publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6IkkaGgB8MfsL8/cB2SP
82o9COxFZPuD2SxZArgSCVLcLX8Iw3zc5Sl5VWLGwJC33Ranm+B1Z27sJoB5fr1E
jZZRWIX4MKjNvQ32OPCwBHn74Sq8TGiR0cCt8TryMMlqX4K5elCjw+dynVbg9rwR
ogxvTB/mKI5f9DPiKpQxQ+u141Yn87jJY4kKzTNPP2sHYVFztxNAGCxL3uXdw3DN
yjf4wjzqX7SkyBIHETk9MB/24JbEq1qGsmQb4uWWC25zwR2rSqau1G32nvTvtrgc
x3Zava4pIMOuwKkVLXj+fW7mYT8TognG3ij1KVhWiRtmJ46KF6P4q098yOneIcrs
LQIDAQAB
-----END PUBLIC KEY-----`
	privateKey = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDoiSRoaAHwx+wv
z9wHZI/zaj0I7EVk+4PZLFkCuBIJUtwtfwjDfNzlKXlVYsbAkLfdFqeb4HVnbuwm
gHl+vUSNllFYhfgwqM29DfY48LAEefvhKrxMaJHRwK3xOvIwyWpfgrl6UKPD53Kd
VuD2vBGiDG9MH+Yojl/0M+IqlDFD67XjVifzuMljiQrNM08/awdhUXO3E0AYLEve
5d3DcM3KN/jCPOpftKTIEgcROT0wH/bglsSrWoayZBvi5ZYLbnPBHatKpq7Ubfae
9O+2uBzHdlq9rikgw67AqRUteP59buZhPxOiCcbeKPUpWFaJG2YnjooXo/irT3zI
6d4hyuwtAgMBAAECggEACqdB0kC1mgzmvrgEqhgK2kmLO03rzfkR7NCx0USSmve/
W5w+0An36s3QH7/TQD5BFjF0A1mwt0jnK8pmBo7wRZZV6rrUiJIrPtaVab1pKwNV
rVayYsPFrnjn3y3mq6KGq5KHIdnw7sg7QYcZ7mPmYOi17qRlfIUtUzhPS8kXSre6
8DwfYbsoW0WHcmKKYZitEhwyauqA/rH0w8vXYV7AIA3fZrpS3pr8bOqxiuIrjdBD
d/EejaSfNEpEfy14lh1a0/AU05nboeCcGf2CCxhM4Q2hly1/NUn45+qRIHdrco58
lNBPMlrvvB1/HPwwk6WiFBcjCmCWS8g+z3HKd69uwQKBgQD/C97xIeIaipbR8E7P
7rfMj7Aqh4Ovq4ITMJtb7DWCtXdcxVcHr3MQ/HGiqKvaGUyS6jdrEzjI5B/PeRTX
DhZCB2IFNBsJ/4eQvrZi2r35YzDUQ3pq3qQ/fe3fpRcXoHHIxXoieAuuRPO1r5m3
Tp3kjRVdpb4NTF2NRH2lPFLwDQKBgQDpZ7lpDThtwXdXLlGyE8Z8GvmzU+Oj7ryZ
TdmwrKHHaP3fw0i9IvI89s0+X/ceHQXXeas726UEBtc5O+eH/PdBhEgCtcHO43WO
HHkJi+3c2YGfebnJBQRY44nLwxTf6pf7mxi2GxpaVM7UFttvc/BrrATdXOmn+3DB
X8tBKzjEoQKBgQDMDOJMR5CPLYwm4L0dTN8OMXN/QzZPSMdjtQLHE39oWOjrdxL/
GhbUYzRDL/F2J8GE1RCLgTBwQVtV8YiD2khigWehxCNR53e9jWd8RYeyS/KYEHiT
ohcEmSrEQF/uTjZaq+vgQe0OeyoElT5FUweuAFY0u1MHbq52RHLFzTKJzQKBgQC2
NMPYD3sCq2oXg9BA3REwpvpRFOb7fY57evu64Tk1629stA1foR1LnDsjO1U1i+CY
kqGrC89pMlHnmy0mysLWwYZZnzwZ3xVRCEcwvazFoIKBVUxEcgcvwQk8KSFtn7xf
rXcACm5rIBOKHAHXosGHvHTbvgGlojMmsjqAuFYLoQKBgQCc0Rv5BR1umylfnLFb
pAk78J1bmTXgSbRSp8gjeTqoNz8Y+XRhTYTszwItCDYP9FESNmbZDCW86BuwGOGX
fQvsWAaj8pKiLgCsW1O2eNt2jhJIp9Ptaoo9WU6Ix5/aFrMJVJpjfn8MUVQGpvOm
hbvsTfC4zscfxmfg23ePMWpJFQ==
-----END PRIVATE KEY-----`

	publicKeyFile  = "/tmp/pub.pem"
	privateKeyFile = "/tmp/pri.pem"
)

func init() {
	fpub, err := os.Create(publicKeyFile)

	if err != nil {
		panic(err)
	}
	defer fpub.Close()

	fpub.WriteString(publicKey)

	fpri, err := os.Create(privateKeyFile)
	if err != nil {
		panic(err)
	}
	defer fpri.Close()
	fpri.WriteString(privateKey)

}

func TestRsaCryptography(t *testing.T) {

	t.Cleanup(func() {
		os.Remove(publicKeyFile)
		os.Remove(privateKeyFile)
	})

	plaintext := "test"

	encrypted, err := EncryptStringWithRsaKey(plaintext, publicKeyFile)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	decrypted, err := DecryptStringWithRsaKey(*encrypted, privateKeyFile)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	if *decrypted != plaintext {
		t.Errorf("%s != %s\n", *decrypted, plaintext)
	}
}
