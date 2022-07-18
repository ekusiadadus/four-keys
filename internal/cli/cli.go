package cli

import (
	"encoding/json"
	"time"

	"github.com/hmiyado/four-keys/internal/releases"
	"github.com/urfave/cli/v2"
)

func DefaultApp() *cli.App {
	return &cli.App{
		Name:   "four-keys",
		Usage:  "analyze four keys metrics",
		Flags:  getCommandReleasesFlags(),
		Action: defaultAction,
		Commands: []*cli.Command{
			GetCommandReleases(),
		},
	}
}

type DefaultCliOutput struct {
	Option                *releases.Option     `json:"option"`
	DeploymentFrequency   float64              `json:"deploymentFrequency"`
	LeadTimeForChanges    DurationWithTimeUnit `json:"leadTimeForChanges"`
	TimeToRestoreServices DurationWithTimeUnit `json:"timeToRestoreServices"`
	ChangeFailureRate     float64              `json:"changeFailureRate"`
}

func defaultAction(ctx *cli.Context) error {
	context := &CliContextWrapper{context: ctx}
	releases, err := QueryReleases(context)
	if err != nil {
		context.Error(err)
		return err
	}
	option, err := context.Option()
	if err != nil {
		context.Error(err)
		return err
	}

	duration := option.Until.Sub(option.Since)
	daysCount := int(duration.Hours() / 24)
	releasesCount := len(releases)
	deploymentFrequency := float64(releasesCount) / float64(daysCount)

	outputJson, err := json.Marshal(&DefaultCliOutput{
		Option:                option,
		DeploymentFrequency:   deploymentFrequency,
		LeadTimeForChanges:    getDurationWithTimeUnit(getMeanLeadTimeForChanges(releases)),
		TimeToRestoreServices: getDurationWithTimeUnit(getTimeToRestoreServices(releases)),
		ChangeFailureRate:     getChangeFailureRate(releases),
	})
	if err != nil {
		context.Error(err)
		return err
	}
	context.Write(outputJson)
	return nil

}

func getMeanLeadTimeForChanges(release []*releases.Release) time.Duration {
	if len(release) == 0 {
		return time.Duration(0)
	}
	sum := time.Duration(0)
	for _, release := range release {
		sum = release.LeadTimeForChanges + sum
	}
	return time.Duration(int64(sum) / int64(len(release)))
}

func getTimeToRestoreServices(releases []*releases.Release) time.Duration {
	sum := time.Duration(0)
	countOfRestoreService := 0
	failedReleaseIndex := -1
	for i := len(releases) - 1; i >= 0; i-- {
		release := releases[i]
		if !release.Result.IsSuccess {
			if failedReleaseIndex < 0 {
				failedReleaseIndex = i
			}
			continue
		}
		if release.Result.IsSuccess && failedReleaseIndex < 0 {
			continue
		}
		sum += release.Date.Sub(releases[failedReleaseIndex].Date)
		countOfRestoreService += 1
		failedReleaseIndex = -1
	}
	if countOfRestoreService == 0 {
		return sum
	}
	return sum / time.Duration(countOfRestoreService)
}

func getChangeFailureRate(releases []*releases.Release) float64 {
	if len(releases) == 0 {
		return 0
	}

	sumOfFailure := 0
	for _, release := range releases {
		if !release.Result.IsSuccess {
			sumOfFailure += 1
		}
	}
	return float64(sumOfFailure) / float64(len(releases))
}
