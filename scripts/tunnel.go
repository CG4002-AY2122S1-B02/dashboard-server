package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
)

type Endpoint struct {
	Host string
	Port int
	User string
}
func NewEndpoint(s string) *Endpoint {
	endpoint := &Endpoint{
		Host: s,
	}
	if parts := strings.Split(endpoint.Host, "@"); len(parts) > 1 {
		endpoint.User = parts[0]
		endpoint.Host = parts[1]
	}
	if parts := strings.Split(endpoint.Host, ":"); len(parts) > 1 {
		endpoint.Host = parts[0]
		endpoint.Port, _ = strconv.Atoi(parts[1])
	}
	return endpoint
}
func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}
type SSHTunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint
	Config *ssh.ClientConfig
	Log    *log.Logger
}
func (tunnel *SSHTunnel) logf(fmt string, args ...interface{}) {
	if tunnel.Log != nil {
		tunnel.Log.Printf(fmt, args...)
	}
}
func (tunnel *SSHTunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()
	tunnel.Local.Port = listener.Addr().(*net.TCPAddr).Port
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		tunnel.logf("accepted connection")
		go tunnel.forward(conn)
	}
}
func (tunnel *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		tunnel.logf("server dial error: %s", err)
		return
	}
	tunnel.logf("connected to %s (1 of 2)\n", tunnel.Server.String())
	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		tunnel.logf("remote dial error: %s", err)
		return
	}
	tunnel.logf("connected to %s (2 of 2)\n", tunnel.Remote.String())
	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			tunnel.logf("io.Copy error: %s", err)
		}
	}
	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}
func PrivateKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}
func NewSSHTunnel(tunnel string, auth ssh.AuthMethod, destination string) *SSHTunnel {
	// A random port will be chosen for us.
	localEndpoint := NewEndpoint("localhost:0")
	server := NewEndpoint(tunnel)
	if server.Port == 0 {
		server.Port = 22
	}
	sshTunnel := &SSHTunnel{
		Config: &ssh.ClientConfig{
			User: server.User,
			Auth: []ssh.AuthMethod{auth},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				// Always accept key.
				return nil
			},
		},
		Local:  localEndpoint,
		Server: server,
		Remote: NewEndpoint(destination),
	}
	return sshTunnel
}


func main() {
	//fmt.Println("x")
	//// Setup the tunnel, but do not yet start it yet.
	//tunnel1 := NewSSHTunnel(
	//	// User and host of tunnel server, it will default to port 22
	//	// if not specified.
	//	"e0325893@sunfire.comp.nus.edu.sg", //good
	//	// Pick ONE of the following authentication methods:
	//	ssh.Password("iLoveCapstoneB02"),                  // 2. password
	//	// The destination host and port of the actual server.
	//	"137.132.86.225:22",
	//)
	//
	//
	////tunnel: ssh_username@TUNNEL_ONE_SSH_ADDR:22
	////password: ssh_password
	////destination: remote_bind_address:22
	//
	//tunnel2 := NewSSHTunnel(
	//	// User and host of tunnel server, it will default to port 22
	//	// if not specified.
	//	"127.0.0.1:22",
	//	// Pick ONE of the following authentication methods:
	//	ssh.Password("cg4002b02"),                  // 2. password
	//	// The destination host and port of the actual server.
	//	"127.0.0.1:8888",
	//)
	//
	//// Start the server in the background. You will need to wait a
	//// small amount of time for it to bind to the localhost port
	//// before you can start sending connections.
	//go tunnel1.Start()
	//go tunnel2.Start()

	c, err := net.Dial("tcp", "127.0.0.1"+":10500")
	for err != nil {
		log.Fatal("Error listening:", err.Error())
	}
	defer c.Close()

	fmt.Println("Successfully Connected to Ultra96 on port:8888")

	br := bufio.NewReader(c)

	for {
		fmt.Println("x")
		packetData, err := br.ReadBytes('\x7F')
		if err != nil {
			fmt.Println("this" + err.Error())

			if err != io.EOF {
				fmt.Println("Error reading:", err.Error())
			}
			break
		}

		fmt.Println(packetData)
	}
	//time.Sleep(100 * time.Millisecond)
	// NewSSHTunnel will bind to a random port so that you can have
	// multiple SSH tunnels available. The port is available through:
	//   tunnel.Local.Port

}
