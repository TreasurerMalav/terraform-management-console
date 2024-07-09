package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type saveData struct {
	Data string
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// call GetListOfApplications function for / endpoint
		apps := GetListOfApplications()
		// convert response to json and return
		for _, app := range apps {
			w.Write([]byte(app + "\n"))
		}
	})
	mux.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var data saveData
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			WriteToFile("/home/ec2-user/my-application-1/main.tf", data.Data)
			SetGitConfig()
			AddCommitAndPushToGit()
			w.Write([]byte("Data saved to file"))
		}
		if r.Method == http.MethodGet {
			data := GetDataFromFile("/home/ec2-user/my-application-1/main.tf")
			w.Write([]byte(data))
		}

	})
	mux.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			cmd := exec.Command("cd", "/home/ec2-user/my-application-1")
			_, err := cmd.Output()
			if err != nil {
				log.Fatal(err)
			}
			cmd = exec.Command("terraform", "init")
			out, err := cmd.Output()
			if err != nil {
				log.Fatal(err)
			}
			//log.Println(out)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(out))
		}

	})
	mux.HandleFunc("/plan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			cmd := exec.Command("cd", "/home/ec2-user/my-application-1")
			_, err := cmd.Output()
			if err != nil {
				log.Fatal(err)
			}
			cmd = exec.Command("terraform", "plan")
			out, err := cmd.Output()
			if err != nil {
				log.Fatal(err)
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(out))
		}

	})
	mux.HandleFunc("/apply", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			cmd := exec.Command("cd", "/home/ec2-user/my-application-1")
			_, err := cmd.Output()
			if err != nil {
				log.Fatal(err)
			}
			cmd = exec.Command("terraform", "apply", "-auto-approve")
			out, err := cmd.Output()
			if err != nil {
				log.Fatal(err)
			}
			// log.Println(out)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(out))
		}

	})

	handler := enableCors(mux)
	http.ListenAndServe(":8080", handler)
}

func WriteToFile(filepath string, data string) {
	var file *os.File

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		file, err = os.Create(filepath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	} else {
		file, err = os.OpenFile(filepath, os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err := file.WriteString(data)
	if err != nil {
		log.Fatal(err)
	}
}

func GetDataFromFile(filepath string) string {
	var file *os.File
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// create empty file
		_, err = os.Create(filepath)
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the content of the file
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the content to string and return
	return string(data)
}

func GetListOfApplications() []string {
	return []string{"my-application-1"}
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SetGitConfig() {
	println("Setting git config")
	cmd := exec.Command("cd", "/home/ec2-user/my-application-1")
	_, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	cmd = exec.Command("git", "config", "user.email", "malav.treasurer@gmail.com") // git config --global user.email
	_, err = cmd.Output()

	if err != nil {
		log.Fatal(err)
	}
	println("done setting git user email")
	cmd = exec.Command("git", "config", "user.name", "Malav Treasurer") // git config --global user.name
	_, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	println("done setting git user name")

}

func AddCommitAndPushToGit() {
	println("Adding, committing and pushing to git")
	cmd := exec.Command("cd", "/home/ec2-user/my-application-1")
	_, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	println("done cd")
	cmd = exec.Command("git", "add", "main.conf") // git add main.conf
	_, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	println("completed git add")
	cmd = exec.Command("git", "commit", "-m", "Update main.conf") // git commit -m "Added new data"
	_, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	println("completed git commit")

	cmd = exec.Command("git", "push") // git push
	_, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	println("completed git push")
}

