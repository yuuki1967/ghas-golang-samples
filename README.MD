# The repo for evaluating the custom query for supporting logrus

## background
The customer is using the 3rd party package logrus (https://github.com/sirupsen/logrus) with GO lang.
When he scanned the codebase with CodeQL, it detected a false positive - “Log entries created from user input” CWE-117

## Root casue
The query -go/log-injection.ql only supports default log framework. So even if you sanitized with "DisableHTMLEscape : false", it would detect the alert.

## workaround
I generated the custom query based on the original go/log-injection

## CLI
```
codeql database create ../codeql-db/go-samples --language="go" --overwrite
codeql database analyze ../codeql-db/go-samples --format=sarif-latest --output=result.sarif custom-queries/code-scanning.qls
```

## Actions
please refer to ./github/workflows/codeql.yml

## how to check
### To detect an alert
```
func init() {
	// Log as JSON instead of the default ASCII formatter.
**	logr.SetFormatter(&logr.JSONFormatter{DisableHTMLEscape: true}) **
  ...
}
...
func hello(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	b, _ := io.ReadAll(req.Body)
	username := string(b)
	logr.WithFields(logr.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("user %s logged in.\n", username)
...
```
### To remidate 
```
func init() {
	// Log as JSON instead of the default ASCII formatter.
**	logr.SetFormatter(&logr.JSONFormatter{DisableHTMLEscape: false}) **
  ...
}
...
func hello(w http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	b, _ := io.ReadAll(req.Body)
	username := string(b)
	logr.WithFields(logr.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("user %s logged in.\n", username)
...
```

## About the sample code
### How to build and run
1) go run sandbox
2) http://localhost:8090/hello on a browser
