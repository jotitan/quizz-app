package filer

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quizz-app/model"
	"github.com/quizz-app/music"
	"github.com/quizz-app/persistence"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type quizzesInfo struct{
	quizzes map[string]*model.LightQuizz
	path string
}

func newQuizzesInfo(path string)quizzesInfo{
	qi := quizzesInfo{path:path,quizzes: make(map[string]*model.LightQuizz)}

	if data,err := ioutil.ReadFile(path) ; err == nil {
		json.Unmarshal(data,&qi.quizzes)
	}
	return qi
}

func (qi quizzesInfo)addOrUpdateQuizz(quizz model.Quizz)error{
	qi.quizzes[quizz.Id] = model.NewLightQuizz(quizz)
	return qi.save()
}

func (qi quizzesInfo)save()error{
	data,err := json.Marshal(qi.quizzes)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(qi.path,data,777)
}

func (qi quizzesInfo) deleteQuizz(id string) error {
	delete(qi.quizzes,id)
	return qi.save()
}

func (qi quizzesInfo) nb() int {
	return len(qi.quizzes)
}

// Store each quizz in filer for now
type filerQuizzStorage struct {
	filer string
	quizzes quizzesInfo
	musicCutter music.Cutter
}

func (fgs filerQuizzStorage) getByFile(filename string) (model.Quizz, error) {
	quizz := model.Quizz{}
	data,err := ioutil.ReadFile(filename)
	if err != nil {
		return quizz,err
	}
	err = json.Unmarshal(data,&quizz)
	return quizz,err
}

func (fgs filerQuizzStorage) Get(id string) (model.Quizz, error) {
	return fgs.getByFile(fgs.generatePath(id))
}

func (fgs filerQuizzStorage) Create(dto model.QuizzDto) (string, error) {
	// Create quizz in a specific json file
	id := generateUniqueId()
	quizz := model.Quizz{Id: id,Name:dto.Name,Description:dto.Description,Questions: []model.Question{}}
	if err := fgs.addImageQuizz(dto,&quizz) ; err != nil {
		return "",err
	}
	return id,fgs.save(quizz)
}

func (fgs filerQuizzStorage) removeImageQuizz(quizz *model.Quizz)  error {
	if quizz.Image {
		path,err := fgs.generateImagePath(*quizz)
		if err != nil{
			return err
		}
		err = os.Remove(path)
		if err != nil{
			return err
		}
		quizz.Image = false
	}
	return nil
}

func (fgs filerQuizzStorage) GetCover(quizz model.Quizz) (io.ReadCloser, error) {
	if !quizz.Image {
		return nil,errors.New("no image")
	}
	path,err := fgs.generateImagePath(quizz)
	if err != nil {
		return nil,err
	}
	return os.Open(path)
}

func (fgs filerQuizzStorage) addImageQuizz(dto model.QuizzDto,quizz *model.Quizz)  error {
	if dto.ImageDescriptionHeader != nil {
		if !strings.HasSuffix(dto.ImageDescriptionHeader.Filename,".jpeg") && !strings.HasSuffix(dto.ImageDescriptionHeader.Filename,".jpg") {
			return errors.New("image must be jpeg")
		}
		path,err := fgs.generateImagePath(*quizz)
		if err != nil {
			return err
		}
		img,err := os.OpenFile(path,os.O_CREATE|os.O_TRUNC|os.O_RDWR,os.ModePerm)
		if err != nil {
			return err
		}
		defer img.Close()
		if _,err = io.Copy(img,dto.ImageDescription) ; err != nil {
			return err
		}
		quizz.Image = true
	}
	return nil
}

func (fgs filerQuizzStorage) Update(id string, dto model.QuizzDto) (string, error) {
	quizz,err := fgs.Get(id)
	if err != nil {
		return "",err
	}
	quizz.Name = dto.Name
	quizz.Description = dto.Description
	if dto.RemoveImage {
		if err := fgs.removeImageQuizz(&quizz) ; err != nil {
			return "",err
		}
	}
	if err := fgs.addImageQuizz(dto,&quizz) ; err != nil {
		return "",err
	}
	return id,fgs.save(quizz)
}

func (fgs filerQuizzStorage)save(quizz model.Quizz)error {
	path := fgs.generatePath(quizz.Id)
	data,_ := json.Marshal(quizz)

	err := ioutil.WriteFile(path,data,os.ModePerm)
	if err != nil {
		return err
	}
	return fgs.quizzes.addOrUpdateQuizz(quizz)
}

