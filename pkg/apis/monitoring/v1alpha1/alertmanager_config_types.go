// Copyright 2020 The prometheus-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"errors"
	"regexp"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	Version = "v1alpha1"

	AlertmanagerConfigKind    = "AlertmanagerConfig"
	AlertmanagerConfigName    = "alertmanagerconfigs"
	AlertmanagerConfigKindKey = "alertmanagerconfig"
)

var (
	opsGenieTypeRe = regexp.MustCompile("^(team|user|escalation|schedule)$")
)

// AlertmanagerConfig defines a namespaced AlertmanagerConfig to be aggregated across multiple namespaces configuring one Alertmanager.
// +genclient
// +k8s:openapi-gen=true
type AlertmanagerConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AlertmanagerConfigSpec `json:"spec"`
}

// AlertmanagerConfigList is a list of AlertmanagerConfig.
// +k8s:openapi-gen=true
type AlertmanagerConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`
	// List of AlertmanagerConfig
	Items []*AlertmanagerConfig `json:"items"`
}

type AlertmanagerConfigSpec struct {
	Route        *Route        `json:"route,omitempty"`
	Receivers    []Receiver    `json:"receivers,omitempty"`
	InhibitRules []InhibitRule `json:"inhibitRules,omitempty"`
}

type Route struct {
	Receiver       string    `json:"receiver,omitempty"`
	GroupBy        []string  `json:"groupBy,omitempty"`
	GroupWait      string    `json:"groupWait,omitempty"`
	GroupInterval  string    `json:"groupInterval,omitempty"`
	RepeatInterval string    `json:"repeatInterval,omitempty"`
	Matchers       []Matcher `json:"matchers,omitempty"`
	Continue       bool      `json:"continue,omitempty"`
	Routes         []Route   `json:"routes,omitempty"`
}

type Receiver struct {
	Name             string            `json:"name"`
	OpsGenieConfigs  []OpsGenieConfig  `json:"opsgenieConfigs,omitempty"`
	PagerDutyConfigs []PagerDutyConfig `json:"pagerDutyConfigs,omitempty"`
	SlackConfigs     []SlackConfig     `json:"slackConfigs,omitempty"`
	WebhookConfigs   []WebhookConfig   `json:"webhookConfigs,omitempty"`
	WeChatConfigs    []WeChatConfig    `json:"weChatConfigs,omitempty"`
}

type PagerDutyConfig struct {
	SendResolved *bool                   `json:"sendResolved,omitempty"`
	RoutingKey   *v1.SecretKeySelector   `json:"routingKey,omitempty"`
	ServiceKey   *v1.SecretKeySelector   `json:"serviceKey,omitempty"`
	URL          *string                 `json:"url,omitempty"`
	Client       *string                 `json:"client,omitempty"`
	ClientURL    *string                 `json:"clientURL,omitempty"`
	Description  *string                 `json:"description,omitempty"`
	Severity     *string                 `json:"severity,omitempty"`
	Class        *string                 `json:"class,omitempty"`
	Group        *string                 `json:"group,omitempty"`
	Component    *string                 `json:"component,omitempty"`
	Details      []PagerDutyConfigDetail `json:"details,omitempty"`
	HTTPConfig   *HTTPConfig             `json:"httpConfig,omitempty"`
}

// SlackConfig configures notifications via Slack.
type SlackConfig struct {
	SendResolved *bool                 `json:"sendResolved,omitempty"`
	APIURL       *v1.SecretKeySelector `json:"apiURL,omitempty"`
	Channel      *string               `json:"channel,omitempty"`
	Username     *string               `json:"username,omitempty"`
	Color        *string               `json:"color,omitempty"`
	Title        *string               `json:"title,omitempty"`
	TitleLink    *string               `json:"titleLink,omitempty"`
	Pretext      *string               `json:"pretext,omitempty"`
	Text         *string               `json:"text,omitempty"`
	Fields       []SlackField          `json:"fields,omitempty"`
	ShortFields  *bool                 `json:"shortFields,omitempty"`
	Footer       *string               `json:"footer,omitempty"`
	Fallback     *string               `json:"fallback,omitempty"`
	CallbackID   *string               `json:"callbackId,omitempty"`
	IconEmoji    *string               `json:"iconEmoji,omitempty"`
	IconURL      *string               `json:"iconURL,omitempty"`
	ImageURL     *string               `json:"imageURL,omitempty"`
	ThumbURL     *string               `json:"thumbURL,omitempty"`
	LinkNames    *bool                 `json:"linkNames,omitempty"`
	MrkdwnIn     []string              `json:"mrkdwnIn,omitempty"`
	Actions      []SlackAction         `json:"actions,omitempty"`
	HTTPConfig   *HTTPConfig           `json:"httpConfig,omitempty"`
}

// Validate ensures SlackConfig is valid
func (sc *SlackConfig) Validate() error {
	for _, action := range sc.Actions {
		if err := action.Validate(); err != nil {
			return err
		}
	}
	for _, field := range sc.Fields {
		if err := field.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// SlackAction configures a single Slack action that is sent with each notification.
// See https://api.slack.com/docs/message-attachments#action_fields and https://api.slack.com/docs/message-buttons
// for more information.
type SlackAction struct {
	Type         string                  `json:"type"`
	Text         string                  `json:"text"`
	URL          string                  `json:"url,omitempty"`
	Style        string                  `json:"style,omitempty"`
	Name         string                  `json:"name,omitempty"`
	Value        string                  `json:"value,omitempty"`
	ConfirmField *SlackConfirmationField `json:"confirm,omitempty"`
}

// Validate ensures SlackAction is valid
func (sa *SlackAction) Validate() error {
	if sa.Type == "" {
		return errors.New("missing type in Slack action configuration")
	}
	if sa.Text == "" {
		return errors.New("missing text in Slack action configuration")
	}
	if sa.URL == "" && sa.Name == "" {
		return errors.New("missing name or url in Slack action configuration")
	}
	if sa.ConfirmField != nil {
		if err := sa.ConfirmField.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// SlackConfirmationField protect users from destructive actions or particularly distinguished decisions
// by asking them to confirm their button click one more time.
// See https://api.slack.com/docs/interactive-message-field-guide#confirmation_fields for more information.
type SlackConfirmationField struct {
	Text        string  `json:"text"`
	Title       *string `json:"title,omitempty"`
	OkText      *string `json:"okText,omitempty"`
	DismissText *string `json:"dismissText,omitempty"`
}

// Validate ensures SlackConfirmationField is valid
func (scf *SlackConfirmationField) Validate() error {
	if scf.Text == "" {
		return errors.New("missing text in Slack confirmation configuration")
	}
	return nil
}

// SlackField configures a single Slack field that is sent with each notification.
// Each field must contain a title, value, and optionally, a boolean value to indicate if the field
// is short enough to be displayed next to other fields designated as short.
// See https://api.slack.com/docs/message-attachments#fields for more information.
type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short *bool  `json:"short,omitempty"`
}

// Validate ensures SlackField is valid
func (sf *SlackField) Validate() error {
	if sf.Title == "" {
		return errors.New("missing title in Slack field configuration")
	}
	if sf.Value == "" {
		return errors.New("missing value in Slack field configuration")
	}
	return nil
}

type WebhookConfig struct {
	SendResolved *bool                 `json:"sendResolved,omitempty"`
	URL          *string               `json:"url,omitempty"`
	URLSecret    *v1.SecretKeySelector `json:"urlSecret,omitempty"`
	HTTPConfig   *HTTPConfig           `json:"httpConfig,omitempty"`
	MaxAlerts    *int32                `json:"maxAlerts,omitempty"`
}

type OpsGenieConfig struct {
	SendResolved *bool                     `json:"sendResolved,omitempty"`
	APIKey       *v1.SecretKeySelector     `json:"apiKey,omitempty"`
	APIURL       *string                   `json:"apiURL,omitempty"`
	Message      *string                   `json:"message,omitempty"`
	Description  *string                   `json:"description,omitempty"`
	Source       *string                   `json:"source,omitempty"`
	Tags         *string                   `json:"tags,omitempty"`
	Note         *string                   `json:"note,omitempty"`
	Priority     *string                   `json:"priority,omitempty"`
	Details      []OpsGenieConfigDetail    `json:"details,omitempty"`
	Responders   []OpsGenieConfigResponder `json:"responders,omitempty"`
	HTTPConfig   *HTTPConfig               `json:"httpConfig,omitempty"`
}

// Validate ensures OpsGenieConfig is valid
func (o *OpsGenieConfig) Validate() error {
	for _, responder := range o.Responders {
		if err := responder.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type OpsGenieConfigResponder struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
	Type     string `json:"type,omitempty"`
}

// Validate ensures OpsGenieConfigResponder is valid
func (r *OpsGenieConfigResponder) Validate() error {
	if r.ID == "" && r.Name == "" && r.Username == "" {
		return errors.New("responder must have at least an ID, a Name or an Username defined")
	}

	if !opsGenieTypeRe.MatchString(r.Type) {
		return errors.New("responder type should match team, user, escalation or schedule")
	}

	return nil
}

type HTTPConfig struct {
	BasicAuth         *monitoringv1.BasicAuth     `json:"basicAuth,omitempty"`
	BearerTokenSecret *v1.SecretKeySelector       `json:"bearerTokenSecret,omitempty"`
	TLSConfig         *monitoringv1.SafeTLSConfig `json:"tlsConfig,omitempty"`
	ProxyURL          *string                     `json:"proxyURL,omitempty"`
}

type OpsGenieConfigDetail struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PagerDutyConfigDetail struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type WeChatConfig struct {
	SendResolved *bool                 `json:"sendResolved,omitempty"`
	APISecret    *v1.SecretKeySelector `json:"apiSecret,omitempty"`
	APIURL       *string               `json:"apiURL,omitempty"`
	CorpID       *string               `json:"corpID,omitempty"`
	AgentID      *string               `json:"agentID,omitempty"`
	ToUser       *string               `json:"toUser,omitempty"`
	ToParty      *string               `json:"toParty,omitempty"`
	ToTag        *string               `json:"toTag,omitempty"`
	Message      *string               `json:"message,omitempty"`
	MessageType  *string               `json:"messageType,omitempty"`
	HTTPConfig   *HTTPConfig           `json:"httpConfig,omitempty"`
}

type InhibitRule struct {
	TargetMatch []Matcher `json:"targetMatch,omitempty"`
	SourceMatch []Matcher `json:"sourceMatch,omitempty"`
	Equal       []string  `json:"equal,omitempty"`
}

type Matcher struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Regex bool   `json:"regex,omitempty"`
}

// DeepCopyObject implements the runtime.Object interface.
func (l *AlertmanagerConfig) DeepCopyObject() runtime.Object {
	return l.DeepCopy()
}

// DeepCopyObject implements the runtime.Object interface.
func (l *AlertmanagerConfigList) DeepCopyObject() runtime.Object {
	return l.DeepCopy()
}
