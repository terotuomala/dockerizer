package prompt

import (
	"log"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/terotuomala/dockerizer/pkg/file"
)

type promptContent struct {
	errorMsg string
	label    string
}

func getLanguages() []string {
	return []string{"Go", "Node.js"}
}

func getGoBaseImageDistributions() []string {
	return []string{"alpine", "bullseye"}
}

func getGoReleaseImageDistributions() []string {
	return []string{"alpine", "bullseye", "scratch"}
}

func getNodejsImageDistributions() []string {
	return []string{"alpine", "bullseye", "slim"}
}

func getGoVersions() []string {
	return []string{"1.19", "1.20"}
}

func getNodeVersions() []string {
	return []string{"16", "18", "19", "20"}
}

func getNodePackageManagers() []string {
	return []string{"npm", "yarn"}
}

func useTypeScript() []string {
	return []string{"yes", "no"}
}

func promptGetSelect(pc promptContent, selectableItems []string) string {
	items := selectableItems
	index := -1
	var result string
	var err error

	for index < 0 {
		templates := &promptui.SelectTemplates{
			Label:    "{{ . | green }}?",
			Active:   "\U0001F433 {{ .| cyan }}",
			Inactive: "{{ . | blue }}",
			Selected: "{{ . | cyan }}",
		}

		prompt := promptui.Select{
			Label:     pc.label,
			Items:     items,
			Templates: templates,
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		log.Fatalf("Promt failed %v", err.Error())
		os.Exit(1)
	}

	return result
}

func StartPrompt() {
	const baseImageErroMsg = "No docker image distribution selected for 'build' stage"
	const baseImageLabel = "Select docker image distribution for 'build' stage"
	const releaseImageErroMsg = "No docker image distribution selected for 'release' stage"
	const releaseImageLabel = "Select docker image distribution for 'release' stage"
	const versionContentErrorMsg = "No version selected"
	const versionContentLabel = "Select version"

	baseLanguageSelection := promptContent{
		errorMsg: "No language selected",
		label:    "Select language",
	}
	language := promptGetSelect(baseLanguageSelection, getLanguages())

	switch language {
	case "Go":
		baseImagePromptContent := promptContent{
			errorMsg: baseImageErroMsg,
			label:    baseImageLabel,
		}
		baseImage := promptGetSelect(baseImagePromptContent, getGoBaseImageDistributions())

		releaseImagePromptContent := promptContent{
			errorMsg: releaseImageErroMsg,
			label:    releaseImageLabel,
		}
		releaseImage := promptGetSelect(releaseImagePromptContent, getGoReleaseImageDistributions())

		versionPromptContent := promptContent{
			errorMsg: versionContentErrorMsg,
			label:    versionContentLabel,
		}
		goVersion := promptGetSelect(versionPromptContent, getGoVersions())

		file.CreateGoDockerfile(language, baseImage, releaseImage, goVersion)

	case "Node.js":
		baseImagePromptContent := promptContent{
			errorMsg: baseImageErroMsg,
			label:    baseImageLabel,
		}
		baseImage := promptGetSelect(baseImagePromptContent, getNodejsImageDistributions())

		releaseImagePromptContent := promptContent{
			errorMsg: releaseImageErroMsg,
			label:    releaseImageLabel,
		}
		releaseImage := promptGetSelect(releaseImagePromptContent, getNodejsImageDistributions())

		nodeVersionPromptContent := promptContent{
			errorMsg: versionContentErrorMsg,
			label:    versionContentLabel,
		}
		nodeVersion := promptGetSelect(nodeVersionPromptContent, getNodeVersions())

		nodePackageManagerPromptContent := promptContent{
			errorMsg: "No package manager selected",
			label:    "Select package manager",
		}
		nodePackageManager := promptGetSelect(nodePackageManagerPromptContent, getNodePackageManagers())

		useTypeScriptPromptContent := promptContent{
			errorMsg: "TypeScript usage not selected",
			label:    "Use TypeScript?",
		}
		useTypeScript := promptGetSelect(useTypeScriptPromptContent, useTypeScript())

		file.CreateNodejsDockerfile(language, baseImage, releaseImage, nodeVersion, nodePackageManager, useTypeScript)
	}
}
