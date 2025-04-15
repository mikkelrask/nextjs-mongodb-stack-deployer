package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
  "path/filepath"
	"syscall"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)


var cwd = os.Getenv("PWD")
var dirName = filepath.Base(cwd)
var frontendContainer = dirName+"-frontend-1"
var mongoContainer = dirName+"-mongodb-1"

var envVars = map[string]string{
	"NEXTJS_PORT":   "3000",
	"MONGO_PORT":    "27017",
	"MONGO_DB_NAME": "production",
	"MONGO_USER":    "admin",
	"MONGO_PASS":    "",
	"DUMP_PATH":  	 "/path/to/dump",
	"REPO_URL":      "",
	"REPO_PRIVATE":  "No",
	"GIT_USER":      "",
	"GIT_PASS":      "",
}
var args = map[string]string{
	"down": 		"Stopping the running Docker containers and removing the volumes",
	"rebuild": 		"Rebuilding the Docker image with current configuration",
	"redeploy": 	"Stopping running container and rebuilding image before re-starting the stack",
	"reset": 		"Deleting configuration and starting fresh",
	"restart": 		"Restarting the Docker containers",
	"start": 		"Starting the Docker containers",
	"up": 			"Starting the Docker containers",
  "logs":     "Tail the logs for the frondend container, the mongodb container or both",
}

func deployerHelp() {
	fmt.Println("🚀 DPLOY - BUILD YOUR WEBSITE")
	fmt.Println("\nUsage: deployer [optionable action]")
	fmt.Println("\nActions:")
	fmt.Printf(" deployer: %s\n", "Run without any arguments to start the setup process")
	for arg, message := range args {
		fmt.Printf(" deployer %s: %s\n", arg, message)
	}
	fmt.Println("\nExamples:")
	fmt.Println("  deployer")
	fmt.Println("  deployer redeploy")
}

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if message, exists := args[arg]; exists {
			fmt.Printf("\n%s\n", message)
			switch arg {
			case "rebuild":
				stopContainers()
				buildImage()
				startContainers()
			case "redeploy":
				stopContainers()
				buildImage()
				startContainers()
			case "stop":
				stopContainers()
			case "down":
				stopContainers()
			case "up", "start":
				startContainers()
			case "restart":
				stopContainers()
				startContainers()
			case "logs":
				container := askWhatContainer()
				showLogs(container)
			case "attach":
				attach()
			case "reset":
				fmt.Println("\n🔄 Resetting configuration...")
				os.Remove(".env")
				fmt.Println("✅ Configuration reset. Exiting...")
			default:
				fmt.Printf("⚠️  Unknown action: %s\n", arg)
			}
			os.Exit(0)
		} else {
			fmt.Printf("⚠️  Unknown argument: %s\n", arg)
			deployerHelp()
			os.Exit(1)
		}
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChannel
		fmt.Println("\n👋️ Interrupted. Cleaning up and exiting...")
		os.Exit(0)
	}()

	fmt.Println("🚀 DPLOY - BUILD YOUR WEBSITE")
	runChecks()
	println("")
	fmt.Println("Fill the needed data to deploy your webapp - Enter keeps the suggested default\n")

	envVars["REPO_URL"] = askQuestion("NextJS webapp repository URL:", envVars["REPO_URL"])
	envVars["REPO_PRIVATE"] = askRepoPrivate()

	if envVars["REPO_PRIVATE"] == "Yes" {
		fmt.Println("🔒 Enter your Github credentials to clone the repo:")
		envVars["GIT_USER"] = askQuestion("Username: ", "")
		var gitPassword string
		prompt := &survey.Password{
			Message: "Access token or password: ",
		}
		err := survey.AskOne(prompt, &gitPassword)
		if err != nil {
			handleSurveyError(err)
		}
		envVars["GIT_PASS"] = gitPassword
	}
	envVars["NEXTJS_PORT"] = askQuestion("Do you want to specify a port for NextJS? ", envVars["NEXTJS_PORT"])

	if askForDataBaseDump() == "Yes" {
		dumpPath := askQuestion("Database dump path: ", envVars["DUMP_PATH"])
		copyDumpToDumpDirectory(dumpPath)
		envVars["DUMP_PATH"] = dumpPath
	}
	envVars["MONGO_DB_NAME"] = askQuestion("MongoDB Database name?", envVars["MONGO_DB_NAME"])
	envVars["MONGO_USER"] = askQuestion("MongoDB Username: ", envVars["MONGO_USER"])
	envVars["MONGO_PASS"] = askQuestion("MongoDB Password: ", envVars["MONGO_PASS"])
	envVars["MONGO_PORT"] = askQuestion("Do you want to specify a port for MongoDB? ", envVars["MONGO_PORT"])

	var saveResponse string
	prompt := &survey.Select{
		Message: "Do you want to save the configuration to `.env`?",
		Options: []string{"Yes", "No"},
		Default: "Yes", 
	}
	err := survey.AskOne(prompt, &saveResponse)
	if err != nil {
		handleSurveyError(err)
	}

	if saveResponse == "Yes" {
		saveEnv()
	} else {
		fmt.Println("\n👋 Configuration was not saved - Exiting ...")
		os.Exit(0)
	}

	containerStatus, err := exec.Command("docker", "ps", "--format", "{{.Names}}").Output()
	if err != nil {
		fmt.Println("Error checking Docker containers:", err)
		return
	}

	containerNames := strings.Split(string(containerStatus), "\n")
	expectedContainerName := "-frontend-1"

	for _, name := range containerNames {
		if strings.Contains(name, expectedContainerName) {
			var stopResponse string
			prompt := &survey.Select{
				Message: "The Docker containers are already running. Do you want to stop them before continuing?",
				Options: []string{"Yes", "No"},
				Default: "No",
			}
			err := survey.AskOne(prompt, &stopResponse)
			if err != nil {
				handleSurveyError(err)
			}
			if stopResponse == "Yes" {
				stopContainers()
			}
			break
		}
	}

	buildResponse := askToBuildImage()
	if buildResponse == "No" {
		os.Exit(0)
	}
	upResponse := askToUpContainers()

	if buildResponse == "Yes" {
		buildImage()
	}

	if upResponse == "Yes" {
		startContainers()
	}
}

