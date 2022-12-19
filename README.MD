# The repo for evaluating the custom query for supporting logrus

## Background
So far out of the box query has not supported the 3rd party package logrus (https://github.com/sirupsen/logrus) with GO lang. Whenever you use SetFormatter with DisableHTMLEscape:false, it still alerts as “Log entries created from user input” CWE-117.

## Root cause
The query -go/log-injection.ql only supports default log framework. So even if you sanitized with "DisableHTMLEscape : false", it would detect the alert.

## Workaround
I generated the custom query based on the original go/log-injection. The query is custom-queries/go/Custom_LogInjection.ql
```
... snippet
  private class SafeFormatArgumentSanitizer extends Sanitizer {
    SafeFormatArgumentSanitizer() {
      exists(StringOps::Formatting::StringFormatCall call, string safeDirective |
        this = call.getOperand(_, safeDirective) and
        // Mark "%q" formats as safe, but not "%#q", which would preserve newline characters.
        safeDirective.regexpMatch("%[^%#]*q")
      ) or
      exists(CallExpr cexpr | cexpr.getAnArgument().getAChild().getChild(1).getChild(0).toString() = "DisableHTMLEscape" and 
        cexpr.getAnArgument().getAChild().getChild(1).getChild(1).toString() = "false")
    }
  }
...
```
### For applying the custom query
- The following files are needed.
	- custom-queries/go/qlpack.yml
	- custom-queries/code-scanning.qls
	- .github/codeql/codeql-config.yml
- you need modify .github/workflows/codeql.yml
```
...snippet
    - name: Initialize CodeQL
      uses: github/codeql-action/init@v2
      with:
        languages: ${{ matrix.language }}
        # If you wish to specify custom queries, you can do so here or in a config file.
        # By default, queries listed here will override any specified in a config file.
        # Prefix the list here with "+" to use these queries and those in the config file.

        # Details on CodeQL's query packs refer to : https://docs.github.com/en/code-security/code-scanning/automatically-scanning-your-code-for-vulnerabilities-and-errors/configuring-code-scanning#using-queries-in-ql-packs
        # queries: security-extended,security-and-quality
        config-file: ./.github/codeql/codeql-config.yml
```
## CLI
```
codeql database create ../codeql-db/go-samples --language="go" --overwrite
codeql database analyze ../codeql-db/go-samples --format=sarif-latest --output=result.sarif custom-queries/code-scanning.qls
```

## Actions
please refer to ./github/workflows/codeql.yml

## how to check
### Vulnerable code
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
### Sanitized code
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
- This sample is a simple web server. When you specify the URL http://localhost:8090/hello on your browser, you can see "hello" on the browser.
- Also you can see some log message on your console

### How to build and run
1) go run sandbox
2) http://localhost:8090/hello on a browser
- In advance you might need the commands "go mod".
