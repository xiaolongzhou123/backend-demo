package myssh

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sso/pkg"
	"time"

	"github.com/gorilla/websocket"
	myuuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"golang.org/x/crypto/ssh"
)

var interval time.Duration

type Conn struct {
	net.Conn
	// ReadTimeout  time.Duration
	// WriteTimeout time.Duration
}

// Cli ...
type Cli struct {
	IP         string      //IP地址
	Username   string      //用户名
	Password   string      //密码
	Privatekey string      //私钥登陆
	Port       int         //端口号
	client     *ssh.Client //ssh客户端
	uuid       string      //uuid唯一标识
	exit       chan bool

	//
	cols      int
	rows      int
	StdinPipe io.WriteCloser
	StdOutput *wsBufferWriter
	Session   *ssh.Session
	backspace bool
}

type MessData struct {
	Op   string
	Cols int
	Rows int
	Data string
}

// New 创建命令行对象
// ip IP地址
// username 用户名
// password 密码
// port 端口号,默认22
func NewConn(ip, username, password, Privatekey string, backspace bool, port int) (*Cli, error) {

	cli := new(Cli)
	cli.IP = ip
	cli.Username = username
	cli.Password = password
	cli.Privatekey = Privatekey
	cli.Port = port
	cli.exit = make(chan bool)
	cli.backspace = backspace

	cli.StdOutput = new(wsBufferWriter)
	uid := myuuid.NewV4()

	// if err != nil {
	//      return nil, fmt.Errorf("uuid create error:%v", err)
	// }
	cli.uuid = uid.String()

	if port <= 0 {
		cli.Port = 22
	}

	return cli, cli.connect()
}

// Run 执行 shell脚本命令
func (c *Cli) Run(shell string) (string, error) {
	session, err := c.newSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	buf, err := session.CombinedOutput(shell)
	return string(buf), err
}

// RunTerminal 执行带交互的命令
func (c *Cli) RunTerminal(shell string) error {
	session, err := c.newSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return c.runTerminalSession(session, shell)
}

// runTerminalSession 执行带交互的命令
func (c *Cli) runTerminalSession(session *ssh.Session, shell string) error {

	session.Stdout = os.Stdout
	session.Stderr = os.Stdin
	session.Stdin = os.Stdin

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", c.cols, c.rows, modes); err != nil {
		return err
	}

	session.Run(shell)
	return nil
}

