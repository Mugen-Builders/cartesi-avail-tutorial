package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"dapp/rollups"
)

var (
	infolog = log.New(os.Stderr, "[ info ]  ", log.Lshortfile)
	errlog  = log.New(os.Stderr, "[ error ] ", log.Lshortfile)
)

// Circle defines the structure for student data
type Circle struct {
	Name           string
	Age            int
	WalletAddress  string
	AttendanceCount int
}

// AdvancePayload defines the structure of the payload received in advance requests
type AdvancePayload struct {
	Method  string `json:"method"`
	Payload string `json:"payload"`
}

// CreateStudentPayload defines the payload for creating a new student
type CreateStudentPayload struct {
	Name          string `json:"name"`
	Age           int    `json:"age"`
	WalletAddress string `json:"wallet_address"`
}

func (c *Circle) TakeAttendance() {
	c.AttendanceCount++
	fmt.Printf("%s has taken attendance. Attendance count: %d\n", c.Name, c.AttendanceCount)
}

var allStudents []Circle
var mu sync.Mutex

func HandleAdvance(data *rollups.AdvanceResponse) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("HandleAdvance: failed to marshal JSON: %w", err)
	}
	infolog.Println("Received advance request data:", string(dataBytes))

	// Process payload
	payloadStr, err := rollups.Hex2Str(data.Payload)
	if err != nil {
		return fmt.Errorf("HandleAdvance: failed to decode payload: %w", err)
	}

	var advancePayload AdvancePayload
	if err := json.Unmarshal([]byte(payloadStr), &advancePayload); err != nil {
		return fmt.Errorf("HandleAdvance: failed to unmarshal advance payload JSON: %w", err)
	}

	// Handle method based on request type
	switch advancePayload.Method {
	case "create":
		var studentData CreateStudentPayload
		if err := json.Unmarshal([]byte(advancePayload.Payload), &studentData); err != nil {
			return fmt.Errorf("HandleAdvance: failed to unmarshal 'payload' field: %w", err)
		}
		createStudent(studentData)
	case "sign_attendance":
		signAttendance(advancePayload.Payload)
	case "delete":
		deleteStudent(advancePayload.Payload)
	default:
		fmt.Println("Unsupported method in payload")
	}

	return nil
}

func createStudent(payload CreateStudentPayload) {
	newStudent := Circle{
		Name:          payload.Name,
		Age:           payload.Age,
		WalletAddress: payload.WalletAddress,
	}

	mu.Lock()
	allStudents = append(allStudents, newStudent)
	mu.Unlock()

	fmt.Printf("New student added: %s, %d, %s\n", newStudent.Name, newStudent.Age, newStudent.WalletAddress)
	rollups.SendNotice(&rollups.NoticeRequest{Payload: rollups.Str2Hex(fmt.Sprintf("New student created: %v", newStudent))})
}

func signAttendance(walletAddress string) {
	mu.Lock()
	defer mu.Unlock()

	for i, student := range allStudents {
		if strings.EqualFold(student.WalletAddress, walletAddress) {
			allStudents[i].TakeAttendance()
			fmt.Printf("%s has signed attendance.\n", student.Name)
			rollups.SendNotice(&rollups.NoticeRequest{Payload: rollups.Str2Hex("Attendance signed by: " + walletAddress)})
			return
		}
	}

	fmt.Printf("Student not found: %s\n", walletAddress)
	rollups.SendReport(&rollups.ReportRequest{Payload: rollups.Str2Hex("Student not found: " + walletAddress)})
}

func deleteStudent(walletAddress string) {
	mu.Lock()
	defer mu.Unlock()

	for i, student := range allStudents {
		if strings.EqualFold(student.WalletAddress, walletAddress) {
			allStudents = append(allStudents[:i], allStudents[i+1:]...)
			fmt.Printf("Student deleted: %s\n", walletAddress)
			rollups.SendNotice(&rollups.NoticeRequest{Payload: rollups.Str2Hex("Student deleted successfully: " + walletAddress)})
			return
		}
	}

	fmt.Printf("Student not found: %s\n", walletAddress)
	rollups.SendReport(&rollups.ReportRequest{Payload: rollups.Str2Hex("Student not found: " + walletAddress)})
}

func HandleInspect(data *rollups.InspectResponse) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("HandleInspect: failed to marshal JSON: %w", err)
	}
	infolog.Println("Received inspect request data:", string(dataBytes))

	payloadStr, err := rollups.Hex2Str(data.Payload)
	if err != nil {
		return fmt.Errorf("HandleInspect: failed to decode payload: %w", err)
	}

	routes := strings.Split(payloadStr, "/")

	switch routes[0] {
	case "all_students":
		mu.Lock()
		studentsList := allStudents
		mu.Unlock()

		response, _ := json.Marshal(map[string]interface{}{"students_list": studentsList})
		rollups.SendReport(&rollups.ReportRequest{Payload: rollups.Str2Hex(string(response))})
	case "student":
		walletAddress := routes[1]
		var student *Circle

		mu.Lock()
		for _, s := range allStudents {
			if strings.EqualFold(s.WalletAddress, walletAddress) {
				student = &s
				break
			}
		}
		mu.Unlock()

		if student != nil {
			response, _ := json.Marshal(student)
			rollups.SendReport(&rollups.ReportRequest{Payload: rollups.Str2Hex(string(response))})
		} else {
			rollups.SendReport(&rollups.ReportRequest{Payload: rollups.Str2Hex("Student not found")})
		}
	default:
		fmt.Println("Unsupported route in payload")
		rollups.SendReport(&rollups.ReportRequest{Payload: rollups.Str2Hex("Function not implemented")})
	}

	return nil
}

func Handler(response *rollups.FinishResponse) error {
	var err error

	switch response.Type {
	case "advance_state":
		data := new(rollups.AdvanceResponse)
		if err = json.Unmarshal(response.Data, data); err != nil {
			return fmt.Errorf("Handler: error unmarshalling advance: %w", err)
		}
		err = HandleAdvance(data)
	case "inspect_state":
		data := new(rollups.InspectResponse)
		if err = json.Unmarshal(response.Data, data); err != nil {
			return fmt.Errorf("Handler: error unmarshalling inspect: %w", err)
		}
		err = HandleInspect(data)
	}
	return err
}

func main() {
	finish := rollups.FinishRequest{Status: "accept"}

	for {
		infolog.Println("Sending finish")
		res, err := rollups.SendFinish(&finish)
		if err != nil {
			errlog.Panicln("Error: error making HTTP request: ", err)
		}
		infolog.Println("Received finish status", strconv.Itoa(res.StatusCode))

		if res.StatusCode == 202 {
			infolog.Println("No pending rollup request, trying again")
		} else {
			resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				errlog.Panicln("Error: could not read response body: ", err)
			}

			var response rollups.FinishResponse
			err = json.Unmarshal(resBody, &response)
			if err != nil {
				errlog.Panicln("Error: unmarshaling body:", err)
			}

			finish.Status = "accept"
			err = Handler(&response)
			if err != nil {
				errlog.Println(err)
				finish.Status = "reject"
			}
		}
	}
}
