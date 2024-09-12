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
	rules, err := a.listEventBridgeRules()
	if err != nil {
		return "", err
	}
	eventbridgeTemplate := template.Must(template.New("eventbridge").Funcs(template.FuncMap{
		"isCronExpression": isCronExpression,
		"latestSchedules":  latestSchedules,
	}).Parse(eventBridgeTmpl))
	var buf bytes.Buffer
	if err := eventbridgeTemplate.Execute(&buf, rules); err != nil {
		return "", err
	}
	glueTriggerSchedules, err := a.listGlueTrigger()
	if err != nil {
		return "", err
	}
	glueTriggerTemplate := template.Must(template.New("glueTrigger").Funcs(template.FuncMap{
		"isCronExpression": isCronExpression,
		"latestSchedules":  latestSchedules,
	}).Parse(glueTriggerTmpl))
	var glueTriggerBuf bytes.Buffer
	if err := glueTriggerTemplate.Execute(&glueTriggerBuf, glueTriggerSchedules); err != nil {
		return "", err
	}
	buf.Write(glueTriggerBuf.Bytes())
	return buf.String(), nil
}

func (a *AwsCronDoc) listEventBridgeRules() ([]*etypes.Rule, error) {
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

const eventBridgeTmpl = `
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

// TODO: Comment in when support for Glue Crawler is added.
// type GlueCrawlerSchedule struct {
// 	Name               *string
// 	Description        *string
// 	ScheduleExpression *string
// 	ScheduleState      string
// }

type GlueTriggerSchedule struct {
	Name               *string
	WorkflowName       *string
	ScheduleExpression *string
	ScheduleState      string
}

func (a *AwsCronDoc) listGlueTrigger() ([]*GlueTriggerSchedule, error) {
	var nextToken *string
	triggerNames := make([]string, 0)
	for {
		output, err := a.g.ListTriggers(context.TODO(),
			&glue.ListTriggersInput{
				NextToken: nextToken,
			})
		if err != nil {
			return nil, err
		}
		triggerNames = append(triggerNames, output.TriggerNames...)
		if output.NextToken == nil {
			break
		}
		*nextToken = *output.NextToken
	}
	glueTriggerSchedules := make([]*GlueTriggerSchedule, 0)
	for _, triggerName := range triggerNames {
		output, err := a.g.GetTrigger(context.TODO(), &glue.GetTriggerInput{
			Name: &triggerName,
		})
		if err != nil {
			return nil, err
		}
		glueTriggerSchedules = append(glueTriggerSchedules, &GlueTriggerSchedule{
			Name:               output.Trigger.Name,
			ScheduleExpression: output.Trigger.Schedule,
			ScheduleState:      string(output.Trigger.State),
		})
	}
	return glueTriggerSchedules, nil
}

// TODO: Comment in when support for Glue Crawler is added.
// func (a *AwsCronDoc) GetCrawlers() ([]*GlueCrawlerSchedule, error) {
// 	var nextToken *string
// 	crawlers := make([]types.Crawler, 0)
// 	for {
// 		output, err := a.g.GetCrawlers(context.TODO(),
// 			&glue.GetCrawlersInput{
// 				NextToken: nextToken,
// 			})
// 		if err != nil {
// 			return nil, err
// 		}
// 		if output.NextToken == nil {
// 			break
// 		}
// 		*nextToken = *output.NextToken
// 		crawlers = append(crawlers, output.Crawlers...)
// 	}
// 	glueCrawlerSchedules := make([]*GlueCrawlerSchedule, 0)
// 	for _, crawler := range crawlers {
// 		glueCrawlerSchedules = append(glueCrawlerSchedules, &GlueCrawlerSchedule{
// 			Name:               crawler.Name,
// 			Description:        crawler.Description,
// 			ScheduleExpression: crawler.Schedule.ScheduleExpression,
// 			ScheduleState:      string(crawler.State),
// 		})

// 	}
// 	return glueCrawlerSchedules, nil
// }

const glueTriggerTmpl = `
# Glue Trigger Schedules
{{ range $i, $r := . }}
	{{- if and (ne $r.ScheduleExpression nil) (isCronExpression $r.ScheduleExpression) }}
{{- if ne $r.Name nil}}## {{ $r.Name }}{{ end}}

{{ if ne $r.WorkflowName nil }}* WorkflowName: {{ $r.WorkflowName }}{{ end }}
{{- if ne $r.ScheduleExpression nil }}
* CronExperssion: {{ $r.ScheduleExpression }}
* Example:
  {{- range $i, $t := $r.ScheduleExpression | latestSchedules }}
  * {{ $t }}
  {{- end }}
* State: {{ $r.ScheduleState }}
{{- end }}

	{{- end }}
{{ end }}
`

// TODO: Comment in when support for Glue Crawler is added.
// const glueCrawlerTmpl = `
// # Glue Crawler Schedules
// {{ range $i, $r := . }}
// 	{{- if and (ne $r.ScheduleExpression nil) (isCronExpression $r.ScheduleExpression) }}
// {{- if ne $r.Name nil}}## $r.Name }}{{ end}}

// {{ if ne $r.Description nil }}* Description: {{ $r.Description }}{{ end }}
// {{- if ne $r.ScheduleExpression nil }}
// * CronExperssion: {{ $r.ScheduleExpression }}
// * Example:
//   {{- range $i, $t := $r.ScheduleExpression | latestSchedules }}
//   * {{ $t }}
//   {{- end }}
// * State: {{ $r.ScheduleState }}
// {{- end }}

// 	{{- end }}
// {{ end }}
// `

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
