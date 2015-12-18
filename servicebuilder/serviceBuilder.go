package servicebuilder

import (
	"log"
	"strings"

	"github.com/bndr/gopencils"
	"github.com/essentier/spickspan/config"
)

func createServiceBuilder(serviceConfig config.Service, providerUrl string, token string) *serviceBuilder {
	nomockApi := gopencils.Api(providerUrl)
	return &serviceBuilder{nomockApi: nomockApi, token: token, serviceConfig: serviceConfig}
}

type serviceBuilder struct {
	nomockApi     *gopencils.Resource
	serviceConfig config.Service
	token         string
}

func (p *serviceBuilder) buildService() error {
	if !p.serviceConfig.IsSourceProject() {
		panic("Could not build a service that is not a source project.")
	}

	log.Printf("going to git push code of service %v", p.serviceConfig.ServiceName)
	err := gitPush(p.serviceConfig.ServiceName, p.serviceConfig.ProjectSrcRoot, p.nomockApi.Api.BaseUrl.String()+"/nomockbuilder", p.token)
	if err != nil {
		return err
	}

	log.Printf("going to build code")
	return p.buildServiceOnNomockBuilder(p.serviceConfig.ServiceName)
}

func (p *serviceBuilder) buildServiceOnNomockBuilder(serviceName string) error {
	log.Printf("building service %v on nomock builder", serviceName)
	builderResource := p.nomockApi.Res("nomockbuilder/build/" + serviceName)
	builderResource.SetHeader("Authorization", "Bearer "+p.token)
	_, err := builderResource.Get()
	return err
}

func gitPush(serviceName string, projectDir string, builderUrl string, token string) error {
	log.Printf("gitpusher.Push project dir: %v", projectDir)
	gitVersion, _ := runCmd("git", "--version")
	log.Printf("git version is %v", gitVersion.String())
	remoteUrl := getEssentierGitRemote(serviceName, builderUrl, token)
	git := &gitProject{projectDir: projectDir, err: nil}
	originalBranch := git.getCurrentBranch()
	git.stashAll()
	if git.err == nil {
		defer git.deferredPopStashed()
	}

	git.branch("nomock")
	if git.err == nil {
		defer git.deferredDeleteBranch("nomock")
	}

	git.checkout("nomock")
	if git.err == nil {
		defer git.deferredCheckout(originalBranch)
	}

	git.pull(remoteUrl, "nomock")
	git.applyStash()
	git.addAll()
	git.commit("'done by nomock'")
	git.push(remoteUrl, "nomock")
	return git.err
}

func getEssentierGitRemote(serviceName string, builderUrl string, token string) string {
	remoteUrl := ""
	if strings.HasPrefix(builderUrl, "git") {
		remoteUrl = builderUrl + ":" + serviceName
	} else if strings.HasPrefix(builderUrl, "http://") {
		if strings.HasSuffix(builderUrl, "/") {
			remoteUrl = "http://" + token + ":@" + builderUrl[7:] + serviceName
		} else {
			remoteUrl = "http://" + token + ":@" + builderUrl[7:] + "/" + serviceName
		}
	}
	return remoteUrl
}