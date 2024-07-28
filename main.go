package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type saveData struct {
	Data string `json:"data"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// call GetListOfApplications function for / endpoint
		apps := GetListOfApplications()
		// convert response to json and return
		for _, app := range apps {
			w.Write([]byte(app + "\n"))
		}
	})

	mux.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodPost {
			var data saveData
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			WriteToFile("/terraform/my-application-1/main.tf", data.Data)
			SetGitConfig()
			AddCommitAndPushToGit()
			w.Write([]byte("Data saved to file"))
		} else if r.Method == http.MethodGet {
			data := GetDataFromFile("/terraform/my-application-1/main.tf")
			w.Write([]byte(data))
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodPost {
			commands := []string{
				"docker",
				"run",
				"-w",
				"/terraform/my-application-1",
				"-v",
				"/home/ec2-user/.aws:/root/.aws",
				"-v",
				"/terraform/my-application-1:/terraform/my-application-1",
				"docker-terraform:latest",
				"init",
			}
			cmd := exec.Command("sh", "-c", strings.Join(commands, " "))
			out, err := cmd.Output()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println(err)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write(out)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/plan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodPost {
                        commands := []string{
                                "docker",
                                "run",
                                "-w",
                                "/terraform/my-application-1",
                                "-v",
                                "/home/ec2-user/.aws:/root/.aws",
                                "-v",
                                "/terraform/my-application-1:/terraform/my-application-1",
                                "docker-terraform:latest",
                                "plan",
                        }
                        cmd := exec.Command("sh", "-c", strings.Join(commands, " "))
			out, err := cmd.Output()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write(out)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/apply", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodPost {
                        commands := []string{
                                "docker",
                                "run",
                                "-w",
                                "/terraform/my-application-1",
                                "-v",
                                "/home/ec2-user/.aws:/root/.aws",
                                "-v",
                                "/terraform/my-application-1:/terraform/my-application-1",
                                "docker-terraform:latest",
                                "apply",
				"-auto-approve",
                        }
                        cmd := exec.Command("sh", "-c", strings.Join(commands, " "))
			out, err := cmd.Output()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write(out)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	handler := enableCors(mux)
	log.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

func WriteToFile(filepath string, data string) {
	var file *os.File
	var err error

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		file, err = os.Create(filepath)
		if err != nil {
			log.Printf("Error creating file: %s\n", err.Error())
			return
		}
		defer file.Close()
	} else {
		file, err = os.OpenFile(filepath, os.O_RDWR|os.O_TRUNC, 0644)
		if err != nil {
			log.Printf("Error opening file: %s\n", err.Error())
			return
		}
		defer file.Close()
	}
	_, err = file.WriteString(data)
	if err != nil {
		log.Printf("Error writing to file: %s\n", err.Error())
	}
}

func GetDataFromFile(filepath string) string {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		file, err := os.Create(filepath)
		if err != nil {
			log.Printf("Error creating file: %s\n", err.Error())
			return ""
		}
		file.Close()
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("Error reading file: %s\n", err.Error())
		return ""
	}

	return string(data)
}

func GetListOfApplications() []string {
	return []string{"my-application-1"}
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SetGitConfig() {
	cmd := exec.Command("sh", "-c", "cd /terraform/my-application-1 && git config user.email 'malav.treasurer@gmail.com' && git config user.name 'Malav Treasurer'")
	if _, err := cmd.Output(); err != nil {
		if strings.Contains("Everything up-to-date", err.Error()) {
			log.Printf("No changes to be pushed.")
		} else {
			log.Printf("Error setting git config: %s\n", err.Error())
		}
	}
}

func AddCommitAndPushToGit() {
	cmd := exec.Command("sh", "-c", "cd /terraform/my-application-1 && git add main.tf && git commit -m 'Update main.tf' && git push")
	if _, err := cmd.Output(); err != nil {
		log.Printf("Error adding, committing, or pushing to git: %s\n", err.Error())
	}
}
