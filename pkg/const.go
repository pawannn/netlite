package pkg

const START_RANGE = 1
const END_RANGE = 65535
const CONCURRENCY = 200
const INTERVAL = 3

const HelpTemplate = `
{{with (or .Long .Short)}}{{.}}{{end}}

Available tools:
{{- range .Commands }}
{{- if .IsAvailableCommand}}

{{ rpad .Name 12 }} {{ .Short }}
{{- if .HasAvailableLocalFlags}}
{{ .LocalFlags.FlagUsages }}

{{- end}}
{{- end}}
{{- end}}

Global flags:
{{ .PersistentFlags.FlagUsages }}

Use "netlite [command] --help" for more information about a command.
`
