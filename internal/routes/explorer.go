package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Explorer(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(explorerHTML))
}

const explorerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
	<title>Kavalife ERP API</title>

	<style>
		:root {
			--bg: #0f172a;
			--card: #020617;
			--border: #1e293b;
			--text: #e5e7eb;
			--muted: #94a3b8;
			--primary: #38bdf8;
			--success: #22c55e;
		}

		* { box-sizing: border-box; }

		body {
			margin: 0;
			font-family: system-ui, -apple-system, BlinkMacSystemFont, sans-serif;
			background: radial-gradient(circle at top, #020617, #000);
			color: var(--text);
			min-height: 100vh;
			display: flex;
			align-items: center;
			justify-content: center;
			padding: 24px;
		}

		.container {
			width: 100%;
			max-width: 760px;
		}

		.header {
			margin-bottom: 24px;
		}

		.header h1 {
			font-size: 2rem;
			margin: 0 0 6px;
		}

		.header p {
			margin: 0;
			color: var(--muted);
		}

		.card {
			background: linear-gradient(180deg, var(--card), #000);
			border: 1px solid var(--border);
			border-radius: 16px;
			padding: 24px;
			box-shadow: 0 20px 40px rgba(0,0,0,.6);
		}

		.grid {
			display: grid;
			grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
			gap: 16px;
		}

		.section h3 {
			margin-top: 0;
			font-size: 0.9rem;
			letter-spacing: .04em;
			text-transform: uppercase;
			color: var(--muted);
		}

		.list {
			list-style: none;
			padding: 0;
			margin: 0;
		}

		.list li {
			margin-bottom: 12px;
		}

		.link {
			display: flex;
			align-items: center;
			justify-content: space-between;
			padding: 14px 16px;
			border-radius: 10px;
			border: 1px solid var(--border);
			text-decoration: none;
			color: var(--text);
			background: #020617;
			transition: all .15s ease;
		}

		.link:hover {
			border-color: var(--primary);
			transform: translateY(-1px);
		}

		.badge {
			font-size: .75rem;
			padding: 4px 8px;
			border-radius: 999px;
			background: #022c22;
			color: var(--success);
			border: 1px solid #14532d;
		}

		.footer {
			margin-top: 24px;
			text-align: center;
			font-size: .85rem;
			color: var(--muted);
		}
	</style>
</head>

<body>
	<div class="container">
		<div class="header">
			<h1>Kavalife ERP API</h1>
		</div>

		<div class="card">
			<div class="grid">

				<div class="section">
					<h3>API Tools</h3>
					<ul class="list">
						<li>
							<a class="link" href="/swagger/index.html">
								<span>Swagger UI</span>
								<span class="badge">Docs</span>
							</a>
						</li>
						<li>
							<a class="link" href="/swagger/doc.json">
								<span>OpenAPI Spec</span>
								<span class="badge">JSON</span>
							</a>
						</li>
					</ul>
				</div>

				<div class="section">
					<h3>Health</h3>
					<ul class="list">
						<li>
							<a class="link" href="/api/health">
								<span>/api/health</span>
								<span class="badge">200 OK</span>
							</a>
						</li>
					</ul>
				</div>

			</div>
		</div>

		<div class="footer">
			Kavalife ERP Backend • Powered by Go & Gin
		</div>
	</div>
</body>
</html>`
