package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jcelliott/lumber"
)

const Version = "1.0.0"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)

	opts := Options{}

	if options != nil {
		opts = *options
	}

	if opts.Logger != nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     options.Logger,
	}

	if _, err := os.Stat(dir); err != nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		fmt.Println("Error: ", err)
		return &driver, nil
	}

	opts.Logger.Debug("Creating the database at '%s'....\n", dir)
	return &driver, os.MkdirAll(dir, 0755)

}

func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing collection - No place to save the records")
	}

	if resource == "" {
		return fmt.Errorf("Missing resource - unable to save the record( no name)")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(d.dir, resource+".json")
	tempPath := fnlPath + ".tmp"

	if err := os.Mkdir(dir, 0755); err != nill {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")

	b = append(b, byte('\n'))

	if err := ioutil.WriteFile(tempPath, b, 0644); err != nil {
		return err
	}
}

func (d *Driver) Read() error {

}

func (d *Driver) ReadAll() error {

}

func (d *Driver) Delete() error {

}

func (d *Driver) getOrCreateMutex() *sync.Mutex {

}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err := os.Stat(path); os.IsNotExist(err) {
		fi, err := os.Stat(path + ".json")
	}
	return
}

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

func main() {
	dir := "./"

	db, err := New(dir)

	if err != nil {
		fmt.Println("Error creating db: ", err)
	}

	rand.Seed(time.Now().UnixNano())

	employees := []User{
		{"John", randomAge(), randomPhone(), "Google", Address{"Bangalore", "Karnataka", "India", "560054"}},
		{"Paul", randomAge(), randomPhone(), "Facebook", Address{"LA", "California", "US", "90001"}},
		{"Joe", randomAge(), randomPhone(), "Amazon", Address{"Noida", "UP", "India", "201301"}},
		{"Ram", randomAge(), randomPhone(), "Adobe", Address{"London", "London", "UK", "EC1A1BB"}},
		{"Sam", randomAge(), randomPhone(), "Flipkart", Address{"Bangalore", "Karnataka", "India", "560034"}},
		{"Peter", randomAge(), randomPhone(), "Netflix", Address{"Los Gatos", "California", "US", "95030"}},
		{"Alice", randomAge(), randomPhone(), "Tesla", Address{"Austin", "Texas", "US", "73301"}},
		{"Bob", randomAge(), randomPhone(), "Microsoft", Address{"Redmond", "Washington", "US", "98052"}},
		{"Charlie", randomAge(), randomPhone(), "Spotify", Address{"Stockholm", "Stockholm", "Sweden", "11122"}},
		{"Diana", randomAge(), randomPhone(), "Uber", Address{"San Francisco", "California", "US", "94103"}},
	}

	for _, value := range employees {
		db.Write("users", value.Name, User{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}

	records, err := db.readAll("users")
	if err != nil {
		fmt.Println("error reading users", err)
	}

	fmt.Println(records)

	allusers := []User{}

	for _, f := range records {
		employeeFounded := User{}
		if err := json.Unmarshal([]byte(f), &employeeFounded); err != nil {
			fmt.Println("error", err)
		}

		allusers = append(allusers, employeeFounded)
	}

	fmt.Println(allusers)

	// if err := db.Delete("user" , "john"); err != nil {
	// 	fmt.Println("error deleting user", err)
	// }

	// if err := db.Delete("user" , ""); err != nil {
	// 	fmt.Println("error deleting all user", err)
	// }

}

func randomPhone() string {
	return fmt.Sprintf("952%07d", rand.Intn(10000000))
}

func randomAge() json.Number {
	return json.Number(fmt.Sprintf("%d", rand.Intn(19)+22))
}
