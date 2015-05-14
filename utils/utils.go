package utils

import (
    "bufio"
    "encoding/json"
    "io"
    "io/ioutil"
    "net"
    "os"
    "strings"
    "math/rand"
    "time"
    "golang.org/x/crypto/ssh"
    "bytes")


func ExecSSHCmd(publicIP string, privateKey string, uname string, command string) (string, error) {

	LogInfo("\nWaiting for server to start before adding ssh keys")
	e := WaitForSSH(publicIP)

	if e != nil {
		return "", e
	}

    k, e := ssh.ParsePrivateKey([]byte(privateKey))
    if e != nil {
        return "", e
    }

	c := &ssh.ClientConfig{
		User: uname,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(k)},
	}

	cn, _ := ssh.Dial("tcp", publicIP+":22", c)
	s, _ := cn.NewSession()
	defer s.Close()

	var sOut bytes.Buffer
	s.Stdout = &sOut
	s.Run(command)

	LogInfo(sOut.String())
	LogInfo("\nCommand Complete")

	return sOut.String(), nil
}


func WaitForSSH(publicIP string) error {
    for i:=0; i < 10; i++ {
        conn, e := net.Dial("tcp", publicIP+":22")
        if e != nil {
            time.Sleep(time.Second * 10)
            continue
        }
        defer conn.Close()
        if _, e = conn.Read(make([]byte, 1)); e != nil {
            continue
        }
        break
    }
    return nil
}

func LoadJsonConfig() (map[string]string, error) {
    var m map[string]string
    c, e := ioutil.ReadFile("./config.json")
    if e != nil {
        return nil, e
    }
    if e = json.Unmarshal(c, &m); e != nil {
        return nil, e
    }
    return m, nil
}

func LoadStdinToEnvAndKeys() error {
    rd := bufio.NewReader(os.Stdin)
    for {
        ln := ""
        ln, e := rd.ReadString('\n')
        if e == io.EOF {
            break
        } else if e != nil {
            return e
        } else if strings.Contains(ln, "=") {
            kv := strings.SplitN(ln, "=", 2)
            SetKey(kv[0], strings.TrimSpace(kv[1]))
            os.Setenv(kv[0], strings.TrimSpace(kv[1]))
        }
    }
    return nil
}

func RandSeq(n int) string {
    var letters = []rune("abcdefghijklmnopqrstuvwxyz")
    rand.Seed(time.Now().UTC().UnixNano())
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}