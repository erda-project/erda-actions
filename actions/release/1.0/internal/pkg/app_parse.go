package pkg

import (
	"archive/zip"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/andrianbdn/iospng"
	"github.com/pkg/errors"
	"github.com/shogo82148/androidbinary/apk"
	"github.com/sirupsen/logrus"
	"howett.net/plist"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/template"
)

const (
	AndroidManifestXML = "AndroidManifest.xml"
)

var (
	reInfoPlist = regexp.MustCompile(`Payload/[^/]+/Info\.plist`)
	ErrNoIcon   = errors.New("icon not found")
)

type IOSPlist struct {
	CFBundleName         string         `plist:"CFBundleName"`
	CFBundleDisplayName  string         `plist:"CFBundleDisplayName"`
	CFBundleVersion      string         `plist:"CFBundleVersion"`
	CFBundleShortVersion string         `plist:"CFBundleShortVersionString"`
	CFBundleIdentifier   string         `plist:"CFBundleIdentifier"`
	CFBundleIcons        *CFBundleIcons `plist:"CFBundleIcons"`
}
type CFBundleIcons struct {
	CFBundlePrimaryIcon *CFBundlePrimaryIcon `plist:"CFBundlePrimaryIcon"`
}

type CFBundlePrimaryIcon struct {
	CFBundleIconFiles []string `plist:"CFBundleIconFiles"`
	CFBundleIconName  string   `plist:"CFBundleIconName"`
}

type IosAppInfo struct {
	Name     string
	BundleId string
	Version  string
	Build    string
	Icon     image.Image
	Size     int64
	IconName string
}

type AndroidAppInfo struct {
	PackageName string
	Version     string
	VersionCode int32
	Icon        image.Image
	DisplayName string
}

func GetAndroidAppInfo(appFilePath string) (*AndroidAppInfo, error) {
	info := &AndroidAppInfo{}
	pkg, err := apk.OpenFile(appFilePath)
	if err != nil {
		return nil, fmt.Errorf("error open file %s %v", appFilePath, err)
	}
	defer pkg.Close()
	info.Version = pkg.Manifest().VersionName.MustString()
	info.VersionCode = pkg.Manifest().VersionCode.MustInt32()
	info.DisplayName = pkg.Manifest().App.Name.MustString()
	info.PackageName = pkg.PackageName()
	icon, err := pkg.Icon(nil)
	if err != nil {
		logrus.Errorf("error extract icon %s %v", appFilePath, err)
	} else {
		info.Icon = icon
	}
	return info, nil
}

// TODO complete this method, extract the package name from the aab file
// Maybe you can refer to https://github.com/chenquincy/app-info-parser/issues/63
func GetAndroidAppBundleInfo(appFilePath string) (*AndroidAppInfo, error) {
	return nil, nil
}

func SaveImageToFile(icon image.Image, logoPath string) error {
	out, err := os.Create(logoPath)
	if err != nil {
		return err
	}
	var opt jpeg.Options
	opt.Quality = 80
	err = jpeg.Encode(out, icon, &opt) // put quality to 80%
	return err
}

func GetIOSAppInfo(filePath string) (*IosAppInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil, err
	}

	var plistFile *zip.File
	for _, f := range reader.File {
		if reInfoPlist.MatchString(f.Name) {
			plistFile = f
			break
		}
	}
	info, err := parseIpaFile(plistFile)
	if err != nil {
		return nil, err
	}

	var iosIconFile *zip.File
	for _, f := range reader.File {
		if strings.Contains(f.Name, info.IconName) {
			iosIconFile = f
			break
		}
	}
	icon, err := parseIpaIcon(iosIconFile)
	if err != nil {
		logrus.Errorf("failed to parse ipa icon " + filePath)
	} else {
		info.Icon = icon
	}

	info.Size = stat.Size()
	return info, nil
}

func parseIpaFile(plistFile *zip.File) (*IosAppInfo, error) {
	if plistFile == nil {
		return nil, errors.New("info.plist not found")
	}

	rc, err := plistFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	buf, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	p := new(IOSPlist)
	decoder := plist.NewDecoder(bytes.NewReader(buf))
	if err := decoder.Decode(p); err != nil {
		return nil, err
	}

	info := new(IosAppInfo)
	if p.CFBundleDisplayName == "" {
		info.Name = p.CFBundleName
	} else {
		info.Name = p.CFBundleDisplayName
	}
	info.BundleId = p.CFBundleIdentifier
	info.Version = p.CFBundleShortVersion
	info.Build = p.CFBundleVersion
	if p.CFBundleIcons != nil &&
		p.CFBundleIcons.CFBundlePrimaryIcon != nil &&
		p.CFBundleIcons.CFBundlePrimaryIcon.CFBundleIconFiles != nil &&
		len(p.CFBundleIcons.CFBundlePrimaryIcon.CFBundleIconFiles) > 0 {
		files := p.CFBundleIcons.CFBundlePrimaryIcon.CFBundleIconFiles
		info.IconName = files[len(files)-1]
	} else {
		info.IconName = "Icon.png"
	}
	return info, nil
}

func parseIpaIcon(iconFile *zip.File) (image.Image, error) {
	if iconFile == nil {
		return nil, ErrNoIcon
	}

	rc, err := iconFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var w bytes.Buffer
	iospng.PngRevertOptimization(rc, &w)

	return png.Decode(bytes.NewReader(w.Bytes()))
}

func GenerateInstallPlist(info *IosAppInfo, downloadUrl string) string {
	plistTemplate := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
   <key>items</key>
   <array>
       <dict>
           <key>assets</key>
           <array>
               <dict>
                   <key>kind</key>
                   <string>software-package</string>
                   <key>url</key>
                   <string>{{appUrl}}</string>
               </dict>
           </array>
           <key>metadata</key>
           <dict>
               <key>bundle-identifier</key>
               <string>{{bundleId}}</string>
               <key>bundle-version</key>
               <string>{{version}}</string>
               <key>kind</key>
               <string>software</string>
               <key>subtitle</key>
               <string>{{displayName}}</string>
               <key>title</key>
               <string>{{displayName}}</string>
           </dict>
       </dict>
   </array>
</dict>
</plist>`
	plistContent := template.Render(plistTemplate, map[string]string{
		"bundleId":    info.BundleId,
		"version":     info.Version,
		"displayName": info.Name,
		"appUrl":      downloadUrl,
	})
	return plistContent
}

// getAppStoreURL 根据bundleID从app store搜索链接，目前只从中国区查找
func getAppStoreURL(bundleID string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	hc := &http.Client{
		Transport: tr,
	}
	var getResp apistructs.AppStoreResponse
	resp, err := hc.Get("https://itunes.apple.com/cn/lookup?bundleId=" + bundleID)
	if err != nil {
		logrus.Errorf("get app store url err: %v", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("get app store url err: %v", err)
		return ""
	}
	err = json.Unmarshal(body, &getResp)
	if err != nil {
		logrus.Errorf("get app store url err: %v", err)
		return ""
	}

	fmt.Println(getResp)

	if getResp.ResultCount == 0 {
		return ""
	}

	return getResp.Results[0].TrackViewURL
}
