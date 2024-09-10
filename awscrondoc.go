package awscrondoc

import (
	"bytes"
	"context"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	etypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/winebarrel/cronplan"
)

type AwsCronDoc struct {
	eb *eventbridge.Client
	g  *glue.Client
}

func New() (*AwsCronDoc, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	eb := eventbridge.NewFromConfig(cfg)
	g := glue.NewFromConfig(cfg)
	return &AwsCronDoc{
		eb: eb,
		g:  g,
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
	t := template.Must(template.New("eventbridge").Funcs(funcMap).Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, rules); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (a *AwsCronDoc) listRules() ([]*etypes.Rule, error) {
	var nextToken *string
	rules := make([]*etypes.Rule, 0)
	for {
		output, err := a.eb.ListRules(context.TODO(),
			&eventbridge.ListRulesInput{
				NextToken: nextToken,
			})
		if err != nil {
			return nil, err
		}
		for _, r := range output.Rules {
			rules = append(rules, &r)
		}
		if output.NextToken == nil {
			break
		}
		*nextToken = *output.NextToken

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
