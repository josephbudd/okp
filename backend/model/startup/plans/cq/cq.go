// Package cq creates a lesson plan of 5 lessons.
// A call sign is constructed from any 5 of the 26 letters of the alphabet but always begins with "K" or "W".
// Each lesson has
//  * no characters,
//  * no words,
//  * 1 sentence like "CQ CQ CQ DE WKRAP"
package cq

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/josephbudd/okp/shared/store"
	"github.com/josephbudd/okp/shared/store/record"
)

const (
	PlanName = "Call Signs Only"
	planDesc = "5 call sign only lessons. Each lesson is for 1 call sign."

	maxSentenceOptionsCount = 15
)

// Create creates the plan in the stores if it does not already exist.
func Create(
	stores *store.Stores,
	keyCodeCharMap map[string]*record.KeyCode,
	keyCodeWordMap map[string]*record.KeyCode,
) (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("cq.Create: %w", err)
		}
	}()

	// Find this plan in the stores.
	var plans []*record.Plan
	if plans, err = stores.Plan.GetAll(); err != nil {
		return
	}
	if len(plans) > 0 {
		for _, plan := range plans {
			if plan.Name == PlanName {
				// This lesson plan has already been created.
				return
			}
		}
	}

	// Create this plan.
	plan := record.NewPlan(PlanName, planDesc)
	if plan.Lessons, err = lessons(stores, keyCodeCharMap, keyCodeWordMap); err != nil {
		return
	}
	err = stores.Plan.Update(plan)
	return
}

func lessons(
	stores *store.Stores,
	keyCodeCharMap map[string]*record.KeyCode,
	keyCodeWordMap map[string]*record.KeyCode,
) (lessons []record.Lesson, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("lessons: %w", err)
		}
	}()

	lessons = make([]record.Lesson, 0, 1024)
	var lessonR record.Lesson
	var lessonNumber uint64
	// Lesson 1.
	lessonNumber++
	if lessonR, err = callSignLesson(lessonNumber, keyCodeCharMap, keyCodeWordMap); err != nil {
		return
	}
	lessons = append(lessons, lessonR)
	// Lesson 2.
	lessonNumber++
	if lessonR, err = callSignLesson(lessonNumber, keyCodeCharMap, keyCodeWordMap); err != nil {
		return
	}
	lessons = append(lessons, lessonR)
	// Lesson 3.
	lessonNumber++
	if lessonR, err = callSignLesson(lessonNumber, keyCodeCharMap, keyCodeWordMap); err != nil {
		return
	}
	lessons = append(lessons, lessonR)
	// Lesson 4.
	lessonNumber++
	if lessonR, err = callSignLesson(lessonNumber, keyCodeCharMap, keyCodeWordMap); err != nil {
		return
	}
	lessons = append(lessons, lessonR)
	// Lesson 5.
	lessonNumber++
	if lessonR, err = callSignLesson(lessonNumber, keyCodeCharMap, keyCodeWordMap); err != nil {
		return
	}
	lessons = append(lessons, lessonR)

	return
}

var aRune = uint64('A')

func callSignLesson(lessonNumber uint64, keyCodeCharMap, keyCodeWordMap map[string]*record.KeyCode) (lesson record.Lesson, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("callSignLesson: %w", err)
		}
	}()

	callSigns := make([][][]*record.KeyCode, maxSentenceOptionsCount)
	for i := 0; i < maxSentenceOptionsCount; i++ {
		callSigns[i] = randomCallSign(keyCodeCharMap, keyCodeWordMap)
	}

	lesson, err = record.BuildSentenceLesson(
		fmt.Sprintf("Lesson %d.", lessonNumber),
		"Call Signs.",
		lessonNumber,
		callSigns,
	)
	return
}

// randomCallSign returns a call sign.
func randomCallSign(keyCodeCharMap, keyCodeWordMap map[string]*record.KeyCode) (call [][]*record.KeyCode) {
	callsign := make([]*record.KeyCode, 5)
	// A call sign starts with "K" or "W"
	callsign[0] = randKW(keyCodeCharMap)
	for i := 1; i < 5; i++ {
		callsign[i] = randAZ(keyCodeCharMap)
	}
	call = make([][]*record.KeyCode, 5)
	cq := keyCodeWordMap["CQ"]
	de := keyCodeWordMap["DE"]
	call[0] = []*record.KeyCode{cq}
	call[1] = []*record.KeyCode{cq}
	call[2] = []*record.KeyCode{cq}
	call[3] = []*record.KeyCode{de}
	call[4] = callsign
	return
}

func randKW(keyCodeCharMap map[string]*record.KeyCode) (kw *record.KeyCode) {

	var err error
	var i int
	defer func() {
		if err != nil {
			err = fmt.Errorf("rangeKW: %w", err)
		}
		if i == 0 {
			kw = keyCodeCharMap["K"]
		} else {
			kw = keyCodeCharMap["W"]
		}
	}()

	var bigMax, bigI *big.Int
	bigMax = big.NewInt(2)
	if bigI, err = rand.Int(rand.Reader, bigMax); err != nil {
		return
	}
	i = int(bigI.Int64())
	return
}

func randAZ(keyCodeCharMap map[string]*record.KeyCode) (az *record.KeyCode) {

	var err error
	var i int64
	defer func() {
		// if err != nil then i == 0 by default.
		azRune := aRune + uint64(i)
		azString := string(rune(azRune))
		az = keyCodeCharMap[azString]
	}()

	var bigMax, bigI *big.Int
	bigMax = big.NewInt(26)
	if bigI, err = rand.Int(rand.Reader, bigMax); err != nil {
		return
	}
	i = bigI.Int64()
	return
}
