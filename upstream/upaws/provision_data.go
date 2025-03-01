// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package upaws

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/apps"
	"github.com/mattermost/mattermost-plugin-apps/apps/path"
	"github.com/mattermost/mattermost-plugin-apps/utils"
)

// DeployData contains all the necessary data for deploying an app.
type DeployData struct {
	// StaticFiles key is the name of the static file in the /static folder
	// Staticfiles value is the S3 Key where file should be deployed
	StaticFiles map[string]AssetData `json:"static_files"`

	// LambdaFunctions key is the name of the lambda function zip bundle
	// LambdaFunctions value contains info for deploying a function in the AWS.
	// LambdaFunctions value's Name field contains functions name in the AWS.
	LambdaFunctions map[string]FunctionData `json:"lambda_functions"`
	Manifest        *apps.Manifest          `json:"-"`
	ManifestKey     string                  `json:"manifest_key"`
}

type FunctionData struct {
	Bundle  io.ReadCloser `json:"-"`
	Name    string        `json:"name"`
	Handler string        `json:"handler"`
	Runtime string        `json:"runtime"`
}

type AssetData struct {
	File io.ReadCloser `json:"-"`
	Key  string        `json:"key"`
}

func GetDeployDataFromFile(path string, log utils.Logger) (*DeployData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read file from  path %s", path)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "can't read file")
	}

	return getDeployData(b, log)
}

// getDeployData takes app bundle zip as a byte slice and returns DeployData
func getDeployData(b []byte, log utils.Logger) (*DeployData, error) {
	bundleReader, bundleErr := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if bundleErr != nil {
		return nil, errors.Wrap(bundleErr, "can't get zip reader")
	}
	bundleFunctions := []FunctionData{}
	var m *apps.Manifest
	assets := []AssetData{}

	// Read all the files from zip archive
	for _, file := range bundleReader.File {
		switch {
		case strings.HasSuffix(file.Name, "manifest.json"):
			manifestFile, err := file.Open()
			if err != nil {
				return nil, errors.Wrap(err, "can't open manifest.json file")
			}
			defer manifestFile.Close()

			data, err := io.ReadAll(manifestFile)
			if err != nil {
				return nil, errors.Wrap(err, "can't read manifest.json file")
			}
			m, err = apps.DecodeCompatibleManifest(data)
			if err != nil {
				return nil, errors.Wrapf(err, "can't unmarshal manifest.json file %s", string(data))
			}
			log.Infow("Found manifest", "file", file.Name)

		case strings.HasSuffix(file.Name, ".zip"):
			lambdaFunctionFile, err := file.Open()
			if err != nil {
				return nil, errors.Wrapf(err, "can't open file %s", file.Name)
			}
			defer lambdaFunctionFile.Close()
			bundleFunctions = append(bundleFunctions, FunctionData{
				Name:   strings.TrimSuffix(file.Name, ".zip"),
				Bundle: lambdaFunctionFile,
			})
			log.Infow("Found lambda function bundle", "file", file.Name)

		case strings.HasPrefix(file.Name, path.StaticFolder+"/"):
			assetName := strings.TrimPrefix(file.Name, path.StaticFolder+"/")
			if assetName == "" {
				continue
			}
			assetFile, err := file.Open()
			if err != nil {
				return nil, errors.Wrapf(err, "can't open file %s", file.Name)
			}
			assets = append(assets, AssetData{
				Key:  assetName,
				File: assetFile,
			})
			log.Infow("Found static asset", "file", file.Name)

		default:
			log.Infow("Ignored unknown file", "file", file.Name)
		}
	}

	if m == nil {
		return nil, errors.New("no manifest found")
	}
	if !m.Contains(apps.DeployAWSLambda) {
		return nil, errors.New(`"aws_lambda" is not present n the manifest`)
	}

	resFunctions := []FunctionData{}

	// Matching bundle functions to the functions listed in manifest
	// O(n^2) code for simplicity
	for _, bundleFunction := range bundleFunctions {
		for _, manifestFunction := range m.AWSLambda.Functions {
			if strings.HasSuffix(bundleFunction.Name, manifestFunction.Name) {
				resFunctions = append(resFunctions, FunctionData{
					Bundle:  bundleFunction.Bundle,
					Name:    manifestFunction.Name,
					Handler: manifestFunction.Handler,
					Runtime: manifestFunction.Runtime,
				})
				continue
			}
		}
	}

	generatedAssets := generateAssetNames(m, assets)
	generatedFunctions := generateFunctionNames(m, resFunctions)

	pd := &DeployData{
		StaticFiles:     generatedAssets,
		LambdaFunctions: generatedFunctions,
		Manifest:        m,
		ManifestKey:     S3ManifestName(m.AppID, m.Version),
	}
	if err := pd.Validate(); err != nil {
		return nil, errors.Wrap(err, "deploy data is not valid")
	}
	return pd, nil
}

func generateAssetNames(manifest *apps.Manifest, assets []AssetData) map[string]AssetData {
	generatedAssets := make(map[string]AssetData, len(assets))
	for _, asset := range assets {
		generatedAssets[asset.Key] = AssetData{
			Key:  S3StaticName(manifest.AppID, manifest.Version, asset.Key),
			File: asset.File,
		}
	}
	return generatedAssets
}

func generateFunctionNames(manifest *apps.Manifest, functions []FunctionData) map[string]FunctionData {
	generatedFunctions := make(map[string]FunctionData, len(functions))
	for _, function := range functions {
		name := LambdaName(manifest.AppID, manifest.Version, function.Name)
		generatedFunctions[function.Name] = FunctionData{
			Name:    name,
			Bundle:  function.Bundle,
			Handler: function.Handler,
			Runtime: function.Runtime,
		}
	}
	return generatedFunctions
}

func (pd *DeployData) Validate() error {
	var result error

	if pd.Manifest == nil || !pd.Manifest.Contains(apps.DeployAWSLambda) {
		result = multierror.Append(result,
			errors.New("no manifest or AWS Lamda metadata"))
	}
	if err := pd.Manifest.Validate(); err != nil {
		return err
	}

	if len(pd.Manifest.AWSLambda.Functions) != len(pd.LambdaFunctions) {
		result = multierror.Append(result,
			errors.New("different number of functions in the manifest and in the bundle"))
	}

	for _, function := range pd.Manifest.AWSLambda.Functions {
		data, ok := pd.LambdaFunctions[function.Name]
		if !ok {
			result = multierror.Append(result,
				errors.Errorf("function %s was not found in the bundle", function))
		}
		if data.Handler != function.Handler {
			result = multierror.Append(result,
				errors.New("mismatched handler"))
		}
		if data.Runtime != function.Runtime {
			result = multierror.Append(result,
				errors.New("mismatched runtime"))
		}
	}

	return result
}
