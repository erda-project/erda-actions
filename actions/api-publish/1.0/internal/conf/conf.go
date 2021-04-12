package conf

// Conf action 入参
type Conf struct {
	// Params

	RegisterId string `env:"ACTION_REGISTER_ID"`
	ItemName   string `env:"ACTION_ITEM_NAME"`
	OwnerEmail string `env:"ACTION_OWNER_EMAIL"`

	// env
	OrgID             int64  `env:"DICE_ORG_ID" required:"true"`
	CiOpenapiToken    string `env:"DICE_OPENAPI_TOKEN" required:"true"`
	DiceOpenapiPrefix string `env:"DICE_OPENAPI_ADDR" required:"true"`
}
