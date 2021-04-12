// Package out dicehub action
package pkg

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/erda-project/erda-actions/actions/release/1.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/release/1.0/internal/diceyml"
	"github.com/erda-project/erda/apistructs"
)

func fillDiceYml(cfg *conf.Conf, storage *StorageURL) (string, error) {
	d, err := composeEnvYml(cfg)
	if err != nil {
		return "", err
	}

	if err = insertImages(d, cfg); err != nil {
		return "", err
	}

	// push db to oss/disk
	if cfg.InitSQL != "" {
		if err = executeSQL(cfg, storage, d); err != nil {
			return "", err
		}
	}
	// 判断是否存在serviceMesh addon，存在的话meta中添加service-mesh:on, service 中增加对应开启的开关
	dCopy := d.Copy()
	err = serviceMeshAddonAdjust(&dCopy)
	if err != nil {
		return "", err
	}

	err = apiGatewayAddonAdjust(&dCopy)
	if err != nil {
		return "", err
	}

	//假如是service模式，需要将service中的cmd塞入dice.yaml中的cmd中
	if cfg.Services != nil {
		insertCommands(&dCopy, cfg)
	}

	if err != nil {
		return "", err
	}
	b, err := yaml.Marshal(dCopy.Obj())
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func apiGatewayAddonAdjust(dice *diceyml.DiceYaml) error {
	addons := dice.Addons()
	if addons != nil {
		for _, addon := range addons {
			plan := diceyml.GetAddonPlan(addon)
			planSplit := strings.Split(plan, ":")
			if len(planSplit) == 0 {
				return errors.Errorf("plan about addon was error")
			}
			if planSplit[0] == apistructs.AddonApiGateway {
				return nil
			}
		}
	}
	// addon not exist
	services := dice.Services()
	for _, service := range services {
		if objI, ok := service["endpoints"]; ok {
			if obj, ok := objI.([]interface{}); ok && len(obj) > 0 {
				apiGatewayAddonInsert(dice)
				return nil
			}
		}
	}
	return nil
}

func apiGatewayAddonInsert(dice *diceyml.DiceYaml) {
	addons := dice.Addons()
	if addons != nil {
		addons["api-gateway"] = map[string]interface{}{
			"plan": apistructs.AddonApiGateway + ":basic",
		}
	}
	dice.SetAddons(addons)
}

func serviceMeshAdjust(dice *diceyml.DiceYaml) {
	meta := dice.Meta()
	meta[apistructs.AddonServiceMesh] = "on"
	services := dice.Services()
	for name, _ := range services {
		if _, ok := services[name]["mesh_enable"]; !ok {
			services[name]["mesh_enable"] = &[]bool{true}[0]
		}
	}
	dice.SetServices(services)
	dice.SetMeta(meta)
}

func serviceMeshAddonInsert(dice *diceyml.DiceYaml) {
	addons := dice.Addons()
	if addons != nil {
		addons["service-mesh"] = map[string]interface{}{
			"plan": apistructs.AddonServiceMesh + ":basic",
		}
	}
	dice.SetAddons(addons)
}

// serviceMeshAddonAdjust addon中如果存在service-mesh，meta中加上service-mesh:on
func serviceMeshAddonAdjust(dice *diceyml.DiceYaml) error {
	addons := dice.Addons()
	if addons != nil {
		for _, addon := range addons {
			plan := diceyml.GetAddonPlan(addon)
			planSplit := strings.Split(plan, ":")
			if len(planSplit) == 0 {
				return errors.Errorf("plan about addon was error")
			}
			if planSplit[0] == apistructs.AddonServiceMesh {
				serviceMeshAdjust(dice)
				return nil
			}
		}
	}
	// addon not exist
	services := dice.Services()
	if services != nil {
		for _, service := range services {
			if objI, ok := service["traffic_security"]; ok {
				if objII, ok := objI.(map[interface{}]interface{}); ok {
					if objII["mode"] != "" {
						serviceMeshAddonInsert(dice)
						serviceMeshAdjust(dice)
						return nil
					}
				}
			}
		}
	}
	return nil
}

func parseURL(fullURL string) (*StorageURL, error) {
	var (
		URL     *url.URL
		err     error
		storage = &StorageURL{}
	)

	if fullURL == "" {
		return nil, errors.New("nil url")
	}

	if URL, err = url.Parse(fullURL); err != nil {
		return nil, errors.Errorf("failed to parse url, url: %s, (%+v)", fullURL, err)
	}

	storage.Scheme = URL.Scheme
	storage.Path = URL.Path
	storage.UserName = URL.User.Username()
	storage.PassWord, _ = URL.User.Password()
	storage.Host = URL.Hostname()

	return storage, nil
}
