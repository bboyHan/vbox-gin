package config

type MQTask struct {
	//    order-wait: false
	//    order-confirm: false
	//    order-callback: false
	//    acc-enable-check: false
	//    pay-code-cd-check: false
	//    pay-code-exp-check: false
	OrderWait       bool `mapstructure:"order-wait" json:"order-wait" yaml:"order-wait"`                         // 环境值
	OrderConfirm    bool `mapstructure:"order-confirm" json:"order-confirm" yaml:"order-confirm"`                // 环境值
	OrderCallback   bool `mapstructure:"order-callback" json:"order-callback" yaml:"order-callback"`             // 环境值
	AccEnableCheck  bool `mapstructure:"acc-enable-check" json:"acc-enable-check" yaml:"acc-enable-check"`       // 环境值
	PayCodeCdCheck  bool `mapstructure:"pay-code-cd-check" json:"pay-code-cd-check" yaml:"pay-code-cd-check"`    // 环境值
	PayCodeExpCheck bool `mapstructure:"pay-code-exp-check" json:"pay-code-exp-check" yaml:"pay-code-exp-check"` // 环境值
	AccCDCheck      bool `mapstructure:"acc-cd-check" json:"acc-cd-check" yaml:"acc-cd-check"`
}
