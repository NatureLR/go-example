package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	key := "LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQkc1dmJtVUFBQUFFYm05dVpRQUFBQUFBQUFBQkFBQUJsd0FBQUFkemMyZ3RjbgpOaEFBQUFBd0VBQVFBQUFZRUEybEZyc1FVamdDd2ZnNzFMWlVBMCtuYnN4b3p6YXkyOHNNSHJtdXp1RkcvOWk5VlR6c2x6CkJZSjNyUWxGMW5UWG1lRE5wNVliMnd0V2w2RlYzdkZNS3VrSmMySTNnWGY2NkdONW5pODZVR05XSWZ4b21kR0pmeUVNMUIKVk8yNUlpbXB2NVhDZmdyOG0wdXVzSjdBWU0vRDZKUHlLeDRtVjY1MjhBVGFUOVRUYVZiUndteGkvRnFmTU5uN0JMa0dOTgpURk1ZTzlHUGR4U2NtVjU5TjJVLzVhZThYR0tWQnZ6N1VFa291dTJpQ0gwRjc3bzNYQkcweGhEMG9JbldSSVBnY0F5YkhNCnpLelAvU2xqSjNWdG0xMStjQUlSL2lDVlRaZmtreFBLYkxKeUNKelJ3cVV1ajVRcnVaRHJsbWNmVk04OGhwQTgwMWZOMzQKRldMelM3Z1FHbGk2bFJZRWlTK0xaMDE5bHM5OVFhdk1PSlhsQmJlZjFrU1prUEx6ZnZucElGcmNLWUJqajlDS05jaE5MRgp1M2RROXhzSEkweU9SNFZpZks5SnpyK1V4Y3A4QmVTVzJHMmQ3dnBKQ05oUlNrVjJwRkNIRTExRW5tMlVKYzd1MTEzR0FRClZEckNSTWtHS0c5RE1OcDBaZkNMbTR0RzBaeWcwMUdQOGpHZkVsalZBQUFGa0E1bzNoME9hTjRkQUFBQUIzTnphQzF5YzIKRUFBQUdCQU5wUmE3RUZJNEFzSDRPOVMyVkFOUHAyN01hTTgyc3R2TERCNjVyczdoUnYvWXZWVTg3SmN3V0NkNjBKUmRaMAoxNW5nemFlV0c5c0xWcGVoVmQ3eFRDcnBDWE5pTjRGMyt1aGplWjR2T2xCalZpSDhhSm5SaVg4aEROUVZUdHVTSXBxYitWCnduNEsvSnRMcnJDZXdHRFB3K2lUOGlzZUpsZXVkdkFFMmsvVTAybFcwY0pzWXZ4YW56RFord1M1QmpUVXhUR0R2UmozY1UKbkpsZWZUZGxQK1dudkZ4aWxRYjgrMUJKS0xydG9naDlCZSs2TjF3UnRNWVE5S0NKMWtTRDRIQU1teHpNeXN6LzBwWXlkMQpiWnRkZm5BQ0VmNGdsVTJYNUpNVHlteXljZ2ljMGNLbExvK1VLN21RNjVabkgxVFBQSWFRUE5OWHpkK0JWaTgwdTRFQnBZCnVwVVdCSWt2aTJkTmZaYlBmVUdyekRpVjVRVzNuOVpFbVpEeTgzNzU2U0JhM0NtQVk0L1FpalhJVFN4YnQzVVBjYkJ5Tk0KamtlRllueXZTYzYvbE1YS2ZBWGtsdGh0bmU3NlNRallVVXBGZHFSUWh4TmRSSjV0bENYTzd0ZGR4Z0VGUTZ3a1RKQmlodgpRekRhZEdYd2k1dUxSdEdjb05OUmovSXhueEpZMVFBQUFBTUJBQUVBQUFHQkFJVGN3RGsrODFmeGdreGVTeUFYYnladWNiCkp6M1VBQTJiQ0lrNlg1UXZyVkhPeVlxeVJSbk5waGlBdWFkUklLa1p0b0lFQTVMa0trSjlLbnNPYTQycTNTbkpuSDBCZk8KdUxmc3Nmcitxdko1UWRYMUVvTnA3YytjZ1g0Z2FabGUyZ2hWbSsvbHBPdldTVkxuNzJYZ1dNNjZFRFNJSE5HM0NKRUlFSgpzd0thZHY3SWcyZGJKdktGQkJScTFFVVBoU05weVloVXNDUWRrcWhoTWdXZnY0ZG1hQktqVGZHZXJpMHQ4Tlp1Zzc3anhVCk9lU200MEg0MUpjbGducWc1L0ZTZjRnbno2MFNLdnZHSit6ZmFyaEtPTU5kKzZIV3I4WlZoblBpelc1ai9Sb3lwN05zWEYKb2VDOUp5a3BjN3h6TXM4bXZQNC9MYlRRUnJMT243TDBaZnNWcWZCdlRxUkZ5OU9weDVXdnh3d2tKQWZLTW5BZ1JGaFpDYgpsbEdYWmdCalpVQVBuV2lPbWhIN0FjRys1MU1oZHBDNG9xUkkyVE1qRlZCbjg4Y1QzYVpZSTdoclVNMnVSNG5xejRwRWJzCkNVR0Z6MXNIcmRUU2E3QlhmWDJTekJFZWxETllVS1JGaEU1OFJNYnNWV1F3S29QRzY2d21vZTlJZ1RYdWxmM0crNWtRQUEKQU1BaHFaKzFHUisxOUNjbXplS0hPV1cxYU1Dbkl5VDl6ZUFlQWRhVG55eHhZUnRNL1RINGNtdFE1Y1FMcHlya05ES25CUwpyL0U4dndYUWwzd2ZqUGRtWllwb0NBTTFYbjdzVkVzWHFwbkwrRHgyejRTWlFSRmx3T1FZdUxnam5wT0tDcmxQYkJUYlFtCmk4Q2pMVGRrWVIrYiszQ0lSL1BxYkxXTmNaUlJRd3VMN2ZmNzZRM3ZudUNUSUVEWkxpRWltaGdGMDNFQnZSYThHVk9LYzAKOFZJWk1Fb3BFTHZLRVlidzhKVUNYS05nYnFBc1Q4RjRxUXdmTVJRbTc2b1pwT1I4MEFBQURCQVBjOHF6RW5zZVR3cE1IeAovVGxySzBvYzlGU1lGb0hhdnVOT3lRMVdySGM3WGlSTHJHaS9Pa2hmN0RQWjM1ZmpaR1JnRmp5d0U4YTRtalpGS25KY2pYCnB2ZStWOCtSSnYzUHVKaUx0V1AybjRoSHl0Y2JrS04yS2wwVU9zeFJOSm51dVpEUUErYTB6V3JBS2pVOFBid09iMVdseGUKbmlxYnBtZENEbXNKT0Z0Umg1Z1NjTGJvQ2IyYjMrdkU4QkZyTjlPa0pLRjE0M0ZrQzZvMFhzSzJ4TENuTjRZWFhvNWQ4cAp6R2l1REFsaWwyV0I3SG9OOWlHa3lMejhjL3g0OGRZd0FBQU1FQTRnNWFWZCs4eGlsbEFuN3Jmc2lGLzlGSzd2OFJkYnBDClpFaVR4YnNCZ2NKL2NwQmJ4d1puejEwYk44WVg2MUord0RtNXlxLzI5eUIwVGJpbEV4U01QTnRMSTI3aERsaWcxQWVXSWYKb3Y0Y3hsbnhQdmR6K2pUenR0M20wdTloQ2l0WnlTWTN5L1ZtNHN5TktBUTFnSFIrTXRkRExURzM4NHFQTzZWZXdZRjc0RgptNm1vZ0cxRkZtU2lYZy9nZXpldUZKMjBxeS8wVUozV2kwVFR0QXJ6ak1vOGYxN1JKWURzMGh2RjcwYXE5VjByVGNrQndKCjNrOEFuUjlKTDJqa0puQUFBQUducDRla0I2ZUhwa1pVMWhZMEp2YjJzdFVISnZMbXh2WTJGcwotLS0tLUVORCBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0K"
	sc, err := newSSHConnect("root", "", key, "101.36.124.102:22")
	if err != nil {
		panic(err)
	}

	if err := scp(sc, "/tmp/xx", "test"); err != nil {
		panic(err)
	}
}

func scp(sc *ssh.Client, remotePath, content string) error {
	sftpClient, err := sftp.NewClient(sc)
	if err != nil {
		return fmt.Errorf("create sftp client failed: %s", err.Error())
	}
	defer sftpClient.Close()

	dstFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("create file on ssh server failed: %s", err.Error())
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("write file on server failed: %s", err.Error())
	}

	return nil
}

func newSSHConnect(user, password, privateKey, address string) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		err          error
	)
	if len(password) > 0 {
		auth = append(auth, ssh.Password(password))
	} else {
		privatekeybyte, err := base64.StdEncoding.DecodeString(privateKey)
		if err != nil {
			return nil, err
		}

		signer, err := ssh.ParsePrivateKey(privatekeybyte)
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 10 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	if client, err = ssh.Dial("tcp", address, clientConfig); err != nil {
		return nil, err
	}

	return client, nil
}
