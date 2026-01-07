package report

import (
	"html/template"
	"strings"

	"github.com/ismailtsdln/HardenaK8s/internal/policy"
)

// HTMLFormatter implements Formatter for HTML output
type HTMLFormatter struct{}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>HardenaK8s Security Audit Report</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;800&display=swap" rel="stylesheet">
    <style>
        :root {
            --primary: #7D56F4;
            --bg: #0f172a;
            --card-bg: #1e293b;
            --text-main: #f8fafc;
            --text-dim: #94a3b8;
            --critical: #ef4444;
            --high: #f97316;
            --medium: #eab308;
            --low: #22c55e;
            --info: #3b82f6;
        }

        body {
            font-family: 'Inter', sans-serif;
            background-color: var(--bg);
            color: var(--text-main);
            margin: 0;
            padding: 2rem;
            line-height: 1.5;
        }

        .container {
            max-width: 1000px;
            margin: 0 auto;
        }

        header {
            text-align: center;
            margin-bottom: 3rem;
            padding: 3rem;
            background: linear-gradient(135deg, var(--primary), #4c1d95);
            border-radius: 1.5rem;
            box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
        }

        h1 { margin: 0; font-size: 3rem; font-weight: 800; }
        .tagline { opacity: 0.8; font-size: 1.1rem; }

        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 1.5rem;
            margin-bottom: 3rem;
        }

        .stat-card {
            background: var(--card-bg);
            padding: 1.5rem;
            border-radius: 1rem;
            text-align: center;
            border: 1px solid rgba(255,255,255,0.05);
            transition: transform 0.2s;
        }
        
        .stat-card:hover { transform: translateY(-5px); }
        .stat-value { font-size: 2rem; font-weight: 800; display: block; }
        .stat-label { font-size: 0.8rem; color: var(--text-dim); text-transform: uppercase; letter-spacing: 0.05em; }

        .issue-card {
            background: var(--card-bg);
            margin-bottom: 1.5rem;
            border-radius: 1rem;
            overflow: hidden;
            border-left: 6px solid #ccc;
            border: 1px solid rgba(255,255,255,0.05);
        }

        .issue-header {
            padding: 1.5rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
            border-bottom: 1px solid rgba(255,255,255,0.05);
        }

        .severity-badge {
            padding: 0.25rem 0.75rem;
            border-radius: 9999px;
            font-size: 0.75rem;
            font-weight: 700;
            text-transform: uppercase;
        }

        .CRITICAL { border-left: 8px solid var(--critical); color: var(--critical); }
        .HIGH { border-left: 8px solid var(--high); color: var(--high); }
        .MEDIUM { border-left: 8px solid var(--medium); color: var(--medium); }
        .LOW { border-left: 8px solid var(--low); color: var(--low); }
        .INFO { border-left: 8px solid var(--info); color: var(--info); }

        .issue-body { padding: 1.5rem; }
        .remediation {
            background: rgba(0,0,0,0.2);
            padding: 1rem;
            border-radius: 0.5rem;
            margin-top: 1rem;
            border: 1px dashed var(--primary);
        }
        
        .footer {
            text-align: center;
            margin-top: 5rem;
            color: var(--text-dim);
            font-size: 0.9rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>HardenaK8s</h1>
            <p class="tagline">Security Audit & Hardening Report</p>
        </header>

        <div class="stats-grid">
            <div class="stat-card">
                <span class="stat-value">{{.Stats.TotalIssues}}</span>
                <span class="stat-label">Total Issues</span>
            </div>
            {{range $sev, $count := .Stats.SeverityCount}}
            <div class="stat-card">
                <span class="stat-value">{{$count}}</span>
                <span class="stat-label">{{$sev}}</span>
            </div>
            {{end}}
        </div>

        <h2>Security Findings</h2>
        {{range .Issues}}
        <div class="issue-card {{.Severity}}">
            <div class="issue-header">
                <h3>{{.Title}}</h3>
                <span class="severity-badge">{{.Severity}}</span>
            </div>
            <div class="issue-body">
                <p><strong>Resource:</strong> {{.Resource}} ({{.Namespace}})</p>
                <p>{{.Description}}</p>
                <div class="remediation">
                    <strong>Remediation:</strong> {{.Remediation}}
                </div>
            </div>
        </div>
        {{else}}
        <p>No security issues found! Your cluster is hardened. 🛡️</p>
        {{end}}

        <div class="footer">
            Generated by HardenaK8s CLI tool.
        </div>
    </div>
</body>
</html>
`

func (f *HTMLFormatter) Format(result *policy.Result) ([]byte, error) {
	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, result); err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}
