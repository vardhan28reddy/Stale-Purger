package consts

type ResourceType string

const (
	DeploymentResource ResourceType = "Deployment"
	ReplicaSetResource ResourceType = "ReplicaSet"
)

type ActionType string

const (
	DeleteActionType          ActionType = "DELETED"
	SkipOnNotStableActionType ActionType = "SKIP_ON_NOT_STABLE"
	SkipOnNoOwnerActionType   ActionType = "SKIP_ON_NO_OWNER"
)

type PodPhase string

const (
	FailedPodPhase  PodPhase = "Failed"
	UnknownPodPhase PodPhase = "Unknown"
)

// Define the HTML template
const HTMLOutputTemplate = `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Stale Pods Report</title>
		<style>
			body { font-family: Arial, sans-serif; margin: 20px; }
			table { border-collapse: collapse; width: 80%; margin-bottom: 40px; }
			th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
			th { background-color: #f2f2f2; }
			h2 { color: #333; }
		</style>
	</head>
	<body>
		<h1>Stale Pods Report</h1>
		{{range $namespace, $pods := .Info}}
			<h2>Namespace: {{$namespace}}</h2>
			<table>
				<tr>
					<th>Pod Name</th>
					<th>Phase</th>
					<th>Reason</th>
				</tr>
				{{range $pod := $pods}}
				<tr>
					<td>{{$pod.PodName}}</td>
					<td>{{$pod.Phase}}</td>
					<td>{{$pod.Reason}}</td>
				</tr>
				{{end}}
			</table>
		{{end}}
	</body>
	</html>`

const HTMLOutputTemplateFileName string = "stale_pods_report.html"