func handleSurveyError(err error) {
	if err == terminal.InterruptErr {
		fmt.Println("\n👋 Interrupted. Exiting...")
		os.Exit(0)
	}
	fmt.Println("Error:", err)
	os.Exit(1)
}

func askQuestion(question, defaultValue string) string {
	var response string
	prompt := &survey.Input{
		Message: fmt.Sprintf("%s", question),
		Default: defaultValue,
	}
	err := survey.AskOne(prompt, &response)
	if err != nil {
		handleSurveyError(err)
	}
	if response == "" {
		response = defaultValue
	}
	return response
}

func askForDataBaseDump() string {
	var response string
	prompt := &survey.Select{
		Message: "Do you want to import a database dump?",
		Options: []string{"Yes", "No"},
		Default: "No",
	}
	err := survey.AskOne(prompt, &response)
	if err != nil {
		handleSurveyError(err)
	}
	return response
}

func askRepoPrivate() string {
	var response string
	prompt := &survey.Select{
		Message: "Is this repository private?",
		Options: []string{"Yes", "No"},
		Default: "No",
	}
	err := survey.AskOne(prompt, &response)
	if err != nil {
		handleSurveyError(err)
	}
	return response
}

func askToBuildImage() string {
	var response string
	prompt := &survey.Select{
		Message: "Do you want to build the Docker image?",
		Options: []string{"Yes", "No"},
		Default: "Yes",
	}
	err := survey.AskOne(prompt, &response)
	if err != nil {
		handleSurveyError(err)
	}
	return response
}

func askToUpContainers() string {
	var response string
	prompt := &survey.Select{
		Message: "Do you want to start the Docker containers?",
		Options: []string{"Yes", "No"},
		Default: "Yes",
	}
	err := survey.AskOne(prompt, &response)
	if err != nil {
		handleSurveyError(err)
	}
	return response
}

func runChecks() {
	_, err := os.Stat(".env")
	if err != nil {
		if os.IsNotExist(err) {
			// .env file does not exist, continue
		} else {
			fmt.Printf("⚠️  Error checking .env file: %v\n", err)
		}
	} else {
		fmt.Println("⚠️  .env file already exists. Continue to overwrite current configuration - Press Ctrl+C to cancel.\n")
	}

	if err == nil {
		for key := range envVars {
			if value := checkEnv(key); value != "" {
				envVars[key] = value
			}
		}
	}
}

func checkEnv(key string) string {
	file, err := os.Open(".env")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && parts[0] == key {
			return parts[1]
		}
	}
	return ""
}

func copyDumpToDumpDirectory(dumpPath string) {
	fmt.Println("\n📁 Copying the database dump to the dump directory...")
	copyCommand := fmt.Sprintf("cp -r %s/* mongo-dump/", dumpPath)
	runCommand(copyCommand)
}

func runCommand(command string) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}

func saveEnv() {
	file, err := os.Create(".env")
	if err != nil {
		fmt.Println("Error creating .env file:", err)
		return
	}
	defer file.Close()

	for key, value := range envVars {
		file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}

  fmt.Println("💾 Configuration saved as `.env`")
}

func buildImage() {
	fmt.Println("\n🔨 Building the Docker image...")
	buildCommand := fmt.Sprintf("docker compose build --no-cache")
	runCommand(buildCommand)
}

func startContainers() {
	fmt.Println("\n🚀 Starting the Docker containers...")
	upCommand := "docker compose up -d --force-recreate"
	runCommand(upCommand)
}

func stopContainers() {
	fmt.Println("\n🛑 Stopping the Docker containers...")
	stopCommand := "docker compose down --volumes"
	runCommand(stopCommand)
}

func showLogs(container string) {
  fmt.Printf("\n👨‍💻 Logs for %s\n", container)
  logsCommand := "" // Declare the variable before switch
  switch container {
  case "NextJS":
    logsCommand = "docker logs -f " + frontendContainer
  case "MongoDB":
    logsCommand = "docker logs -f " + mongoContainer
  case "":
    logsCommand = "docker compose logs -f"
  default:
    fmt.Println("❌ Unknown container:", container)
    return
  }
  runCommand(logsCommand)
}

func askWhatContainer() string{
	var response string
	prompt := &survey.Select{
		Message: "What container?",
		Options: []string{"NextJS", "MongoDB", ""},
		Default: "",
	}
	err := survey.AskOne(prompt, &response)
	if err != nil {
		handleSurveyError(err)
	}
	return response
}
