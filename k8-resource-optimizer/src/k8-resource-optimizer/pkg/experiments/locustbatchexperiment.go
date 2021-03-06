package experiments

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"k8-resource-optimizer/pkg/locustwrap"
	"k8-resource-optimizer/pkg/models"
)

//LocustBatchExperiment : Locust experiment for all tenants of all namespaces
type LocustBatchExperiment struct {
	Iteration    int
	Sample       int
	Locustconfig locustwrap.Config
	OutputDir    string
	Type         string
}

const experimentDurationInSeconds = 60

func locustBatchFromNamespaces(slas []models.SLA) models.Experiment {
	// var e models.Experiment
	e := &LocustBatchExperiment{}
	e.Locustconfig.DurationInSeconds = experimentDurationInSeconds
	users := []locustwrap.User{}
	for _, sla := range slas {
		for u := 0; u < sla.NbOfTenants; u++ {
			jobsize, err := sla.GetSLO("jobsize")
			if err != nil {
				log.Printf("LocustBatchFromNamespaces: %v", err)
				panic(err)
			}
			usr := createLocustBatchExperimentParameters(sla, jobsize.(int), u)
			users = append(users, usr)
			e.Locustconfig.Users++
			e.Locustconfig.HatchRate++
		}
	}
	// create script
	e.Locustconfig.ScriptPath = locustwrap.CreateRunScript(users, "locust_exp")

	e.Type = "LocustBatchExperiment"

	return e
}

func createLocustBatchExperimentParameters(s models.SLA, jobsize int, tenantNb int) locustwrap.User {
	user := locustwrap.User{}
	user.TenantID = s.Name + strconv.Itoa(tenantNb)
	user.Name = s.Name + "-" + strconv.Itoa(tenantNb)
	user.URL = "http://demo." + s.Name + ".svc.cluster.local:80/pushJob/" + strconv.Itoa(jobsize)
	user.Amount = 1
	user.MaxWait = 0
	user.MinWait = 0
	return user
}

//Type returns the type of the experiment
func (lbe *LocustBatchExperiment) GetType() string {
	return lbe.Type
}

// SetOutputDir set the output directory for raw locust experiment results
func (lbe *LocustBatchExperiment) SetOutputDir(directory string) {
	if !strings.HasSuffix(directory, "/") {
		directory = directory + "/"
	}
	lbe.OutputDir = directory
}

// SetIterationAndSample sets the current iteration and sample being evaluated
func (lbe *LocustBatchExperiment) SetIterationAndSample(iter int, sample int) {
	lbe.Iteration = iter
	lbe.Sample = sample
}

//Run : runs the experiment
func (lbe *LocustBatchExperiment) Run() (models.ExperimentResult, error) {
	results := lbe.Locustconfig.Run(fmt.Sprintf("%v%v-%v", lbe.OutputDir, lbe.Iteration, lbe.Sample))
	result := CreateLocustBarchExperimentResult(results)
	return result, nil
}
