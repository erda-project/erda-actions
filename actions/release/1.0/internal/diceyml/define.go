package diceyml

// release action 会对 dice.yml 进行反序列化之后再重新序列化，
// 如果对 dice.yml 的理解不一致，就会丢失信息，
// release action，可以不用知道 dice.yml 中的具体内容，比如 Services 中具体有哪些字段，
// 这样 dice.yml 扩展之后，action 也无需重新编译了。
// 所以这里没有直接用 dice parser 库里的解析方式，而是重写了一个。

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type EnvType int

const (
	BaseEnv EnvType = iota
	DevEnv
	TestEnv
	StagingEnv
	ProdEnv
)

func (e EnvType) String() string {
	switch e {
	case BaseEnv:
		return "base"
	case DevEnv:
		return "development"
	case TestEnv:
		return "test"
	case StagingEnv:
		return "staging"
	case ProdEnv:
		return "production"
	default:
		panic("should not be here!")
	}
}

type Object map[string]interface{}

type DiceYaml struct {
	data []byte
	obj  *Object
}

func New(b []byte) (*DiceYaml, error) {
	var obj Object
	d := &DiceYaml{
		data: b,
	}
	if err := yaml.Unmarshal(d.data, &obj); err != nil {
		return nil, errors.Wrap(err, "fail to yaml unmarshal")
	}
	d.obj = &obj
	return d, nil
}

func override(src, dst interface{}) {
	r, err := yaml.Marshal(src)
	if err != nil {
		panic("deepcopy marshal")
	}
	if err := yaml.Unmarshal(r, dst); err != nil {
		panic("deepcopy unmarshal")
	}
}

func (d *DiceYaml) Copy() DiceYaml {
	dst := new(Object)
	override(d.obj, dst)
	return DiceYaml{
		obj: dst,
	}
}
func (d *DiceYaml) Obj() *Object {
	return d.obj
}

func (d *DiceYaml) Meta() map[string]string {
	if d.obj == nil {
		return nil
	}
	obj := *d.obj
	var (
		metaI  interface{}
		metaII map[interface{}]interface{}
		ok     bool
	)
	meta := map[string]string{}
	if metaI, ok = obj["meta"]; !ok {
		return meta
	}
	if metaII, ok = metaI.(map[interface{}]interface{}); !ok {
		if meta, ok = metaI.(map[string]string); ok {
			return meta
		}
		return nil
	}
	for keyI, valueI := range metaII {
		var key, value string
		if key, ok = keyI.(string); !ok {
			return nil
		}
		if value, ok = valueI.(string); !ok {
			return nil
		}
		meta[key] = value
	}
	return meta
}

func (d *DiceYaml) Envs() map[string]string {
	if d.obj == nil {
		return nil
	}
	obj := *d.obj
	var (
		envsI  interface{}
		envsII map[interface{}]interface{}
		ok     bool
	)
	envs := map[string]string{}
	if envsI, ok = obj["envs"]; !ok {
		return envs
	}
	if envsII, ok = envsI.(map[interface{}]interface{}); !ok {
		if envs, ok = envsI.(map[string]string); ok {
			return envs
		}
		return nil
	}
	for keyI, valueI := range envsII {
		var key, value string
		if key, ok = keyI.(string); !ok {
			return nil
		}
		if value, ok = valueI.(string); !ok {
			return nil
		}
		envs[key] = value
	}
	return envs
}

func (d *DiceYaml) extractSpecialMap(key string) map[string]map[string]interface{} {
	if d.obj == nil {
		return nil
	}
	obj := *d.obj
	var (
		itemsI  interface{}
		itemsII map[interface{}]interface{}
		ok      bool
	)
	items := map[string]map[string]interface{}{}
	if itemsI, ok = obj[key]; !ok {
		return items
	}
	if itemsII, ok = itemsI.(map[interface{}]interface{}); !ok {
		if items, ok = itemsI.(map[string]map[string]interface{}); ok {
			return items
		}
		return nil
	}
	for nameI, itemI := range itemsII {
		var (
			name   string
			itemII map[interface{}]interface{}
		)
		item := map[string]interface{}{}
		if name, ok = nameI.(string); !ok {
			return nil
		}
		if itemII, ok = itemI.(map[interface{}]interface{}); !ok {
			return nil
		}
		for keyI, valueI := range itemII {
			var key string
			if key, ok = keyI.(string); !ok {
				return nil
			}
			item[key] = valueI
		}
		items[name] = item
	}
	return items
}

func (d *DiceYaml) Services() map[string]map[string]interface{} {
	return d.extractSpecialMap("services")
}

func (d *DiceYaml) Jobs() map[string]map[string]interface{} {
	return d.extractSpecialMap("jobs")
}

func (d *DiceYaml) Addons() map[string]map[string]interface{} {
	return d.extractSpecialMap("addons")
}