// EnterTerminal 完全进入终端
func (c *Cli) Enter() error {

	session, err := c.newSession()
	if err != nil {
		return err
	}
	defer session.Close()
	c.Session = session

	session.Stdout = c.StdOutput
	session.Stderr = c.StdOutput

	c.StdinPipe, err = session.StdinPipe()
	if err != nil {
		return nil
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm-256color", c.cols, c.rows, modes)
	if err != nil {
		return err
	}

	err = session.Shell()
	if err != nil {
		return err
	}

	return session.Wait()
}

// 连接
func (c *Cli) connect() error {
	conf := pkg.Conf()
	timeout := 5 * time.Second
	var config ssh.ClientConfig
	if c.Password != "" {
		config = ssh.ClientConfig{
			User:            c.Username,
			Auth:            []ssh.AuthMethod{ssh.Password(c.Password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         timeout,
		}
	} else {
		signer, err := ssh.ParsePrivateKey([]byte(c.Privatekey))
		if err != nil {
			return err
		}
		config = ssh.ClientConfig{
			User: c.Username,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         timeout,
		}
	}
	config.KeyExchanges = conf.Ssh.KeyExchanges
	config.Ciphers = conf.Ssh.Ciphers
	config.MACs = conf.Ssh.Macs

	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	sshClient, err := ssh.Dial("tcp", addr, &config)
	// sshClient, err := SSHDialTimeout("tcp", addr, &config, 60*time.Second)
	// // sshClient, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		fmt.Println("ssh connect err:", err)
		return err
	}
	c.client = sshClient

	return nil
}
func SSHDialTimeout(network, addr string, config *ssh.ClientConfig, timeout time.Duration) (*ssh.Client, error) {
	conn, err := net.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}

	timeoutConn := &Conn{conn}
	c, chans, reqs, err := ssh.NewClientConn(timeoutConn, addr, config)
	if err != nil {
		return nil, err
	}
	client := ssh.NewClient(c, chans, reqs)

	// this sends keepalive packets every 10 seconds
	// there's no useful response from these, so we can just abort if there's an error
	go func() {
		t := time.NewTicker(10 * time.Second)
		defer t.Stop()
		for range t.C {
			_, _, err := client.Conn.SendRequest("hello zenki jiuweihu", true, nil)
			if err != nil {
				return
			}
		}
	}()
	return client, nil
}

func (c *Cli) SetUUID(id string) {
	c.uuid = id
}
func (c *Cli) GetUUID() string {
	return c.uuid
}

// newSession new session
func (c *Cli) newSession() (*ssh.Session, error) {
	if c.client == nil {
		if err := c.connect(); err != nil {
			return nil, err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (c *Cli) Recv(wsConn *websocket.Conn, logBuff *bytes.Buffer, exitCh chan bool) {
	defer func() {
		fmt.Println("recv exit")
		c.Close()
	}()

	for {
		select {
		case <-exitCh:
			return
		default:
			//read websocket msg
			_, wsData, err := wsConn.ReadMessage()
			if err != nil {
				logrus.WithError(err).Error("reading webSocket message failed")
				return
			}
			var mess MessData
			if err := json.Unmarshal(wsData, &mess); err != nil {
				logrus.WithError(err).WithField("wsData", string(wsData)).Error("unmarshal websocket message failed")
				return
			}

			//			fmt.Println(mess)
			dlen := len(mess.Data)

			switch mess.Op {
			case "resize":
				//handle xterm.js size change
				if mess.Cols > 0 && mess.Rows > 0 {
					if err := c.Session.WindowChange(mess.Rows, mess.Cols); err != nil {
						logrus.WithError(err).Error("ssh pty change windows size failed")
					}
				}
			case "stdin":
				bs := []byte(mess.Data)
				// fmt.Println("============", dlen)
				// for k, v := range mess.Data {
				// 	fmt.Println(k, int(v))
				// }

				if err != nil {
					logrus.WithError(err).Error("websock cmd string base64 decoding failed")
				}
				//替换127 del 退格。换成ctrl+h，这个是ascii为8
				if c.backspace == true && dlen == 1 && int(bs[0]) == 127 {
					if _, err := c.StdinPipe.Write([]byte{8}); err != nil {
						logrus.WithError(err).Error("ws cmd bytes write to ssh.stdin pipe failed")
						//                                              setQuit(exitCh)
					}
				} else {

					if _, err := c.StdinPipe.Write(bs); err != nil {
						logrus.WithError(err).Error("ws cmd bytes write to ssh.stdin pipe failed")
						//                                              setQuit(exitCh)
					}
				}
				//write input cmd to log buffer
				if _, err := logBuff.Write(bs); err != nil {
					logrus.WithError(err).Error("write received cmd into log buffer failed")
				}
			}
		}
	}
}
func (c *Cli) Send(wsConn *websocket.Conn, exitCh chan bool) {

	defer func() {
		fmt.Println("send exit")
	}()

	//every 120ms write combine output bytes into websocket response
	tick := time.NewTicker(time.Second * time.Duration(3))
	// send ping message
	pingTick := time.NewTimer(interval)
	//for range time.Tick(120 * time.Millisecond){}
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			//write combine output bytes into websocket response
			if err := flushComboOutput(c.StdOutput, wsConn); err != nil {
				logrus.WithError(err).Error("ssh sending combo output to webSocket failed")
				return
			}
		case <-pingTick.C:
			if err := wsConn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				logrus.WithError(err).Error("ssh sending ping to webSocket failed")
				return
			}

		case <-exitCh:
			return
		}
	}
}
func (c *Cli) Wait(exitCh chan bool) {
	if err := c.Enter(); err != nil {
		fmt.Println("wait enter return err:", err)
	}
	fmt.Println("==wait final return")
	exitCh <- true

}
func (c *Cli) Close() {
	if c.Session != nil {
		c.Session.Close()
		c.client.Close()
	}
}

func flushComboOutput(w *wsBufferWriter, wsConn *websocket.Conn) error {
	if w.buffer.Len() != 0 {
		encodeString := base64.StdEncoding.EncodeToString(w.buffer.Bytes())
		err := wsConn.WriteMessage(websocket.TextMessage, []byte(encodeString))

		if err != nil {
			return err
		}
		w.buffer.Reset()
	}
	return nil
}
