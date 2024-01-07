{{if .Server}}# if your repository is private, please login...
# docker login {{ .Server }} --username='your username' --password='your password' [server]
{{end}}
{{- range .Outputs }}
docker pull {{ .Target }}
docker tag {{ .Target }} {{ .Source }}
{{- end }}