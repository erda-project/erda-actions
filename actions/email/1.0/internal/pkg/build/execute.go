package build

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/smtp"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/erda-project/erda-actions/actions/email/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

func Execute() error {
	fmt.Println("start send email")
	defer fmt.Println("end send email")
	var cfg conf.Conf
	envconf.MustLoad(&cfg)

	// customize the contrast structure
	var traces = []map[string]interface{}{}
	err := json.Unmarshal([]byte(cfg.EmailTemplateObject), &traces)
	if err == nil {
		var newTraces = []map[string]interface{}{}
		for _, trace := range traces {
			trace["applicationName"] = os.Getenv("DICE_APPLICATION_NAME")
			if trace["discovered"] != nil {
				discoveredInt, ok := trace["discovered"].(float64)
				if ok {
					tm := time.Unix(int64(discoveredInt/1000), 0)
					trace["discoveredTime"] = tm.Format("2006-01-02 15:04:05")
				}
			}
			newTraces = append(newTraces, trace)
		}
		objects, err := json.Marshal(newTraces)
		if err != nil {
			return fmt.Errorf("json Marshal object error %v", err)
		}
		cfg.EmailTemplateObject = string(objects)
	}

	templateData, err := ioutil.ReadFile(cfg.EmailTemplateAddr)
	if err != nil {
		return fmt.Errorf("template file read failed %v", err)
	}

	tmpl, err := template.New("test").Parse(string(templateData))
	if err != nil {
		return fmt.Errorf("template render error %v", err)
	}

	var objects interface{}
	err = json.Unmarshal([]byte(cfg.EmailTemplateObject), &objects)
	if err != nil {
		return fmt.Errorf("json Unmarshal email_template_object error %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"email_template_object": objects,

		"HOST":                  os.Getenv("HOST"),
		"PIPELINE_ID":           os.Getenv("PIPELINE_ID"),
		"DICE_ORG_ID":           cfg.OrgID,
		"DICE_ORG_NAME":         cfg.OrgName,
		"DICE_PROJECT_ID":       cfg.ProjectID,
		"DICE_PROJECT_NAME":     cfg.ProjectName,
		"DICE_APPLICATION_ID":   cfg.AppID,
		"DICE_APPLICATION_NAME": cfg.AppName,
		"PIPELINE_TASK_NAME":    cfg.TaskName,
		"DICE_CLUSTER_NAME":     cfg.ClusterName,
		"DICE_VERSION":          cfg.DiceVersion,
	})
	if err != nil {
		return fmt.Errorf("render template error %v", err)
	}

	var host, port, email, password string
	var isSSL = true
	data, err := erdaSmtpInfo()
	if err != nil {
		fmt.Printf("get erda smtp info error %v", err)
	} else {
		host = data.Host
		port = data.Port
		email = data.User
		password = data.Password
		isSSL = data.IsSSL
	}

	if os.Getenv("SMTP_HOST") != "" {
		host = os.Getenv("SMTP_HOST")
	}

	if os.Getenv("SMTP_PORT") != "" {
		port = os.Getenv("SMTP_PORT")
	}

	if os.Getenv("SMTP_EMAIL") != "" {
		email = os.Getenv("SMTP_EMAIL")
	}

	if os.Getenv("SMTP_PASSWORD") != "" {
		password = os.Getenv("SMTP_PASSWORD")
	}

	if os.Getenv("SMTP_NOT_SSL") != "" {
		isSSL = false
	}

	fmt.Printf("smtp_host %v, smtp_port %v, smtp_email %v \n", host, port, email)

	header := make(map[string]string)
	header["From"] = "test" + "<" + email + ">"
	header["To"] = cfg.ToMail[0]
	header["Subject"] = "dice get safety production inspection report"
	header["Content-Type"] = "text/html; charset=UTF-8"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + buf.String()

	auth := smtp.PlainAuth(
		"",
		email,
		password,
		host,
	)
	if os.Getenv("SMTP_GSSAPI") != "" {
		os.Setenv("python_smtp_host", host)
		os.Setenv("python_smtp_sender", email)
		os.Setenv("python_smtp_receivers", strings.Join(cfg.ToMail, ","))
		os.Setenv("python_smtp_message", message)

		err = simpleRun("python2.7", "/app/files/sendmail.py")
		if err != nil {
			return err
		}
	} else {
		if isSSL {
			err = SendMailUsingTLS(fmt.Sprintf("%s:%s", host, port), auth, email, cfg.ToMail, true, []byte(message))
		} else {
			err = smtp.SendMail(fmt.Sprintf("%s:%s", host, port), auth, email, cfg.ToMail, []byte(message))
		}
	}

	if err != nil {
		return fmt.Errorf("error to send email, error %v", err)
	}
	return nil
}

func simpleRun(name string, arg ...string) error {
	fmt.Fprintf(os.Stdout, "Run: %s, %v\n", name, arg)
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type MailSubscriberInfo struct {
	Host               string
	Port               string
	User               string
	Password           string
	DisplayUser        string
	IsSSL              bool
	IsSSLStr           string
	InsecureSkipVerify bool
}

func erdaSmtpInfo() (*MailSubscriberInfo, error) {

	var resp MailSubscriberInfo
	response, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(os.Getenv("DICE_OPENAPI_PUBLIC_URL")).
		Path("/api/dice/eventbox/actions/get-smtp-info").
		Header("Authorization", os.Getenv("DICE_OPENAPI_TOKEN")).Do().JSON(&resp)

	if err != nil {
		return nil, fmt.Errorf("failed to request (%s)", err.Error())
	}

	if !response.IsOK() {
		return nil, fmt.Errorf(fmt.Sprintf("failed to request, status-code: %d, content-type: %s", response.StatusCode(), response.ResponseHeader("Content-Type")))
	}

	return &resp, nil
}

func DialTLS(addr string, insecureSkipVerify bool) (*smtp.Client, error) {
	tlsconfig := &tls.Config{}
	if insecureSkipVerify {
		tlsconfig.InsecureSkipVerify = true
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func SendMailUsingTLS(addr string, auth smtp.Auth, from string, tos []string, insecureSkipVerify bool, msg []byte) (err error) {
	c, err := DialTLS(addr, insecureSkipVerify)
	if err != nil {
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range tos {
		if err = c.Rcpt(addr); err != nil {
			fmt.Print(err)
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
