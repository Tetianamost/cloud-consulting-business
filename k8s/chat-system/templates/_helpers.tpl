{{/*
Expand the name of the chart.
*/}}
{{- define "ai-consultant-chat.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "ai-consultant-chat.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "ai-consultant-chat.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "ai-consultant-chat.labels" -}}
helm.sh/chart: {{ include "ai-consultant-chat.chart" . }}
{{ include "ai-consultant-chat.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "ai-consultant-chat.selectorLabels" -}}
app.kubernetes.io/name: {{ include "ai-consultant-chat.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "ai-consultant-chat.serviceAccountName" -}}
{{- if .Values.security.serviceAccount.create }}
{{- default (include "ai-consultant-chat.fullname" .) .Values.security.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.security.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the PostgreSQL service
*/}}
{{- define "ai-consultant-chat.postgresql.fullname" -}}
{{- include "ai-consultant-chat.fullname" . }}-postgresql
{{- end }}

{{/*
Create the name of the Redis service
*/}}
{{- define "ai-consultant-chat.redis.fullname" -}}
{{- include "ai-consultant-chat.fullname" . }}-redis-master
{{- end }}

{{/*
Create database URL
*/}}
{{- define "ai-consultant-chat.databaseUrl" -}}
postgres://{{ .Values.postgresql.auth.username }}:{{ .Values.postgresql.auth.password }}@{{ include "ai-consultant-chat.postgresql.fullname" . }}:5432/{{ .Values.postgresql.auth.database }}?sslmode=disable
{{- end }}

{{/*
Create Redis URL
*/}}
{{- define "ai-consultant-chat.redisUrl" -}}
redis://{{ include "ai-consultant-chat.redis.fullname" . }}:6379/0
{{- end }}