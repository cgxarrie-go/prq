package createpr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cgxarrie-go/prq/internal/azure/origin"
	"github.com/cgxarrie-go/prq/internal/models"
	"github.com/cgxarrie-go/prq/internal/ports"
	"github.com/cgxarrie-go/prq/internal/utils"
)

type service struct {
	pat       string
	originSvc ports.OriginSvc
}

// NewService return new instnce of azure service
func NewService(pat string, originSvc ports.OriginSvc) ports.PRCreator {
	return service{
		pat:       fmt.Sprintf("`:%s", pat),
		originSvc: originSvc,
	}
}

// Create .
func (svc service) Run(req ports.CreatePRRequest) (
	pr models.CreatedPullRequest, err error) {

	src, err := utils.GitCurrentBranchName()
	if err != nil {
		return pr, fmt.Errorf("getting current branch name: %w", err)
	}

	source := svc.originSvc.NewBranch(src)
	destination := svc.originSvc.DefaultTargetBranch()
	if req.Destination != "" {
		destination = svc.originSvc.NewBranch(req.Destination)
	}

	title := req.Title
	if req.Title == "" {
		title = fmt.Sprintf("PR from %s to %s",
			source.Name(), destination.Name())
	}

	desc := []byte("")
	if req.PRTemplate != "" {
		desc, err = os.ReadFile(req.PRTemplate)
		if err != nil {
			desc = []byte("")
		}
	}

	o, err := utils.CurrentFolderRemote()
	if err != nil {
		return pr, fmt.Errorf("getting repository origin: %w", err)
	}

	azOrigin := origin.NewAzureOrigin(o)

	svcResp := Response{}
	err = svc.doPOST(source, destination, title, string(desc), req.IsDraft,
		azOrigin, &svcResp)

	if err != nil {
		return pr, fmt.Errorf("creating PR: %w", err)
	}

	pr = svcResp.ToPullRequest(azOrigin.Organization())
	pr.Link, _ = svc.originSvc.PRLink(o, pr.ID, "open")

	return pr, nil
}

func (svc service) doPOST(src, dest models.Branch, ttl, desc string, draft bool,
	o origin.AzureOrigin, resp *Response) (err error) {

	url, err := svc.originSvc.CreatePRsURL(o.Remote)
	if err != nil {
		return fmt.Errorf("getting url: %w", err)
	}

	b64PAT := base64.RawStdEncoding.EncodeToString([]byte(svc.pat))
	bearer := fmt.Sprintf("Basic %s", b64PAT)

	pullRequest := map[string]interface{}{
		"sourceRefName": src.FullName(),  // Source branch
		"targetRefName": dest.FullName(), // Target branch
		"title":         ttl,             // Title of PR
		"isDraft":       draft,           // Draft PR
		"description":   desc,            // Description of PR
	}

	body, err := json.Marshal(pullRequest)
	if err != nil {
		return fmt.Errorf("marshalling request body: %w", err)
	}

	azReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("creating http request: %w", err)
	}

	azReq.Header.Add("Authorization", bearer)
	azReq.Header.Add("Content-Type", "application/json")
	azReq.Header.Add("Host", "dev.azure.com")
	azReq.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/117.0")
	azReq.Header.Add("Accept", "application/json;api-version=5.0-preview.1;excludeUrls=true;enumsAsNumbers=true;msDateFormat=true;noArrayWrap=true")
	azReq.Header.Add("Accept-Encoding", "gzip,deflate,br")
	azReq.Header.Add("Referer", fmt.Sprintf("https://dev.azure.com/%s/%s/_git/"+
		"%s/pullrequestcreate?sourceRef=%s&targetRef=%s"+
		"&sourceRepositoryId=%s&targetRepositoryId=%s", o.Organization(),
		o.Project(), o.Repository(), src.Name(), dest.Name(), o.Repository(),
		o.Repository()))
	azReq.Header.Add("Origin", "https://dev.azure.com")
	azReq.Header.Add("Connection", "keep-alive")
	azReq.Header.Add("Sec-Fetch-Dest", "empty")
	azReq.Header.Add("Sec-Fetch-Mode", "cors")
	azReq.Header.Add("Sec-Fetch-Site", "same-origin")

	client := &http.Client{}
	azResp, err := client.Do(azReq)
	if err != nil {
		return fmt.Errorf("creating PR via client: %w\nurl: %s\n%+v",
			err, url, pullRequest)
	}

	if azResp.StatusCode != http.StatusCreated {
		respBody, err := io.ReadAll(azResp.Body)
		if err != nil {
			respBody = []byte("cannot read response body content")
		}

		return fmt.Errorf("response code: %d\n"+
			"response body: %+v\n"+
			"pull request: %+v\n"+
			"url: %s\n"+
			"request: %+v",
			azResp.StatusCode, string(respBody), pullRequest, url, azReq)

	}

	defer azResp.Body.Close()
	return json.NewDecoder(azResp.Body).Decode(&resp)

}