func (fgs filerQuizzStorage) AddQuestion(id string, question model.Question) error{
	quizz,err := fgs.Get(id)
	if err != nil {
		return err
	}
	if strings.EqualFold(question.Id,"") {
		question.Id = quizz.GetNextId()
		quizz.Questions = append(quizz.Questions, question)
	}else{
		fgs.updateQuestion(&quizz,question)
	}
	return fgs.save(quizz)
}

func (fgs filerQuizzStorage) updateQuestion(quizz *model.Quizz, question model.Question){
	// search position
	pos := -1
	for i,q := range quizz.Questions {
		if strings.EqualFold(q.Id,question.Id) {
			pos = i
			break
		}
	}
	if pos != -1 {
		// Keep music
		quizz.Questions[pos] = question
	}
}

func (fgs filerQuizzStorage)generatePath(id string)string{
	return filepath.Join(fgs.filer, fmt.Sprintf("quizz.%s.json",id))
}

func (fgs filerQuizzStorage) GetAll() []model.LightQuizz {
	quizzes := make([]model.LightQuizz,0,fgs.quizzes.nb())
	for _,q := range fgs.quizzes.quizzes {
		quizzes = append(quizzes,*q)
	}
	return quizzes
}

func (fgs filerQuizzStorage) ReadMusic(quizz model.Quizz,questionId string, writer io.Writer) error{
	pos,err := quizz.GetPositionQuestionById(questionId)
	if err != nil || strings.EqualFold("",quizz.Questions[pos].MusicPath){
		return err
	}
	data,err := ioutil.ReadFile(quizz.Questions[pos].MusicPath)
	if err != nil {
		return err
	}
	_,err = writer.Write(data)
	return err
}

func (fgs filerQuizzStorage) DeleteQuestion(quizz model.Quizz,questionId string)error {
	pos,err := quizz.GetPositionQuestionById(questionId)
	if err == nil {
		quizz.Questions = append(quizz.Questions[0:pos],quizz.Questions[pos+1:]...)
		return fgs.save(quizz)
	}
	return err
}

func (fgs filerQuizzStorage) DeleteMusic(idQuizz,idQuestion string) error {
	quizz,err := fgs.Get(idQuizz)
	if err != nil {
		return err
	}
	pos,err := quizz.GetPositionQuestionById(idQuestion)
	if err != nil {
		return err
	}
	if strings.EqualFold("",quizz.Questions[pos].MusicPath){
		return errors.New("impossible to find path")
	}
	return os.Remove(quizz.Questions[pos].MusicPath)
}

func (fgs filerQuizzStorage) StoreMusic(quizz model.Quizz, musicFile string, from, to int) (string, error) {
	path,err := fgs.generateMusicPath(quizz)
	if err != nil {
		return "",err
	}
	return path,fgs.musicCutter.Cut(musicFile,path,from,to)
}

func (fgs filerQuizzStorage)generateMusicPath(quizz model.Quizz)(string,error){
	folder := filepath.Join(fgs.filer,"assets",quizz.Id)
	if err := os.MkdirAll(folder,os.ModePerm) ; err != nil {
		return "",err
	}
	return filepath.Join(folder,fmt.Sprintf("%s.mp3",generateUniqueIdAsString())),nil
}

func (fgs filerQuizzStorage)generateImagePath(quizz model.Quizz)(string,error){
	folder := filepath.Join(fgs.filer,"assets",quizz.Id)
	if err := os.MkdirAll(folder,os.ModePerm) ; err != nil {
		return "",err
	}
	return filepath.Join(folder,fmt.Sprintf("cover_%s.jpeg",quizz.Id)),nil
}

func generateUniqueIdAsString()string{
	r := rand.Intn(10000)
	return fmt.Sprintf("%s-%d",time.Now().Format("20060102150405"),r)
}

func generateUniqueId()string{
	// Base on timestamp and random id
	r := rand.Intn(10000)
	value := fmt.Sprintf("%s-%d",time.Now().Format("20060102150405"),r)
	return base64.RawStdEncoding.EncodeToString([]byte(value))
}

func (fgs filerQuizzStorage) DeleteQuizz(quizz model.Quizz) error {
	// Remove all music of questions
	for _,q := range quizz.Questions {
		if !strings.EqualFold("",q.MusicPath){
			if err := os.Remove(q.MusicPath) ; err != nil {
				return err
			}
		}
	}
	err := os.Remove(fgs.generatePath(quizz.Id))
	if err != nil {
		return err
	}
	return fgs.quizzes.deleteQuizz(quizz.Id)
}

func NewQuizzStorage(filer, ffmpegPath string)persistence.QuizzStorage {
	return filerQuizzStorage{
		filer:filer,
		quizzes: newQuizzesInfo(filepath.Join(filer,"list_quizzes.json")),
		musicCutter: music.NewCutter(ffmpegPath),
	}
}
