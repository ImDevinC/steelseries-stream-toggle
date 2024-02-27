package steelseries

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/imdevinc/steelseries-stream-toggle/internals/models"
)

const ModeClassic string = "classic"
const ModeStream string = "stream"

type steelseries struct {
	httpClient *http.Client
}

func New(client *http.Client) *steelseries {
	return &steelseries{
		httpClient: client,
	}
}

func getCoreProps() (models.CoreProps, error) {
	props, err := os.ReadFile("C:\\ProgramData\\SteelSeries\\GG\\coreProps.json")
	if err != nil {
		return models.CoreProps{}, fmt.Errorf("could not read coreProps.json: %w", err)
	}
	var coreProps models.CoreProps
	err = json.Unmarshal(props, &coreProps)
	if err != nil {
		return models.CoreProps{}, fmt.Errorf("could not unmarshal coreProps.json: %w", err)
	}
	return coreProps, nil
}

func (s *steelseries) getSubApps() (models.SubApps, error) {
	coreProps, err := getCoreProps()
	if err != nil {
		return models.SubApps{}, fmt.Errorf("could not get coreProps: %w", err)
	}
	rawSubApps, err := s.httpClient.Get("https://" + coreProps.GgEncryptedAddress + "/subApps")
	if err != nil {
		return models.SubApps{}, fmt.Errorf("could not get subApps: %w", err)
	}
	var subApps models.SubApps
	err = json.NewDecoder(rawSubApps.Body).Decode(&subApps)
	if err != nil {
		return models.SubApps{}, fmt.Errorf("could not decode subApps: %w", err)
	}
	return subApps, nil
}

func (s *steelseries) getWebAddress() (string, error) {
	subApps, err := s.getSubApps()
	if err != nil {
		return "", fmt.Errorf("could not get subApps: %w", err)
	}
	return subApps.SubApps.Sonar.Metadata.WebServerAddress, nil
}

func (s *steelseries) getCurrentSonarMode() (string, error) {
	webAddress, err := s.getWebAddress()
	if err != nil {
		return "", fmt.Errorf("could not get web address: %w", err)
	}
	modeResp, err := s.httpClient.Get(webAddress + "/mode")
	if err != nil {
		return "", fmt.Errorf("could not get mode: %w", err)
	}
	rawMode, err := io.ReadAll(modeResp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read mode: %w", err)
	}
	mode := strings.Trim(string(rawMode), "\"")
	return mode, nil
}

func (s *steelseries) setMode(mode string) error {
	webAddress, err := s.getWebAddress()
	if err != nil {
		return fmt.Errorf("could not get web address: %w", err)
	}
	req, err := http.NewRequest(http.MethodPut, webAddress+"/mode/"+mode, nil)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	_, err = s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not set mode: %w", err)
	}
	return nil
}

func (s *steelseries) ToggleStreamerMode() (string, error) {
	mode, err := s.getCurrentSonarMode()
	if err != nil {
		return "", fmt.Errorf("could not get current mode: %w", err)
	}
	newMode := ""
	switch mode {
	case ModeClassic:
		newMode = ModeStream
	case ModeStream:
		newMode = ModeClassic
	}
	if newMode == "" {
		return "", fmt.Errorf("could not determine new mode from current mode: %s", mode)
	}
	err = s.setMode(newMode)
	if err != nil {
		return "", fmt.Errorf("could not set mode: %w", err)
	}
	return newMode, nil
}