func (d *DiceYaml) Environments() map[string]map[string]interface{} {
	return d.extractSpecialMap("environments")
}

func (d *DiceYaml) SetServices(services map[string]map[string]interface{}) {
	if services != nil {
		(*d.obj)["services"] = services
	}
}

func (d *DiceYaml) SetJobs(jobs map[string]map[string]interface{}) {
	if jobs != nil {
		(*d.obj)["jobs"] = jobs
	}
}

func (d *DiceYaml) SetAddons(addons map[string]map[string]interface{}) {
	if addons != nil {
		(*d.obj)["addons"] = addons
	}
}

func (d *DiceYaml) SetEnvironments(environments map[string]map[string]interface{}) {
	if environments != nil {
		(*d.obj)["environments"] = environments
	}
}

func (d *DiceYaml) SetMeta(meta map[string]string) {
	if meta != nil {
		(*d.obj)["meta"] = meta
	}
}

func (d *DiceYaml) SetEnvs(envs map[string]string) {
	if envs != nil {
		(*d.obj)["envs"] = envs
	}
}

func (d *DiceYaml) SetEnv(key, value string) error {
	if d.obj == nil {
		return errors.New("modify dice.yml base on raw bytes is not allowed")
	}
	envs := d.Envs()
	if envs == nil {
		envs = map[string]string{}
		(*d.obj)["envs"] = envs
	}
	envs[key] = value
	d.SetEnvs(envs)
	return nil
}

func (d *DiceYaml) InsertImage(images map[string]string) error {
	if d.obj == nil {
		return errors.New("modify dice.yml base on raw bytes is not allowed")
	}
	services := d.Services()
	jobs := d.Jobs()
	for name, image := range images {
		if services != nil {
			if service, ok := services[name]; ok {
				service["image"] = image
			}
		}
		if jobs != nil {
			if job, ok := jobs[name]; ok {
				job["image"] = image
			}
		}
		delete(images, name)
	}
	d.SetJobs(jobs)
	d.SetServices(services)
	return nil
}

func GetAddonPlan(addon map[string]interface{}) string {
	var (
		planI interface{}
		plan  string
		ok    bool
	)

	if planI, ok = addon["plan"]; !ok {
		return ""
	}
	if plan, ok = planI.(string); !ok {
		return ""
	}
	return plan
}

func insertAddonOptions(addons map[string]map[string]interface{}, addonPlan string, moreOptions map[string]string) {
	for _, addon := range addons {
		plan := GetAddonPlan(addon)
		splitted := strings.Split(plan, ":")
		if len(splitted) < 1 {
			continue
		}
		if strings.TrimSpace(splitted[0]) == addonPlan {
			options := map[string]string{}
			if optionsI, ok := addon["options"]; ok {
				if existOptions, ok := optionsI.(map[interface{}]interface{}); ok {
					for keyI, valueI := range existOptions {
						if key, ok := keyI.(string); ok {
							if value, ok := valueI.(string); ok {
								options[key] = value
							}
						}
					}
				}
				if existOptions, ok := optionsI.(map[string]string); ok {
					options = existOptions
				}
				for key, value := range moreOptions {
					options[key] = value
				}
			}
			addon["options"] = options
		}
	}
}

func (d *DiceYaml) InsertAddonOptions(env EnvType, addonPlan string, moreOptions map[string]string) error {
	if d.obj == nil {
		return errors.New("modify dice.yml base on raw bytes is not allowed")
	}
	addons := d.Addons()
	environments := d.Environments()
	if addons != nil {
		insertAddonOptions(addons, addonPlan, moreOptions)
		d.SetAddons(addons)
	}
	if environments != nil {
		var (
			subObject Object
			ok        bool
		)
		if subObject, ok = environments[env.String()]; !ok {
			return nil
		}
		subDiceYml := DiceYaml{
			obj: &subObject,
		}
		addons = subDiceYml.Addons()
		if addons != nil {
			insertAddonOptions(addons, addonPlan, moreOptions)
			subDiceYml.SetAddons(addons)
		}
		d.SetEnvironments(environments)
	}
	return nil
}

func (d *DiceYaml) Compose(env string, yml *DiceYaml) error {
	if d.obj == nil {
		return errors.New("modify dice.yml base on raw bytes is not allowed")
	}
	environments := d.Environments()
	if environments == nil {
		environments = map[string]map[string]interface{}{}
		(*d.obj)["environments"] = environments
	}
	if e, ok := environments[env]; !ok || len(e) == 0 {
		environments[env] = map[string]interface{}{}
	}
	envs := yml.Envs()
	if envs != nil && len(envs) > 0 {
		environments[env]["envs"] = envs
	}
	services := yml.Services()
	if services != nil && len(services) > 0 {
		environments[env]["services"] = services
	}
	addons := yml.Addons()
	if addons != nil && len(addons) > 0 {
		environments[env]["addons"] = addons
	}
	d.SetEnvironments(environments)
	return nil
}
