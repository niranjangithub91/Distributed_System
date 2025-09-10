package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
	controller_helper "user_entry/helper/Controller"
	db "user_entry/helper/Database"
	faulttolerance "user_entry/helper/Fault_tolerance"
	"user_entry/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/klauspost/reedsolomon"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var a model.Data
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, "Request error", http.StatusBadRequest)
		return
	}
	status := controller_helper.Validate_signup(a)
	if !status {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	var a1 model.User
	if name, ok := a["name"].(string); ok {
		a1.Name = name
	}
	if pass, ok := a["password"].(string); ok {
		a1.Password = pass
	}
	o := db.Add_User(a1)
	if !o {
		http.Error(w, "Data entry not done", http.StatusInternalServerError)
	}
}
func Login(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	var jwtKey = []byte(os.Getenv("Keys"))
	w.Header().Set("Content-Type", "application/json")
	var a model.Data
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Invalid", http.StatusBadRequest)
		return
	}
	var a1 model.User
	if name, ok := a["name"].(string); ok {
		a1.Name = name
	}
	if pass, ok := a["password"].(string); ok {
		a1.Password = pass
	}
	status := db.Find_User(a1)
	if !status {
		fmt.Println("Hi")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	expirationTime := time.Now().Add(time.Minute * 10)
	claims := &model.Claims{
		Data: a1.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Server fails", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
	return
}
func Upload(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	w.Header().Set("Content-Type", "application/json")
	claims, ok := r.Context().Value("claims").(*model.Claims)
	if !ok {
		http.Error(w, "No claims found", http.StatusUnauthorized)
		return
	}
	name := claims.Data
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File is larger then 10 Mb", http.StatusBadRequest)
		log.Fatal(err)
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fmt.Println(handler.Filename)
	fmt.Println(handler.Size)
	data, err := io.ReadAll(file)
	enc, err := reedsolomon.New(2, 1)
	if err != nil {
		http.Error(w, "Failed to create encoder", http.StatusInternalServerError)
		return
	}
	shards, err := enc.Split(data)
	if err != nil {
		http.Error(w, "Failed to split", http.StatusInternalServerError)
		return
	}
	err = enc.Encode(shards)
	if err != nil {
		http.Error(w, "Failed to encode", http.StatusInternalServerError)
		return
	}
	o := faulttolerance.Health_checker()
	var t []int
	for key, value := range o {
		if value {
			t = append(t, key)
		}
	}
	fmt.Print(name)
	length := len(t)
	counter := 0
	m1 := make(map[int]int)
	fmt.Println("Ykjnjkefjre")
	for n, shard := range shards {
		r := counter % length
		p := fmt.Sprintf("http://localhost:%d/data_receive", t[r])
		m1[n] = t[r]
		wg.Add(1) // increment first
		fmt.Println(p)
		go func(p string, shard []byte, c int, m map[int]int, counter int) {
			defer wg.Done()
			controller_helper.SendChunk(p, shard, handler.Filename, c, name, counter)
		}(p, shard, t[r], m1, counter) // pass variables explicitly
		counter++
	}
	wg.Wait()
	return
}

func Download(w http.ResponseWriter, r *http.Request) {
	// var mt sync.Mutex
	var wg sync.WaitGroup
	var a model.Data
	var mt sync.Mutex
	var total = make([][]byte, 3) // since you have 3 servers
	fmt.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, "Invalid decode", http.StatusBadRequest)
		return
	}

	file_name := a["file_name"].(string)
	claims, ok := r.Context().Value("claims").(*model.Claims)
	if !ok {
		http.Error(w, "No claims found", http.StatusUnauthorized)
		return
	}
	name := claims.Data
	counter := 0
	for i := 0; i < 3; i++ {
		p := fmt.Sprintf("%s_%s_%d", file_name, name, i)
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			q := db.Get_cache_mem(p)
			if q != nil {
				mt.Lock()
				total[i] = q
				counter++
				mt.Unlock()
			}
		}(p)
	}
	wg.Wait()

	if counter < 2 {
		z := db.Collection_retreive_details(file_name, name)

		for i := 3001; i < 3004; i++ {
			wg.Add(1)
			p := fmt.Sprintf("http://localhost:%d/data_retreive", i)
			go func(p string, file_name string, name string, port int) {
				defer wg.Done()
				_, resp := controller_helper.Send_Chunk_request(p, file_name, name, port, z)
				mt.Lock()
				for key, value := range resp {
					total[key] = value
				}
				mt.Unlock()

			}(p, file_name, name, i)
		}
		wg.Wait()
	}

	// Reed-Solomon decoder (2 data, 1 parity)
	enc, err := reedsolomon.New(2, 1)
	if err != nil {
		http.Error(w, "Failed to create decoder", http.StatusInternalServerError)
		return
	}

	// Reconstruct missing shards if needed
	if ok, _ := enc.Verify(total); !ok {
		if err := enc.Reconstruct(total); err != nil {
			http.Error(w, "Failed to reconstruct shards", http.StatusInternalServerError)
			return
		}
	}

	for key, value := range total {
		var wg sync.WaitGroup
		p := fmt.Sprintf("%s_%s_%d", file_name, name, key)
		wg.Add(1)
		go func(p string, x []byte) {
			defer wg.Done()
			db.Set_cache_mem(p, x)
		}(p, value)
	}
	wg.Wait()

	originalSize := len(total[0]) + len(total[1])

	buf := make([]byte, originalSize)

	copy(buf[:len(total[0])], total[0])
	copy(buf[len(total[0]):], total[1])

	data := buf

	w.Header().Set("Content-Disposition", "attachment; filename="+file_name)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(data)
}
