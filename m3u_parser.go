package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/dhowden/tag"
)

func main() {
	os.RemoveAll("img/")
	os.MkdirAll("img", 0700)
	os.RemoveAll("render/")
	os.MkdirAll("render", 0700)

	file, err := os.Open("m3u.m3u")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		renderFile(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func renderFile(file string) {
	fmt.Println("abrindo - ", file)
	data, err := os.Open(file)
	if err != nil {
		fmt.Println("skippando ", file, " err: ", err)
		return
	}

	m, err := tag.ReadFrom(data)
	if err != nil {
		fmt.Println("skippando ", file, " err: ", err)
		return
	}

	var artist = string(m.Artist())
	// ai q nojo odeio regex
	artist = regexp.MustCompile(`[^\x00-\x7F]+`).ReplaceAllString(artist, "")
	// nao sei regex, e a porra do gpt n sabe tbm entao vai ser replaceall separado
	artist = strings.ReplaceAll(artist, "\\", "")
	artist = strings.ReplaceAll(artist, "/", "")
	artist = strings.ReplaceAll(artist, ":", "")
	artist = strings.ReplaceAll(artist, "*", "")
	artist = strings.ReplaceAll(artist, "?", "")
	artist = strings.ReplaceAll(artist, "\"", "")
	artist = strings.ReplaceAll(artist, "<", "")
	artist = strings.ReplaceAll(artist, ">", "")
	artist = strings.ReplaceAll(artist, "|", "")
	artist = strings.TrimSpace(artist)
	if len(artist) == 0 {
		artist = "BLANK"
	}

	var path = regexp.MustCompile(`[^\x00-\x7F]+`).ReplaceAllString(artist+" - "+m.Title(), "")
	path = strings.ReplaceAll(path, "  ", " ")
	path = strings.ReplaceAll(path, "__", "_")

	if m.Picture() == nil {
		fmt.Println("erro pegando imagem")
		return
	}

	f, err := os.Create("img/" + path + ".jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	write, err := f.Write(m.Picture().Data)

	if err != nil {
		fmt.Println(err)
		return
	}

	// como q eu so passo o err eu n quero essa var
	write = write

	// var renderName = m.Artist() + " - " + m.Title()
	// var renderName = randSeq(10)
	cmd := exec.Command("ffmpeg", "-r", "1", "-loop", "1", "-i", "./img/"+path+".jpg", "-i", file, "-c:a", "copy", "-shortest", "-pix_fmt", "yuv420p", "render/"+path+".mp4")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Run()
	fmt.Println("out:", outb.String(), "err:", errb.String())
}
