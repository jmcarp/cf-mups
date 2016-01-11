package main

import (
        "os"
        "fmt"
        "errors"
        "strings"
        "encoding/json"

        "github.com/cloudfoundry/cli/plugin"
)

type MupsPlugin struct {}

func getServiceApp(conn plugin.CliConnection, name string) (string, error) {
        services, err := conn.GetServices()
        if err != nil {
              return "", err
        }
        for _, service := range services {
                if service.Name == name {
                        return service.ApplicationNames[0], nil
                }
        }
        return "", errors.New(fmt.Sprintf("No service %s bound to an app", name))
}

func getServiceCreds(conn plugin.CliConnection, name string) (map[string]interface{}, error) {
        app, err := getServiceApp(conn, name)
        if err != nil {
                return nil, err
        }

        env, err := conn.CliCommandWithoutTerminalOutput("env", app)
        if err != nil {
                return nil, err
        }

        var raw string;
        for _, line := range env {
                if strings.Contains(line, "VCAP_SERVICES") {
                        raw = line
                        break
                }
        }
        if raw == "" {
                return nil, errors.New("VCAP_SERVICES not found in env")
        }
        var unmarshaled map[string]interface{}
        err = json.Unmarshal([]byte(raw), &unmarshaled)
        if err != nil {
                return nil, err
        }

        vcapServices := unmarshaled["VCAP_SERVICES"].(map[string]interface{})
        userProvided := vcapServices["user-provided"].([]interface{})
        var credentials map[string]interface{}
        for _, service := range userProvided {
                service := service.(map[string]interface{})
                if service["name"] == name {
                        credentials = service["credentials"].(map[string]interface{})
                        break
                }
        }
        return credentials, nil
}

func setServiceCreds(conn plugin.CliConnection, name string, creds map[string]interface{}) error {
        marshaled, err := json.Marshal(creds)
        if err != nil {
                return err
        }

        _, err = conn.CliCommandWithoutTerminalOutput("uups", name, "-p", string(marshaled[:]))
        return err
}

func (p *MupsPlugin) Run(cliConnection plugin.CliConnection, args []string) {
        if len(args) > 0 && args[0] == "CLI-MESSAGE-UNINSTALL" {
              return
        }

        serviceName := args[1]
        credName := args[2]

        var credValue interface{}
        err := json.Unmarshal([]byte(args[3]), &credValue)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

        creds, err := getServiceCreds(cliConnection, serviceName)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

        creds[credName] = credValue

        err = setServiceCreds(cliConnection, serviceName, creds)
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
}

func (p *MupsPlugin) GetMetadata() plugin.PluginMetadata {
        return plugin.PluginMetadata{
                Name: "modify-user-provided-service",
                Version: plugin.VersionType{
                        Major: 1,
                        Minor: 0,
                        Build: 0,
                },
                MinCliVersion: plugin.VersionType{
                        Major: 6,
                        Minor: 7,
                        Build: 0,
                },
                Commands: []plugin.Command{
                        plugin.Command{
                                Name: "modify-user-provided-service",
                                Alias: "mups",
                                HelpText: "Set a credential within a user-provided service",
                                UsageDetails: plugin.Usage{
                                        Usage: "mups SERVICE_NAME CREDENTIAL_NAME CREDENTIAL_VALUE",
                                },
                        },
                },
        }
}

func main() {
        plugin.Start(new(MupsPlugin))
}
