package main

type Sonar struct {
	Auth SonarAuth
}

type SonarAuth struct {
	HostURL  string
	Login    string
	Password string
}

func NewSonar(hostURL, login, password string) *Sonar {
	auth := SonarAuth{HostURL: hostURL, Login: login, Password: password}
	sonar := Sonar{Auth: auth}
	return &sonar
}
