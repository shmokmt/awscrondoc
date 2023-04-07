package awscrondoc

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/winebarrel/cronplan"
)

type AwsCronDoc struct {
	svc *eventbridge.EventBridge
}

func New() (*AwsCronDoc, error) {
	sessOpts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	sess, err := session.NewSessionWithOptions(sessOpts)
	if err != nil {
		return nil, err
	}
	svc := eventbridge.New(sess)
	return &AwsCronDoc{
		svc: svc,
	}, nil
}

func (a *AwsCronDoc) MarkdownString() (string, error) {
	rules, err := a.listRules()
	if err != nil {
		return "", err
	}
	funcMap := template.FuncMap{
		"isCronExpression": isCronExpression,
		"latestSchedules":  latestSchedules,
	}
	t := template.Must(template.New("doc").Funcs(funcMap).Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, rules); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (a *AwsCronDoc) listRules() ([]*eventbridge.Rule, error) {
	var nextToken *string
	rules := make([]*eventbridge.Rule, 0)
	for {
		resp, err := a.svc.ListRules(&eventbridge.ListRulesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		for _, r := range resp.Rules {
			rules = append(rules, r)
		}
		if resp.NextToken == nil {
			break
		}
		*nextToken = *resp.NextToken

	}
	return rules, nil
}

const tmpl = `
{{ range $i, $r := . }}
	{{- if and (ne $r.ScheduleExpression nil) (isCronExpression $r.ScheduleExpression) }}
## {{ $r.Name }}

{{ if ne $r.Description nil }}* Description: {{ $r.Description }}{{ end }}
{{- if ne $r.ScheduleExpression nil }}
* CronExperssion: {{ $r.ScheduleExpression }}
* Example:
  {{- range $i, $t := $r.ScheduleExpression | latestSchedules }} 
  * {{ $t }}
  {{- end }}
* State: {{ $r.State }}
{{- end }}

	{{- end }}
{{ end }}
`

func latestSchedules(exp string) []time.Time {
	cp, _ := cronplan.Parse(trimCronBracket(exp))
	nextN := cp.NextN(time.Now().UTC(), 4)
	for i, t := range nextN {
		localTime := t.Local()
		nextN[i] = localTime
	}
	return nextN
}

func trimCronBracket(exp string) string {
	exp = strings.Replace(exp, "cron(", "", 1)
	exp = strings.Replace(exp, ")", "", 1)
	return exp
}

func isCronExpression(exp string) bool {
	return strings.HasPrefix(exp, "cron")
}
