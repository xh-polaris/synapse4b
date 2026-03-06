package conf

type Email struct {
	MaxInPeriod  int               // 周期内最多发送的次数
	Period       int               // 以分钟为单位
	Provider     string            // 供应商
	Account      string            // 账号
	Token        string            // 鉴权token
	FromAddress  string            // 发件人地址
	Extra        map[string]string // 不同渠道的额外信息，包括发件人名称和地址
	TemplateConf EmailTemplateConf `json:",default={}"` // 以AppName为名获取到到对应App的配置, 可从其中通过原因获取对应的邮件模板
	SubjectConf  EmailSubjectConf  `json:",default={}"` // 同上，获取对应邮件主题
}

type EmailTemplateConf = map[string]EmailCauseToTemplate

type EmailCauseToTemplate = map[string]uint64

type EmailSubjectConf = map[string]EmailCauseToSubject

type EmailCauseToSubject = map[string]string
