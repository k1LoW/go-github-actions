package attest

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/known/structpb"
)

const (
	slsaPredicateType = "https://slsa.dev/provenance/v1"
	buildType         = "https://actions.github.io/buildtypes/workflow/v1"
)

// buildProvenancePredicate builds the same predicate as buildSLSAProvenancePredicate of @actions/attest.
func buildProvenancePredicate(c *claims, serverURL string) (*structpb.Struct, error) {
	// workflow_ref is "{owner}/{repo}/{path}@{ref}".
	rest, ok := strings.CutPrefix(c.WorkflowRef, c.Repository+"/")
	if !ok {
		return nil, fmt.Errorf("workflow_ref claim does not belong to repository %s: %s", c.Repository, c.WorkflowRef)
	}
	path, _, ok := strings.Cut(rest, "@")
	if !ok {
		return nil, fmt.Errorf("invalid workflow_ref claim: %s", c.WorkflowRef)
	}
	return structpb.NewStruct(map[string]any{
		"buildDefinition": map[string]any{
			"buildType": buildType,
			"externalParameters": map[string]any{
				"workflow": map[string]any{
					// Not the ref of workflow_ref, which differs from the ref of the run for events on ref-less commits.
					"ref":        c.Ref,
					"repository": fmt.Sprintf("%s/%s", serverURL, c.Repository),
					"path":       path,
				},
			},
			"internalParameters": map[string]any{
				"github": map[string]any{
					"event_name":          c.EventName,
					"repository_id":       c.RepositoryID,
					"repository_owner_id": c.RepositoryOwnerID,
					"runner_environment":  c.RunnerEnvironment,
				},
			},
			"resolvedDependencies": []any{
				map[string]any{
					"uri": fmt.Sprintf("git+%s/%s@%s", serverURL, c.Repository, c.Ref),
					"digest": map[string]any{
						"gitCommit": c.SHA,
					},
				},
			},
		},
		"runDetails": map[string]any{
			"builder": map[string]any{
				"id": fmt.Sprintf("%s/%s", serverURL, c.JobWorkflowRef),
			},
			"metadata": map[string]any{
				"invocationId": fmt.Sprintf("%s/%s/actions/runs/%s/attempts/%s", serverURL, c.Repository, c.RunID, c.RunAttempt),
			},
		},
	})
}
