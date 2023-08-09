// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fluentbit

const BaseConfigName = "fluent-bit.conf"
const UpstreamConfigName = "upstream.conf"
const CustomParsersConfigName = "custom-parsers.conf"
const StockConfigPath = "/fluent-bit/etc"
const StockBinPath = "/fluent-bit/bin/fluent-bit"
const OperatorConfigPath = "/fluent-bit/etc-operator"

var fluentBitConfigTemplate = `
[SERVICE]
    Flush        {{ .Flush }}
    Grace        {{ .Grace }}
    Daemon       Off
    Log_Level    {{ .LogLevel }}
    Parsers_File {{ .DefaultParsers }}
    {{- if .CustomParsers }}
    Parsers_File {{ .CustomParsers }}
    {{- end }}
    Coro_Stack_Size    {{ .CoroStackSize }}
    {{- if .Monitor.Enabled }}
    HTTP_Server  On
    HTTP_Listen  0.0.0.0
    HTTP_Port    {{ .Monitor.Port }}
    {{- end }}
    {{- range $key, $value := .BufferStorage }}
    {{- if $value }}
    {{ $key }}  {{$value}}
    {{- end }}
    {{- end }}

[INPUT]
    Name         tail
    {{- range $key, $value := .Input.Values }}
    {{- if $value }}
    {{ $key }}  {{$value}}
    {{- end }}
    {{- end }}
    {{- range $id, $v := .Input.ParserN }}
    {{- if $v }}
    Parse_{{ $id}} {{$v}}
    {{- end }}
    {{- end }}
    {{- if .Input.MultilineParser }}
    multiline.parser {{- range $i, $v := .Input.MultilineParser }}{{ if $i }},{{ end}} {{ $v }}{{ end }}
    {{- end }}

{{- if not .DisableKubernetesFilter }}
[FILTER]
    Name        kubernetes
    {{- range $key, $value := .KubernetesFilter }}
    {{- if $value }}
    {{ $key }}  {{$value}}
    {{- end }}
    {{- end }}
{{- end}}

{{- if .AwsFilter }}
[FILTER]
    Name        aws
    {{- range $key, $value := .AwsFilter }}
    {{- if $value }}
    {{ $key }}  {{$value}}
    {{- end }}
    {{- end }}
{{- end}}

{{- range $modify := .FilterModify }}

[FILTER]
    Name modify
    Match *
    {{- range $condition := $modify.Conditions }}
    {{- $operation :=  $condition.Operation }}
    Condition {{ $operation.Op }} {{ $operation.Key }} {{ if $operation.Value }}{{ $operation.Value }}{{ end }}
    {{- end }}

    {{- range $rule := $modify.Rules }}
    {{- $operation :=  $rule.Operation }}
    {{ $operation.Op }} {{ $operation.Key }} {{ if $operation.Value }}{{ $operation.Value }}{{ end }}
    {{- end }}
{{- end}}

{{ with $out := .FluentForwardOutput }}
{{- range $target := $out.Targets }}
[OUTPUT]
    Name          forward
    {{- if eq $target.Namespace "*" }}
    Match *
    {{- else }}
    Match *_{{ $target.Namespace }}_*
    {{- end }}
    {{- if $out.Upstream.Enabled }}
    Upstream      {{ $out.Upstream.Config.Path }}
    {{- else }}
    Host          {{ $target.Host }}
    Port          {{ $target.Port }}
    {{- end }}
    {{ if $out.TLS.Enabled }}
    tls           On
    tls.verify    Off
    tls.ca_file   /fluent-bit/tls/ca.crt
    tls.crt_file  /fluent-bit/tls/tls.crt
    tls.key_file  /fluent-bit/tls/tls.key
    {{- if $out.TLS.SharedKey }}
    Shared_Key    {{ $out.TLS.SharedKey }}
    {{- else }}
    Empty_Shared_Key true
    {{- end }}
    {{- end }}
    {{- if $out.Network.ConnectTimeoutSet }}
    net.connect_timeout {{ $out.Network.ConnectTimeout }}
    {{- end }}
    {{- if $out.Network.ConnectTimeoutLogErrorSet }}
    net.connect_timeout_log_error {{ $out.Network.ConnectTimeoutLogError }}
    {{- end }}
    {{- if $out.Network.DNSMode }}
    net.dns.mode {{ $out.Network.DNSMode }}
    {{- end }}
    {{- if $out.Network.DNSPreferIPV4Set }}
    net.dns.prefer_ipv4 {{ $out.Network.DNSPreferIPV4 }}
    {{- end }}
    {{- if $out.Network.DNSResolver }}
    net.dns.resolver {{ $out.Network.DNSResolver }}
    {{- end }}
    {{- if $out.Network.KeepaliveSet}}
    net.keepalive {{if $out.Network.Keepalive }}on{{else}}off{{end}}
    {{- end }}
    {{- if $out.Network.KeepaliveIdleTimeoutSet }}
    net.keepalive_idle_timeout {{ $out.Network.KeepaliveIdleTimeout }}
    {{- end }}
    {{- if $out.Network.KeepaliveMaxRecycleSet }}
    net.keepalive_max_recycle {{ $out.Network.KeepaliveMaxRecycle }}
    {{- end }}
    {{- if $out.Network.SourceAddress }}
    net.source_address {{ $out.Network.SourceAddress }}
    {{- end }}
    {{- with $out.Options }}
    {{- range $key, $value := . }}
    {{- if $value }}
    {{ $key }}  {{$value}}
    {{- end }}
    {{- end }}
    {{- end }}
{{- end }}
{{- end }}

{{- with $out := .SyslogNGOutput }}
{{- range $target := $out.Targets }}
[OUTPUT]
    Name tcp
    {{- if eq $target.Namespace "*" }}
    Match *
    {{- else }}
    Match *_{{ $target.Namespace }}_*
    {{- end }}
    Host {{ $target.Host }}
    Port {{ $target.Port }}
    Format json_lines
    {{- with $out.JSONDateKey }}
    json_date_key {{ . }}
    {{- end }}
    {{- with $out.JSONDateFormat }}
    json_date_format {{ . }}
    {{- end }}
    {{- with $out.Workers }}
    Workers {{ . }}
    {{- end }}
    {{- if $out.Network.ConnectTimeoutSet }}
    net.connect_timeout {{ $out.Network.ConnectTimeout }}
    {{- end }}
    {{- if $out.Network.ConnectTimeoutLogErrorSet }}
    net.connect_timeout_log_error {{ $out.Network.ConnectTimeoutLogError }}
    {{- end }}
    {{- if $out.Network.DNSMode }}
    net.dns.mode {{ $out.Network.DNSMode }}
    {{- end }}
    {{- if $out.Network.DNSPreferIPV4Set }}
    net.dns.prefer_ipv4 {{ $out.Network.DNSPreferIPV4 }}
    {{- end }}
    {{- if $out.Network.DNSResolver }}
    net.dns.resolver {{ $out.Network.DNSResolver }}
    {{- end }}
    {{- if $out.Network.KeepaliveSet}}
    net.keepalive {{if $out.Network.Keepalive }}on{{else}}off{{end}}
    {{- end }}
    {{- if $out.Network.KeepaliveIdleTimeoutSet }}
    net.keepalive_idle_timeout {{ $out.Network.KeepaliveIdleTimeout }}
    {{- end }}
    {{- if $out.Network.KeepaliveMaxRecycleSet }}
    net.keepalive_max_recycle {{ $out.Network.KeepaliveMaxRecycle }}
    {{- end }}
    {{- if $out.Network.SourceAddress }}
    net.source_address {{ $out.Network.SourceAddress }}
    {{- end }}
{{- end }}
{{- end }}
`

var upstreamConfigTemplate = `
[UPSTREAM]
    Name {{ .Config.Name }}
{{- range $idx, $element:= .Config.Nodes}}
[NODE]
    Name {{.Name}}
    Host {{.Host}}
    Port {{.Port}}
{{- end}}
`
