package file

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/terotuomala/dockerizer/pkg/http"
)

type Dockerfile struct {
	BaseImage          string
	BaseImageDigest    string
	ReleaseImage       string
	ReleaseImageDigest string
	Version            string
	PackageManager     string
	TypeScript         string
}

const startMsg = "Starting to bake Dockerfile. \n"
const fetchBuildDigestMsg = "Fetching digest for the 'build' stage image"
const fetchReleaseDigestMsg = "Fetching digest for the 'release' stage image"

const (
	goBuildTemplate = `

	FROM golang:{{ .Version }}-{{ .BaseImage }}@{{ .BaseImageDigest }} as build

	WORKDIR /app

	RUN useradd -u 1001 -m app

	COPY go.mod ./

	RUN go mod download

	COPY *.go ./

	RUN go build -o /app

	{{ if eq .ReleaseImage "scratch" }}FROM {{ .ReleaseImage }} as release{{ else }}FROM golang:{{ .Version }}-{{ .ReleaseImage }}@{{ .ReleaseImageDigest }} as release{{ end }}

	USER 1001

	WORKDIR /app
	
	COPY --from=build /app /app
	COPY --from=build /etc/passwd /etc/passwd
	{{ if eq .ReleaseImage "scratch" }}COPY --from=golang:{{ .Version }}-{{ .BaseImage }}@{{ .BaseImageDigest }} /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/{{ end }}

	ENTRYPOINT ["/app"]
`
)

const (
	nodejsBuildTemplate = `
	FROM node:{{ .Version }}-{{ .BaseImage }}@{{ .BaseImageDigest }} as build

	WORKDIR /app

	{{ if eq .PackageManager "npm" }}COPY package.json package-lock.json ./{{ else }}COPY package.json yarn.lock ./{{ end }}

	{{ if and (eq .PackageManager "npm") (eq .TypeScript "no")}}RUN npm ci --production{{ else if and (eq .PackageManager "npm") (eq .TypeScript "yes")}}RUN npm ci{{ else if and (eq .PackageManager "yarn") (eq .TypeScript "no")}}RUN yarn install --production --frozen-lockfile{{ else if and (eq .PackageManager "yarn") (eq .TypeScript "yes")}}RUN yarn install --frozen-lockfile{{ end }}
	
	COPY . .

	{{if and (eq .PackageManager "npm") (eq .TypeScript "no")}}{{ end }}{{ if and (eq .PackageManager "npm") (eq .TypeScript "yes")}}RUN npm run build{{ end }}{{ if and (eq .PackageManager "yarn") (eq .TypeScript "no")}}{{ end }}{{if and (eq .PackageManager "yarn") (eq .TypeScript "yes")}}RUN yarn run build{{ end }}

	FROM node:{{ .Version }}-{{ .ReleaseImage }}@{{ .ReleaseImageDigest }} as release

	USER node

	ENV NPM_CONFIG_LOGLEVEL=warn
	ENV NODE_ENV=production

	WORKDIR /home/node

	COPY --chown=node:node --from=build /app .

	CMD ["node", "src/index.js"]
`
)

func CreateGoDockerfile(language, baseImage, releaseImage, goVersion string) {
	fmt.Println(startMsg)

	var dt Dockerfile

	baseImageDigest := http.GetDockerImageDigest("library/golang", goVersion+"-"+baseImage)
	fmt.Printf(fetchBuildDigestMsg+" golang:%v \n \n", goVersion+"-"+baseImage)

	switch releaseImage {
	case "scratch":
		dt = Dockerfile{
			BaseImage:       baseImage,
			BaseImageDigest: baseImageDigest,
			ReleaseImage:    releaseImage,
			Version:         goVersion,
		}

	default:
		releaseImageDigest := http.GetDockerImageDigest("library/golang", goVersion+"-"+releaseImage)
		fmt.Printf(fetchReleaseDigestMsg+" golang:%v \n \n", goVersion+"-"+releaseImage)

		dt = Dockerfile{
			BaseImage:          baseImage,
			BaseImageDigest:    baseImageDigest,
			ReleaseImage:       releaseImage,
			ReleaseImageDigest: releaseImageDigest,
			Version:            goVersion,
		}
	}

	buildTemplate(language, dt)

}

func CreateNodejsDockerfile(language, baseImage, releaseImage, nodeVersion, packageManager, typeScript string) {
	fmt.Println(startMsg)

	baseImageDigest := http.GetDockerImageDigest("library/node", nodeVersion+"-"+baseImage)
	fmt.Printf(fetchBuildDigestMsg+" node:%v \n \n", nodeVersion+"-"+baseImage)

	releaseImageDigest := http.GetDockerImageDigest("library/node", nodeVersion+"-"+releaseImage)
	fmt.Printf(fetchReleaseDigestMsg+" node:%v \n \n", nodeVersion+"-"+releaseImage)

	dt := Dockerfile{
		BaseImage:          baseImage,
		BaseImageDigest:    baseImageDigest,
		ReleaseImage:       releaseImage,
		ReleaseImageDigest: releaseImageDigest,
		Version:            nodeVersion,
		PackageManager:     packageManager,
		TypeScript:         typeScript,
	}

	buildTemplate(language, dt)
}

func buildTemplate(language string, dockerFileTemplate Dockerfile) {
	var t *template.Template

	switch language {
	case "Go":
		t = template.Must(template.New("buildTemplate").Parse(goBuildTemplate))

	case "Node.js":
		t = template.Must(template.New("buildTemplate").Parse(nodejsBuildTemplate))
	}

	f, err := os.Create("Dockerfile")

	if err != nil {
		log.Fatalf("Error creating Dockerfile %v", err.Error())
		return
	}
	defer f.Close()

	t.Execute(f, dockerFileTemplate)

	fmt.Println("Dockerfile baked succesfully!")
}
