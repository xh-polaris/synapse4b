package conf

type Email struct {
	MaxInPeriod  int               // 周期内最多发送的次数
	Period       int               // 以分钟为单位
	Provider     string            // 供应商
	Account      string            // 账号
	Token        string            // 鉴权token
	Extra        map[string]string // 不同渠道的额外信息
	TemplateConf EmailTemplateConf `json:",default={}"` // 以AppName为名获取到到对应App的配置, 可从其中通过原因获取对应的模板
	SubjectConf  EmailSubjectConf  `json:",default={}"`
}

type EmailTemplateConf = map[string]EmailCauseToTemplate

type EmailCauseToTemplate = map[string]uint64

type EmailSubjectConf = map[string]EmailCauseToSubject

type EmailCauseToSubject = map[string]string
