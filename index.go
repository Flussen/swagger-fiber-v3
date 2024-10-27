package swagger

const indexTmpl string = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
    <link href="https://fonts.googleapis.com/css?family=Open+Sans:400,700|Source+Code+Pro:300,600|Titillium+Web:400,600,700" rel="stylesheet">
    <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.1.3/swagger-ui.css">
    <link rel="icon" type="image/png" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.1.3/favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.1.3/favicon-16x16.png" sizes="16x16" />
    {{- if .CustomStyle}}
      <style>
        body { margin: 0; }
        {{.CustomStyle}}
      </style>
    {{- end}}
    {{- if .CustomScript}}
      <script>
        {{.CustomScript}}
      </script>
    {{- end}}
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.1.3/swagger-ui-bundle.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.1.3/swagger-ui-standalone-preset.js"></script>
    <script>
    window.onload = function() {
      const config = {{.}};
      config.dom_id = '#swagger-ui';
      config.plugins = [
        {{- range $plugin := .Plugins }}
          {{$plugin}},
        {{- end}}
      ];
      config.presets = [
        {{- range $preset := .Presets }}
          {{$preset}},
        {{- end}}
      ];
      config.filter = {{.Filter.Value}};
      config.syntaxHighlight = {{.SyntaxHighlight.Value}};
      {{if .TagsSorter}} config.tagsSorter = {{.TagsSorter}}; {{end}}
      {{if .OnComplete}} config.onComplete = {{.OnComplete}}; {{end}}
      {{if .RequestInterceptor}} config.requestInterceptor = {{.RequestInterceptor}}; {{end}}
      {{if .ResponseInterceptor}} config.responseInterceptor = {{.ResponseInterceptor}}; {{end}}
      {{if .ModelPropertyMacro}} config.modelPropertyMacro = {{.ModelPropertyMacro}}; {{end}}
      {{if .ParameterMacro}} config.parameterMacro = {{.ParameterMacro}}; {{end}}

      const ui = SwaggerUIBundle(config);

      {{if .OAuth}} ui.initOAuth({{.OAuth}}); {{end}}
      {{if .PreauthorizeBasic}} ui.preauthorizeBasic({{.PreauthorizeBasic}}); {{end}}
      {{if .PreauthorizeApiKey}} ui.preauthorizeApiKey({{.PreauthorizeApiKey}}); {{end}}

      window.ui = ui;
    }
    </script>
  </body>
</html>
`
